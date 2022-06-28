terraform {
  required_providers {
    solace = {
      source  = "TelusAg/solace"
      version = ">= 0.6.3"
    }
  }

  required_version = "> 0.14"
}

provider "solace" {
  hostname = "solace:8080"
  username = "admin"
  password = "admin"
  // default_msgvpn = "bar"
}
