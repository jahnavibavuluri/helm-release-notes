variable "istio_version" {
  description = "The version of the Istio Helm charts to install"
  type        = string
  default     = "1.21.0"
}

variable "istio_namespace" {
  description = "Kubernetes namespace where Istio will be installed"
  type        = string
  default     = "istio-system"
}

variable "istio_helm_repo" {
  description = "Helm repository URL for Istio charts"
  type        = string
  default     = "https://istio-release.storage.googleapis.com/charts"
}

variable "enable_mtls" {
  description = "Enable mutual TLS in Istio"
  type        = bool
  default     = true
}

variable "proxy_log_level" {
  description = "Log level for Istio proxy (e.g. warning, info, debug)"
  type        = string
  default     = "debug"
}

variable "trace_sampling_rate" {
  description = "Percentage of traces to sample (0.0 - 100.0)"
  type        = number
  default     = 100.0
}

variable "ingress_autoscaling_enabled" {
  description = "Enable autoscaling for the ingress gateway"
  type        = bool
  default     = true
}

variable "ingress_autoscaling_min_replicas" {
  description = "Minimum number of ingress gateway pods"
  type        = number
  default     = 2
}

variable "ingress_autoscaling_max_replicas" {
  description = "Maximum number of ingress gateway pods"
  type        = number
  default     = 5
}

variable "ingress_service_type" {
  description = "Type of Kubernetes service for the Istio ingress gateway"
  type        = string
  default     = "LoadBalancer"
}

variable "access_log_file" {
  description = "Path to the access log file for the proxy"
  type        = string
  default     = "/dev/stdout"
  sensitive   = true
}
