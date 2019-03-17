resource "aws_lambda_function" "classifier" {
  function_name = "classifier"
  filename = "${pathexpand("./main.zip")}"

  handler = "main"
  runtime = "go1.x"

  role = "${aws_iam_role.resume_lambda_backend.arn}"
}

resource "aws_api_gateway_resource" "commit" {
  rest_api_id = "${aws_api_gateway_rest_api.resume.id}"
  parent_id = "${aws_api_gateway_rest_api.resume.root_resource_id}"
  path_part = "commit"
}

resource "aws_api_gateway_method" "commit" {
  rest_api_id = "${aws_api_gateway_rest_api.resume.id}"
  resource_id = "${aws_api_gateway_resource.commit.id}"
  http_method = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "commit" {
  rest_api_id = "${aws_api_gateway_rest_api.resume.id}"
  resource_id = "${aws_api_gateway_resource.commit.id}"
  http_method = "${aws_api_gateway_method.commit.http_method}"

  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = "${aws_lambda_function.classifier.invoke_arn}"
}

resource "aws_api_gateway_deployment" "commit" {
  depends_on = [
    "aws_api_gateway_integration.commit"
  ]

  rest_api_id = "${aws_api_gateway_rest_api.resume.id}"
  stage_name = "1"
}

resource "aws_lambda_permission" "commit" {
  statement_id = "AllowAPIGatewayInvoke"
  action  = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.classifier.arn}"
  principal = "apigateway.amazonaws.com"
  source_arn = "${aws_api_gateway_deployment.commit.execution_arn}/*/*"
}

output "base_url" {
  url = "${aws_api_gateway_deployment.commit.invoke_url}"
}

