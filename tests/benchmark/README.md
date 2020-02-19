# Azure AD Benchmarking

|Script|Parallelism|Create (sec)|Delete (sec)|Notes|
|---|---:|---:|---:|---|
|msgraphgo|1|82.119|42.548||
|msgraphtf|1|103.473|107.952|
|msgraphtf|10|12.568|12.632|Default parallelism|
|msgraphtf|100|7.550|7.943|
|azureadtf|1|>1000?|?|Not yet measured|
|azureadtf|10|125.191|12.832|Default parallelism|
|azureadtf|100|30.266|10.745||

This folder contains 3 benchmarking scripts:

- [msgraphgo](msgraphgo) ... [msgraph.go](https://github.com/yaegashi/msgraph.go) program performing sequential creation/deletion
  - Build: `go build .`
  - Create: `time ./msgraphgo`
  - Delete: `time ./msgraphgo -clean`
- [msgraphtf](msgraphtf) ... Terraform script using [msgraph provider](https://github.com/yaegashi/terraform-provider-msgraph)
  - Create: `time terraform apply -auto-approve -parallelism N`
  - Delete: `time terraform destroy -auto-approve -parallelism N`
- [azureadtf](azureadtf) ... Terraform script using [azuread provider](https://www.terraform.io/docs/providers/azuread/)
  - Create: `time terraform apply -auto-approve -parallelism N`
  - Delete: `time terraform destroy -auto-approve -parallelism N`

Each script creates the following resources in Azure AD tenant:

- 100 users
- 10 groups
- 100 group-user member relationships
- 9 group-group member relationships

![](diag.png)
