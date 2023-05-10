package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type Granularity = types.Granularity

const (
	GranularityDaily   Granularity = "DAILY"
	GranularityMonthly Granularity = "MONTHLY"
)

type Provider interface {
	GetCosts(ctx context.Context, params ParamGetCosts) (Costs, error)
}

type ParamGetCosts struct {
	// Services
	Services []string
	// Regions
	Regions []string
	// Date start
	DateStart *time.Time
	// Date end
	DateEnd *time.Time
	// Granularity
	Granularity Granularity
}

type Costs struct {
	// Total costs
	TotalCosts float64
	// Costs by service
	ByService map[string]float64
	// Costs by region
	ByRegion map[string]float64
}
