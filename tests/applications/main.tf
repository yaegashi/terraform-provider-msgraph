
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

resource "msgraph_application" "demoapplication" {
  display_name                 = "Demo Application"
  sign_in_audience             = "AzureADMyOrg"
  identifier_uris              = ["http://localhost"]
  home_page_url                = "http://localhost"
  logout_url                   = "http://localhost/logout"
  redirect_uris                = ["http://localhost:8080"]
  enable_id_token_issuance     = true
  enable_access_token_issuance = true
}

resource "msgraph_application_password" "demoapplicationpassword" {
  application_id  = msgraph_application.demoapplication.id
  display_name    = "Demo Application Password"
  end_date_time   = "2100-01-01T00:00:00Z"
}