resource "awx_organization" "example" {
  name        = "example"
  description = "example"
}

resource "awx_project" "example-git" {
  name         = "example_git"
  organization = awx_organization.example.id
  scm_type     = "git"
  scm_url      = "git@github.com:user/repo.git"
}

resource "awx_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = awx_organization.example.id
}

resource "awx_inventory_source" "github_inventory_source" {
  name             = "example"
  inventory        = awx_inventory.example.id
  source           = "scm"
  source_project   = awx_project.example_git.id
  source_path      = "inventory"
  overwrite        = true
  overwrite_vars   = true
  update_on_launch = true
}
