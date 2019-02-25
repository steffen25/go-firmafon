workflow "New workflow" {
  on = "push"
  resolves = ["HTTP client"]
}

action "Call httpbin" {
  uses = "swinton/httpie.action@master"
  args = ["POST", "httpbin.org/anything", "hello=world"]
}
