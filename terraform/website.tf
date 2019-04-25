resource "aws_s3_bucket" "resume" {
  bucket = "${var.s3_bucket}"
  acl = "public-read"

  website {
    index_document = "index.html"
  }
}

resource "aws_s3_bucket_policy" "resume_bucket_policy" {
  bucket = "${aws_s3_bucket.resume.id}"
  policy = "${data.aws_iam_policy_document.resume_bucket_policy_document.json}"
}

data "aws_iam_policy_document" "resume_bucket_policy_document" {
  statement {
    actions = [
      "s3:GetObject"
    ]

    resources = [
      "${aws_s3_bucket.resume.arn}/*"
    ]

    principals {
      type = "AWS"
      identifiers = [ "*" ]
    }
  }
}
