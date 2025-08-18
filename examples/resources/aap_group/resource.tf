resource "aap_organization" "example" {
  name        = "example"
  description = "example"
}

resource "aap_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = aap_organization.example.id
}

resource "aap_group" "example-variables" {
  name        = "group-name-example"
  description = "Example with jsonencoded variables."
  inventory   = aap_inventory.example.id
  variables = jsonencode(
    {
      foo = "bar"
      baz = "qux"
    }
  )
}
