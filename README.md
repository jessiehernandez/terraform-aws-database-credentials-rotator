# Database Credentials Rotator

This module creates a Lambda function that is able to rotate database
credentials, using the 4-step process described in detail at https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets-lambda-function-overview.html.
Only PostgreSQL and SQL Server are supported at this moment. This module will
usually be used in tandem with the `autorotated-database-credentials` module.

## Contents

* [File Structure](#file-structure)
* [Inputs](#inputs)
* [Outputs](#outputs)
* [Expected Secret Format](#expected-secret-format)
* [Usage](#usage)
* [Making Changes to the Rotator](#making-changes-to-the-rotator)
* [Author(s)](#authors)

## File Structure

```
build/    (Compiled Lambda assets will be stored here)
src/      (Rotator source code. Dependencies are vendored)
build.bsh (Compiles the Lambda assets)
```

## Inputs

| Name           | Description                          | Type     | Default | Required |
|----------------|--------------------------------------|----------|---------|----------|
| account        | Name of the account.                 | `string` |         | yes      |
| application    | Value for the application tag.       | `string` |         | yes      |
| description    | Value for the description tag.       | `string` |         | yes      |
| environment    | Environment.                         | `string` |         | yes      |
| name           | Name of the rotator Lambda function. | `string` | database-password-rotator | no |
| organization   | Organization.                        | `string` | palig                     | no |
| vpc_subnet_ids | ID of the VPC subnets to use for the rotator Lambda. | `list(string)` | | yes |
| vpc_id         | ID of the VPC in which to place the rotator Lambda.  | `string`       | | yes |

## Outputs

| Name                | Description                                                       |
|---------------------|-------------------------------------------------------------------|
| rotation_lambda_arn | ARN to the rotator Lambda function.                               |
| security_group_id   | ID of the security group assigned to the rotator Lambda function. |

## Expected Secret Format

This Lambda rotator expects the Secrets Manager secret to have the following
JSON format:

```json
{
  "dbname":   "Database name",
  "engine":   "Either 'postgres' or 'sqlserver'",
  "host":     "Database host",
  "instance": "Instance name (only useful for certain SQL Server installations)",
  "password": "Password",
  "port":     "Database port. Can be left blank for SQL Server instance installations.",
  "username": "Database username"
}
```

## Usage

The value of the `rotation_lambda_arn` output can be used as the input value for
the `rotation_lambda_arn` input in either the `autorotated-database-credentials`
module or the `aws_secretsmanager_secret_rotation` Terraform resource, as in the
examples below:

```
# When using the autorotated-database-credentials module
module "db_master_credentials" {
  source  = "jessiehernandez/autorotated-database-credentials/aws"
  version = "v1.0.0"

  ...
  rotation_lambda_arn = module.db_credentials_rotator.rotation_lambda_arn
  ...
}

# When using the Terraform resource directly
resource "aws_secretsmanager_secret_rotation" "database_credentials" {
  rotation_lambda_arn = module.rotator.rotation_lambda_arn
  secret_id           = aws_secretsmanager_secret.database_credentials.id

  rotation_rules {
    automatically_after_days = 60
  }
}
```

## Making Changes to the Rotator

If changes need to be made to the rotator source code, please ensure you first
have the `go` tools and `zip` package installed. Make your changes to the source
code, then compile the code by running the `build.bsh` script. If compilation is
successful, then the compiled assets will be placed in the `build/` directory.
The Lambda function is built from the `package.zip` file in this directory, so
it is critical that the code is compiled before deploying the Lambda function.

## Author(s)

Module was created by `Jessie Hernandez`.
