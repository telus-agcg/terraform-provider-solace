---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "solace_aclprofile_publish_exception Resource - terraform-provider-solace"
subcategory: ""
description: |-
  MsgVpnAclProfilePublishException
---

# solace_aclprofile_publish_exception (Resource)

MsgVpnAclProfilePublishException



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `acl_profile_name` (String) The name of the ACL Profile. Deprecated since 2.14. Replaced by publishTopicExceptions.
- `msg_vpn_name` (String) The name of the Message VPN. Deprecated since 2.14. Replaced by publishTopicExceptions.
- `publish_exception_topic` (String) The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by publishTopicExceptions.
- `topic_syntax` (String) The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> "smf" - Topic uses SMF syntax. "mqtt" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by publishTopicExceptions.
