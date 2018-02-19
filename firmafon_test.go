package firmafon

import (
	"net/http"
	"fmt"
	"os"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"net/url"
)

const (
	baseURLPath = "/api/v2"
)

func setup() (c *Client, mux *http.ServeMux, serverURL string, teardown func())  {
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

func TestNewRequest(t *testing.T)  {
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

