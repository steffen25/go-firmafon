package firmafon

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestCallsService_All(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/calls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{
		  "calls": [
			{
			  "call_uuid": "e54f5820-386d-0132-5bc3-14dae9edd21d",
			  "company_id": 1,
			  "endpoint": "Reception#1",
			  "from_number": "4512345678",
			  "to_number": "4571999999",
			  "from_contact": {
				"id": 1,
				"number": "4512345678",
				"name": "Kim Kontakt",
				"email": "kimkontakt@example.com"
			  },
			  "to_contact": null,
			  "direction": "incoming",
			  "started_at": "2014-03-21T13:59:04Z",
			  "answered_at": "2014-03-21T13:59:07Z",
			  "answered_by": {
				"id": 2,
				"name": "Karsten Kollega",
				"number": "4587654321"
			  },
			  "ended_at": "2014-03-21T13:59:59Z",
			  "status": "answered"
			}
		  ]
		}`)
	})

	calls, _, err := client.Calls.GetAll(nil)
	if err != nil {
		t.Errorf("Get all calls returned error: %v", err)
	}

	layout := time.RFC3339
	startedStr := "2014-03-21T13:59:04Z"
	started, _ := time.Parse(layout, startedStr)

	answeredStr := "2014-03-21T13:59:07Z"
	answered, _ := time.Parse(layout, answeredStr)

	endedStr := "2014-03-21T13:59:59Z"
	ended, _ := time.Parse(layout, endedStr)

	call := &Call{
		CallUUID:   "e54f5820-386d-0132-5bc3-14dae9edd21d",
		CompanyID:  1,
		Endpoint:   "Reception#1",
		FromNumber: "4512345678",
		ToNumber:   "4571999999",
		FromContact: &CallFromContact{
			ID:     1,
			Number: "4512345678",
			Name:   "Kim Kontakt",
			Email:  "kimkontakt@example.com",
		},
		ToContact:  nil,
		Direction:  "incoming",
		StartedAt:  started,
		AnsweredAt: answered,
		AnsweredBy: &CallAnsweredBy{
			ID:     2,
			Name:   "Karsten Kollega",
			Number: "4587654321",
		},
		EndedAt: ended,
		Status:  "answered",
	}

	want := []*Call{call}
	if !reflect.DeepEqual(calls, want) {
		t.Errorf("Get all calls returned %+v, want %+v", calls, want)
	}
}

func TestCallsService_All_Invalid_Request(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	baseURL, _ := url.Parse("https://app.firmafon.dk/api/v2")
	client.BaseURL = baseURL

	mux.HandleFunc("/calls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{
		  "calls": [
			{
			  "call_uuid": "e54f5820-386d-0132-5bc3-14dae9edd21d",
			  "company_id": 1,
			  "endpoint": "Reception#1",
			  "from_number": "4512345678",
			  "to_number": "4571999999",
			  "from_contact": {
				"id": 1,
				"number": "4512345678",
				"name": "Kim Kontakt",
				"email": "kimkontakt@example.com"
			  },
			  "to_contact": null,
			  "direction": "incoming",
			  "started_at": "2014-03-21T13:59:04Z",
			  "answered_at": "2014-03-21T13:59:07Z",
			  "answered_by": {
				"id": 2,
				"name": "Karsten Kollega",
				"number": "4587654321"
			  },
			  "ended_at": "2014-03-21T13:59:59Z",
			  "status": "answered"
			}
		  ]
		}`)
	})

	_, _, err := client.Calls.GetAll(nil)
	if err == nil {
		t.Error("Get all calls expected an error but got none")
	}
}

func TestCallsService_All_Invalid_URL(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{
		  "calls": [
			{
			  "call_uuid": "e54f5820-386d-0132-5bc3-14dae9edd21d",
			  "company_id": 1,
			  "endpoint": "Reception#1",
			  "from_number": "4512345678",
			  "to_number": "4571999999",
			  "from_contact": {
				"id": 1,
				"number": "4512345678",
				"name": "Kim Kontakt",
				"email": "kimkontakt@example.com"
			  },
			  "to_contact": null,
			  "direction": "incoming",
			  "started_at": "2014-03-21T13:59:04Z",
			  "answered_at": "2014-03-21T13:59:07Z",
			  "answered_by": {
				"id": 2,
				"name": "Karsten Kollega",
				"number": "4587654321"
			  },
			  "ended_at": "2014-03-21T13:59:59Z",
			  "status": "answered"
			}
		  ]
		}`)
	})

	_, _, err := client.Calls.GetAll(nil)
	if err == nil {
		t.Error("Get all calls expected an error but got none")
	}
}

func TestCallsService_All_Invalid_Options(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/calls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{
		  "calls": [
			{
			  "call_uuid": "e54f5820-386d-0132-5bc3-14dae9edd21d",
			  "company_id": 1,
			  "endpoint": "Reception#1",
			  "from_number": "4512345678",
			  "to_number": "4571999999",
			  "from_contact": {
				"id": 1,
				"number": "4512345678",
				"name": "Kim Kontakt",
				"email": "kimkontakt@example.com"
			  },
			  "to_contact": null,
			  "direction": "incoming",
			  "started_at": "2014-03-21T13:59:04Z",
			  "answered_at": "2014-03-21T13:59:07Z",
			  "answered_by": {
				"id": 2,
				"name": "Karsten Kollega",
				"number": "4587654321"
			  },
			  "ended_at": "2014-03-21T13:59:59Z",
			  "status": "answered"
			}
		  ]
		}`)
	})

	opts := &CallsListOptions{
		Endpoint:        "Reception#1",
		Direction:       "",
		Status:          "answered",
		Number:          "",
		Limit:           "",
		StartedAtGtOrEq: "",
		StartedAtLtOrEq: "",
		EndedAtGtOrEq:   "",
		EndedAtLtOrEq:   "",
	}
	// set a invalid endpoint
	client.Calls.Endpoint = ":"
	_, _, err := client.Calls.GetAll(opts)
	if err == nil {
		t.Errorf("Get all calls did not expect an error but got %v", err)
	}
}

