variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "zones" {
  type = list(string)
}

variable "name_prefix" {
  type = string
}

variable "network_self_link" {
  type = string
}

variable "subnetwork_self_link" {
  type = string
}

variable "pods_secondary_range" {
  type = string
}

variable "services_secondary_range" {
  type = string
}

variable "gke_service_account" {
  type = string
}
