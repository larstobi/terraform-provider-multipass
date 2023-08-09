Example of the use of The Terraform Provider for Multipass Hypervisor.
------------------------------------------------------------------------------------------
 Url - https://registry.terraform.io/providers/larstobi/multipass/1.4.2
 By - Lars Tobias Skjong-Børsting

 Multipass Hypervisor created by Cononical [Ubuntu] 
 Url - https://multipass.run/

 Terraform example provided by Robert Weaver

What is this Example
---------------------

This Terraform Module will create 2 example VM's within the multipass hypervisor
And will use the cloudinit to install the below package on each of the two servers apon bootup
As well as ensuringin the server is updated, and SSH keys are installed. 

Packages to be installed. 
---------------------------
- Apache2
- tmux
- nmon 

Terraform Module
-----------------------

── README.md                       < This Readme>
├── main.tf                             < main calling module>  
├── multipass_module             < This is the main module folder>
│   ├── main.tf                  < Main terraform Module >
│   ├── provider.tf              < Link to provider in the main root folder>
│   └── vars.tf                  < Variables passed to the main multipass module> 
├── provider.tf                  < Main provider for multipass version 1.4.2>
├── user_data.cfg                < Bootstrap installation of packages , ssh keys and upgrade of VM>
└── variables.tf                 < Variables used by the module , and default settings>

Module Configuration
----------------------
`resource "multipass_instance" "multipass_vm" {
    count           = var.instance_count
    cloudinit_file  = "${path.module}/user_data.cfg"
    name            = vm.name
    cpus            = var.cpus
    memory          = var.memory
    disk            = var.disks
    image           = var.image_name
}
`
## running the terraform plan

- cd into the main root directory
- terraform init
- terraform validate
- terraform plan
- terraform apply  < answer yes to proceed >
