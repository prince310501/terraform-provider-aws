---
subcategory: "IAM (Identity & Access Management)"
layout: "aws"
page_title: "AWS: aws_iam_role_policies"
description: |-
  Data source to list inline policies for an IAM role.
---

# Data Source: aws_iam_role_policies

Use this data source to get the names of all the inline policies of an IAM role.

## Example Usage

```terraform
data "aws_iam_role_policies" "example" {
  role_name = "my_role_name"
}
```

## Argument Reference

This data source supports the following arguments:

* `role_name` (Required) - Name of the IAM role to list inline policies for.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `policy_names` - Set of names of all the inline policies attached to the role.