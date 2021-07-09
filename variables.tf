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

variable "description" {
    default     = "Database password rotator."
    description = "Description of the rotator Lambda function."
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
