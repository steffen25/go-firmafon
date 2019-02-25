workflow "New workflow" {
  on = "push"
  resolves = ["go"]
}

action "Call httpbin" {
  uses = "swinton/httpie.action@master"
  args = ["POST", "httpbin.org/anything", "hello=world"]
}

action "go" {
  uses = "go"
  needs = ["Call httpbin"]
  args = "test"
}
