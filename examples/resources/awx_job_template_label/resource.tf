resource "awx_job_template_label" "example" {
  label_ids       = [1, 2, 3]
  job_template_id = 100
}
