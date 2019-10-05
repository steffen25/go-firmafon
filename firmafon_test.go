package firmafon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	baseURLPath = "/api/v2"
)

func setup() (c *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	c = NewClient("")
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	c.BaseURL = url

	return c, mux, server.URL, server.Close
}

func TestNewClient(t *testing.T) {
	token := "123"
	c := NewClient(token)

	if actual, expected := c.BaseURL.String(), defaultBaseURL; actual != expected {
		t.Errorf("NewClient BaseURL is %v, want %v", actual, expected)
	}

	if actual, expected := c.AccessToken, token; actual != expected {
		t.Errorf("NewClient AccessToken is %v, want %v", actual, expected)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient("")
	inURL, outURL := "users", defaultBaseURL+"users"
	inBody, outBody := &Employee{Name: "Steffen"}, `{"name":"Steffen"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body is %v, want %v", inBody, got, want)
	}
}

func TestNewRequest_badBaseURL(t *testing.T) {
	c := NewClient("")
	c.BaseURL, _ = url.Parse("https://app.firmafon.dk/api/v2")
	_, err := c.NewRequest("GET", "users", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient("")
	_, err := c.NewRequest("GET", ":", nil)
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewRequest_emptyBody(t *testing.T) {
	c := NewClient("")
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("HTTP request contains a non-nil Body")
	}
}

func TestNewRequest_invalidMethod(t *testing.T) {
	c := NewClient("")
	req, err := c.NewRequest("🍆", ".", nil)
	if err == nil {
		t.Error("Expected error to be returned")
	}
	if req != nil {
		t.Fatalf("Expected request to be nil")
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient("")
	type MyType struct {
		Test map[interface{}]interface{}
	}

	req, err := c.NewRequest("POST", ".", &MyType{})

	if err == nil {
		t.Error("Expected error to be returned")
	}

	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}

	if req != nil {
		t.Fatalf("Expected request to be nil")
	}
}

func TestCheckResponse(t *testing.T) {
	tests := []struct {
		res          *http.Response
		errorMessage string
		errorStatus  string
		wantError    bool
	}{
		{res: &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"success": true}`)),
		},
			wantError: false,
		},
		{res: &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"success": false}`)),
		},
			wantError: true,
		},
	}

	for _, test := range tests {
		err := CheckResponse(test.res)
		want := &ErrorResponse{
			Response: test.res,
			Message:  test.errorMessage,
			Status:   test.errorStatus,
		}

		if err == nil && test.wantError {
			t.Errorf("Expected error response.")
		}

		if err != nil && !reflect.DeepEqual(err, want) {
			t.Errorf("Error = %#v, want %#v", err, want)
		}
	}
}

func TestCheckResponse_AuthError(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusUnauthorized,
		Body:       ioutil.NopCloser(strings.NewReader(`{"success": false, "message": "unauthorized", "status": "401 unauthorized"}`)),
	}
	err := CheckResponse(res).(*AuthError)

	want := &AuthError{
		Response: res,
		Message:  "unauthorized",
		Status:   "401 unauthorized",
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestNewResponse(t *testing.T) {
	tests := []struct {
		res *http.Response
	}{
		{res: &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusOK,
			Body:       nil,
		},
		},
		{res: &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"success": false}`)),
		},
		},
	}

	for _, test := range tests {
		want := newResponse(test.res)
		if !reflect.DeepEqual(test.res, want.Response) {
			t.Errorf("Error = %#v, want %#v", test.res, want.Response)
		}
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		url, want string
	}{
		{"users/1337", "users/1337"},
		{"users/1337?access_token=secret", "users/1337?access_token=REDACTED"},
		// The Encode call will sort the params therefore since a comes before i it will be the first param
		{"users?id=1&access_token=secret", "users?access_token=REDACTED&id=1"},
		{":", ":"},
	}

	for _, test := range tests {
		inURL, _ := url.Parse(test.url)
		want, _ := url.Parse(test.want)

		if got := sanitizeURL(inURL); !reflect.DeepEqual(got, want) {
			t.Errorf("sanitizeURL(%v) returned %v, want %v", test.url, got, want)
		}
	}
}

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		res          *http.Response
		errorMessage string
		errorStatus  string
	}{
		{res: &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"success": true}`)),
		},
		},
		{res: &http.Response{
			Request:    &http.Request{},
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"success": false}`)),
		},
		},
		{res: &http.Response{
			Request: &http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "example.com",
					Path:     "foo/bar",
					RawQuery: "access_token=secret&id=1",
				},
			},
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"success": false, "message": "unauthorized", "status": "401 unauthorized"}`)),
		},
			errorMessage: "unauthorized",
			errorStatus:  "401 unauthorized",
		},
	}

	for _, test := range tests {
		err := &ErrorResponse{Response: test.res, Message: test.errorMessage, Status: test.errorStatus}
		want := fmt.Sprintf("%v %v: %d %v",
			test.res.Request.Method, sanitizeURL(test.res.Request.URL),
			test.res.StatusCode, test.errorMessage)

		if !reflect.DeepEqual(err.Error(), want) {
			t.Errorf("Error = %#v, want %#v", err, want)
		}
	}
}

func TestErrorResponse_AuthError(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := &AuthError{Response: res, Message: "unauthorized", Status: "401 unauthorized"}
	want := fmt.Sprintf("%v %v: %d %v",
		res.Request.Method, sanitizeURL(res.Request.URL),
		res.StatusCode, err.Message)

	if !reflect.DeepEqual(err.Error(), want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type User struct {
		Name string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"Name":"John"}`)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	body := new(User)
	client.Do(req, body)

	want := &User{"John"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_reqError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	req := &http.Request{}
	_, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".", nil)
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestDo_ioWriter(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test")
	})

	var b bytes.Buffer

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(req, &b)
	got := b.String()
	want := "test"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response body = %v, want %v", got, want)
	}

	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestAddOptions_Parse(t *testing.T) {
	_, err := addOptions(":", nil)
	if err == nil {
		t.Fatal("Addoptions Parse did not return an error")
	}
}

func TestAddOptions_Values(t *testing.T) {
	type test interface{}
	var tester test
	_, err := addOptions("/", &tester)
	if err == nil {
		t.Fatal("Addoptions Values did not return an error")
	}
}
