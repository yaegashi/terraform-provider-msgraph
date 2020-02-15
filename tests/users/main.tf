
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

resource "msgraph_user" "demousers" {
  count                              = 5
  user_principal_name                = "demouser${count.index}@${var.tenant_domain}"
  display_name                       = "Demo User ${count.index}"
  given_name                         = "User ${count.index}"
  surname                            = "Demo"
  mail_nickname                      = "demouser${count.index}"
  other_mails                        = ["demouser${count.index}@example.com"]
  account_enabled                    = true
}

data "msgraph_user" "demouser0" {
  id = msgraph_user.demousers[0].id
}

data "msgraph_user" "demouser1" {
  user_principal_name = msgraph_user.demousers[1].user_principal_name
}

data "msgraph_user" "demouser2" {
  mail_nickname = msgraph_user.demousers[2].mail_nickname
}
