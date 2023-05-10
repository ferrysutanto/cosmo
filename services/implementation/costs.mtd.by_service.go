package implementation

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/aws/aws-sdk-go/aws"
)

func getMtdCostsByService(costs *costexplorer.GetCostAndUsageOutput, date time.Time, keys []string) (*types.Group, error) {
	resp := &types.Group{
		Keys: keys,
		Metrics: map[string]types.MetricValue{
			"UnblendedCost": {
				Amount: aws.String("0"),
				Unit:   aws.String("USD"),
			},
		},
	}

	for _, cost := range costs.ResultsByTime {
		for _, group := range cost.Groups {
			if group.Keys[0] == keys[0] && group.Keys[1] == keys[1] {
				for name, metric := range group.Metrics {
					if _, exist := resp.Metrics[name]; !exist {
						resp.Metrics[name] = types.MetricValue{
							Amount: metric.Amount,
							Unit:   metric.Unit,
						}
					} else {
						// convert monthToDateCosts.Metrics[name].Amount to float64
						monthToDateCostsAmount, err := strconv.ParseFloat(*resp.Metrics[name].Amount, 64)
						if err != nil {
							return nil, err
						}

						// convert metric.Amount to float64
						metricAmount, err := strconv.ParseFloat(*metric.Amount, 64)
						if err != nil {
							return nil, err
						}

						// combine both values
						combinedAmount := monthToDateCostsAmount + metricAmount

						// convert combinedAmount back to string
						combinedAmountAsStr := strconv.FormatFloat(combinedAmount, 'f', -1, 64)

						// set monthToDateCosts.Metrics[name].Amount to the combined value
						// monthToDateCosts.Metrics[name].Amount = &combinedAmountAsStr
						resp.Metrics[name] = types.MetricValue{
							Amount: &combinedAmountAsStr,
							Unit:   metric.Unit,
						}
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
