locals {
  engine = "postgres"
  engine_version = "11.10"
  family = "postgres11"
  major_engine_version = "11"
  instance_class = "db.t3.large"
  username = "go_counter_online"
  aws_rds_credentials = {
    instance_address: module.db.db_instance_address,
    instance_name: module.db.db_instance_name,
    username: module.db.db_instance_username,
    password: module.db.db_instance_password
  }
}

module "security_group_database" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4"

  name        = format("%s-sg", module.labels.id)
  description = "Go Counter Online database security group"
  vpc_id      = module.vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [for s in module.vpc.private_subnets_cidr_blocks : {
    from_port = 5432
    to_port   = 5432
    protocol  = "tcp"
    cidr_blocks = s
    description = "PostgreSQL access from within VPC private subnets"
  }]

  tags = module.labels.tags
}

resource "random_password" "db" {
  length           = 16
  special          = false
}

resource "random_password" "secret" {
  length    = 4
  min_upper = 0
  special   = false
}

module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 3.0"

  identifier = format("%s-db", module.labels.id)

  # All available versions: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html#PostgreSQL.Concepts
  engine               = local.engine
  engine_version       = local.engine_version
  family               = local.family # DB parameter group
  major_engine_version = local.major_engine_version # DB option group
  instance_class       = local.instance_class

  allocated_storage     = 10
  max_allocated_storage = 50
  storage_encrypted     = false

  # NOTE: Do NOT use 'user' as the value for 'username' as it throws:
  # "Error creating DB Instance: InvalidParameterValue: MasterUsername
  # user cannot be used as it is a reserved word used by the engine"
  name     = replace(module.labels.id, "-", "_")
  username = local.username
  password = random_password.db.result
  port     = 5432

  multi_az               = true
  subnet_ids             = module.vpc.database_subnets
  vpc_security_group_ids = [module.security_group_database.security_group_id]

  maintenance_window              = "Mon:00:00-Mon:03:00"
  backup_window                   = "03:00-06:00"
  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  backup_retention_period = 7
  skip_final_snapshot     = true
  deletion_protection     = false

  performance_insights_enabled          = true
  performance_insights_retention_period = 7
  create_monitoring_role                = true
  monitoring_interval                   = 60
  monitoring_role_name                  = format("%s-monitoring-role", module.labels.id)
  monitoring_role_description           = "Monitoring role for Go Counter Online database"

  parameters = [
    {
      name  = "autovacuum"
      value = 1
    },
    {
      name  = "client_encoding"
      value = "utf8"
    }
  ]

  tags = module.labels.tags

  db_option_group_tags = {
    "Sensitive" = "low"
  }
  db_parameter_group_tags = {
    "Sensitive" = "low"
  }
  db_subnet_group_tags = {
    "Sensitive" = "high"
  }
}

# change name every time because I need to create and destroy many times and We can't create with the same name
resource "aws_secretsmanager_secret" "db" {
  name = format("%s-db-%s", module.labels.id, random_password.secret.result)
}

resource "aws_secretsmanager_secret_version" "db" {
  secret_id     = aws_secretsmanager_secret.db.id
  secret_string = jsonencode(local.aws_rds_credentials)
}