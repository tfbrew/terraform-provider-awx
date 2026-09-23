data "awx_inventory" "example" {
  id = "1"
}


data "awx_inventory" "example-name" {
  name         = "Demo Inventory"
  organization = 1
}
