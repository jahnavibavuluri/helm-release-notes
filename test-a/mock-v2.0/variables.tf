# Version 2 - variables.tf

variable "dns_provider" {
  description = "DNS provider for external-dns"
  type        = string
  default     = "aws"
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "production-cluster"  # CHANGED: Different default value
}

variable "log_level" {
  description = "Log level for external-dns"
  type        = string
  default     = "debug"  # CHANGED: Different default value
}

variable "domain_filters" {
  description = "List of domains to manage"
  type        = list(string)
  default     = ["example.com", "app.example.com"]  # CHANGED: Added additional domain
}

# ADDED: New variable
variable "zone_type" {
  description = "AWS Route53 zone type"
  type        = string
  default     = "public"
}

# ADDED: New variable
variable "sync_interval" {
  description = "Sync interval for external-dns"
  type        = string
  default     = "30s"
}

# ADDED: New variable
variable "sync_policy" {
  description = "Sync policy for external-dns"
  type        = string
  default     = "upsert-only"
}

variable "external_dns_role_arn" {
  description = "IAM role ARN for external-dns service account"
  type        = string
}