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

resource "aap_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = aap_organization.example.id
}

resource "aap_inventory_source" "github_inventory_source" {
  name             = "example"
  inventory        = aap_inventory.example.id
  source           = "scm"
  source_project   = aap_project.example_git.id
  source_path      = "inventory"
  overwrite        = true
  overwrite_vars   = true
  update_on_launch = true
}
