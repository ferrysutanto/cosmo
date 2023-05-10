package es

import (
	"context"
	"time"
)

type Provider interface {
	// Feed AWS costs into Observability platform
	FeedAwsCosts(ctx context.Context, params ParamFeedAwsCosts) error
}

type ParamFeedAwsCosts struct {
	Index     string    `json:"index" validate:"required" example:"aws-costs" format:"string"`
	AccountID string    `json:"account.id" validate:"required" example:"123456789012" format:"string"`
	Region    *string   `json:"region,omitempty" validate:"required" example:"us-east-1" format:"string"`
	Service   *string   `json:"service,omitempty" validate:"required" example:"Amazon Elastic Compute Cloud - Compute" format:"string"`
	Date      time.Time `json:"date,omitempty" validate:"required" example:"2021-01-01" format:"date"`
	Cost      float64   `json:"cost,omitempty" validate:"required" example:"100.00" format:"float"`
	Currency  string    `json:"currency,omitempty" validate:"required" example:"USD" format:"string"`
}
