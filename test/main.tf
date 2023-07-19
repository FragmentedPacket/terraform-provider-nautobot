terraform {
  required_providers {
    nautobot = {
      source = "nautobot/nautobot"
    }
  }
}

provider "nautobot" {
  url   = "https://demo.nautobot.com/api/"
  token = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}

// resource "nautobot_manufacturer" "new" {
//   description = "Created with Terraform"
//   name        = "New Vendor"
// }
