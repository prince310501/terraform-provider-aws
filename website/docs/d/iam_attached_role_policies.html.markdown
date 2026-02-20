---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_attached_role_policies"
description: |-
  Data source to list managed policies attached to an IAM role.
---

# Data Source: aws_iam_attached_role_policies

Use this data source to get the ARNs and names of all the managed policies attached to an IAM role.

## Example Usage

### All managed policies attached to a role

```terraform
data "aws_iam_attached_role_policies" "example" {
  role_name = "my_role_name"
}
```
### Filter by path prefix

```terraform
data "aws_iam_attached_role_policies" "example" {
  role_name   = "my_role_name"
  path_prefix = "/custom-path/"
}
```

## Argument Reference

This data source supports the following arguments:

* `role_name` (Required) - Name of the IAM role to list attached policies for.
* `path_prefix` - (Optional)  Path prefix for filtering the results.
  Defaults to a slash (`/`).

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `policy_names` - Set of names of the managed policies attached to the role.
* `policy_arns` - Set of ARNs of the managed policies attached to the role.