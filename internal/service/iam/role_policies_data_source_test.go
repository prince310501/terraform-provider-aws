// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iam_test

import (
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccIAMRolePoliciesDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_iam_role_policies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.IAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRolePoliciesDataSourceConfig_basic(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("role_name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("policy_names"), knownvalue.ListSizeExact(2)),
				},
			},
		},
	})
}

func TestAccIAMRolePoliciesDataSource_nonExistentRole(t *testing.T) {
	ctx := acctest.Context(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.IAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRolePoliciesDataSourceConfig_nonExistentRole,
				ExpectError: regexache.MustCompile(`(?i)NoSuchEntity`),
			},
		},
	})
}

func testAccRolePoliciesDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name = "%[1]s"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "policyA" {
  name = "policyA-%[1]s"
  role = aws_iam_role.test.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = "s3:ListBucket"
      Resource = "*"
    }]
  })
}

resource "aws_iam_role_policy" "policyB" {
  name = "policyB-%[1]s"
  role = aws_iam_role.test.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = "ec2:DescribeInstances"
      Resource = "*"
    }]
  })
}

data "aws_iam_role_policies" "test" {
  role_name = aws_iam_role.test.name
  depends_on = [aws_iam_role_policy.policyA, aws_iam_role_policy.policyB]
}
`, rName)
}

const testAccRolePoliciesDataSourceConfig_nonExistentRole = `
data "aws_iam_role_policies" "test" {
  role_name = "non-existent-role"
}
`
