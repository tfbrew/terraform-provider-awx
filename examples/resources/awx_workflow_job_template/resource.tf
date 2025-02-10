resource "awx_workflow_job_template" "example" {
  allow_simultaneous       = false
  ask_inventory_on_launch  = false
  ask_labels_on_launch     = false
  ask_limit_on_launch      = false
  ask_scm_branch_on_launch = false
  ask_skip_tags_on_launch  = false
  ask_tags_on_launch       = false
  ask_variables_on_launch  = false
  description              = null
  extra_vars               = "---"
  inventory                = 1
  job_tags                 = null
  limit                    = null
  name                     = "example"
  organization             = 1
  scm_branch               = "main"
  skip_tags                = null
  survey_enabled           = false
  webhook_credential       = null
  webhook_service          = null
}
