data "awx_instance_group" "example-id" {
  id = "2"
}

data "awx_instance_group" "example-name" {
  name = "controlplane"
}
