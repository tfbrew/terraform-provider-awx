resource "aap_organization" "example" {
  name        = "example"
  description = "example"
}

resource "aap_inventory" "example" {
  name         = "example"
  description  = "example"
  organization = aap_organization.example.id
}

resource "aap_group" "group-example" {
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

resource "aap_host" "host-1" {
  name      = "host-1"
  inventory = aap_inventory.example.id
}


resource "aap_host" "host-2" {
  name      = "host-2"
  inventory = aap_inventory.example.id
}

resource "aap_group_host" "grp-host-link" {
  group_id = aap_group.group-example.id
  host_id  = aap_host.host-1.id
}

resource "aap_group_host" "grp-host-link-2" {
  group_id = aap_group.group-example.id
  host_id  = aap_host.host-2.id
}

resource "aap_group" "group-example-2" {
  name        = "group-name-example-2"
  description = "A second group example."
  inventory   = aap_inventory.example.id
}

resource "aap_group_host" "grp2-host-link" {
  group_id = aap_group.group-example-2.id
  host_id  = aap_host.host-2.id
}