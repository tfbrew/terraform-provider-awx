resource "awx_job_template_credential" "example" {
  credential_ids  = [1, 2, 3]
  job_template_id = 100
}
