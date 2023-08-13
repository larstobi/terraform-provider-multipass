variable "user_data" {
  description = "cloudinit_file that contains bootstrap provision commands"
  type        = string
  default     = "./user_data"
}

variable "name" {
  description = "Name of the VM your creating"
  type        = string
  default     = "dev"
}

variable "image_name" {
  description = "ubuntu image name default jammy Lts"
  type        = string
  default     = "jammy"
}

variable "cpus" {
  description = "virtual cpu count"
  type        = number
  default     = 4
}

variable "memory" {
  description = "virtual Vm memory allocation"
  type        = string
  default     = "4G"
}

variable "disks" {
  description = "Thin provisioned disk size"
  type        = string
  default     = "20G"
}

variable "instance_count" {
  description = "Number of instances to create"
  type        = number
  default     = 2
}

variable "name_prefix" {
  description = "Instance name prefix"
  type        = string
  default     = "instance"
}
