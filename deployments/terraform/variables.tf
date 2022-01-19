# #########################################
# Misc
# #########################################
variable "stage" {
  description = "Environment i.e. dev"
}
variable "full_name" {
  description = "Fullname, will add as a tag to resources"
}
variable "name" {
  description = "Name which allows identifying resources"
}

# #########################################
# Networking
# #########################################

variable "public_subnets_cidr" {
  description = "Public subnet CIDRs"
  default     = ["10.0.128.0/24", "10.0.129.0/24"]
}

variable "private_subnets_cidr" {
  description = "Private subnet CIDRs"
  default     = ["10.0.0.0/18", "10.0.64.0/18"]
}

variable "database_subnets_cidr" {
  description = "Database subnet CIDRs"
  default     = ["10.0.145.0/24", "10.0.146.0/24"]
}

variable "vpc_cidr" {
  description = "VPC CIDR"
  default     = "10.0.0.0/16"
}

# #########################################
# Kubernetes
# #########################################
variable "kubernetes_cluster_version" {
  description = "kubernetes cluster version"
}
variable "kubernetes_instance_type" {
  type        = set(string)
  description = "List of instace type for ASG"
}
variable "kubernetes_asg_max_size" {
  description = "Maximum size for ASG"
}

variable "capacity_type" {
  description = "Capacity type"
  default     = "SPOT"
}

# #########################################
# Github
# #########################################
variable "github_repository_name" {
  description = "Github repository name"
}
