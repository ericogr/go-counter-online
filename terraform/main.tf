terraform {
  backend "s3" {
    bucket         = "online-counter-terraform"
    key            = "production/state/online-counter"
    region         = "us-east-2"
    dynamodb_table = "terraform-state-lock"
    encrypt        = true
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}