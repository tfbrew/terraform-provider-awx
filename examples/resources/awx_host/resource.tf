resource "awx_organization" "example" {
  name        = "example"
  description = "example"
}

resource "awx_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = awx_organization.example.id
}

resource "awx_host" "example-variables" {
  name        = "localhost"
  description = "Example with jsonencoded variables for localhost"
  inventory   = awx_inventory.example.id
  variables = jsonencode(
    {
      foo = "bar"
    }
  )
}
