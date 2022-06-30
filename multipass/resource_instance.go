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
                Description: "Name for the instance. If it is 'primary' " +
                    "(the configured primary instance name), the user's " +
                    "home directory is mounted inside the newly launched " +
                    "instance, in 'Home'.",
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "image": {
                Description: "Optional image to launch. If omitted, then " +
                    "the default Ubuntu LTS will be used. <remote> can be " +
                    "either ‘release’ or ‘daily‘. If <remote> is " +
                    "omitted, ‘release’ will be used. <image> can be a " +
                    "partial image hash or an Ubuntu release version, " +
                    "codename or alias. <url> is a custom image URL " +
                    "that is in http://, https://, or file:// format.",
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
            },
            "cpus": {
                Description: "Number of CPUs to allocate. Minimum: 1, default: 1.",
                Type:        schema.TypeInt,
                Optional:    true,
                ForceNew:    true,
            },
            "memory": {
                Description: "Amount of memory to allocate. Positive integers, " +
                    "in bytes, or with K, M, G suffix. Minimum: 128M, default: 1G.",
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
            },
            "disk": {
                Description: "Disk space to allocate. Positive integers, in bytes, " +
                    "or with K, M, G suffix. Minimum: 512M, default: 5G.",
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
            },
            "cloudinit_file": {
                Description: "Path to a user-data cloud-init configuration.",
                Type:        schema.TypeString,
                Optional:    true,
                ForceNew:    true,
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
    d.SetId(name)

    _, err := multipass.Launch(&multipass.LaunchReq{
        Name:          name,
        Image:         d.Get("image").(string),
        CPU:           d.Get("cpus").(int),
        Memory:        d.Get("memory").(string),
        Disk:          d.Get("disk").(string),
        CloudInitFile: d.Get("cloudinit_file").(string),
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
