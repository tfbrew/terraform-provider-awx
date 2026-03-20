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
  job_type = "run"
  name     = "example"
  project  = awx_organization.example.id
  playbook = "example.yml"
}

resource "awx_workflow_job_template" "example" {
  name         = "example"
  inventory    = awx_inventory.example.id
  organization = awx_organization.example.id
}

resource "awx_workflow_job_template_job_node" "awx_workflow_job_template_job_node" {
  unified_job_template     = awx_job_template.example.id
  workflow_job_template_id = awx_workflow_job_template.example.id
  inventory                = awx_inventory.example.id
}

resource "awx_workflow_job_template_job_node_credential" "example" {
  credential_ids = [1, 2, 3]
  id             = awx_workflow_job_template_job_node.example.id
}
