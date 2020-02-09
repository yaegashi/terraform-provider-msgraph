# WIP: Terraform Provider for Microsoft Graph

## Introduction

The POC implementation of [Terraform](https://terraform.io) Provider for [Microsoft Graph](https://developer.microsoft.com/en-us/graph) using [msgraph.go](https://github.com/yaegashi/msgraph.go).

## How to test

You need Terraform v0.12 and Go v1.13, and an Azure AD tenant with the admin privilege.

Clone the repository, then move to one of [the test directories](tests) and build `terraform-provider-msgraph` executable there:

```console
$ git clone https://github.com/yaegashi/terraform-provider-msgraph
$ cd terraform-provider-msgraph/tests/users
$ go build ../..
```

Configure the provider and variable by editing `main.tf`:

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

## Authorization

When the provider configuration `client_secret` is empty,
it requests you for [the device code authorization](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code)
in debug log output as follows on the first invocation of `terraform plan`:

```
2020-02-09T03:55:33.204+0900 [DEBUG] plugin.terraform-provider-msgraph: 2020/02/09 03:55:33 To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code GEXSRT5LT to authenticate.
```

Open https://microsoft.com/devicelogin with your web browser and enter the code to proceed the authorization steps.
After completing authorization it stores auth tokens in a file specified by `token_cache_path`.
On subsequent terraform invocations it can skip the authorization steps above with this file.

## Todo

- [ ] Support various graph resources
  - [ ] [User](https://docs.microsoft.com/en-us/graph/api/resources/user)
  - [ ] [Group](https://docs.microsoft.com/en-us/graph/api/resources/group)
  - [ ] [Calendar](https://docs.microsoft.com/en-us/graph/api/resources/calendar)
  - [ ] Licenses
  - [ ] [Team](https://docs.microsoft.com/en-us/graph/api/resources/teams-api-overview)
  - [ ] ~~[Site](https://docs.microsoft.com/en-us/graph/api/resources/sharepoint)~~ (no ability to create new sites)
- [ ] Support importing
- [ ] Code auto-generation based on the API metadata
- [ ] Persist OAuth2 tokens in backend storage?
- [ ] Better device auth grant experience (no `TF_LOG=DEBUG`)
- [ ] Unit testing
- [ ] CI
