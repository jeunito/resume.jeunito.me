data "aws_iam_policy_document" "resume_lambda_assume_role_policy" {
  statement {
    actions = [ "sts:AssumeRole" ]
    principals {
      type = "Service"
      identifiers = [ "lambda.amazonaws.com" ]
    }
  }
}

resource "aws_iam_role" "resume_lambda_backend" {
  name = "ResumeLambdaBackend"
  assume_role_policy = "${data.aws_iam_policy_document.resume_lambda_assume_role_policy.json}"
}

data "aws_iam_policy_document" "resume_backend" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["*"]
  }
  statement {
    actions = [
      "dynamodb:PutItem",
      "dynamodb:Scan"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "resume_backend_policy" {
  name = "ResumeLambdaBackendPolicy"
  role = "${aws_iam_role.resume_lambda_backend.id}"
  policy = "${data.aws_iam_policy_document.resume_backend.json}"
}

