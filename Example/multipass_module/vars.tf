variable "user_data" {
  type    = string
}

variable "name" {
  type    = string
}

variable "image_name" {
    type    = string
}

variable "cpus" {
    type    = number
}

variable "memory" {
    type    = string
}

variable "disks" {
    type     = string
}

variable "instance_count" {
  description = "Number of instances to create"
  type        = number
  default     = 1
}

variable "name_prefix" {
  description = "Instance name prefix"
  type        = string
  default     = "instance"
}