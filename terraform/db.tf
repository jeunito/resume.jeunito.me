resource "aws_dynamodb_table" "commits" {
  name = "commits"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "Id"
  range_key = "Timestamp"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }

  stream_enabled = "true"
  stream_view_type = "KEYS_ONLY"
}
