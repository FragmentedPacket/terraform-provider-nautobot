terraform {
  required_providers {
    nautobot = {
      source = "nautobot/nautobot"
    }
  }
}
provider "nautobot" {}

data "nautobot_manufacturers" "test" {}

output "manus" {
  value = data.nautobot_manufacturers.test
}
