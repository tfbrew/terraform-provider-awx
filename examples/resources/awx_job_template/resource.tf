resource "awx_organization" "example" {
  name        = "example"
  description = "example"
}

resource "awx_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = awx_organization.example.id
}

resource "awx_job_template" "example" {
  job_type  = "run"
  name      = "test"
  inventory = awx_inventory.example.id
  project   = awx_organization.example.id
  playbook  = "test.yml"
}
