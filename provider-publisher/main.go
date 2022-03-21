package main

import (
	"context"
	"fmt"
	"log"

	"telusag/provider-publisher/tfclient"

	rt "github.com/go-openapi/runtime/client"
)

const org string = "telusagriculture"

var ctx context.Context
var client *tfclient.APIClient

func main() {
	ctx = context.WithValue(context.Background(), tfclient.ContextAccessToken, tok)

	config := tfclient.NewConfiguration()
	config.Scheme = "https"
	config.Host = "app.terraform.io"
	config.UserAgent = "Terraform Provider Publisher"

	httpClient, err := rt.TLSClient(
		rt.TLSClientOptions{InsecureSkipVerify: true})
	if err != nil {
		fmt.Println("Unable to create HTTPS client")
	}
	config.HTTPClient = httpClient

	client = tfclient.NewAPIClient(config)

	CreateProvider("solace")
	CreateProviderVersion("solace", "0.1.0")
}

func CreateProviderVersion(name string, version string) {
	getProviderVersionsRes, _, err := client.DefaultApi.GetProviderVersions(
		ctx, org, name).Execute()
	if err != nil {
		log.Fatalf("Unable to get provider versions: %s\n", err.Error())
	}

	// Check to see if the requested version already exists
	found := false
	for _, v := range getProviderVersionsRes.Data {
		if *v.Attributes.Version == version {
			found = true
		}
	}
	if found {
		fmt.Printf("Provider %s version %s already exists.\n", name, version)
		return
	}

	fmt.Printf("Creating provider %s version %s\n", name, version)
	providerVersion := tfclient.CreateRegistryProviderVersion{
		Data: tfclient.CreateRegistryProviderVersionData{
			Type: "registry-provider-versions",
			Attributes: tfclient.CreateRegistryProviderVersionDataAttributes{
				Version:   version,
				KeyId:     "",
				Protocols: []string{"6.0"},
			},
		},
	}
	createProviderVersionRes, _, err := client.DefaultApi.CreateProviderVersion(
		ctx, org, name).CreateRegistryProviderVersion(providerVersion).Execute()
	if err != nil {
		log.Fatalf("Error creating %s, version %s: %s\n", name, version, err.Error())
	}

	fmt.Printf("%+v", createProviderVersionRes)
}

func CreateProvider(name string) {
	getProvidersRes, _, err := client.DefaultApi.GetProviders(
		ctx, org).Execute()
	if err != nil {
		log.Fatalf("Unable to get providers: %s\n", err.Error())
	}

	if len(getProvidersRes.Data) == 0 {
		fmt.Println("Creating provider")

		// Create the Solace provider
		provider := tfclient.RegistryProvider{
			Data: tfclient.RegistryProviderData{
				Type: "registry-providers",
				Attributes: tfclient.RegistryProviderDataAttributes{
					Name:         "solace",
					Namespace:    org,
					RegistryName: "private",
				},
			},
		}
		_, _, err := client.DefaultApi.
			CreateProvider(ctx, org).RegistryProvider(provider).
			Execute()
		if err != nil {
			fmt.Printf("Error creating provider: %s\n", err.Error())
		}
	}
}
