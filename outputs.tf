output "rotation_lambda_arn" {
  description = "ARN to the rotator Lambda function."
  value       = module.rotator_lambda.lambda_function_arn
}
