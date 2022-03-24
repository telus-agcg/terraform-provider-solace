generate:
	rm -rf tfclient
	openapi-generator generate \
		-g go \
		-i tf-private-registry.yaml \
		--skip-validate-spec \
		--output tfclient \
		--package-name tfclient
	gofmt -w tfclient
	cd tfclient; rm -rf go.mod go.sum git_push.sh docs .travis.yml api

provider-publisher: *.go
	go build

run:
	cd ../terraform-provider-solace/dist &&	../../provider-publisher/provider-publisher \
		-organization telusagriculture	\
		-name solace					\
		-version 0.3.0					\
		-gpg-key 14F5FEA3B2DA691F		

all: generate provider-publisher run
