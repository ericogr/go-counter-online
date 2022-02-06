module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~>3.11.0"

  name = module.labels.id
  cidr = var.vpc_cidr

  azs = [
    data.aws_availability_zones.available.names[0],
    data.aws_availability_zones.available.names[1]
  ]

  private_subnets  = var.private_subnets_cidr
  public_subnets   = var.public_subnets_cidr
  database_subnets = var.database_subnets_cidr

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
  }
  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
  }

  tags = module.labels.tags
}

resource "aws_security_group" "vpce" {
  name   = format("%s-endpoints", module.labels.id)
  vpc_id = module.vpc.vpc_id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }
  tags = module.labels.tags
}