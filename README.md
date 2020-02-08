# WIP: Terraform Provider for Microsoft Graph

## Introduction

The POC implementation of [Terraform](https://terraform.io) Provider for [Microsoft Graph](https://developer.microsoft.com/en-us/graph) using [msgraph.go](https://github.com/yaegashi/msgraph.go).

## How to test

You need Terraform v0.12 and Go v1.13.

Clone the repository and build `terraform-provider-msgraph` executable:

```console
$ git clone https://github.com/yaegashi/terraform-provider-msgraph
$ cd terraform-provider-msgraph
$ go build .
```

Edit [test_users.tf](test_users.tf) before running terraform:

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

Run terraform with an environment variable `TF_LOG=DEBUG` to see the debug log output:

```console
$ TF_LOG=DEBUG terraform init
$ TF_LOG=DEBUG terraform plan
$ TF_LOG=DEBUG terraform apply
```

On the first `terraform plan` invocation, you'll see a device code authorization message like the following.  Open https://microsoft.com/devicelogin with your web browser and enter the code to authenticate.

```
2020-02-09T03:55:33.204+0900 [DEBUG] plugin.terraform-provider-msgraph: 2020/02/09 03:55:33 To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code GEXSRT5LT to authenticate.
```

Authentication tokens are save in token_cache.json in the current directory.  With this file you can bypass authentication  in subsequent terraform invocations.

## Todo

- [ ] Support various graph resources
  - [ ] [User](https://docs.microsoft.com/en-us/graph/api/resources/user)
  - [ ] [Group](https://docs.microsoft.com/en-us/graph/api/resources/group)
  - [ ] [Calendar](https://docs.microsoft.com/en-us/graph/api/resources/calendar)
  - [ ] Licenses
  - [ ] [Team](https://docs.microsoft.com/en-us/graph/api/resources/teams-api-overview)
  - [ ] ~~[Site](https://docs.microsoft.com/en-us/graph/api/resources/sharepoint)~~ (no ability to create new sites)
- [ ] Auto-generate code based on the API metadata
- [ ] Persist OAuth2 tokens in backend storage?
- [ ] Better device auth grant experience
- [ ] Unit testing
- [ ] CI
