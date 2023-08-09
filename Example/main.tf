module "multipass_vm" {
  source = "./multipass_module"

  instance_count = var.instance_count
  user_data      = "${path.module}/user_data.cfg"
  name_prefix    = "vmtf"
  name           = var.name
  image_name     = var.image_name
  cpus           = var.cpus
  memory         = var.memory
  disks          = var.disks
}
