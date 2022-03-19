PKG_NAME=provider

.PHONY:

openapi-provider-generator:
	mvn -f provider-generator/pom.xml package

generate-provider: openapi-provider-generator
	java -cp "provider-generator/target/terraform-provider-openapi-generator-1.0.0.jar:/usr/local/Cellar/openapi-generator/5.4.0/libexec/openapi-generator-cli.jar" \
		-Dmodels=MsgVpn,MsgVpnQueue,MsgVpnClientUsername,MsgVpnAclProfile,MsgVpnClientProfile \
	org.openapitools.codegen.OpenAPIGenerator generate \
		-g terraform-provider \
		-i terraform-provider-solace/sempv2-config.json \
		--skip-validate-spec \
		--output terraform-provider-solace/$(PKG_NAME) \
		--package-name $(PKG_NAME)
	gofmt -w terraform-provider-solace/$(PKG_NAME)

install:
	cd terraform-provider-solace; go install