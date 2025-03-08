package database

import (
	"database/sql"
	"jobgolangcrawl/config"
	"log"
	"os"
	"strings"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	_ "github.com/go-sql-driver/mysql"
)

func Initialize(config *config.Config) *sql.DB {
	// 데이터베이스 연결 설정
	dsn := config.DB.Url

	env := os.Getenv("APP_ENV")
	if env == "aws_lambda" {
		secret, err := getRDSSecret()
		if err != nil {
			log.Fatal(err)
		}
		dsn = strings.ReplaceAll(dsn, "<password>", secret)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 데이터베이스 연결 확인
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getRDSSecret() (string, error) {
	secretName := "rds!db-f47dff62-66bf-44f5-b0d9-d3b3d35b934b"
	region := "ap-northeast-2"

	config, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString

	return secretString, nil
}
