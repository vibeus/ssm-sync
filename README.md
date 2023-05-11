# Sync AWS SSM Parameters with Terraform Config

This is an internal tool used by Vibe team to manage secrets stored in AWS SSM using Terraform.

## Install

```shell
go install github.com/vibeus/ssm-sync@main
```

## Usage

Make sure `ssm-sync` is in your $PATH.  Then, in the `ssm` folders of Terraform configuration, run:

```shell
ssm-sync
```

Current secrets stored in SSM will be stored in `terraform.tfvars` file. Now you can update them
and use Terraform to push them back to SSM.
