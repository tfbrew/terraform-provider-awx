data "awx_credential_type" "example-id" {
  id = "5"
}

data "awx_credential_type" "example-name" {
  name = "Amazon Web Services"
  kind = "cloud"
}
