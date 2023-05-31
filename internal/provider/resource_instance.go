package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/larstobi/go-multipass/multipass"
)

var _ tfsdk.ResourceType = instanceResourceType{}
var _ tfsdk.Resource = instanceResource{}
var _ tfsdk.ResourceWithImportState = instanceResource{}

type instanceResourceType struct{}

func (r instanceResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        MarkdownDescription: "Multipass instance resource.",
        Version:             0,
        Attributes: map[string]tfsdk.Attribute{
            "name": {
                MarkdownDescription: "Name for the instance. If it is 'primary' " +
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
                MarkdownDescription: "Optional image to launch. If omitted, then " +
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
                MarkdownDescription: "Number of CPUs to allocate. Minimum: 1, default: 1.",
                Type:                types.NumberType,
                Optional:            true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "memory": {
                MarkdownDescription: "Amount of memory to allocate. Positive integers, " +
                    "in KiB, MiB, GiB or TiB suffix. Minimum: 128MiB, default: 1GiB.",
                Type:     types.StringType,
                Optional: true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "disk": {
                MarkdownDescription: "Disk space to allocate. Positive integers, " +
                    "in KiB, MiB, GiB or TiB suffix. Minimum: 512MiB, default: 5GiB.",
                Type:     types.StringType,
                Optional: true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
            "cloudinit_file": {
                MarkdownDescription: "Path to a user-data cloud-init configuration.",
                Type:                types.StringType,
                Optional:            true,
                PlanModifiers: []tfsdk.AttributePlanModifier{
                    tfsdk.RequiresReplace(),
                },
            },
        },
    }, nil
}

func (t instanceResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
    provider, diags := convertProviderType(in)

    return instanceResource{
        provider: provider,
    }, diags
}

type Instance struct {
    Name          types.String `tfsdk:"name"`
    Image         types.String `tfsdk:"image"`
    CPUS          types.Number `tfsdk:"cpus"`
    Memory        types.String `tfsdk:"memory"`
    Disk          types.String `tfsdk:"disk"`
    CloudInitFile types.String `tfsdk:"cloudinit_file"`
}

type instanceResource struct {
    provider provider
}

func (r instanceResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
    // Retrieve values from plan
    var plan Instance
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    tflog.Info(ctx, "Multipass instanceResource", map[string]interface{}{
        "name": plan.Name.String(),
    })

    var cpus string
    if plan.CPUS.Null {
        cpus = ""
    } else {
        cpus = plan.CPUS.Value.String()
    }

    _, err := multipass.LaunchV2(&multipass.LaunchReqV2{
        Name:          plan.Name.Value,
        Image:         plan.Image.Value,
        CPUS:          cpus,
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

func (r instanceResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
    // Get current state
    var state Instance
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    instanceInfo, infoErr := multipass.Info(&multipass.InfoRequest{Name: state.Name.Value})
    if instanceInfo == nil || infoErr != nil {
        tflog.Warn(ctx, "Multipass instance not found, removing from state.", map[string]interface{}{
            "name": state.Name.Value,
        })
        resp.State.RemoveResource(ctx)
        return
    }

    result, err := QueryInstance(state)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error from multipass",
            "Could not query instance: "+err.Error(),
        )
        return
    }

    // Set state
    diags = resp.State.Set(ctx, result)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r instanceResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
    // multipass stop instance
    // multipass set local.instance.cpus etc
    // multipass start instance
}

func (r instanceResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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
func (r instanceResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
    // Save the import identifier in the id attribute
    tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
