package implementation

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go/aws"
)

func getServiceAvgCost(costs *costexplorer.GetCostAndUsageOutput, keys []string) (*types.Group, error) {
	totalDays := float64(0)
	totalCost := float64(0)

	// loop costs and append the cost within Metrics["UnblendedCost"].Amount to totalCost
	for _, cost := range costs.ResultsByTime {
		for _, group := range cost.Groups {
			if group.Keys[0] == keys[0] && group.Keys[1] == keys[1] {
				cost, err := strconv.ParseFloat(*group.Metrics["UnblendedCost"].Amount, 64)
				if err != nil {
					return nil, err
				}

				totalCost += cost
				totalDays++
			}
		}
	}

	// get average cost per day
	avgCostPerDay := float64(0)

	if totalCost > 0 {
		avgCostPerDay = totalCost / totalDays
	}

	// declare resp
	resp := &types.Group{
		Keys: keys,
		Metrics: map[string]types.MetricValue{
			"UnblendedCost": {
				Amount: aws.String(strconv.FormatFloat(avgCostPerDay, 'f', 2, 64)),
				Unit:   aws.String("USD"),
			},
		},
	}

	return resp, nil
}
