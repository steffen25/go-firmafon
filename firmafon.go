package firmafon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	apiVersion     = "2"
	defaultBaseURL = "https://app.firmafon.dk/api/v" + apiVersion + "/"
	mediaTypeJSON  = "application/json"
)

// A Client manages communication with the Firmafon API.
type Client struct {
	AccessToken string
	client      *http.Client
	BaseURL     *url.URL

	common service

	// Services used for talking to different parts of the Firmafon API
	Employees *EmployeesService
	Calls     *CallsService
}

type service struct {
	client *Client
}

type CallsListOptions struct {
	Endpoint        string `url:"endpoint"`
	Direction       string `url:"direction"`
	Status          string `url:"status"`
	Number          string `url:"number"`
	Limit           string `url:"limit"`
	StartedAtGtOrEq string `url:"started_at_gt_or_eq"`
	StartedAtLtOrEq string `url:"started_at_lt_or_eq"`
	EndedAtGtOrEq   string `url:"ended_at_gt_or_eq"`
	EndedAtLtOrEq   string `url:"ended_at_lt_or_eq"`
}

type Response struct {
	*http.Response
}

type ErrorResponse struct {
	Response *http.Response
	Status   string `json:"status"`  // error message returned from api
	Message  string `json:"message"` // error message returned from api
}

type AuthError ErrorResponse

func (r *AuthError) Error() string { return (*ErrorResponse)(r).Error() }

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message)
}

func NewClient(token string) *Client {
	httpClient := http.DefaultClient
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, AccessToken: token}
	c.common.client = c
	c.Employees = (*EmployeesService)(&c.common)
	callSrv := &CallsService{
		service:  &c.common,
		Endpoint: "calls",
	}
	c.Calls = callSrv

	return c
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	defer func() {
		resp.Body.Close()
	}()

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", mediaTypeJSON)
	}
	req.Header.Set("Accept", mediaTypeJSON)

	return req, nil
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	if r.StatusCode == 401 {
		data, err := ioutil.ReadAll(r.Body)
		if err == nil && data != nil {
			json.Unmarshal(data, errorResponse)
		}

		return (*AuthError)(errorResponse)
	}

	return errorResponse
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("access_token")) > 0 {
		params.Set("access_token", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}
