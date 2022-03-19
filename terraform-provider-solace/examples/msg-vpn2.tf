resource "solace_msgvpn" "vpn2" {
  msg_vpn_name                         = "vpn2"
  enabled                      = true
  authentication_basic_enabled = true
  max_connection_count         = 75
  max_egress_flow_count        = 100
  max_endpoint_count           = 10
  max_ingress_flow_count       = 100
  max_msg_spool_usage              = 1500
  max_subscription_count       = 100000
  max_transacted_session_count = 50
  max_transaction_count        = 100
  replication_enabled          = false
}

resource "solace_clientprofile" "vpn2_client_profile" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  client_profile_name    = "my-client-profile"
}

resource "solace_aclprofile" "vpn2_aclprofile" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  acl_profile_name                             = "aclprofile"
  client_connect_default_action = "disallow"
  publish_topic_default_action     = "disallow"
  subscribe_topic_default_action   = "disallow"
}

/*
resource "solace_aclprofile_clientconnexception" "vpn2_allow_10dot" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  acl_profile = solace_aclprofile.vpn2_aclprofile.name
  address     = "10.2.0.0/24"
}

resource "solace_aclprofile_clientconnexception" "vpn2_allow_fd44" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  acl_profile = solace_aclprofile.vpn2_aclprofile.name
  address     = "fd44:5249:4553:10::/64"
}

resource "solace_aclprofile_publishexception" "vpn2_allow_pub_abc" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  acl_profile  = solace_aclprofile.vpn2_aclprofile.name
  topic_syntax = "smf"
  topic        = "a/b/c"
}

resource "solace_aclprofile_subscribeexception" "vpn2_allow_sub_abc" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  acl_profile  = solace_aclprofile.vpn2_aclprofile.name
  topic_syntax = "smf"
  topic        = "a/b/d"
}
*/

resource "solace_clientusername" "vpn2_user1" {
  msg_vpn_name = solace_msgvpn.vpn2.msg_vpn_name
  client_username           = "user1"
  enabled        = true
  client_profile_name = solace_clientprofile.vpn2_client_profile.client_profile_name
}
