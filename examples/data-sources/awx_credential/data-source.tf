data "awx_credential" "example" {
  id = "1"
}

data "awx_credential" "example-name" {
  name            = "Demo Credential"
  credential_type = 1
  organization    = 1
}
