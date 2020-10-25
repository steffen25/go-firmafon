# go-firmafon
Go library for accessing the Firmafon API.

[![Build Status](https://github.com/steffen25/go-firmafon/workflows/golangci-lint/badge.svg)](https://github.com/steffen25/go-firmafon/actions?query=workflow%3Agolangci-lint)
[![Test suite Status](https://github.com/steffen25/go-firmafon/workflows/test-suite/badge.svg)](https://github.com/steffen25/go-firmafon/actions?query=workflow%3Atest-suite)
[![Go Report Card](https://goreportcard.com/badge/github.com/steffen25/go-firmafon)](https://goreportcard.com/report/github.com/steffen25/go-firmafon)
[![codecov](https://codecov.io/gh/steffen25/go-firmafon/branch/master/graph/badge.svg)](https://codecov.io/gh/steffen25/go-firmafon)

## Installation
`go get github.com/steffen25/go-firmafon`

## Usage ##

```go
import "github.com/steffen25/go-firmafon"
```

Construct a new Firmafon client using an access token which you can generate here [Generate Token](https://app.firmafon.dk/account/authorized_applications)
```go
client := firmafon.NewClient("token")
```
#### List all employees ####

```go
client := firmafon.NewClient("token")

// list all employees for your organization
users, _, err := client.Employees.All()
if err != nil {
	// Handle error
}

// print each employee's name
for _, u := range users {
	fmt.Println(u.Name)
}
```
