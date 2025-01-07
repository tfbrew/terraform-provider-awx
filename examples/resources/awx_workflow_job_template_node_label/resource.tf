resource "awx_workflow_job_template_node_label" "example_node_label" {
  node_id   = 201
  label_ids = [322, 121]
}