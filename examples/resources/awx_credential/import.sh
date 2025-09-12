# The inputs field contains values that are not returned by the automation controller API and thus not available in state.
# The first plan/apply after import will result in a modification to the inputs so that the state can be updated.

terraform import awx_credential.example 5

# If you have an inputs_with_import field defined with secrets, you need to specify them here as comma-separated field name and secret-value pairs.
# The string at the end of this example command below would correlate to:
#   resource "awx_credential" "example_with_input" {
#      id = 5   
#      inputs_with_import = {
#         password = "12345"
#         token = "token,a1b2c3-d4e5"
#         non-secret = "do not include in import cli command"
#      }
#   }

terraform import awx_credential.example_with_input "5,password,12345,token,a1b2c3-d4e5"
