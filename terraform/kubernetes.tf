locals {
  kubeconfig = <<KUBECONFIG
apiVersion: v1
clusters:
- cluster:
    server: ${data.aws_eks_cluster.eks.endpoint}
    certificate-authority-data: ${data.aws_eks_cluster.eks.certificate_authority.0.data}
  name: ${data.aws_eks_cluster.eks.name}
contexts:
- context:
    cluster: ${data.aws_eks_cluster.eks.name}
    user: ${data.aws_eks_cluster.eks.name}
  name: ${data.aws_eks_cluster.eks.name}
current-context: ${data.aws_eks_cluster.eks.name}
kind: Config
preferences: {}
users:
- name: ${data.aws_eks_cluster.eks.name}
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws
      args:
        - "eks"
        - "get-token"
        - "--cluster-name"
        - "${data.aws_eks_cluster.eks.name}"
        - "--region"
        - "${data.aws_region.current.name}"
KUBECONFIG
}

data "aws_region" "current" {
}

data "aws_eks_cluster" "eks" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "eks" {
  name = module.eks.cluster_id
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.eks.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.eks.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.eks.token
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~>18.0.6"

  cluster_name                    = format("%s-eks", module.labels.id)
  cluster_version                 = var.kubernetes_cluster_version
  cluster_endpoint_private_access = true
  cluster_endpoint_public_access  = true

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  enable_irsa = true

  eks_managed_node_group_defaults = {
    ami_type       = "AL2_x86_64" #https://docs.aws.amazon.com/eks/latest/APIReference/API_Nodegroup.html
    disk_size      = 50
    instance_types = var.kubernetes_instance_type
  }
  eks_managed_node_groups = {
    workers_spot = {
      min_size     = 2
      max_size     = var.kubernetes_asg_max_size
      desired_size = 2

      instance_types = var.kubernetes_instance_type
      capacity_type  = var.capacity_type
      k8s_labels = {
        Environment = var.stage
      }
      tags = {
        format("k8s.io/cluster-autoscaler/%s-eks", module.labels.id) = "owned"
        "k8s.io/cluster-autoscaler/enabled"                          = "TRUE"
      }
    }
  }

  node_security_group_additional_rules = {
    ingress_cluster_9443 = {
      description                   = "Cluster API to node groups webhook"
      protocol                      = "tcp"
      from_port                     = 9443
      to_port                       = 9443
      type                          = "ingress"
      source_cluster_security_group = true
    }

    ingress_cluster_8443 = {
      description                   = "Cluster API to node groups webhook"
      protocol                      = "tcp"
      from_port                     = 8443
      to_port                       = 8443
      type                          = "ingress"
      source_cluster_security_group = true
    }

    ingress_cluster_80 = {
      description = "Internal communcation 80"
      protocol    = "tcp"
      from_port   = 80
      to_port     = 80
      type        = "ingress"
      self        = true
    }

    engress_cluster_80 = {
      description = "Internal communcation 80"
      protocol    = "tcp"
      from_port   = 80
      to_port     = 80
      type        = "egress"
      self        = true
    }

    engress_cluster_5432 = {
      description              = "Internal communcation to postgres"
      protocol                 = "tcp"
      from_port                = 5432
      to_port                  = 5432
      type                     = "egress"
      source_security_group_id = module.security_group_database.security_group_id
    }
  }

  tags = module.labels.tags
}