# ------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# ------------------------------------------------------------------------------

variable "vpc_subnet_ids" {
  description = "ID of the VPC subnets to use for the rotator Lambda."
  type        = list(string)
}

variable "vpc_id" {
  description = "ID of the VPC in which to place the rotator Lambda."
  type        = string
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
