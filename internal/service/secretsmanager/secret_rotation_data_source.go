// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_secretsmanager_secret_rotation")
func DataSourceSecretRotation() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSecretRotationRead,

		Schema: map[string]*schema.Schema{
			"rotation_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rotation_lambda_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"automatically_after_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"secret_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(1, 2048),
				Required:     true,
			},
		},
	}
}

func dataSourceSecretRotationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerConn()

	secretID := d.Get("secret_id").(string)
	output, err := FindSecretByID(ctx, conn, secretID)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Secrets Manager Secret (%s): %s", secretID, err)
	}

	d.SetId(aws.StringValue(output.ARN))
	d.Set("rotation_enabled", output.RotationEnabled)
	d.Set("rotation_lambda_arn", output.RotationLambdaARN)
	if err := d.Set("rotation_rules", flattenRotationRules(output.RotationRules)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting rotation_rules: %s", err)
	}

	return diags
}
