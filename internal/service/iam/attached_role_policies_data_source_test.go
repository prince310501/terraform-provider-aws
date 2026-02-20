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

func TestAccIAMAttachedRolePoliciesDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	datasourceName := "data.aws_iam_attached_role_policies.test"
	policyName := "test-policy-" + rName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.IAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachedRolePoliciesDataSourceConfig_basic(rName, policyName, "/"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("role_name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("policy_names"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.StringExact(policyName),
					})),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("policy_arns"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("path_prefix"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccIAMAttachedRolePoliciesDataSource_pathPrefix(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	pathPrefix := fmt.Sprintf("/custom-path-%s/", rName)
	datasourceName := "data.aws_iam_attached_role_policies.test"
	policyName := "test-policy-" + rName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.IAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachedRolePoliciesDataSourceConfig_pathPrefix(rName, policyName, pathPrefix),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("role_name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("policy_names"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.StringExact(policyName),
					})),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("policy_arns"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("path_prefix"), knownvalue.StringExact(pathPrefix)),
				},
			},
		},
	})
}

func TestAccIAMAttachedRolePoliciesDataSource_nonExistentRole(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	policyName := "test-policy-" + rName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.IAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccAttachedRolePoliciesDataSourceConfig_nonExistentRole(rName, policyName, "/"),
				ExpectError: regexache.MustCompile(`(?i)NoSuchEntity`),
			},
		},
	})
}

func TestAccIAMAttachedRolePoliciesDataSource_nonExistentPathPrefix(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	policyName := "test-policy-" + rName
	datasourceName := "data.aws_iam_attached_role_policies.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.IAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAttachedRolePoliciesDataSourceConfig_nonExistentPathPrefix(rName, policyName, "/"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("policy_names"), knownvalue.ListExact([]knownvalue.Check{})),
					statecheck.ExpectKnownValue(datasourceName, tfjsonpath.New("policy_arns"), knownvalue.ListExact([]knownvalue.Check{})),
				},
			},
		},
	})
}

func testAccAttachedRolePoliciesBaseDataSourceConfig(rName, policyName, path string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name = %q
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "test" {
  name = %q
  path = %q
  description = "Test policy for attachment"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = "ec2:DescribeInstances"
      Resource = "*"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "test" {
  role       = aws_iam_role.test.name
  policy_arn = aws_iam_policy.test.arn
}
`, rName, policyName, path)
}

func testAccAttachedRolePoliciesDataSourceConfig_basic(rName, policyName, pathPrefix string) string {
	return acctest.ConfigCompose(
		testAccAttachedRolePoliciesBaseDataSourceConfig(rName, policyName, pathPrefix),
		`
data "aws_iam_attached_role_policies" "test" {
  role_name = aws_iam_role.test.name
  depends_on = [aws_iam_role_policy_attachment.test]
}
`)
}

func testAccAttachedRolePoliciesDataSourceConfig_pathPrefix(rName, policyName, pathPrefix string) string {
	return acctest.ConfigCompose(
		testAccAttachedRolePoliciesBaseDataSourceConfig(rName, policyName, pathPrefix),
		fmt.Sprintf(`
data "aws_iam_attached_role_policies" "test" {
  role_name   = aws_iam_role.test.name
  path_prefix = %q
  depends_on = [aws_iam_role_policy_attachment.test]
}
`, pathPrefix),
	)
}

func testAccAttachedRolePoliciesDataSourceConfig_nonExistentRole(rName, policyName, pathPrefix string) string {
	return acctest.ConfigCompose(
		testAccAttachedRolePoliciesBaseDataSourceConfig(rName, policyName, pathPrefix),
		`
data "aws_iam_attached_role_policies" "test" {
  role_name = "non-existent-role"
  depends_on = [aws_iam_role_policy_attachment.test]
}
`,
	)
}

func testAccAttachedRolePoliciesDataSourceConfig_nonExistentPathPrefix(rName, policyName, pathPrefix string) string {
	return acctest.ConfigCompose(
		testAccAttachedRolePoliciesBaseDataSourceConfig(rName, policyName, pathPrefix),
		`
data "aws_iam_attached_role_policies" "test" {
  role_name   = aws_iam_role.test.name
  path_prefix = "/nonexistent-path/"
  depends_on = [aws_iam_role_policy_attachment.test]
}
`,
	)
}
