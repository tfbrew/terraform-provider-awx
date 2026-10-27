resource "awx_organization" "example" {
  name        = "Example Organization"
  description = "An example organization"
}

resource "awx_inventory" "example" {
  name         = "Example Inventory"
  description  = "An example inventory"
  organization = awx_organization.example.id
}

resource "awx_constructed_inventory" "example" {
  name         = "Example Constructed Inventory"
  description  = "An example constructed inventory"
  organization = awx_organization.example.id
  limit        = "dev,&no_mgmt"
  source_vars  = <<-EOT
---
plugin: ansible.builtin.constructed
strict: true
groups:
  dev: az_env_tag == "Development"
  no_mgmt: not ("mgmt" in inventory_hostname)
EOT
}

resource "awx_inventory_input" "example" {
  constructed_inventory_id = awx_constructed_inventory.example.id
  input_inventory_id       = awx_inventory.example.id
}  