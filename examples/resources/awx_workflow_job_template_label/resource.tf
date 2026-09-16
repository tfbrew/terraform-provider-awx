resource "awx_organization" "example" {
  name = "example-organization"
}

resource "awx_workflow_job_template" "example" {
  name         = "example-workflow-job-template"
  organization = awx_organization.example.id
}

resource "awx_label" "example-1" {
  name         = "example-label-1"
  organization = awx_organization.example.id
}

resource "awx_label" "example-2" {
  name         = "example-label-2"
  organization = awx_organization.example.id
}

resource "awx_workflow_job_template_label" "example" {
  label_ids                = [awx_label.example-1.id, awx_label.example-2.id]
  workflow_job_template_id = awx_workflow_job_template.example.id
}
