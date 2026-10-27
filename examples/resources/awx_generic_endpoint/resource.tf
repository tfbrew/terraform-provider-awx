resource "awx_generic_endpoint" "test-credential-type" {
  api_path     = "credential_types"
  api_endpoint = "controller"
  data_json = jsonencode(
    {
      "description" : "",
      "injectors" : {},
      "inputs" : { "fields" : [] },
      "kind" : "cloud",
      "name" : "test-credential-type"
    }
  )
  ignore_read = ["created", "created_by", "id", "managed", "modified", "modified_by", "namespace", "related", "summary_fields", "type", "url"]
}
