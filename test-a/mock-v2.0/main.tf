# Version 2 - main.tf
terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.8"
    }
  }
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "helm_release" "external_dns" {
  name       = "external-dns"
  repository = "https://kubernetes-sigs.github.io/external-dns/"
  chart      = "external-dns"
  version    = "1.14.3"  # CHANGED: Updated version
  namespace  = "external-dns"  # CHANGED: Different namespace
  
  create_namespace = true  # ADDED: Create namespace

  set {
    name  = "provider"
    value = var.dns_provider
  }

  set {
    name  = "aws.zoneType"
    value = var.zone_type  # CHANGED: Made configurable
  }

  set {
    name  = "logLevel"
    value = var.log_level
  }

  set {
    name  = "interval"
    value = var.sync_interval  # CHANGED: Made configurable
  }

  set {
    name  = "registry"
    value = "txt"
  }

  set {
    name  = "txtOwnerId"
    value = var.cluster_name
  }

  # ADDED: Domain filters configuration
  set {
    name  = "domainFilters"
    value = "{${join(",", var.domain_filters)}}"
  }

  # ADDED: Policy configuration
  set {
    name  = "policy"
    value = var.sync_policy
  }
}

# ADDED: New nginx-ingress helm release
resource "helm_release" "nginx_ingress" {
  name       = "nginx-ingress"
  repository = "https://kubernetes.github.io/ingress-nginx"
  chart      = "ingress-nginx"
  version    = "4.7.1"
  namespace  = "ingress-nginx"
  
  create_namespace = true
  
  set {
    name  = "controller.service.type"
    value = "LoadBalancer"
  }
}

resource "kubernetes_config_map" "external_dns_config" {
  metadata {
    name      = "external-dns-config"
    namespace = "external-dns"  
  }

  data = {
    "config.yaml" = "aws-zone-type: public\nlog-level: debug\nsync-interval: 30s"  # CHANGED: Updated content
  }
}

# ADDED: New service account resource
resource "kubernetes_service_account" "external_dns" {
  metadata {
    name      = "external-dns"
    namespace = "external-dns"
    
    annotations = {
      "eks.amazonaws.com/role-arn" = var.external_dns_role_arn
    }
  }
}