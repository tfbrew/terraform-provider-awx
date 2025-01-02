resource "awx_job_template" "example" {
  job_type  = "run"
  name      = "test"
  inventory = 1
  project   = 1
  playbook  = "test.yml"
}
