resource "awx_job_templates_credentials" "default" {
  credential_ids  = [1, 2, 3]
  job_template_id = 1
}
