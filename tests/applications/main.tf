
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

resource "msgraph_application" "demo_provider_app" {
  display_name                 = "Demo Provider App"
  sign_in_audience             = "AzureADMyOrg"
  identifier_uris              = ["http://localhost/provider"]
  home_page_url                = "http://localhost/provider"
  logout_url                   = "http://localhost/logout"
  redirect_uris                = ["http://localhost:8080"]
  enable_id_token_issuance     = true
  enable_access_token_issuance = true

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Hoge description"
      admin_consent_display_name = "Hoge display name"
      id                         = "f5e0be5f-ef28-40ba-8cc9-ec0beedc59d5"
      is_enabled                 = false
      type                       = "User"
      value                      = "Hoge"
    }
    oauth2_permission_scope {
      admin_consent_description  = "Moge description"
      admin_consent_display_name = "Moge display name"
      id                         = "e0c64d0f-4752-476e-ab3b-b71281bec8a2"
      is_enabled                 = false
      type                       = "Admin"
      value                      = "Moge"
    }
  }

  app_role {
    allowed_member_types = ["Application"]
    description          = "Write description"
    display_name         = "Write display name"
    id                   = "32bc1ceb-418e-4b68-9080-a9fda28d2071"
    is_enabled           = false
    value                = "Write"
  }
  app_role {
    allowed_member_types = ["User"]
    description          = "Read description"
    display_name         = "Read display name"
    id                   = "ad936841-3b89-480d-8cf1-a617cf880e6c"
    is_enabled           = false
    value                = "Read"
  }

  required_resource_access {
    resource_app_id = "00000003-0000-0000-c000-000000000000"
    resource_access {
      id   = "205e70e5-aba6-4c52-a976-6d2d46c48043"
      type = "Scope"
    }
    resource_access {
      id   = "df85f4d6-205c-4ac5-a5ea-6bf408dba283"
      type = "Scope"
    }
  }
}

resource "msgraph_application_password" "demo_provider_app_password" {
  application_id = msgraph_application.demo_provider_app.id
  display_name   = "Demo Provider App Password"
  end_date_time  = "2100-01-01T00:00:00Z"
}

resource "msgraph_application" "demo_consumer_app" {
  display_name                 = "Demo Consumer App"
  sign_in_audience             = "AzureADMyOrg"
  identifier_uris              = ["http://localhost/consumer"]
  home_page_url                = "http://localhost/consumer"
  logout_url                   = "http://localhost/logout"
  redirect_uris                = ["http://localhost:8080"]
  enable_id_token_issuance     = true
  enable_access_token_issuance = true

  // XXX: api{} is needed to make the tfstate consistent
  api {}

  required_resource_access {
    resource_app_id = msgraph_application.demo_provider_app.app_id
    resource_access {
      id   = [for a in [for a in msgraph_application.demo_provider_app.api : a][0].oauth2_permission_scope : a.id if a.value == "Hoge"][0]
      type = "Scope"
    }
    resource_access {
      id   = [for a in msgraph_application.demo_provider_app.app_role : a.id if a.value == "Write"][0]
      type = "Role"
    }
  }
}

resource "msgraph_application_password" "demo_consumer_app_password" {
  application_id = msgraph_application.demo_consumer_app.id
  display_name   = "Demo Consumer App Password"
  end_date_time  = "2100-01-01T00:00:00Z"
}
