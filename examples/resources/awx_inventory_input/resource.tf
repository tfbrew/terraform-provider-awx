resource "awx_organization" "example" {
  name = "example"
}

resource "awx_inventory" "regular" {
  name         = "example"
  organization = awx_organization.example.id
}

resource "awx_constructed_inventory" "constructed" {
  name         = "example"
  organization = awx_organization.example.id
}

resource "awx_inventory_input" "example" {
  constructed_inventory_id = awx_constructed_inventory.constructed.id
  input_inventory_id       = awx_inventory.regular.id
}