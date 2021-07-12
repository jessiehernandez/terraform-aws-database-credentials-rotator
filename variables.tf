# ------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# ------------------------------------------------------------------------------

variable "vpc_subnet_ids" {
  description = "ID of the VPC subnets to use for the rotator Lambda."
  type        = list(string)
}

# ------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# ------------------------------------------------------------------------------

variable "cloudwatch_logs_kms_key_id" {
  default     = null
  description = "The ARN of the KMS key to use when encrypting log data."
  type        = string
}

variable "cloudwatch_logs_retention_in_days" {
  default     = null
  description = "Specifies the number of days you want to retain log events in the specified log group. Possible values are: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, and 3653."
  type        = number
}

variable "cloudwatch_logs_tags" {
  default     = {}
  description = "A map of tags to assign to the CloudWatch logs."
  type        = map(string)
}

variable "create_role" {
  default     = true
  description = "Controls whether an IAM role for the rotator Lambda should be created."
  type        = bool
}

variable "description" {
  default     = "Database password rotator."
  description = "Description of the rotator Lambda function."
  type        = string
}

variable "kms_key_arn" {
  default     = null
  description = "The ARN of the KMS key to use for the rotator Lambda function."
  type        = string
}

variable "lambda_role" {
  default     = ""
  description = "IAM role ARN attached to the Lambda Function. This governs both who / what can invoke your Lambda Function, as well as what resources our Lambda Function has access to. See Lambda Permission Model for more details."
  type        = string
}

variable "name" {
  default     = "database-password-rotator"
  description = "Name of the rotator Lambda function."
  type        = string
}

variable "role_description" {
  default     = null
  description = "Description of the IAM role to use for the rotator Lambda function when creating the role via this module."
  type        = string
}

variable "role_name" {
  default     = null
  description = "Name of the IAM role to use for the rotator Lambda function when creating the role via this module."
  type        = string
}

variable "role_path" {
  default     = null
  description = "Path of IAM role to use for the rotator Lambda function when creating the role via this module."
  type        = string
}

variable "role_permissions_boundary" {
  description = "The ARN of the policy that is used to set the permissions boundary for the IAM role used by the rotator Lambda function when creating the role via this module."
  type        = string
  default     = null
}

variable "role_tags" {
  default     = {}
  description = "A map of tags to assign to the IAM role when creating the role via this module."
  type        = map(string)
}

variable "tags" {
  default     = {}
  description = "Tags to apply to the Lambda function."
  type        = map(string)
}

variable "use_existing_cloudwatch_log_group" {
  default     = false
  description = "Whether to use an existing CloudWatch log group or create a new one."
  type        = bool
}

variable "vpc_security_group_ids" {
  default     = null
  description = "List of security group IDs when the rotator function should run inside a VPC."
  type        = list(string)
}
