resource "awx_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = 0
}

resource "awx_inventory" "example-variables" {
  name         = "example_with_json_variables"
  description  = "Example with jsonencoded variables"
  organization = 0
  variables = jsonencode(
    {
      alex = "test"
    }
  )
}
