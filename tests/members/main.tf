
provider "msgraph" {
  tenant_id        = "common"
  client_id        = "82492584-8587-4e7d-ad48-19546ce8238f"
  client_secret    = "" // empty for device code authorization
  token_cache_path = "token_cache.json"
}

variable "config" {
  type = object({
    tenant_domain = string
    user_count    = number
    group_count   = number
  })
  default = {
    tenant_domain = "l0wdev.onmicrosoft.com"
    user_count    = 100
    group_count   = 10
  }
}

resource "msgraph_user" "demo_users" {
  count               = var.config.user_count
  user_principal_name = "demouser${count.index}@${var.config.tenant_domain}"
  display_name        = "Demo User ${count.index}"
  given_name          = "User ${count.index}"
  surname             = "Demo"
  mail_nickname       = "demouser${count.index}"
  other_mails         = ["demouser${count.index}@example.com"]
  account_enabled     = true
}

resource "msgraph_group" "demo_groups" {
  count         = var.config.group_count
  display_name  = "Demo Group ${count.index}"
  mail_nickname = "demogroup${count.index}"
}

resource "msgraph_group_member" "demo_group_user_members" {
  count     = var.config.user_count
  group_id  = msgraph_group.demo_groups[floor(count.index / (var.config.user_count / var.config.group_count))].id
  member_id = msgraph_user.demo_users[count.index].id
}

resource "msgraph_group_member" "demo_group_group_members" {
  count     = var.config.group_count - 1
  group_id  = msgraph_group.demo_groups[count.index].id
  member_id = msgraph_group.demo_groups[count.index + 1].id
}