package main

import(
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "log"
    "github.com/aws/aws-lambda-go/lambda"
    "fmt"
    "os"
    "strings"
    "text/template"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Resume struct {
    Commits map[string]int
}

func HandleRequest() (string, error) {
    var b strings.Builder
    bucket := os.Getenv("WEBSITE_BUCKET")

    commits, err := Commits()
    if err != nil {
        fmt.Println(err)
        return "", err
    }

    resume := Resume{Commits: commits}

    tmpl, _ := template.New("resume").Parse(string(files["files/resume.html"]))
    tmpl.Execute(&b, resume)

    sess, err := session.NewSession()
    uploader := s3manager.NewUploader(sess)

    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String("index.html"),
        Body: strings.NewReader(b.String()),
        ContentType: aws.String("text/html"),
    })

    if err != nil {
        log.Println(("Unable to upload index.html to %q, %v"), bucket, err)
        os.Exit(1)
    }

    return b.String(), nil
}

func main() {
    lambda.Start(HandleRequest)
}
