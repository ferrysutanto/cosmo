package implementation

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ferrysutanto/cosmo/services"
)

type Provider struct {
	awsClient *costexplorer.Client
	esClient  *elasticsearch.Client
}

type AwsCredential struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
}

type EsCredential struct {
	Host          string
	Port          string
	Username      string
	Password      string
	SkipTLSVerify bool
}

type ProviderOption struct {
	AwsCredential AwsCredential
	EsCredential  EsCredential
}

func newAwsClient(ctx context.Context, opt AwsCredential) (*costexplorer.Client, error) {
	awsCred := credentials.NewStaticCredentialsProvider(opt.AccessKeyID, opt.SecretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(awsCred))
	if err != nil {
		return nil, err
	}

	awsClient := costexplorer.NewFromConfig(cfg)

	return awsClient, nil
}

func newEsClient(ctx context.Context, opt EsCredential) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("%s:%s", opt.Host, opt.Port)},
		Username:  opt.Username,
		Password:  opt.Password,
	}

	if opt.SkipTLSVerify {
		cfg.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return esClient, nil
}

func NewProvider(ctx context.Context, opts ProviderOption) (*Provider, error) {
	awsClient, err := newAwsClient(ctx, opts.AwsCredential)
	if err != nil {
		return nil, err
	}

	esClient, err := newEsClient(ctx, opts.EsCredential)
	if err != nil {
		return nil, err
	}

	resp := &Provider{
		awsClient: awsClient,
		esClient:  esClient,
	}

	return resp, nil
}

func (p *Provider) getCostByServicesPerDay(ctx context.Context, param services.CostSource) (*costexplorer.GetCostAndUsageOutput, error) {
	return p.awsClient.GetCostAndUsage(ctx, &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityDaily,
		GroupBy: []types.GroupDefinition{
			{
				Key:  aws.String("LINKED_ACCOUNT"),
				Type: types.GroupDefinitionTypeDimension,
			},
			{
				Key:  aws.String("SERVICE"),
				Type: types.GroupDefinitionTypeDimension,
			},
		},
		Metrics: []string{"UnblendedCost", "BlendedCost", "UsageQuantity", "NormalizedUsageAmount", "AmortizedCost"},
		TimePeriod: &types.DateInterval{
			Start: aws.String(param.DateStart.Format("2006-01-02")),
			End:   aws.String(param.DateEnd.Format("2006-01-02")),
		},
	})
}
