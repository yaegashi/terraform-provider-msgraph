
provider "msgraph" {
  // Define environment variables:
  // ARM_TENANT_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET
}

variable "config" {
  type    = object({ domain = string, name = string, user_count = number, group_count = number })
  default = { domain = "l0wdev.onmicrosoft.com", name = "msgraphtf", user_count = 100, group_count = 10 }
}

resource "msgraph_user" "demo_users" {
  count               = var.config.user_count
  user_principal_name = "${var.config.name}user${count.index}@${var.config.domain}"
  display_name        = "${var.config.name} user ${count.index}"
  mail_nickname       = "${var.config.name}user${count.index}"
  account_enabled     = true
}

resource "msgraph_group" "demo_groups" {
  count         = var.config.group_count
  display_name  = "${var.config.name} group ${count.index}"
  mail_nickname = "${var.config.name}group${count.index}"
}

resource "msgraph_group_member" "demo_group_user_members" {
  count     = var.config.user_count
  group_id  = msgraph_group.demo_groups[count.index % var.config.group_count].id
  member_id = msgraph_user.demo_users[count.index].id
}

resource "msgraph_group_member" "demo_group_group_members" {
  count     = var.config.group_count - 1
  group_id  = msgraph_group.demo_groups[count.index].id
  member_id = msgraph_group.demo_groups[count.index + 1].id
}
