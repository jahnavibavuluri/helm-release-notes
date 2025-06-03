provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

locals {
  istio_version   = "1.21.0"
  istio_namespace = "istio-system"
  chart_repo      = "https://istio-release.storage.googleapis.com/charts"
}

resource "helm_release" "istio_base" {
  name       = "istio-base"
  namespace  = local.istio_namespace
  repository = local.chart_repo
  chart      = "base"
  version    = local.istio_version
  create_namespace = false
}

resource "helm_release" "istiod" {
  name       = "istiod"
  namespace  = local.istio_namespace
  repository = local.chart_repo
  chart      = "istio-control/istio-discovery"
  version    = local.istio_version
  depends_on = [helm_release.istio_base]

  set {
    name  = "global.proxy.logLevel"
    value = "debug"
  }

  set {
    name  = "pilot.traceSampling"
    value = "100.0"
  }

  set {
    name  = "global.mtls.enabled"
    value = "true"
  }
}

resource "helm_release" "istio_ingress" {
  name       = "istio-ingress"
  namespace  = local.istio_namespace
  repository = local.chart_repo
  chart      = "gateways/istio-ingress"
  version    = local.istio_version
  depends_on = [helm_release.istiod]

  set {
    name  = "service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "autoscaling.enabled"
    value = "true"
  }

  set {
    name  = "autoscaling.minReplicas"
    value = "2"
  }

  set {
    name  = "autoscaling.maxReplicas"
    value = "5"
  }

  set_sensitive {
    name  = "meshConfig.accessLogFile"
    value = "/dev/stdout"
  }
}
