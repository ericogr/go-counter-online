#!/bin/bash
# Create initial aws resources to be used by terraform
# Tested with aws cli 2.4.7

set -e

RESOURCE_NAME="online-counter"
AWS_RESOURCE_NAME="$RESOURCE_NAME-terraform"
AWS_USERNAME="$RESOURCE_NAME-terraform"
AWS_BUCKET_NAME="$RESOURCE_NAME-terraform"
AWS_DYNAMODB_TABLE_NAME="terraform-state-lock"

get_account_id() {
	echo "getting account id..."

    get_account_id=$(aws sts get-caller-identity|jq -r .Account)
}

get_access_key_id_from_file() {
	echo "getting key_id from file..."

    filename=$1
    get_access_key_id_from_file=$(cat $filename|jq -r .AccessKey.AccessKeyId)
}

create_aws_user() {
	echo "creating aws user..."

    username=$1
	aws iam create-user \
		--user-name $username >/dev/null
}

delete_aws_user() {
	echo "delete aws user..."

    username=$1
	aws iam delete-user \
		--user-name $username
}

create_aws_userkey() {
	echo "creating aws user key..."

    username=$1
    key_filename=$2
	aws iam create-access-key \
		--user-name $username>$key_filename
}

delete_aws_userkey() {
	echo "deleting aws user key..."

    username=$1
	user_key_filename=$2
	get_access_key_id_from_file $user_key_filename
	aws iam delete-access-key \
		--access-key-id $get_access_key_id_from_file \
		--user-name $username
}

create_aws_policy() {
	echo "creating aws policy..."

    policy_name=$1
    policy_filename=$2
	aws iam create-policy \
		--policy-name $policy_name \
		--policy-document file://$policy_filename >/dev/null
}

delete_aws_policy() {
	echo "deleting aws policy..."

	policy_name=$1
	account_id=$2
	aws iam delete-policy \
		--policy-arn arn:aws:iam::$account_id:policy/$policy_name
}

attach_aws_user_policy() {
	echo "attaching aws user policy..."

    policy_name=$1
    username=$2
    account_id=$3
	aws iam attach-user-policy \
	  --user-name $username \
	  --policy-arn arn:aws:iam::$account_id:policy/$policy_name
}

detach_aws_user_policy() {
	echo "detaching aws user policy..."

    policy_name=$1
	username=$2
	account_id=$3
	aws iam detach-user-policy \
		--user-name $username \
		--policy-arn arn:aws:iam::$account_id:policy/$policy_name >/dev/null
}

create_aws_s3() {
	echo "creating aws s3..."

    bucket_name=$1
	aws s3api create-bucket \
		--bucket $bucket_name \
		--region us-east-2 \
		--create-bucket-configuration LocationConstraint=us-east-2 >/dev/null
	aws s3api put-public-access-block \
		--bucket $bucket_name \
		--public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"
	aws s3api put-bucket-versioning \
		--bucket $bucket_name \
		--versioning-configuration Status=Enabled
}

delete_aws_s3() {
	echo "deleting aws s3..."

	bucket_name=$1
	aws s3api delete-bucket \
		--bucket $bucket_name
}

create_aws_dynamodb_state() {
	echo "creating aws dynamodb state..."

    table_name=$1
	aws dynamodb create-table \
		--table-name $table_name \
		--attribute-definitions AttributeName=LockID,AttributeType=S \
		--key-schema AttributeName=LockID,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST >/dev/null
}

delete_aws_dynamodb_state() {
	echo "deleting aws dynamodb state..."
	
	table_name=$1
	aws dynamodb delete-table \
        --table-name $table_name >/dev/null
}

cleanup_aws() {
	echo "Cleaning up AWS resources"

	delete_aws_dynamodb_state $AWS_DYNAMODB_TABLE_NAME
	delete_aws_s3 $AWS_BUCKET_NAME
}

create_aws() {
	echo "Creating AWS resources..."

	create_aws_s3 $AWS_BUCKET_NAME
	create_aws_dynamodb_state $AWS_DYNAMODB_TABLE_NAME
}

echo "starting..."

if [[ $1 == "create" ]]; then
	create_aws
elif [[ $1 == "cleanup" ]]; then
	cleanup_aws
else
	echo "Usage: $0 create|cleanup"
fi