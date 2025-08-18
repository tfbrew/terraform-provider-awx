resource "aap_organization" "example" {
  name        = "example"
  description = "example"
}

resource "aap_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = aap_organization.example.id
}

resource "aap_job_template" "example" {
  job_type  = "run"
  name      = "test"
  inventory = aap_inventory.example.id
  project   = aap_organization.example.id
  playbook  = "test.yml"
}
