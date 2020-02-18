provider "azuread" {
  // Define environment variables
  // ARM_TENANT_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_SUBSCRIPTION_ID (dummy)
}

variable "config" {
  type    = object({ domain = string, name = string, user_count = number, group_count = number })
  default = { domain = "l0wdev.onmicrosoft.com", name = "azureadtf", user_count = 100, group_count = 10 }
}

provider random {}

resource "random_password" "password" {
  length = 16
}

resource "azuread_user" "demo_users" {
  count               = var.config.user_count
  user_principal_name = "${var.config.name}user${count.index}@${var.config.domain}"
  display_name        = "${var.config.name} user ${count.index}"
  mail_nickname       = "${var.config.name}user${count.index}"
  password            = random_password.password.result
}

resource "azuread_group" "demo_groups" {
  count = var.config.group_count
  name  = "${var.config.name} group ${count.index}"
}

resource "azuread_group_member" "demo_group_user_members" {
  count            = var.config.user_count
  group_object_id  = azuread_group.demo_groups[count.index % var.config.group_count].id
  member_object_id = azuread_user.demo_users[count.index].id
}

resource "azuread_group_member" "demo_group_group_members" {
  count            = var.config.group_count - 1
  group_object_id  = azuread_group.demo_groups[count.index].id
  member_object_id = azuread_group.demo_groups[count.index + 1].id
}
