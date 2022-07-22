package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/larstobi/go-multipass/multipass"
)

var _ tfsdk.DataSourceType = instanceDataSourceType{}
var _ tfsdk.DataSource = instanceDataSource{}

type instanceDataSourceType struct{}

func (t instanceDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Instance data source",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Instance name",
				Type:                types.StringType,
				Required:            true,
			},
			"ipv4": {
				MarkdownDescription: "The IPv4 address of the instance",
				Type:                types.StringType,
				Computed:            true,
			},
			"state": {
				MarkdownDescription: "The state of the instance",
				Type:                types.StringType,
				Computed:            true,
			},
			"image": {
				MarkdownDescription: "The image of the instance",
				Type:                types.StringType,
				Computed:            true,
			},
			"image_hash": {
				MarkdownDescription: "The image_hash of the instance",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t instanceDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return instanceDataSource{
		provider: provider,
	}, diags
}

type instanceDataSourceData struct {
	Name      types.String `tfsdk:"name"`
	IPv4      types.String `tfsdk:"ipv4"`
	State     types.String `tfsdk:"state"`
	Image     types.String `tfsdk:"image"`
	ImageHash types.String `tfsdk:"image_hash"`
}

type instanceDataSource struct {
	provider provider
}

func (d instanceDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data instanceDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := multipass.Info(&multipass.InfoRequest{Name: data.Name.Value})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance, got error: %s", err))
		return
	}

	data.IPv4 = types.String{Value: instance.IP}
	data.State = types.String{Value: instance.State}
	data.Image = types.String{Value: instance.Image}
	data.ImageHash = types.String{Value: instance.ImageHash}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
