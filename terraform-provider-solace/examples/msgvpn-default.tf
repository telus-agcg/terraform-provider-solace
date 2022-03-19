resource "solace_queue" "default_queue" {
  for_each        = toset(["one", "two", "three", "four", "five"])
  queue_name      = "queue-${each.value}"
  msg_vpn_name    = "default"
  ingress_enabled = true
  egress_enabled  = true
  access_type     = "non-exclusive"
  max_ttl         = 60
}
