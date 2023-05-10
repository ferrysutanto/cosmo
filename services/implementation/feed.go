package implementation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ferrysutanto/cosmo/services"
)

func (p *Provider) FeedAwsCosts(ctx context.Context, params services.ParamFeedCosts) error {
	date := params.DateStart
	end := params.DateEnd

	for date.Before(end) || date.Equal(end) {
		if err := p.FeedAwsCost(ctx, services.ParamFeedCost{
			Date:   date,
			Target: params.Target,
		}); err != nil {
			return err
		}

		log.Printf("Success feeding cost per %s\n", date.Format("2006-01-02"))
		date = date.AddDate(0, 0, 1)
	}

	return nil
}

func (p *Provider) FeedAwsCost(ctx context.Context, params services.ParamFeedCost) error {
	// targeted Date
	// if params.Date is not supplied, then it will be set to current date
	targetedDate := params.Date
	currMonthCosts, err := p.getMtdCosts(ctx, targetedDate)
	if err != nil {
		return fmt.Errorf("Failed to get cost by services per day: %w", err)
	}

	endOfLastMonthDate := getFinalDateOfLastMonth(targetedDate)
	lastMonthCosts, err := p.getMtdCosts(ctx, endOfLastMonthDate)
	if err != nil {
		return fmt.Errorf("Failed to get cost by services per day: %w", err)
	}

	// get last index of currMonthCosts.ResultsByTime
	lastIndex := len(currMonthCosts.ResultsByTime) - 1
	// get currMonthCosts.ResultsByTime[lastIndex]
	targetedDateCosts := currMonthCosts.ResultsByTime[lastIndex]

	// loop per services in form of targetedDateCosts.Groups
	for _, service := range targetedDateCosts.Groups {
		// cost converted to float64
		cost, err := strconv.ParseFloat(*service.Metrics["UnblendedCost"].Amount, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse cost to float64: %w", err)
		}
		cost = roundTo2DecimalPlaces(cost)

		// get lastMonthSameDateCost for the service {
		sameDateLastMonth := getEquivalentDateLastMonth(targetedDate)
		sameDateLastMonthCost, err := getServiceExactDateCost(lastMonthCosts, sameDateLastMonth, service.Keys)
		if err != nil {
			return fmt.Errorf("Failed to get service exact date cost: %w", err)
		}

		// sameDateLastMonthCost converted to float64
		sameDateLastMonthCostFloat := float64(0)
		if sameDateLastMonthCost.Metrics != nil {
			sameDateLastMonthCostFloat, err = strconv.ParseFloat(*sameDateLastMonthCost.Metrics["UnblendedCost"].Amount, 64)
			if err != nil {
				return fmt.Errorf("Failed to parse same date last month cost to float64: %w", err)
			}
		}
		sameDateLastMonthCostFloat = roundTo2DecimalPlaces(sameDateLastMonthCostFloat)

		// get percentage change between sameDateLastMonthCostFloat and cost
		costDiff := getPercentageChange(sameDateLastMonthCostFloat, cost)

		// get monthToDateCosts for the service
		mtdCost, err := getMtdCostsByService(currMonthCosts, targetedDate, service.Keys)
		if err != nil {
			return fmt.Errorf("Failed to get service mtd costs: %w", err)
		}

		// mtdCost converted to float64
		mtdCostFloat, err := strconv.ParseFloat(*mtdCost.Metrics["UnblendedCost"].Amount, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse mtd cost to float64: %w", err)
		}
		mtdCostFloat = roundTo2DecimalPlaces(mtdCostFloat)

		// get lastMtdCost for the service
		lastMtdCost, err := getMtdCostsByService(lastMonthCosts, sameDateLastMonth, service.Keys)
		if err != nil {
			return fmt.Errorf("Failed to get service mtd costs: %w", err)
		}

		// lastMonthToDateCost converted to float64
		lastMtdCostFloat, err := strconv.ParseFloat(*lastMtdCost.Metrics["UnblendedCost"].Amount, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse last mtd cost to float64: %w", err)
		}
		lastMtdCostFloat = roundTo2DecimalPlaces(lastMtdCostFloat)

		// get percentage change between lastMtdCostFloat and mtdCostFloat
		mtdCostDiff := getPercentageChange(lastMtdCostFloat, mtdCostFloat)

		// get lastMonth averageCostPerDay for the service
		lastMonthAvgCostPerDay, err := getServiceAvgCost(lastMonthCosts, service.Keys)
		if err != nil {
			return fmt.Errorf("Failed to get service avg cost: %w", err)
		}

		// lastMonthAvgCostPerDay converted to float64
		lastMonthAvgCostPerDayFloat, err := strconv.ParseFloat(*lastMonthAvgCostPerDay.Metrics["UnblendedCost"].Amount, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse last month avg cost per day to float64: %w", err)
		}
		lastMonthAvgCostPerDayFloat = roundTo2DecimalPlaces(lastMonthAvgCostPerDayFloat)

		// get this month currMonthAvgCostPerDay for the service
		currMonthAvgCostPerDay, err := getServiceAvgCost(currMonthCosts, service.Keys)
		if err != nil {
			return fmt.Errorf("Failed to get service avg cost: %w", err)
		}

		// currMonthAvgCostPerDay converted to float64
		currMonthAvgCostPerDayFloat, err := strconv.ParseFloat(*currMonthAvgCostPerDay.Metrics["UnblendedCost"].Amount, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse curr month avg cost per day to float64: %w", err)
		}
		currMonthAvgCostPerDayFloat = roundTo2DecimalPlaces(currMonthAvgCostPerDayFloat)

		// get percentage change between lastMonthAvgCostPerDayFloat and currMonthAvgCostPerDayFloat
		avgCostDiff := getPercentageChange(lastMonthAvgCostPerDayFloat, currMonthAvgCostPerDayFloat)

		payload := Payload{
			Timestamp:        targetedDate,
			LinkedAccount:    service.Keys[0],
			ServiceName:      service.Keys[1],
			LastMonthDate:    sameDateLastMonth,
			Cost:             cost,
			LastMonthCost:    sameDateLastMonthCostFloat,
			CostDiff:         costDiff,
			MtdCost:          mtdCostFloat,
			LastMonthMtdCost: lastMtdCostFloat,
			MtdCostDiff:      mtdCostDiff,
			AvgCost:          currMonthAvgCostPerDayFloat,
			LastMonthAvgCost: lastMonthAvgCostPerDayFloat,
			AvgCostDiff:      avgCostDiff,
		}

		// if service.Keys[1] == "Amazon Relational Database Service" {
		// 	b, err := json.MarshalIndent(payload, "", "  ")
		// 	if err != nil {
		// 		log.Fatalf("Failed to marshal payload: %v", err)
		// 	}

		// 	log.Println(string(b))
		// }

		b, err := json.Marshal(payload)
		if err != nil {
			log.Println("error payload", payload)
			return fmt.Errorf("Failed to marshal payload: %w", err)
		}

		indexName := params.Target.Index

		if params.Target.TimestampSuffix != nil {
			indexName = fmt.Sprintf("%s-%s", indexName, targetedDate.Format(*params.Target.TimestampSuffix))
		}

		// inject payload to ES
		if _, err := p.esClient.Index(indexName, bytes.NewReader(b)); err != nil {
			return fmt.Errorf("Failed to index to ES: %w", err)
		}
	}

	return nil
}

type Payload struct {
	Timestamp        time.Time `json:"@timestamp"`
	LinkedAccount    string    `json:"linked_account"`
	ServiceName      string    `json:"service_name"`
	LastMonthDate    time.Time `json:"last_month_date"`
	Cost             float64   `json:"cost"`
	LastMonthCost    float64   `json:"last_month_cost"`
	CostDiff         float64   `json:"cost_diff"`
	MtdCost          float64   `json:"mtd_cost"`
	LastMonthMtdCost float64   `json:"last_month_mtd_cost"`
	MtdCostDiff      float64   `json:"mtd_cost_diff"`
	AvgCost          float64   `json:"avg_cost"`
	LastMonthAvgCost float64   `json:"last_month_avg_cost"`
	AvgCostDiff      float64   `json:"avg_cost_diff"`
}
