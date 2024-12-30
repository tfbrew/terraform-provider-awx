resource "awx_job_template" "default" {
  job_type  = "run"
  name      = "test"
  inventory = 1
  project   = 1
  playbook  = "test.yml"
}
