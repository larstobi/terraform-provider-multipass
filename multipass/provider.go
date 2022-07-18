package multipass

import (
    "context"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
    return &provider{}
}

type provider struct{}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
    return tfsdk.Schema{
        Attributes: map[string]tfsdk.Attribute{},
    }, nil
}

type providerData struct{}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
    // Retrieve provider data from configuration
    var config providerData
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
    return map[string]tfsdk.ResourceType{
        "multipass_instance": resourceInstanceType{},
    }, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
    return map[string]tfsdk.DataSourceType{}, nil
}
