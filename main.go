package main

import (
	"context"
	"log"

	"telusag/terraform-provider-solace/provider"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	opts := tfsdk.ServeOpts{
		Name: "registry.terraform.io/telusag/solace",
	}

	err := tfsdk.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}

}
