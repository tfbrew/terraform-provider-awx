resource "awx_workflow_job_template_node_always" "example_node_always" {
  node_id          = 201
  success_node_ids = [241, 914]
}