variable "generator_key" {
  type = "string"
}

resource "aws_lambda_function" "generator" {
  function_name = "generator"

  handler = "main"
  runtime = "go1.x"

  role = "${aws_iam_role.resume_lambda_backend.arn}"

  s3_bucket =  "${var.binary_bucket}"
  s3_key = "${var.generator_key}"

  environment {
    variables = {
      "WEBSITE_BUCKET" = "${var.s3_bucket}"
    }
  }
}
