resource "aap_organization" "example" {
  name        = "example"
  description = "example"
}

resource "aap_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = aap_organization.example.id
}

resource "aap_inventory" "example-variables" {
  name         = "example_with_json_variables"
  description  = "Example with jsonencoded variables"
  organization = aap_organization.example.id
  variables = jsonencode(
    {
      foo = "bar"
    }
  )
}
