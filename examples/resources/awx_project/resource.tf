data "awx_organization" "example" {
  name = "Default"
}

resource "awx_project" "example-git" {
  name         = "example_git"
  organization = data.awx_organization.example
  scm_type     = "git"
  scm_url      = "git@github.com:user/repo.git"
}

resource "awx_project" "example-svn" {
  name         = "example_svn"
  organization = data.awx_organization.example
  scm_type     = "svn"
  scm_url      = "svn://<your_ip>/<repository_name>"
}

resource "awx_project" "example-manual" {
  name         = "example_manual"
  organization = data.awx_organization.example
  scm_type     = "manual"
  local_path   = "playbook_directory"
}
