# Terraform Provider for Microsoft Graph

![Test](https://github.com/yaegashi/terraform-provider-msgraph/workflows/Test/badge.svg)
![Release](https://github.com/yaegashi/terraform-provider-msgraph/workflows/Release/badge.svg)

## Introduction

The POC implementation of [Terraform](https://terraform.io) provider
for [Microsoft Graph](https://developer.microsoft.com/en-us/graph)
using [msgraph.go](https://github.com/yaegashi/msgraph.go).

One of the main purposes of this provider is to become an alternative
to [the official Azure Active Directory provider](https://www.terraform.io/docs/providers/azuread/).

You need Terraform v0.12 and an Azure AD tenant with the admin privilege.

## Supported resources

- Data sources
  - data_group
  - data_user
- Resources
  - msgraph_application
  - msgraph_application_password
  - msgraph_group
  - msgraph_group_member
  - msgraph_user

## Provider configuration

The provider has the configuration with the following default values.
You can modify the default values with the corresponding environment variables.

```hcl
provider "msgraph" {
  tenant_id           = "common"                               // env:ARM_TENANT_ID
  client_id           = "82492584-8587-4e7d-ad48-19546ce8238f" // env:ARM_CLIENT_ID
  client_secret       = ""                                     // env:ARM_CLIENT_SECRET
  token_cache_path    = "token_cache.json"                     // env:ARM_TOKEN_CACHE_PATH
  console_device_path = "/dev/tty"                             // env:ARM_CONSOLE_DEVICE_PATH
}
```

The default configuration above is
to use the public client defined in `l0w.dev` tenant with the permission `Directory.AccessAsUser.All`.
You can use it to make terraform to access your tenant's directory with the delegated privilege.

When `client_secret` is empty,
the provider attempts [the device code authorization](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code).
You can see the following message on the first invocation of `terraform plan`:

```console
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code GNATKX4J8 to authenticate.
```

Open https://microsoft.com/devicelogin with your web browser and enter the code to proceed the authorization steps.
After completing authorization it stores auth tokens in a file specified by `token_cache_path`.
On subsequent terraform invocations it can skip the authorization steps above with this file.

You can also specify an Azure Blob URL with SAS for `token_cache_path`.
It's recommended to pass it via `ARM_TOKEN_CACHE_PATH` envvar
since the SAS is considered sensitive information that should be hidden.

The provider opens `console_device_path` to prompt the instruction of the device code authorization.
It might have no acccess to `/dev/tty` in the restricted environment like GitLab CI runner.
You can workaround it by fd number device and redirection with the shell as follows:

```console
$ 99>&2 ARM_CONSOLE_DEVICE_PATH=/dev/fd/99 terraform plan
```

## How to test

Terraform v0.12 and Go v1.14 are required.
It's strongly recommended to acquire a developer sandbox tenant
by joining [the Office 365 developer program](https://developer.microsoft.com/en-us/office/dev-program).

Clone the repository, then move to one of [the test directories](tests) and build `terraform-provider-msgraph` executable there:

```console
$ git clone https://github.com/yaegashi/terraform-provider-msgraph
$ cd terraform-provider-msgraph/tests/users
$ go build ../..
```

Edit `provider` and `variable` in `main.tf` for your environment:

```hcl
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
```

Run terraform with an environment variable `TF_LOG=DEBUG` to enable debug log output:

```console
$ terraform init
$ TF_LOG=DEBUG terraform plan
$ TF_LOG=DEBUG terraform apply
```

## Todo

- [ ] Support various graph resources (`resource`/`data`)
  - [ ] [User](https://docs.microsoft.com/en-us/graph/api/resources/user)
  - [ ] [Group](https://docs.microsoft.com/en-us/graph/api/resources/group)
  - [ ] [Application](https://docs.microsoft.com/en-us/graph/api/resources/application)
  - [ ] [Team](https://docs.microsoft.com/en-us/graph/api/resources/teams-api-overview)
  - [ ] [Site](https://docs.microsoft.com/en-us/graph/api/resources/sharepoint) (no ability to create new sites)
- [ ] Support importing
- [ ] Code auto-generation based on the API metadata
- [x] Persist OAuth2 tokens in backend storage (Azure Blob Storage)
- [x] Better device auth grant experience (no `TF_LOG=DEBUG`)
- [ ] Unit testing
- [x] CI/CD (GoReleaser)
- [ ] Manuals
- [ ] Publish to the Terraform registry (#1)
