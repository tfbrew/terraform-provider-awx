resource "awx_jobtemplate_instance_groups" "default" {
  instance_group_ids = [1]
  job_template_id    = 1
}
