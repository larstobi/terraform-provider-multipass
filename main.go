package main

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "log"
    "terraform-provider-multipass/multipass"
)

func main() {
    opts := providerserver.ServeOpts{
        Address: "registry.terraform.io/larstobi/multipass",
    }

    err := providerserver.Serve(context.Background(), multipass.New, opts)

    if err != nil {
        log.Fatal(err.Error())
    }
}
