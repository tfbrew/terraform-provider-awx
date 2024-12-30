resource "awx_jobtemplate_survey" "default" {
  question_description = "(String) Description of survey question."
  question_name        = "(String) Name of survey question."
  type                 = "text"
  variable             = "default"
}
