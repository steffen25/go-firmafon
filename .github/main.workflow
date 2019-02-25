workflow "New workflow" {
  on = "push"
  resolves = ["go"]
}

action "Call httpbin" {
  uses = "swinton/httpie.action@master"
  args = ["POST", "httpbin.org/anything", "hello=world"]
}

action "my action" {
  uses = "./my-action/"
  needs = ["Call httpbin"]
  args = "test"
}
