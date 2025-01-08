resource "awx_workflow_job_template_node_success" "example_node_success" {
  node_id          = 201
  success_node_ids = [241, 914]
}