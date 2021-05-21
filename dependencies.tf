data "aws_iam_policy_document" "rotator_lambda" {
  statement {
    actions = [
      "secretsmanager:DescribeSecret",
      "secretsmanager:GetSecretValue",
      "secretsmanager:PutSecretValue",
      "secretsmanager:UpdateSecretVersionStage"
    ]
    effect    = "Allow"
    resources = ["*"]
    sid       = "AllowAccessToSecretsManagedByThisRotator"

    condition {
      test     = "StringEquals"
      variable = "secretsmanager:resource/AllowRotationLambdaArn"
      values   = [module.rotator_lambda.lambda_function_arn]
    }
  }

  statement {
    actions   = ["secretsmanager:GetRandomPassword"]
    effect    = "Allow"
    resources = ["*"]
    sid       = "AllowGeneratingRandomPasswords"
  }
}

data "aws_iam_policy_document" "rotator_lambda_trust" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    sid     = "LambdaTrust"

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }
}
