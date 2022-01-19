#!/bin/bash
export AWS_DEFAULT_REGION="us-east-2" 
terraform "$@" -var-file "env-vars/project.tfvars" -var-file "env-vars/project-prd.tfvars"