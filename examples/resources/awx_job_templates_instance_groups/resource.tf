resource "awx_job_templates_instance_groups" "default" {
  instance_group_ids = [1]
  job_template_id    = 1
}
