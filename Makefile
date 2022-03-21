PKG_NAME=provider
TF_ORG=telusagriculture
TF_REGISTRY=https://app.terraform.io/api/v2/organizations/$(TF_ORG)/registry-providers
CURL=https_proxy=localhost:8080 curl --insecure \
	--header "Authorization: Bearer $(TF_TOKEN)"
JSON=--header "Content-Type: application/vnd.api+json"
GPG_PUBLIC_KEY=$(shell cat gpg-public-key.asc | sed 's/\n/\\n/g' )


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

release:
	GITHUB_TOKEN=$(GITHUB_TOKEN) \
	GPG_FINGERPRINT=$(GPG_FINGERPRINT) \
		goreleaser release --snapshot --rm-dist

upload-key:
	$(CURL) $(JSON) https://app.terraform.io/api/registry/private/v2/gpg-keys -d '{	\
		"data": {									\
			"type": "gpg-keys",						\
			"attributes": {							\
				"namespace": "telusagriculture",	\
				"ascii-armor": "$(GPG_PUBLIC_KEY)"	\
			}										\
		}											\
	}'

create-provider:
	$(CURL) $(JSON) $(TF_REGISTRY) -d '{			\
		"data": {									\
			"type": "registry-providers",			\
			"attributes": {							\
				"name": "solace",					\
				"namespace": "$(TF_ORG)",			\
				"registry-name": "private"			\
			}										\
		}											\
	}'
	

create-version:
	$(CURL) $(JSON) $(TF_REGISTRY)/private/$(TF_ORG)/solace/versions -d '{	\
		"data": {									\
			"type": "registry-provider-versions", 	\
			"attributes": {							\
				"version": "0.1.0", 				\
				"key-id": "14F5FEA3B2DA691F",		\
				"protocols": ["6.0"] 				\
			} 										\
		} 											\
	}'

upload-sigs:
	$(CURL) -T terraform-provider-solace/dist/terraform-provider-solace_*_SHA256SUMS \
		https://archivist.terraform.io/v1/object/dmF1bHQ6djE6S1BmOFhwSFNxYnRmeDJPanJUamJoZzhhY3VxTTlDeG1lc3VSaHVQZFVBa29TNUgxUHBHb0RpUW5od0QvSmdGTHgzdTIrTzhINnhUT095VDRLaFR2MUhpWmJLaTNOSWpiektJaHBtTzF5ZktMMlAvNmFKR1AxRW0wbHN2RG1VS05rVnU2bXQ0Q3ZTdHp6SW0zL2c3TGJCVHlFSjVjdG5uMUFDMENkTWdscS9zd3JxWnluSzRjYTVRb3VubnNyOGhwZmx5elBlbUpZdGc2VG9CNGN2TkJmZFFYUTczMkRmazBBdk16RGR6eTZmbThZNlJXRUYwOGpNbngxaGNjYTVmUlNlY21IWmszbVBZR2I2dW40cldvY3pwZ2pQOWlBK0FpVUgwblp5R0Zza3M9

	$(CURL) -T terraform-provider-solace/dist/terraform-provider-solace_*_SHA256SUMS.sig \
		https://archivist.terraform.io/v1/object/dmF1bHQ6djE6dXc4YW1PMkVrenV6akZ6c2ptTkNDQ25wYk9iU3FJODZZUGo3R3NhRkZVSThua2phazFTKzZNSzRwRy93YWdRenBDNWhXZGxDQUJGeXFjeU1FUzc0eXlQWmJTSk5tNG12cXJ6Qm1uR25JOFVmOHNKcjNaRGNZUHJzbmZ6UUtGeGhmbmFmaTZFSEZSL1BNSWtKVDZKUHJPZnVGSDBEdjFxd2dDREsrckFwaFRvVzVlR1YxcjFEV1dGRG1SeFZiOEp5UGFDNldpd0dBd0xFTnZmVXFtaWtTdUl1RTY2aUV3b3p1VCtxbjRkdVQ4eGx4WGdsY0t2aldUZFpwMElUNitZM0pDWno2Wk1MRHpQazQ2V3JQS2NsYTRndjNPeTdodmltdUthVUs0YXdkQVM1OGl5Rw

