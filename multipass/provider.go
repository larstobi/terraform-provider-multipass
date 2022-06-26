package multipass

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
    return &schema.Provider{
        ResourcesMap: map[string]*schema.Resource{
            "multipass_instance": resourceInstance(),
        },
        DataSourcesMap: map[string]*schema.Resource{},
    }
}
