# Sync AWS SSM Parameters with Terraform Config

This is an internal tool used by Vibe team to manage secrets stored in
AWS SSM using Terraform.

# Install

```shell
go install github.com/vibeus/ssm-sync@v1.0.2
```

# Usage

Make sure `ssm-sync` is in your $PATH.  Commands below should be run
in the `ssm` folders of Terraform configuration.

## Initial Pull

This is useful if you want to keep your local .tfvars file up to date.

```shell
ssm-sync
```

Current secrets stored in SSM will be stored in `terraform.tfvars`
file. Now you can update them and use Terraform to push them back to
SSM.

## Add New SSM Values

To add a new SSM value, edit the resource file (default `main.tf`) by
adding new `aws_ssm_parameter` resources.  Do not edit variable file
(default `variables.tf`) or .tfvars file (default `terraform.tfvars`)
as they should be auto generated.

Then, run the command below.

```shell
ssm-sync
```

You should find new entries being added to .tfvars file with dummy
initial values. Now update these values with real value and run the
command below.

```shell
terraform apply
```

## Update Existing SSM Values

To update SSM values, edit the `terraform.tfvars` file and run:

```shell
terraform apply
```
