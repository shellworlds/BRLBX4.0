variable "project_id" {
  description = "GCP project ID."
  type        = string
  default     = "borel-sigma-prod"
}

variable "region" {
  description = "Primary GCP region."
  type        = string
  default     = "us-central1"
}

variable "zones" {
  description = "Three zones for regional cluster and HA resources."
  type        = list(string)
  default     = ["us-central1-a", "us-central1-b", "us-central1-c"]
}

variable "environment" {
  description = "Deployment environment."
  type        = string
  default     = "prod"
}
