package firmafon

import (
	"fmt"
	"time"
)

type CallsService struct {
	*service
	Endpoint string
}

type Call struct {
	CallUUID    string           `json:"call_uuid"`
	CompanyID   int              `json:"company_id"`
	Endpoint    string           `json:"endpoint"`
	FromNumber  string           `json:"from_number"`
	ToNumber    string           `json:"to_number"`
	FromContact *CallFromContact `json:"from_contact"`
	ToContact   interface{}      `json:"to_contact"`
	Direction   string           `json:"direction"`
	StartedAt   time.Time        `json:"started_at"`
	AnsweredAt  time.Time        `json:"answered_at"`
	AnsweredBy  *CallAnsweredBy  `json:"answered_by"`
	EndedAt     time.Time        `json:"ended_at"`
	Status      string           `json:"status"`
}

type CallFromContact struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type CallAnsweredBy struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

type firmafonCalls struct {
	Calls []*Call `json:"calls"`
}

type firmafonCall struct {
	Call *Call `json:"call"`
}

// GetAll returns a slice of calls to or from one or more numbers
func (s *CallsService) GetAll(opt *CallsListOptions) ([]*Call, *Response, error) {
	url, err := addOptions(s.Endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	calls := &firmafonCalls{}
	resp, err := s.client.Do(req, &calls)
	if err != nil {
		return nil, resp, err
	}

	return calls.Calls, resp, nil
}

func (s *CallsService) Get(uuid string) (*Call, *Response, error) {
	url := s.Endpoint + fmt.Sprintf("/%s", uuid)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	call := &firmafonCall{}
	resp, err := s.client.Do(req, &call)
	if err != nil {
		return nil, resp, err
	}

	return call.Call, resp, nil
}
