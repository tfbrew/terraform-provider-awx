# Terraform Provider AWX (AWX, AAP2.4, AAP2.5)

This is a terraform provider for AWX, AAP2.4 & AAP2.5 built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and based on the [Terraform Provider Scafolding Framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework).

If you find any bugs or have a feature request, please open a GitHub issue.

# Build Tags

This repo has 2 build tags: repoAWX and repoAAP. This is so that this code can be used for 2 differently named repositories & Terraform Providers.

The scaffold template's GNUmakefile has been altered to include refencing these tags. Therefore, use the `make` commands to self-compile instead of just using `go` raw. For example, run `make install` instead fo `go install`.