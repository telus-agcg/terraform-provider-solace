package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"

	"telusag/provider-publisher/tfclient"

	rt "github.com/go-openapi/runtime/client"
)

const org string = "telusagriculture"
const tok string = "VL6OpSEGpgaSUg.atlasv1.FPF3fIrNcFMVAOk3EFnBQBsDxy0nXHMRJmAX9OJiVbjYUX9ZemhBDaFy7YDd2LCcy4U"
const nam string = "solace"
const ver string = "0.1.0"
const key string = "14F5FEA3B2DA691F"

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
		log.Fatal("Unable to create HTTPS client")
	}
	config.HTTPClient = httpClient

	client = tfclient.NewAPIClient(config)

	CreateProvider("solace")
	uploadLinks := CreateProviderVersion(nam, ver, key)

	nameAndVersion := fmt.Sprintf("terraform-provider-%s_%s", nam, ver)

	UploadFile(nameAndVersion+"_SHA256SUMS", "SHA256SUMS", "text/plain", *uploadLinks.ShasumsUpload)
	UploadFile(nameAndVersion+"_SHA256SUMS.sig", "SHA256SUMS.sig", "application/octet-stream", *uploadLinks.ShasumsSigUpload)

	CreatePlatforms(nameAndVersion+"_SHA256SUMS", nam, ver)
}

func CreatePlatforms(sumsfile string, name string, version string) {
	content, err := ioutil.ReadFile(sumsfile)
	if err != nil {
		log.Fatalf("Unable to open sums file %s", sumsfile)
	}

	// Get the sums and create a platform for each
	re := regexp.MustCompile("^([0-9a-fA-F]+)  ([^_]+)_([^_]+)_([^_]+)_([^_]+).zip$")
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		for _, match := range re.FindAllStringSubmatch(line, -1) {
			sum := match[1]
			os := match[4]
			arch := match[5]

			CreatePlatform(name, version, os, arch, sum)
		}
	}
}

func CreatePlatform(name string, version string, os string, arch string, sum string) {
	filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", name, version, os, arch)
	createPlatform := tfclient.CreateProviderPlatform{
		Data: &tfclient.CreateProviderPlatformData{
			Attributes: tfclient.CreateProviderPlatformDataAttributes{
				Os:       os,
				Arch:     arch,
				Shasum:   sum,
				Filename: filename,
			},
		},
	}
	req := client.DefaultApi.CreateProviderPlatform(ctx, org, name, version).CreateProviderPlatform(createPlatform)
	createPlatformRes, _, err := req.Execute()
	if err != nil {
		log.Fatalf("Error creating platform %s/%s: %s", os, arch, err.Error())
	}

	// Now upload the binary
	UploadFile(filename, filename, "application/octet-stream", *createPlatformRes.Data.Links.ProviderBinaryUpload)
}

func UploadFile(filename string, remoteName string, contentType string, url string) {
	log.Printf("Uploading %s", filename)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(contentType, remoteName)

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("PUT", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got %d: %s", response.StatusCode, response.Status)

	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func CreateProviderVersion(name string, version string, keyid string) *tfclient.CreateRegistryProviderVersionResponseDataLinks {
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
		log.Fatalf("Provider %s version %s already exists.\n", name, version)
	}

	log.Printf("Creating provider %s version %s\n", name, version)
	providerVersion := tfclient.CreateRegistryProviderVersion{
		Data: tfclient.CreateRegistryProviderVersionData{
			Type: "registry-provider-versions",
			Attributes: tfclient.CreateRegistryProviderVersionDataAttributes{
				Version:   version,
				KeyId:     keyid,
				Protocols: []string{"6.0"},
			},
		},
	}
	createProviderVersionRes, _, err := client.DefaultApi.CreateProviderVersion(
		ctx, org, name).CreateRegistryProviderVersion(providerVersion).Execute()
	if err != nil {
		log.Fatalf("Error creating %s, version %s: %s\n", name, version, err.Error())
	}

	return createProviderVersionRes.Data.Links
}

func CreateProvider(name string) {
	getProvidersRes, _, err := client.DefaultApi.GetProviders(
		ctx, org).Execute()
	if err != nil {
		log.Fatalf("Unable to get providers: %s\n", err.Error())
	}

	if len(getProvidersRes.Data) == 0 {
		log.Printf("Creating provider %s\n", name)

		// Create the Solace provider
		provider := tfclient.RegistryProvider{
			Data: tfclient.RegistryProviderData{
				Type: "registry-providers",
				Attributes: tfclient.RegistryProviderDataAttributes{
					Name:         name,
					Namespace:    org,
					RegistryName: "private",
				},
			},
		}
		_, _, err := client.DefaultApi.
			CreateProvider(ctx, org).RegistryProvider(provider).
			Execute()
		if err != nil {
			log.Printf("Error creating provider: %s\n", err.Error())
		}
	}
}
