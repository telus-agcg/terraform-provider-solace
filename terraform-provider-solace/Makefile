
format-examples:
	terraform fmt -recursive ./examples/

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

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
