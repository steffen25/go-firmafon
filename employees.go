package firmafon

import (
	"fmt"
	"time"
)

type EmployeesService service

type Employee struct {
	Admin            bool        `json:"admin"`
	CloakReception   interface{} `json:"cloak_reception"`
	CompanyID        int         `json:"company_id"`
	DndTimeoutAt     time.Time   `json:"dnd_timeout_at"`
	DoNotDisturb     bool        `json:"do_not_disturb"`
	EmployeeGroupIds []int       `json:"employee_group_ids"`
	Features         []string    `json:"features"`
	ID               int         `json:"id"`
	LivePresence     string      `json:"live_presence"`
	Name             string      `json:"name"`
	Number           string      `json:"number"`
	SpeedDial        struct {
		Digit int `json:"digit"`
	} `json:"speed_dial"`
}

type firmafonEmployees struct {
	Employees []*Employee `json:"employees"`
}

type firmafonEmployee struct {
	Employee *Employee `json:"employee"`
}

func (s *EmployeesService) All() ([]*Employee, *Response, error) {
	url := "employees"
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var emps *firmafonEmployees
	resp, err := s.client.Do(req, &emps)
	if err != nil {
		return nil, resp, err
	}

	return emps.Employees, resp, nil
}

func (s *EmployeesService) GetById(id int) (*Employee, *Response, error) {
	url := fmt.Sprintf("employees/%d", id)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var e *firmafonEmployee
	resp, err := s.client.Do(req, &e)
	if err != nil {
		return nil, resp, err
	}

	return e.Employee, resp, nil
}

func (s *EmployeesService) Update(e *Employee) (*Employee, *Response, error) {
	url := fmt.Sprintf("employees/%d", e.ID)
	em := firmafonEmployee{e}
	req, err := s.client.NewRequest("PUT", url, em)
	if err != nil {
		return nil, nil, err
	}

	emp := new(firmafonEmployee)
	resp, err := s.client.Do(req, &emp)
	if err != nil {
		return nil, resp, err
	}

	return emp.Employee, resp, nil
}
