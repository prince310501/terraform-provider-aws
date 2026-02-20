// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iam

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	awstypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	inttypes "github.com/hashicorp/terraform-provider-aws/internal/types"
)

// @SDKDataSource("aws_iam_attached_role_policies", name="Attached Role Policies")
func dataSourceAttachedRolePolicies() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceAttachedRolePoliciesRead,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validRolePolicyRole,
			},
			"path_prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validPolicyPath,
			},
			"policy_names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policy_arns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAttachedRolePoliciesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).IAMClient(ctx)
	roleName := d.Get("role_name").(string)

	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	if v, ok := d.GetOk("path_prefix"); ok {
		input.PathPrefix = aws.String(v.(string))
	}

	var results []awstypes.AttachedPolicy
	paginator := iam.NewListAttachedRolePoliciesPaginator(conn, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "listing IAM Role attached Policies (%s): %s", roleName, err)
		}
		for _, policy := range page.AttachedPolicies {
			if p := &policy; inttypes.IsZero(p) {
				continue
			}
			results = append(results, policy)
		}
	}

	var policyNames, policyArns []string
	for _, r := range results {
		policyArns = append(policyArns, aws.ToString(r.PolicyArn))
		policyNames = append(policyNames, aws.ToString(r.PolicyName))
	}

	if err := d.Set("policy_arns", policyArns); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting policy_arns: %s", err)
	}
	if err := d.Set("policy_names", policyNames); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting policy_names: %s", err)
	}
	d.SetId(roleName)
	return diags
}
