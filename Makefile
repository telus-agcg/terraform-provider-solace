PKG_NAME=provider
OPENAPI_GENERATOR_JAR=/usr/local/Cellar/openapi-generator/6.2.0/libexec/openapi-generator-cli.jar
GO_POST_PROCESS_FILE="$(shell which gofmt) -w"
MODELS=MsgVpn,$\
	MsgVpnQueue,$\
	MsgVpnQueueSubscription,$\
	MsgVpnClientUsername,$\
	MsgVpnAclProfile,$\
	MsgVpnAclProfileClientConnectException,$\
	MsgVpnAclProfileSubscribeException,$\
	MsgVpnAclProfilePublishException,$\
	MsgVpnClientProfile,$\
	MsgVpnAuthenticationOauthProfile,$\
	ClientCertAuthority

format-examples:
	terraform fmt -recursive ./examples/

generate-docs:
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
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

release:
	GITHUB_TOKEN=$(GITHUB_TOKEN) \
	GPG_FINGERPRINT=$(GPG_FINGERPRINT) \
		goreleaser release --rm-dist

openapi-provider-generator:
	mvn -f provider-generator/pom.xml clean compile

generate-provider: openapi-provider-generator
	GO_POST_PROCESS_FILE=$(GO_POST_PROCESS_FILE) \
	java -cp "provider-generator/target/classes:$(OPENAPI_GENERATOR_JAR)" \
		-Dmodels=$(MODELS) \
	org.openapitools.codegen.OpenAPIGenerator generate \
		-g terraform-provider \
		-i sempv2-config.json \
		--skip-validate-spec \
		--output $(PKG_NAME) \
		--package-name $(PKG_NAME) \
 		--enable-post-process-file
