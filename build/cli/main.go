package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ferrysutanto/cosmo/services"
	svcImp "github.com/ferrysutanto/cosmo/services/implementation"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func validate() {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		log.Fatalln("ERROR: AWS_ACCESS_KEY_ID is not set")
		return
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatalln("ERROR: AWS_SECRET_ACCESS_KEY is not set")
		return
	}

	if os.Getenv("ES_HOST") == "" {
		log.Fatalln("ERROR: ES_HOST is not set")
		return
	}

	if os.Getenv("ES_PORT") == "" {
		log.Fatalln("ERROR: ES_PORT is not set")
		return
	}

	if os.Getenv("ES_USERNAME") == "" {
		log.Fatalln("ERROR: ES_USER is not set")
		return
	}

	if os.Getenv("ES_PASSWORD") == "" {
		log.Fatalln("ERROR: ES_PASSWORD is not set")
		return
	}
}

func main() {
	validate()

	ctx := context.Background()

	command := os.Args[1]

	switch command {
	case "help":
		break
	case "feed":
		feedCosts(ctx)
		break
	default:
		log.Println("ERROR: command not found")
		return
	}
}

func feedCosts(ctx context.Context) {
	var provider services.Provider

	var err error
	provider, err = svcImp.NewProvider(ctx, svcImp.ProviderOption{
		AwsCredential: svcImp.AwsCredential{
			AccountID:       os.Getenv("AWS_ACCOUNT_ID"),
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
		EsCredential: svcImp.EsCredential{
			Host:          os.Getenv("ES_HOST"),
			Port:          os.Getenv("ES_PORT"),
			Username:      os.Getenv("ES_USERNAME"),
			Password:      os.Getenv("ES_PASSWORD"),
			SkipTLSVerify: true,
		},
	})
	if err != nil {
		log.Fatalf("ERROR: Creating Provider: %v", err)
	}

	if err := provider.FeedAwsCosts(ctx, services.ParamFeedCosts{
		// date start should be 2023-04-01
		DateStart: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
		// date end should be 2023-04-30
		DateEnd: time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
		Target: services.ObservabilityTarget{
			Index: "aws-costs",
			TimestampSuffix: func() *string {
				format := "2006.01"
				return &format
			}(),
		},
	}); err != nil {
		log.Fatalf("ERROR: Feeding AWS Costs: %v", err)
	}

	log.Println("SUCCESS: Feeding AWS Costs")
}
