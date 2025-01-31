resource "awx_job_template" "example" {
  job_type  = "run"
  name      = "test"
  inventory = 1
  project   = 1
  playbook  = "test.yml"
}

resource "awx_job_template_survey_spec" "example" {
  description = "example description"
  id          = awx_job_template.example.id
  name        = ""
  spec = [
    {
      choices              = ["stop", "start", "status", "restart"]
      default              = "status"
      max                  = 1024
      min                  = 0
      question_description = "example question 1"
      question_name        = "example_question_1"
      required             = true
      type                 = "multiplechoice"
      variable             = "examplevar1"
    },
    {
      default              = jsonencode(15)
      max                  = 1024
      min                  = 1
      question_description = "example question 2"
      question_name        = "Example question 1"
      required             = true
      type                 = "integer"
      variable             = "example_2_var"
    },
  ]
}
