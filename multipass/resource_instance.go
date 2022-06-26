package multipass

import (
    "context"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/larstobi/go-multipass/multipass"
    "log"
)

func resourceInstance() *schema.Resource {
    return &schema.Resource{
        Description: "Multipass instance resource.",

        CreateContext: resourceInstanceCreate,
        ReadContext:   resourceInstanceRead,
        DeleteContext: resourceInstanceDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "image": {
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
            },
            "cpus": {
                Type:     schema.TypeInt,
                Optional: true,
                ForceNew: true,
            },
        },
    }
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    name := d.Get("name").(string)
    d.SetId(name)

    info, err := multipass.Info(&multipass.InfoRequest{Name: name})
    if err != nil {
        log.Fatal(err)
    }

    tflog.Debug(ctx, "Multipass resourceInstanceRead", map[string]interface{}{
        "name": info.Name,
    })

    return nil
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    name := d.Get("name").(string)
    image := d.Get("image").(string)
    cpus := d.Get("cpus").(int)

    d.SetId(name)

    _, err := multipass.Launch(&multipass.LaunchReq{
        Image: image,
        Name:  name,
        CPU:   cpus,
    })

    if err != nil {
        log.Fatal(err)
    }

    return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    if err := multipass.Delete(&multipass.DeleteRequest{Name: d.Id()}); err != nil {
        log.Fatal(err)
    }

    return nil
}
