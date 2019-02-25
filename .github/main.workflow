workflow "New workflow" {
  on = "push"
  resolves = ["Call httpbin"]
}

action "Call httpbin" {
  uses = "swinton/httpie.action@master"
  args = ["POST", "httpbin.org/anything", "hello=world"]
}
