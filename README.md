# Solace Terraform Provider

## Introduction

This is a Terraform Provider for the Solace Event Broker. It allows you to configure various resources in Solace using [HCL](https://www.terraform.io/language/syntax/configuration), like other cloud resources compatible with Terraform.

## Installation

To install this provider, use the following snippet in your Terraform configuration:
```hcl
terraform {
  required_providers {
    solace = {
      source  = "telusag/solace"
      version = ">= 0.6.3"
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

## Usage

For documentation on the available resources and their properties, see [docs](docs/index.md)

## Development

To enable provider development, check out the provider and install the binary into your local `~/go/bin` path:
```
$ go install
```

In order for Terraform to use the provider installed locally, rather than downloading a fresh copy from the registry, enable a development override in your `~/terraformrc`:

```hcl
provider_installation {

  dev_overrides {
      "telusag/solace" = "your ~/go/bin path"
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
$ go build
```

## Publishing

1. [Create and upload a GPG key](https://www.terraform.io/cloud-docs/registry/publish-providers#publishing-a-provider-and-creating-a-version) to Terraform Cloud.
1. Set the GPG secret key as an action secret on the repository called GPG_SECRET_KEY, the passphrase should be set as an action secret called PASSPHRASE.
1. Create a tag on the repository called vx.y.z (v1.0.0) and push the tag to Github.
1. Github Actions will build and publish the new release