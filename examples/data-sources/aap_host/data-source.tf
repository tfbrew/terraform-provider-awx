data "aap_host" "example-id" {
  id = "1"
}

data "aap_host" "example-name" {
  name      = "localhost"
  inventory = 1
}
