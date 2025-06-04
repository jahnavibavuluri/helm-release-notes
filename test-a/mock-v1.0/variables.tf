# Version 1 - variables.tf

variable "dns_provider" {
  description = "DNS provider for external-dns"
  type        = string
  default     = "aws"
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "my-cluster"
}

variable "log_level" {
  description = "Log level for external-dns"
  type        = string
  default     = "info"
}

variable "domain_filters" {
  description = "List of domains to manage"
  type        = list(string)
  default     = ["example.com"]
}