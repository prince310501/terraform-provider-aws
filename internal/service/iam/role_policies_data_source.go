// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iam

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// @SDKDataSource("aws_iam_role_policies", name="Role Policies")
func dataSourceRolePolicies() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRolePoliciesRead,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validRolePolicyRole,
			},
			"policy_names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRolePoliciesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).IAMClient(ctx)
	roleName := d.Get("role_name").(string)

	input := iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	var policyNames []string
	paginator := iam.NewListRolePoliciesPaginator(conn, &input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "listing IAM Role Policies (%s): %s", roleName, err)
		}
		policyNames = append(policyNames, page.PolicyNames...)
	}

	if err := d.Set("policy_names", policyNames); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting policy_names: %s", err)
	}
	d.SetId(roleName)
	return diags
}
