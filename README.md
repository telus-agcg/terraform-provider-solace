# Solace Terraform Provider

## Introduction

This is a Terraform Provider for the Solace message broker. It allows you to configure various resources in Solace using [HCL](https://www.terraform.io/language/syntax/configuration), like other cloud resources compatible with Terraform.

## Installation

To install this provider, use the following snippet in your Terraform configuration:
```hcl
terraform {
  required_providers {
    solace = {
      source  = "app.terraform.io/telusagriculture/solace"
      version = "0.3.0"
    }
  }
}

provider "solace" {
  username = "username"
  password = "password"
  scheme   = "http or https"
  hostname = "solace hostname"
}
```

In order to access this provider, your local Terraform process must either be authenticated to the "[telusagriculture](https://app.terraform.io/app/telusagriculture)" organization in Terraform Cloud (see [`terraform login`](https://www.terraform.io/cli/commands/login)) or be running in Terraform Cloud.

## Usage

For documentation on the available resources and their properties, see [terraform-provider-solace/docs](terraform-provider-solace/docs/index.md)

## Development

To enable provider development, check out the provider and install the binary into your local `~/go/bin` path:
```
$ cd terraform-provider-solace
$ go install
```

In order for Terraform to use the provider installed locally, rather than downloading a fresh copy from the registry, enable a development override in your `~/terraformrc`:

```hcl
provider_installation {

  dev_overrides {
      "app.terraform.io/telusagriculture/solace" = "your ~/go/bin path"
  }

  direct {}

}
```

Always run `gofmt -w -s .` before committing to make sure the diffs don't contain minor formatting differences.

## Compiling

### Generate schemas and models
1. Install [OpenAPI generator](https://github.com/OpenAPITools/openapi-generator)
1. Generate the resource schema and models from the Solace OpenAPI description:
   ```
   $ make generate-provider
   ```
   This will generate the various `model_*.go` files. The list of models to generate is specified in the [Makefile](Makefile) (`-Dmodels=MsgVpn,...`). The generated files will sometimes not compile due to unused imports, open all the files and save them using VS Code to automatically remove the unused imports.

### Compile the provider

Once the models and schemas are prepared, compile the provider:
```
$ cd terraform-provider-solace
$ go build
```

## Publishing

1. Set the following environment variables:
   1. `GITHUB_TOKEN` - Your Github [personal access token](https://github.com/settings/tokens)
   1. `GPG_FINGERPRINT` - Your GPG key fingerprint without spaces
   1. `TF_API_TOKEN` - Terraform Cloud [team or personal token](https://www.terraform.io/cloud-docs/users-teams-organizations/api-tokens)
1. Create a release
   ```
   $ git tag v1.0.0
   $ cd terraform-provider-solace
   $ goreleaser release --rm-dist
   ```
1. [Create and upload a GPG key](https://www.terraform.io/cloud-docs/registry/publish-providers#publishing-a-provider-and-creating-a-version) to Terraform Cloud, then put the key-id in the [Makefile](Makefile)
1. Use [terraform-publisher](https://github.com/TelusAg/terraform-publisher) to upload to Terraform Cloud private registry
   ```
   $ make publish
   ```