create-platforms:
	$(CURL) $(JSON) $(TF_REGISTRY)/private/$(TF_ORG)/solace/versions/0.1.0/platforms -d '{		\
		"data": {																				\
			"type": "registry-provider-version-platforms",										\
			"attributes": {																		\
				"os": "darwin",																	\
				"arch": "amd64",																\
				"shasum": "93875bd71a76581977449631e108f3e674cabf1905d314f602ef9c7838a8bf4c",	\
				"filename": "terraform-provider-solace_0.1.0_darwin_amd64.zip"					\
			}																					\
		}																						\
	}'
	$(CURL) $(TF_REGISTRY)/private/$(TF_ORG)/solace/versions/0.1.0/platforms -d '{				\
		"data": {																				\
			"type": "registry-provider-version-platforms",										\
			"attributes": {																		\
				"os": "linux",																	\
				"arch": "amd64",																\
				"shasum": "de78e2e009e3e468d62745e0d6defa1747279d0bd45ef4b5a4a7be2de7731b96",	\
				"filename": "terraform-provider-solace_0.1.0_linux_amd64.zip"					\
			}																					\
		}																						\
	}'
	$(CURL) $(TF_REGISTRY)/private/$(TF_ORG)/solace/versions/0.1.0/platforms -d '{				\
		"data": {																				\
			"type": "registry-provider-version-platforms",										\
			"attributes": {																		\
				"os": "windows",																\
				"arch": "amd64",																\
				"shasum": "c2a25e5ede6f83057b1e64d5f332f44829585a4c6cd9e7f8df7fdf22b1554007",	\
				"filename": "terraform-provider-solace_0.1.0_windows_amd64.zip"					\
			}																					\
		}																						\
	}'

upload-binaries:
	$(CURL) -T terraform-provider-solace/dist/terraform-provider-solace_*_darwin_amd64.zip \
		https://archivist.terraform.io/v1/object/dmF1bHQ6djE6SkZtNllZNHpQeDVVaU9vNEZxRDg1S0NPNVlZUjQ4ZnA5bTBpUnZsVVFHUm1uODFlVnpKYnlQSzBGTWNQU2g1a0ZpVkFpN2g2QmV2TTVqUEpaY0FCK2dPTjYrRVZNdTJrbERGVVJTZjFZVytpRmUrNlBvZDJkdlFzS2dyM0dLcWNzSTZhczdSUWp2dnFIOHNYN3hFVEFtdTJwWHowNU55dGoybnpvZEgrdXdQdDRBVkowbGtrajBYbHp0VjdTdnB1YzZFaVB5cjZ5V3IrV3dmSG5HK3hZQkVMWUI3R25FT2ZZVjV6VWxSc3NRVmR2aHFIOUplYlhjL05aK1YvUmc9PQ

	$(CURL) -T terraform-provider-solace/dist/terraform-provider-solace_*_linux_amd64.zip \
		https://archivist.terraform.io/v1/object/dmF1bHQ6djE6bjA0cGhBK0tyalBEKzB6T2lNa1A5aFFwZ3lyc24zQTJHd1pMYjBvOVo3WHFYaUlhOU5lcGt0VW1Wa29rYUwvcnZmOWtDclZGWTRzTDFubWtoY2YvaXdkb1N5bWxqQmRtaXZnajN6UHlVTUVIUHoxWGFCdVlXbFBSQWg5cE41NmpRc3ZQRjZFTUJhQk96c3FIR000dHcyMTlLUFNqWUc5YWRaUE1EVDE2SVJzVktKbFl1QlJTa0hhaW9TNHRVR1MwcnJzajg0MXA4bitaYUFhYkxObWFpVTlkTVR5L001eTI0YTBOMDBWREJrMklvUUtlMWlhaGtJMzVRWm82RlE9PQ

	$(CURL) -T terraform-provider-solace/dist/terraform-provider-solace_*_windows_amd64.zip \
		https://archivist.terraform.io/v1/object/dmF1bHQ6djE6Rkx4T2pNcElHOHBWYTA4NEh1dDhia2dtMzR1cjhZbmcrTllLbEF2N2k0VXZYZ1d4ZVE1Z2Nna2VCM0VqM2VxZ0IwR0JNQjVZclZKdzhEWmlqbEhJbXpDbHRWZEUvQXFhc2JoRXFPUkJHWkpVNC9yTHRYS0xFV3ByeE54eWJ0eEtBWXMxK2JOVURhVXNxbHVIdUFGVTcxd21WS3ZYalhTZHBmMjNyd3IrMXJQNENVbjZNOWxja011cWM1YVZHOS9aZnFZQWFGVTVKd05tazB5NkFhdksrWUpNQ3hqM2RQOWwzSVdLbEpmSkkyVzNsYjJNYkV6YnpRQ21YR2QrR2c9PQ
