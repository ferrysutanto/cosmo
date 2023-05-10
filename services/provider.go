package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type Provider interface {
	// Feed AWS costs into Observability platform
	FeedAwsCosts(ctx context.Context, params ParamFeedCosts) error
	FeedAwsCost(ctx context.Context, params ParamFeedCost) error
}

type Granularity = types.Granularity

type CostSource struct {
	AccountID   []string
	Regions     []string
	DateStart   time.Time
	DateEnd     time.Time
	Granularity Granularity
}

type ObservabilityTarget struct {
	Host            string
	Index           string
	TimestampSuffix *string
}

type ParamFeedCost struct {
	Date   time.Time
	Target ObservabilityTarget
}

type ParamFeedCosts struct {
	DateStart time.Time
	DateEnd   time.Time
	Target    ObservabilityTarget
}
