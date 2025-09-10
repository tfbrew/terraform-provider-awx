resource "awx_role_definition" "example" {
  name         = "Example Credential Admin"
  description  = "Example credential admin role definition"
  content_type = "awx.credential"
  permissions  = ["awx.change_credential", "awx.delete_credential", "awx.use_credential", "awx.view_credential"]
}
