# Example of Using the Terraform Provider for Multipass Hypervisor

**URL:** [Terraform Provider for Multipass Hypervisor](https://registry.terraform.io/providers/larstobi/multipass/1.4.2)  
**Author:** Lars Tobias Skjong-Børsting

Multipass Hypervisor created by Canonical [Ubuntu]  
**URL:** [Multipass](https://multipass.run/)

Terraform example provided by Robert Weaver

## What is this Example

This Terraform module will create 2 example VMs within the Multipass hypervisor and will use cloud-init to install the following packages on each of the two servers upon bootup, as well as ensuring that the server is updated and SSH keys are installed.

Packages to be installed:
- Apache2
- tmux
- nmon

## Terraform Module Structure

| File/Folder                 | Description                                     |
|-----------------------------|-------------------------------------------------|
| README.md                   | This README                                    |
| main.tf                      | Main calling module                            |
| multipass_module/           | Main module folder                             |
| multipass_module/main.tf     | Main Terraform module                          |
| multipass_module/provider.tf | Link to the provider in the main root folder   |
| multipass_module/vars.tf     | Variables passed to the main Multipass module  |
| provider.tf                 | Main provider for Multipass version 1.4.2     |
| user_data.cfg               | Bootstrap installation of packages, SSH keys, and VM upgrade |
| variables.tf                | Variables used by the module and default settings |


## Module Configuration

```hcl
resource "multipass_instance" "multipass_vm" {
    count          = var.instance_count
    cloudinit_file = "${path.module}/user_data.cfg"
    name           = var.vm_name
    cpus           = var.cpus
    memory         = var.memory
    disk           = var.disks
    image          = var.image_name
} ```

## Running the Terraform Plan

    Change directory to the main root directory.
    Run terraform init.
    Run terraform validate.
    Run terraform plan.
    Run terraform apply (answer yes to proceed).
