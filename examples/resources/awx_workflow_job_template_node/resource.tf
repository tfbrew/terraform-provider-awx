resource "awx_workflow_job_template_node" "example_node" {
  all_parents_must_converge = false
  diff_mode                 = false
  identifier                = "0ddd9e82-e447-4ef2-9c96-20932fef4456"
  extra_data = jsonencode({
    current_version = "101"
    update_zipfile  = "example.zip"
    variable1       = 1
  })
  inventory                = 100
  job_tags                 = null
  job_type                 = null
  limit                    = null
  scm_branch               = null
  skip_tags                = null
  unified_job_template     = 1001
  verbosity                = null
  workflow_job_template_id = 1002
}