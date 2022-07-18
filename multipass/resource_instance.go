package multipass

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-go/tftypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/larstobi/go-multipass/multipass"
)

type resourceInstanceType struct{}

// Instance Resource schema
func (r resourceInstanceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Description: "Multipass instance resource.",
        Attributes: map[string]tfsdk.Attribute{
            "name": {
                Description: "Name for the instance. If it is 'primary' " +
                    "(the configured primary instance name), the user's " +
                    "home directory is mounted inside the newly launched " +
                    "instance, in 'Home'.",
                Type:     types.StringType,
                Required: true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "image": {
                Description: "Optional image to launch. If omitted, then " +
                    "the default Ubuntu LTS will be used. <remote> can be " +
                    "either ‘release’ or ‘daily‘. If <remote> is " +
                    "omitted, ‘release’ will be used. <image> can be a " +
                    "partial image hash or an Ubuntu release version, " +
                    "codename or alias. <url> is a custom image URL " +
                    "that is in http://, https://, or file:// format.",
                Type:     types.StringType,
                Optional: true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "cpus": {
                Description: "Number of CPUs to allocate. Minimum: 1, default: 1.",
                Type:        types.NumberType,
                Optional:    true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "memory": {
                Description: "Amount of memory to allocate. Positive integers, " +
                    "in bytes, or with K, M, G suffix. Minimum: 128M, default: 1G.",
                Type:     types.StringType,
                Optional: true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "disk": {
                Description: "Disk space to allocate. Positive integers, in bytes, " +
                    "or with K, M, G suffix. Minimum: 512M, default: 5G.",
                Type:     types.StringType,
                Optional: true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "cloudinit_file": {
                Description: "Path to a user-data cloud-init configuration.",
                Type:        types.StringType,
                Optional:    true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
        },
    }, nil
}

// New resource instance
func (r resourceInstanceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    return resourceInstance{
        p: *(p.(*provider)),
    }, nil
}

type Instance struct {
    Name          types.String `tfsdk:"name"`
    Image         types.String `tfsdk:"image"`
    CPUS          types.Number `tfsdk:"cpus"`
    Memory        types.String `tfsdk:"memory"`
    Disk          types.String `tfsdk:"disk"`
    CloudInitFile types.String `tfsdk:"cloudinit_file"`
}

type resourceInstance struct {
    p provider
}

// Create a new resource
func (r resourceInstance) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

    // Retrieve values from plan
    var plan Instance
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    tflog.Info(ctx, "Multipass resourceInstance", map[string]interface{}{
        "name": plan.Name.String(),
    })

    _, err := multipass.Launch(&multipass.LaunchReq{
        Name:          plan.Name.Value,
        Image:         plan.Image.Value,
        CPUS:          plan.CPUS.Value.String(),
        Memory:        plan.Memory.Value,
        Disk:          plan.Disk.Value,
        CloudInitFile: plan.CloudInitFile.Value,
    })

    if err != nil {
        resp.Diagnostics.AddError(
            "Error from multipass",
            "Could not create instance, unexpected error: "+err.Error(),
        )
        return
    }

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r resourceInstance) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
}

func (r resourceInstance) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

func (r resourceInstance) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

    var state Instance
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Delete Instance by calling API
    err := multipass.Delete(&multipass.DeleteRequest{Name: state.Name.Value})
    if err != nil {
        resp.Diagnostics.AddError(
            "Error from multipass",
            "Could not delete instance "+state.Name.Value+": "+err.Error(),
        )
        return
    }

    // Remove resource from state
    resp.State.RemoveResource(ctx)

}

// Import resource
func (r resourceInstance) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
    // Save the import identifier in the id attribute
    tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
