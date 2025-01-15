resource "awx_organization" "example" {
  name        = "example"
  description = "example"
}

resource "awx_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = awx_organization.example.id
}

resource "awx_inventory" "example-variables" {
  name         = "example_with_json_variables"
  description  = "Example with jsonencoded variables"
  organization = awx_organization.example.id
  variables = jsonencode(
    {
      foo = "bar"
    }
  )
}
