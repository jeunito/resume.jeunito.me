resource "aws_s3_bucket" "resume" {
  bucket = "${var.s3_bucket}"
  acl = "private"
}
