package firmafon

import (
	"fmt"
	"time"
)

type EmployeesService service

type Employee struct {
	Admin            bool        `json:"admin,omitempty"`
	CloakReception   interface{} `json:"cloak_reception,omitempty"`
	CompanyID        int         `json:"company_id,omitempty"`
	DndTimeoutAt     *time.Time   `json:"dnd_timeout_at,omitempty"`
	DoNotDisturb     bool        `json:"do_not_disturb,omitempty"`
	EmployeeGroupIds []int       `json:"employee_group_ids,omitempty"`
	Features         []string    `json:"features,omitempty"`
	ID               int         `json:"id,omitempty"`
	LivePresence     string      `json:"live_presence,omitempty"`
	Name             string      `json:"name,omitempty"`
	Number           string      `json:"number,omitempty"`
	SpeedDial        *SpeedDial	 `json:"speed_dial,omitempty"`
}

type SpeedDial struct {
	Digit int `json:"digit,omitempty"`
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
