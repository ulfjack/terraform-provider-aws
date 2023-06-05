// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vpclattice

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// Function annotations are used for resource registration to the Provider. DO NOT EDIT.
// @SDKResource("aws_vpclattice_auth_policy")
func ResourceAuthPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAuthPolicyPut,
		ReadWithoutTimeout:   resourceAuthPolicyRead,
		UpdateWithoutTimeout: resourceAuthPolicyPut,
		DeleteWithoutTimeout: resourceAuthPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_identifier": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

const (
	ResNameAuthPolicy = "Auth Policy"
)

func resourceAuthPolicyPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()
	resourceId := d.Get("resource_identifier").(string)

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))
	if err != nil {
		return diag.Errorf("policy (%s) is invalid JSON: %s", policy, err)
	}

	in := &vpclattice.PutAuthPolicyInput{
		Policy:             aws.String(policy),
		ResourceIdentifier: aws.String(resourceId),
	}

	log.Printf("[DEBUG] Putting VPCLattice Auth Policy for resource: %s", resourceId)

	_, err = conn.PutAuthPolicy(ctx, in)
	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionCreating, ResNameAuthPolicy, d.Get("policy").(string), err)
	}

	d.SetId(resourceId)

	return resourceAuthPolicyRead(ctx, d, meta)
}

func resourceAuthPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()
	resourceId := d.Id()

	log.Printf("[DEBUG] Reading VPCLattice Auth Policy for resource: %s", resourceId)

	policy, err := findAuthPolicy(ctx, conn, resourceId)
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] VPCLattice AuthPolicy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameAuthPolicy, d.Id(), err)
	}

	if policy == nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameAuthPolicy, d.Id(), err)
	}

	d.Set("resource_identifier", resourceId)

	policyToSet, err := verify.PolicyToSet(d.Get("policy").(string), aws.ToString(policy.Policy))

	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameAuthPolicy, aws.ToString(policy.Policy), err)
	}

	d.Set("policy", policyToSet)

	return nil
}

func resourceAuthPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	log.Printf("[INFO] Deleting VPCLattice AuthPolicy: %s", d.Id())
	_, err := conn.DeleteAuthPolicy(ctx, &vpclattice.DeleteAuthPolicyInput{
		ResourceIdentifier: aws.String(d.Id()),
	})

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.VPCLattice, create.ErrActionDeleting, ResNameAuthPolicy, d.Id(), err)
	}

	return nil
}

func findAuthPolicy(ctx context.Context, conn *vpclattice.Client, id string) (*vpclattice.GetAuthPolicyOutput, error) {
	in := &vpclattice.GetAuthPolicyInput{
		ResourceIdentifier: aws.String(id),
	}

	out, err := conn.GetAuthPolicy(ctx, in)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}

	return out, nil
}
