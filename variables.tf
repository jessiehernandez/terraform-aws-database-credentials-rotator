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

variable "tags" {
  default     = {}
  description = "Tags to apply to the Lambda function."
  type        = map(string)
}

variable "vpc_security_group_ids" {
  default     = null
  description = "List of security group IDs when the rotator function should run inside a VPC."
  type        = list(string)
}
