workflow "New workflow" {
  on = "push"
  resolves = [
    "Call httpbin",
    "Go test",
  ]
}

action "Call httpbin" {
  uses = "swinton/httpie.action@master"
  args = ["POST", "httpbin.org/anything", "hello=world"]
}

action "Go test" {
  uses = "./my-action/"
  needs = ["Call httpbin"]
  args = "test -v"
}
