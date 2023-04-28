resource "solace_msgvpn" "vpn1" {
  msg_vpn_name                 = "vpn1"
  enabled                      = true
  authentication_basic_enabled = true
  max_connection_count         = 75
  max_egress_flow_count        = 100
  max_endpoint_count           = 10
  max_ingress_flow_count       = 100
  max_msg_spool_usage          = 1500
  max_subscription_count       = 100000
  max_transacted_session_count = 50
  max_transaction_count        = 100
  replication_enabled          = false

  event_egress_msg_rate_threshold = {
    set_value   = 100
    clear_value = 10
  }
}

resource "solace_queue" "vpn1_queue" {
  for_each        = toset(["one", "two", "three", "four", "five"])
  msg_vpn_name    = solace_msgvpn.vpn1.msg_vpn_name
  queue_name      = "queue-${each.value}"
  ingress_enabled = true
  egress_enabled  = true
  access_type     = "non-exclusive"
  max_ttl         = 60
}

resource "solace_queue" "vpn1_another_queue" {
  msg_vpn_name    = solace_msgvpn.vpn1.msg_vpn_name
  queue_name      = "queue-b"
  ingress_enabled = true
  egress_enabled  = true
  access_type     = "non-exclusive"
  max_ttl         = 60
}
