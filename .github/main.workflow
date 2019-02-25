workflow "New workflow" {
  on = "push"
  resolves = [
    "Call httpbin",
    "my action",
  ]
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

action "HTTP client" {
  uses = "swinton/httpie.action@8ab0a0e926d091e0444fcacd5eb679d2e2d4ab3d"
}
