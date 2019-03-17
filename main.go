package main

import(
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "net/http"
    "encoding/json"
    "fmt"
    "gopkg.in/src-d/enry.v1"
    "time"
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "os"
)

type Commit struct {
    Url string
    Id string
    Timestamp int64
    Languages map[string]int
}

type GitCommit struct {
    Url string `json:"url"`
    Id string `json:"id"`
    Timestamp string `json:"timestamp"`
    Modified []string `json:"modified"`
    Added []string `json:"added"`
}

type GitPush struct {
    Commits []GitCommit `json:"commits"`
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    push := GitPush{}
    if err := json.Unmarshal([]byte(req.Body), &push); err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusBadRequest,
            Body: fmt.Sprintf("{\"error\": \"%s\"}", err.Error()),
        }, nil
    }

    commits := []Commit{}

    for _, commit := range push.Commits {
        t, _ := time.Parse(time.RFC3339, commit.Timestamp)
        newCommit := Commit{Url: commit.Url, Id: commit.Id, Timestamp: t.Unix(), Languages: map[string]int{} }
        for _, fileModified := range commit.Modified {
            if lang, _ := enry.GetLanguageByExtension(fileModified); lang != "" {
                newCommit.Languages[lang] += 1
            }
        }
        for _, fileAdded := range commit.Added {
            if lang, _ := enry.GetLanguageByExtension(fileAdded); lang != "" {
                newCommit.Languages[lang] += 1
            }
        }

        commits = append(commits, newCommit)
    }

   bytes, _ := json.Marshal(commits)

    sess, err := session.NewSession()

    if err != nil {
        log.Println("Error creating session:")
        log.Println(err.Error())
        os.Exit(1)
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    for _, commit := range commits {
        av, err := dynamodbattribute.MarshalMap(commit)

        if err != nil {
            log.Println("Got error marshalling map:")
            log.Println(err.Error())
            os.Exit(1)
        }

        // Create item in table Movies
        input := &dynamodb.PutItemInput{
            Item: av,
            TableName: aws.String("commits"),
        }

        _, err = svc.PutItem(input)

        if err != nil {
            log.Println("Got error calling PutItem:")
            log.Println(av)
            log.Println(err.Error())
            os.Exit(1)
        }
    }

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body: string(bytes),
    }, nil
}

func main() {
    lambda.Start(HandleRequest)
}
