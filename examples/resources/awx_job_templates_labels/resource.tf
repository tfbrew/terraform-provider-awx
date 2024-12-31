resource "awx_job_templates_labels" "default" {
  label_ids       = [1, 2, 3]
  job_template_id = 1
}
