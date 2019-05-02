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

resource "aws_lambda_event_source_mapping" "generator" {
  event_source_arn = "${aws_dynamodb_table.commits.stream_arn}"
  function_name = "${aws_lambda_function.generator.arn}"
  starting_position = "LATEST"
}
