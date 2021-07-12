locals {
  description  = "Database password rotator."
  package_file = "${path.module}/build/package.zip"
}

module "rotator_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  attach_network_policy             = true
  cloudwatch_logs_kms_key_id        = var.cloudwatch_logs_kms_key_id
  cloudwatch_logs_retention_in_days = var.cloudwatch_logs_retention_in_days
  cloudwatch_logs_tags              = var.cloudwatch_logs_tags
  create_package                    = false
  create_role                       = var.create_role
  description                       = var.description
  function_name                     = var.name
  handler                           = "notneeded"
  kms_key_arn                       = var.kms_key_arn
  lambda_role                       = var.lambda_role
  local_existing_package            = local.package_file
  memory_size                       = 128
  role_description                  = var.role_description
  role_name                         = var.role_name
  role_path                         = var.role_path
  role_permissions_boundary         = var.role_permissions_boundary
  role_tags                         = var.role_tags
  runtime                           = "provided.al2"
  tags                              = var.tags
  timeout                           = 180
  use_existing_cloudwatch_log_group = var.use_existing_cloudwatch_log_group
  vpc_security_group_ids            = var.vpc_security_group_ids
  vpc_subnet_ids                    = var.vpc_subnet_ids
}

resource "aws_iam_role_policy" "rotator_lambda" {
  name   = "secrets-manager-access"
  policy = data.aws_iam_policy_document.rotator_lambda.json
  role   = module.rotator_lambda.lambda_role_name
}

resource "aws_lambda_permission" "allow_secrets_manager" {
  action        = "lambda:InvokeFunction"
  function_name = module.rotator_lambda.lambda_function_name
  principal     = "secretsmanager.amazonaws.com"
  statement_id  = "AllowExecutionFromSecretsManager"
}
