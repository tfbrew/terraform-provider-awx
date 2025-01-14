data "awx_organization" "example" {
  name = "Default"
}

resource "awx_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = data.awx_organization.example
}

resource "awx_inventory" "example-variables" {
  name         = "example_with_json_variables"
  description  = "Example with jsonencoded variables"
  organization = data.awx_organization.example
  variables = jsonencode(
    {
      foo = "bar"
    }
  )
}
