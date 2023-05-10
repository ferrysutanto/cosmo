package implementation

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go/aws"
)

func getServiceExactDateCost(costs *costexplorer.GetCostAndUsageOutput, date time.Time, keys []string) (*types.Group, error) {
	// declare resp

	resp := &types.Group{
		Keys: keys,
		Metrics: map[string]types.MetricValue{
			"UnblendedCost": {
				Amount: aws.String("0"),
				Unit:   aws.String("USD"),
			},
		},
	}

	// loop costs until the exact date
	for _, cost := range costs.ResultsByTime {
		for _, group := range cost.Groups {
			if group.Keys[0] == keys[0] && group.Keys[1] == keys[1] {
				for name, metric := range group.Metrics {
					if _, exist := resp.Metrics[name]; !exist {
						resp.Metrics[name] = metric
					}
				}
			}
		}

		if *cost.TimePeriod.Start == date.Format("2006-01-02") {
			break
		}
	}

	return resp, nil
}
