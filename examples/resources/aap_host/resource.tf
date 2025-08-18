resource "aap_organization" "example" {
  name        = "example"
  description = "example"
}

resource "aap_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = aap_organization.example.id
}

resource "aap_host" "example-variables" {
  name        = "localhost"
  description = "Example with jsonencoded variables for localhost"
  inventory   = aap_inventory.example.id
  variables = jsonencode(
    {
      foo = "bar"
    }
  )
}
