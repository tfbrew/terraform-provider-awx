resource "awx_jobtemplate_credential" "default" {
  credential_ids  = [1, 2, 3]
  job_template_id = 1
}
