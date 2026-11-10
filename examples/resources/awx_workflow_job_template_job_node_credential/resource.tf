resource "awx_workflow_job_template_job_node" "example" {
  inventory                = 100
  unified_job_template     = 1001
  workflow_job_template_id = 1002
}

resource "awx_workflow_job_template_job_node_credential" "example" {
  id             = awx_workflow_job_template_job_node.example.id
  credential_ids = [1, 2, 3]
}
