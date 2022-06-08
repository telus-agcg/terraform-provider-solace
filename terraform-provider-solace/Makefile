
format-examples:
	terraform fmt -recursive ./examples/

docs:
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs --rendered-website-dir=../docs

openapi: sempv2-config.json
	rm -rf sempv2
	openapi-generator generate \
		-g go \
		-i sempv2-config.json \
		--skip-validate-spec \
		--output sempv2 \
		--package-name sempv2 \
		-p hideGenerationTimestamp=true \
		-p structPrefix=true \
		-p isGoSubmodule=true
	rm -rf sempv2/go.mod sempv2/go.sum sempv2/git_push.sh sempv2/docs

mkdocs:
	tfplugindocs

release:
	GITHUB_TOKEN=$(GITHUB_TOKEN) \
	GPG_FINGERPRINT=$(GPG_FINGERPRINT) \
		goreleaser release --rm-dist

PUBLISHER=github.com/TelusAg/terraform-publisher
terraform-publisher:
	GOPRIVATE=$(PUBLISHER) go install $(PUBLISHER)@latest

publish: terraform-publisher
	terraform-publisher					\
		-organization telusagriculture	\
		-name solace					\
		-version 0.4.0					\
		-gpg-key 14F5FEA3B2DA691F		
