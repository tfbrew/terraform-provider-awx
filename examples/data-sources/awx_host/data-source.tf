data "awx_host" "example-id" {
  id = "1"
}

data "awx_host" "example-name" {
  name      = "localhost"
  inventory = 1
}
