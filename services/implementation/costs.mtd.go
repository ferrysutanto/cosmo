package implementation

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/ferrysutanto/cosmo/services"
)

func (p *Provider) getMtdCosts(ctx context.Context, date time.Time) (*costexplorer.GetCostAndUsageOutput, error) {
	beginningDate := getBeginningDateOfTheMonth(date)
	endDate := getBeginningOfNextDate(date)

	costSource := services.CostSource{
		DateStart:   beginningDate,
		DateEnd:     endDate,
		Granularity: "DAILY",
	}

	resp, err := p.getCostByServicesPerDay(ctx, costSource)
	if err != nil {
		return nil, fmt.Errorf("Failed to get cost by services per day: %w", err)
	}

	return resp, nil
}
