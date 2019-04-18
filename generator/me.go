package main

import (
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


type Commit struct {
    Url string
    Id string
    Timestamp int64
    Languages map[string]int
}

func Commits() (map[string]int, error) {
//    region := "us-west-2"
//    sess, err := session.NewSession(&aws.Config{
//        Endpoint: aws.String("http://localhost:8000"),
//        Region: &region,
//    })

    sess, err := session.NewSession()

    if err != nil {
        log.Println(err)
        return nil, err
    }

    dbSvc := dynamodb.New(sess)

    params := &dynamodb.ScanInput{
		TableName: aws.String("commits"),
        Limit: aws.Int64(100),
	}

    result, err := dbSvc.Scan(params)
	if err != nil {
        log.Println("failed to make Query API call, %v", err)
        return nil, err
	}

    commits := []Commit{}

    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &commits)
	if err != nil {
        log.Println("failed to unmarshal Query result items, %v", err)
	}

    commitStats := map[string]int{}

    for _, commit := range commits {
        for language, count := range commit.Languages {
            commitStats[language] += count
        }
    }

    return commitStats, nil
}
