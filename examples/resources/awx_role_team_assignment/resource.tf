resource "awx_organization" "example" {
  name = "Example Organization"
}
resource "awx_role_definition" "example" {
  name         = "Example Role Definition"
  description  = "Example role definition"
  content_type = "shared.organization"
  permissions  = ["shared.member_organization", "shared.view_organization"]
}
resource "awx_team" "example" {
  name         = "Example Team"
  organization = awx_organization.example.id
  description  = "Example team description"
}
resource "awx_role_user_assignment" "test" {
  object_id       = awx_organization.example.id
  role_definition = awx_role_definition.example.id
  user            = awx_user.example.id
}
