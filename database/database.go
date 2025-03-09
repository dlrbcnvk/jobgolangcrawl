package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func Initialize(cfg *config.Config) *sql.DB {
	// 데이터베이스 연결 설정
	dsn := cfg.DB.Url

	env := os.Getenv("APP_ENV")
	if env == "aws_lambda" {
		secret, err := GetRDSSecret()
		if err != nil {
			log.Fatal(err)
		}
		dsn = strings.ReplaceAll(dsn, "<password>", secret)
		log.Printf("Loading cfg from %s\n", dsn)
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

func GetRDSSecret() (string, error) {
	secretName := os.Getenv("RDS_SECRET_NAME")
	if secretName == "" {
		return "", fmt.Errorf("RDS_SECRET_NAME environment variable not set")
	}
	region := "ap-northeast-2"

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(awsConfig)

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
	var secretJson string = *result.SecretString

	// 빈 맵 생성
	var jsonMap map[string]interface{}

	// JSON 문자열을 맵으로 변환
	err = json.Unmarshal([]byte(secretJson), &jsonMap)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	password := jsonMap["password"].(string)
	return password, nil
}
