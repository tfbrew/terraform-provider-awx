resource "awx_organization" "example" {
  name        = "tf_example"
  description = "example"
}

// Example machine credential

// Inputs options for machine credentials:
// username, password, ssh_key_data, ssh_public_key_data, become_method, become_password, ssh_key_unlock
// Setting a value to "ASK" is equal to choosing "Prompt at Launch"

data "awx_credential_type" "machine" {
  name = "Machine"
  kind = "ssh"
}

resource "awx_credential" "example-machine" {
  name            = "example_machine"
  organization    = awx_organization.example.id
  credential_type = data.awx_credential_type.machine.id
  inputs = jsonencode({
    "username" : "awx",
    "password" : "test1234", // code should not contain secrets, example only
    "become_method" : "sudo",
    "become_password" : "ASK" // ASK = Prompt at Launch checkbox
  })
}

// Example source control credential

// Inputs options for source control credentials:
// username, password, ssh_key_data, ssh_key_unlock
// Setting a value to "ASK" is equal to choosing "Prompt at Launch"

data "awx_credential_type" "source-control" {
  name = "Source Control"
  kind = "scm"
}

resource "awx_credential" "example-source-control" {
  name            = "example_machine"
  organization    = awx_organization.example.id
  credential_type = data.awx_credential_type.source-control.id
  inputs = jsonencode({
    "username" : "awx",
    "ssh_key_data" : file("${path.module}/id_rsa"), // code should not contain secrets, example only
    "ssh_key_unlock" : "test1234"                   // code should not contain secrets, example only
  })
}

// Example container registry credential

// Inputs options for source control credentials:
// host, username, password, verify_ssl
// Setting a value to "ASK" is equal to choosing "Prompt at Launch"

data "awx_credential_type" "container-registry" {
  name = "Container Registry"
  kind = "registry"
}

resource "awx_credential" "example-container-registry" {
  name            = "example_container_registry"
  organization    = awx_organization.example.id
  credential_type = data.awx_credential_type.container-registry.id
  inputs = jsonencode({
    "host" : "quay.io",
    "username" : "test",
    "password" : "test1234", // code should not contain secrets, example only
    "verify_ssl" : true
  })
}
