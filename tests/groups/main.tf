provider "msgraph" {
  tenant_id        = "common"
  client_id        = "82492584-8587-4e7d-ad48-19546ce8238f"
  client_secret    = "" // empty for device code authorization
  token_cache_path = "token_cache.json"
}

variable "tenant_domain" {
  type    = string
  default = "l0wdev.onmicrosoft.com"
}

resource "msgraph_group" "demo_security_group" {
  display_name  = "Demo Security Group"
  mail_nickname = "demo_security_group"
}

resource "msgraph_group" "demo_office365_group" {
  display_name  = "Demo Office365 Group"
  mail_nickname = "demo_office365_group"
  group_types   = ["Unified"]
  visibility    = "Private"
}

data "msgraph_group" "demo_security_group" {
  id = msgraph_group.demo_security_group.id
}

data "msgraph_group" "demo_office365_group" {
  mail_nickname = msgraph_group.demo_office365_group.mail_nickname
}
