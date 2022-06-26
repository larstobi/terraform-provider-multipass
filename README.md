# Terraform Provider for Canonical Multipass


The Terraform Multipass provider is a plugin for Terraform that allows for the full lifecycle management of instances using the Multipass virtual machine manager.

## Example

```hcl
resource "multipass_instance" "test" {
  name  = "instance-1234"
  cpus  = 2
  image = "jammy"
}

provider "multipass" {}

terraform {
  required_providers {
    multipass = {
      source  = "larstobi/multipass"
      version = "~> 1.0.0"
    }
  }
}
```

## Contributing

I appreciate your help! To contribute, please fork the repository, push your changes and make a Pull Request.
