provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "kubernetes_namespace_change" "istio_system" {
  metadata {
    name = var.istio_namespace
  }
}

locals {
  chart_repo = var.istio_helm_repo
}

resource "helm_release" "istio_base" {
  name       = "istio-base"
  namespace  = var.istio_namespace
  repository = local.chart_repo
  chart      = "base"
  version    = var.istio_version
  create_namespace = false
}

resource "helm_release" "istiod" {
  name       = "istiod"
  namespace  = var.istio_namespace
  repository = local.chart_repo
  chart      = "istio-control/istio-discovery"
  version    = var.istio_version
  depends_on = [helm_release.istio_base]

  set {
    name  = "global.proxy.logLevel"
    value = var.proxy_log_level
  }

  set {
    name  = "pilot.traceSampling"
    value = var.trace_sampling_rate
  }

  set {
    name  = "global.mtls.enabled"
    value = var.enable_mtls
  }
}

resource "helm_release" "istio_ingress" {
  name       = "istio-ingress"
  namespace  = var.istio_namespace
  repository = local.chart_repo
  chart      = "gateways/istio-ingress"
  version    = var.istio_version
  depends_on = [helm_release.istiod]

  set {
    name  = "service.type"
    value = var.ingress_service_type
  }

  set {
    name  = "autoscaling.enabled"
    value = var.ingress_autoscaling_enabled
  }

  set {
    name  = "autoscaling.minReplicas"
    value = var.ingress_autoscaling_min_replicas
  }

  set {
    name  = "autoscaling.maxReplicas"
    value = var.ingress_autoscaling_max_replicas
  }

  set {
    name  = "podAnnotations.env"
    value = "production"
  }

  set_sensitive {
    name  = "meshConfig.accessLogFile"
    value = var.access_log_file
  }
}

resource "helm_release" "istio_egress" {
  name       = "istio-egress"
  namespace  = var.istio_namespace
  repository = local.chart_repo
  chart      = "gateways/istio-egress"
  version    = var.istio_version
  depends_on = [helm_release.istiod]

  set {
    name  = "service.type"
    value = "ClusterIP"
  }

  set {
    name  = "autoscaling.enabled"
    value = true
  }

  set {
    name  = "autoscaling.minReplicas"
    value = 1
  }

  set {
    name  = "autoscaling.maxReplicas"
    value = 3
  }
}