func TestCallsService_All_Options(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/calls", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{
		  "calls": [
			{
			  "call_uuid": "e54f5820-386d-0132-5bc3-14dae9edd21d",
			  "company_id": 1,
			  "endpoint": "Reception#1",
			  "from_number": "4512345678",
			  "to_number": "4571999999",
			  "from_contact": {
				"id": 1,
				"number": "4512345678",
				"name": "Kim Kontakt",
				"email": "kimkontakt@example.com"
			  },
			  "to_contact": null,
			  "direction": "incoming",
			  "started_at": "2014-03-21T13:59:04Z",
			  "answered_at": "2014-03-21T13:59:07Z",
			  "answered_by": {
				"id": 2,
				"name": "Karsten Kollega",
				"number": "4587654321"
			  },
			  "ended_at": "2014-03-21T13:59:59Z",
			  "status": "answered"
			}
		  ]
		}`)
	})

	opts := &CallsListOptions{
		Endpoint:        "Reception#1",
		Direction:       "",
		Status:          "answered",
		Number:          "",
		Limit:           "",
		StartedAtGtOrEq: "",
		StartedAtLtOrEq: "",
		EndedAtGtOrEq:   "",
		EndedAtLtOrEq:   "",
	}
	_, _, err := client.Calls.GetAll(opts)
	if err != nil {
		t.Errorf("Get all calls did not expect an error but got %v", err)
	}
}

func TestCallsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/calls/e54f5820-386d-0132-5bc3-14dae9edd21d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{
		  "call": {
			  "call_uuid": "e54f5820-386d-0132-5bc3-14dae9edd21d",
			  "company_id": 1,
			  "endpoint": "Reception#1",
			  "from_number": "4512345678",
			  "to_number": "4571999999",
			  "from_contact": {
				"id": 1,
				"number": "4512345678",
				"name": "Kim Kontakt",
				"email": "kimkontakt@example.com"
			  },
			  "to_contact": null,
			  "direction": "incoming",
			  "started_at": "2014-03-21T13:59:04Z",
			  "answered_at": "2014-03-21T13:59:07Z",
			  "answered_by": {
				"id": 2,
				"name": "Karsten Kollega",
				"number": "4587654321"
			  },
			  "ended_at": "2014-03-21T13:59:59Z",
			  "status": "answered"
			}
		}`)
	})

	call, _, err := client.Calls.Get("e54f5820-386d-0132-5bc3-14dae9edd21d")
	if err != nil {
		t.Errorf("Get call returned error: %v", err)
	}

	layout := time.RFC3339
	startedStr := "2014-03-21T13:59:04Z"
	started, _ := time.Parse(layout, startedStr)

	answeredStr := "2014-03-21T13:59:07Z"
	answered, _ := time.Parse(layout, answeredStr)

	endedStr := "2014-03-21T13:59:59Z"
	ended, _ := time.Parse(layout, endedStr)

	want := &Call{
		CallUUID:   "e54f5820-386d-0132-5bc3-14dae9edd21d",
		CompanyID:  1,
		Endpoint:   "Reception#1",
		FromNumber: "4512345678",
		ToNumber:   "4571999999",
		FromContact: &CallFromContact{
			ID:     1,
			Number: "4512345678",
			Name:   "Kim Kontakt",
			Email:  "kimkontakt@example.com",
		},
		ToContact:  nil,
		Direction:  "incoming",
		StartedAt:  started,
		AnsweredAt: answered,
		AnsweredBy: &CallAnsweredBy{
			ID:     2,
			Name:   "Karsten Kollega",
			Number: "4587654321",
		},
		EndedAt: ended,
		Status:  "answered",
	}

	if !reflect.DeepEqual(call, want) {
		t.Errorf("Get call returned %+v, want %+v", call, want)
	}
}
