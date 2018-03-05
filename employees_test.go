package firmafon

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestEmployeesService_All(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employees":[{"id": 1}, {"id": 2}]}`)
	})

	emps, _, err := client.Employees.All()
	if err != nil {
		t.Errorf("Get all employees returned error: %v", err)
	}

	want := []*Employee{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(emps, want) {
		t.Errorf("Get all employees returned %+v, want %+v", emps, want)
	}
}

func TestEmployeesService_All_InvalidURL(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employees":[{"id": 1}, {"id": 2}]}`)
	})

	_, _, err := client.Employees.All()
	if err == nil {
		t.Error("Get all employees expected error to be returned but gone none")
	}
}

func TestEmployeesService_All_InvalidRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	baseURL, _ := url.Parse("https://app.firmafon.dk/api/v2")
	client.BaseURL = baseURL

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employees":[{"id": 1}, {"id": 2}]}`)
	})

	_, _, err := client.Employees.All()
	if err == nil {
		t.Error("Get all employees expected error to be returned but gone none")
	}
}

func TestEmployeesService_GetById(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/employees/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1}}`)
	})

	emp, _, err := client.Employees.GetById(1)
	if err != nil {
		t.Errorf("Get employee by id returned error: %v", err)
	}

	want := &Employee{ID: 1}
	if !reflect.DeepEqual(emp, want) {
		t.Errorf("Get employee by id returned %+v, want %+v", emp, want)
	}
}

func TestEmployeesService_GetById_InvalidURL(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1}}`)
	})

	_, _, err := client.Employees.GetById(1)
	if err == nil {
		t.Error("Get employee by id expected error to be returned but gone none")
	}
}

func TestEmployeesService_GetById_InvalidRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	baseURL, _ := url.Parse("https://app.firmafon.dk/api/v2")
	client.BaseURL = baseURL

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1}}`)
	})

	_, _, err := client.Employees.GetById(1)
	if err == nil {
		t.Error("Get employee by id expected error to be returned but gone none")
	}
}

func TestEmployeesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/employees/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1, "name": "Steffen"}}`)
	})
	emp := &Employee{ID: 1, Name: "John"}
	emp.Name = "Steffen"

	uEmp, _, err := client.Employees.Update(emp)
	if err != nil {
		t.Errorf("Update employee returned error: %v", err)
	}

	want := &Employee{ID: 1, Name: "Steffen"}
	if !reflect.DeepEqual(uEmp, want) {
		t.Errorf("Update employee returned %+v, want %+v", uEmp, emp)
	}
}

func TestEmployeesService_Update_InvalidURL(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1, "name": "Steffen"}}`)
	})
	emp := &Employee{ID: 1, Name: "John"}
	emp.Name = "Steffen"

	_, _, err := client.Employees.Update(emp)
	if err == nil {
		t.Error("Update employee expected error to be returned but gone none")
	}
}

func TestEmployeesService_Update_InvalidRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	baseURL, _ := url.Parse("https://app.firmafon.dk/api/v2")
	client.BaseURL = baseURL

	mux.HandleFunc("/employees/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1, "name": "Steffen"}}`)
	})
	emp := &Employee{ID: 1, Name: "John"}
	emp.Name = "Steffen"

	_, _, err := client.Employees.Update(emp)
	if err == nil {
		t.Error("Update employee expected error to be returned but gone none")
	}
}

func TestEmployeesService_Authenticated(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1, "name": "John"}}`)
	})

	emp, _, err := client.Employees.Authenticated()
	if err != nil {
		t.Errorf("Update employee returned error: %v", err)
	}

	want := &Employee{ID: 1, Name: "John"}
	if !reflect.DeepEqual(emp, want) {
		t.Errorf("Get authenticated employee returned %+v, want %+v", emp, want)
	}
}

func TestEmployeesService_Authenticated_InvalidURL(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(":", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1, "name": "John"}}`)
	})

	_, _, err := client.Employees.Authenticated()
	if err == nil {
		t.Error("Get authenticated employee expected error to be returned but gone none")
	}
}

func TestEmployeesService_Authenticated_InvalidRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	baseURL, _ := url.Parse("https://app.firmafon.dk/api/v2")
	client.BaseURL = baseURL

	mux.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeJSON)
		fmt.Fprint(w, `{"employee":{"id": 1, "name": "John"}}`)
	})

	_, _, err := client.Employees.Authenticated()
	if err == nil {
		t.Error("Get authenticated employee expected error to be returned but gone none")
	}
}
