resource "awx_job_templates_survey_spec" "default" {
  question_description = "(String) Description of survey question."
  question_name        = "(String) Name of survey question."
  type                 = "text"
  variable             = "default"
}
