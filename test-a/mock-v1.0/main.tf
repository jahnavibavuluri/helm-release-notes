# Version 1 - main.tf
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
  version    = "1.13.0"
  namespace  = "kube-system"

  set {
    name  = "provider"
    value = var.dns_provider
  }

  set {
    name  = "aws.zoneType"
    value = "public"
  }

  set {
    name  = "logLevel"
    value = var.log_level
  }

  set {
    name  = "interval"
    value = "1m"
  }

  set {
    name  = "registry"
    value = "txt"
  }

  set {
    name  = "txtOwnerId"
    value = var.cluster_name
  }
}

resource "helm_release" "cert_manager" {
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.12.2"
  namespace  = "cert-manager"
  
  create_namespace = true
  
  set {
    name  = "installCRDs"
    value = "true"
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