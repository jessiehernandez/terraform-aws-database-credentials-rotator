output "rotation_lambda_arn" {
  description = "ARN to the rotator Lambda function."
  value       = module.rotator_lambda.function_arn
}

output "security_group_id" {
  description = "ID of the security group assigned to the rotator Lambda function."
  value       = module.rotator_lambda.security_group_id
}
