# Terraform Provider for Microsoft Graph

implementation of Terraform provider for Microsoft Graph using msgraph.go.
One of the main purposes of this provider is to become an alternative to the official Azure Active Directory provider.
You need Terraform v0.12 and an Azure AD tenant with the admin privilege.

## Example Usage

```hcl
provider "msgraph" {
  tenant_id        = "common"
  client_id        = "82492584-8587-4e7d-ad48-19546ce8238f"
  client_secret    = "" // empty for device code authorization
  token_cache_path = "token_cache.json"
}

resource "msgraph_group" "demo_office365_group" {
  display_name  = "Demo Office365 Group"
  mail_nickname = "demo_office365_group"
  group_types   = ["Unified"]
  visibility    = "Private"
}
```

## Provider configuration

The provider has the configuration with the following default values. You can modify the default values with the corresponding environment variables.

```hcl
provider "msgraph" {
  tenant_id           = "common"                               // env:ARM_TENANT_ID
  client_id           = "82492584-8587-4e7d-ad48-19546ce8238f" // env:ARM_CLIENT_ID
  client_secret       = ""                                     // env:ARM_CLIENT_SECRET
  token_cache_path    = "token_cache.json"                     // env:ARM_TOKEN_CACHE_PATH
  console_device_path = "/dev/tty"                             // env:ARM_CONSOLE_DEVICE_PATH
}
```
The default configuration above is to use the public client defined in l0w.dev tenant with the permission Directory.AccessAsUser.All. You can use it to make terraform to access your tenant's directory with the delegated privilege.

When client_secret is empty, the provider attempts the device code authorization. You can see the following message on the first invocation of terraform plan:

```bash
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code GNATKX4J8 to authenticate.
```

Open https://microsoft.com/devicelogin with your web browser and enter the code to proceed the authorization steps. After completing authorization it stores auth tokens in a file specified by token_cache_path. On subsequent terraform invocations it can skip the authorization steps above with this file.

You can also specify an Azure Blob URL with SAS for `token_cache_path`. It's recommended to pass it via `ARM_TOKEN_CACHE_PATH` envvar since the SAS is considered sensitive information that should be hidden.

The provider opens `console_device_path` to prompt the instruction of the device code authorization. It might have no acccess to /dev/tty in the restricted environment like GitLab CI runner. You can workaround it by fd number device and redirection with the shell as follows:

## Supported resources

* Data sources
  * data_group
  * data_user
* Resources
  * msgraph_application
  * msgraph_application_password
  * msgraph_group
  * msgraph_group_member
  * msgraph_user
