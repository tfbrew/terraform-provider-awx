resource "aap_organization" "example" {
  name        = "example"
  description = "example"
}

resource "aap_project" "example-git" {
  name         = "example_git"
  organization = aap_organization.example.id
  scm_type     = "git"
  scm_url      = "git@github.com:user/repo.git"
}

resource "aap_project" "example-svn" {
  name         = "example_svn"
  organization = aap_organization.example.id
  scm_type     = "svn"
  scm_url      = "svn://<your_ip>/<repository_name>"
}

resource "aap_project" "example-archive" {
  name         = "example_archive"
  organization = aap_organization.example.id
  scm_type     = "archive"
  scm_url      = "https://github.com/user/repo"
}

resource "aap_project" "example-manual" {
  name         = "example_manual"
  organization = aap_organization.example.id
  scm_type     = "manual"
  local_path   = "directory/on/controller"
}

data "aap_credential" "example-insights" {
  id = "1"
}

resource "aap_project" "example-insights" {
  name         = "example_insights"
  organization = aap_organization.example.id
  scm_type     = "insights"
  credential   = data.aap_credential.example_insights.id
}
