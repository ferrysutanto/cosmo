package implementation

import (
	"reflect"
	"testing"
	"time"
)

func Test_getFinalDateOfLastMonth(t *testing.T) {
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"success", args{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}, time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFinalDateOfLastMonth(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFinalDateOfLastMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBeginningDateOfTheMonth(t *testing.T) {
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"success", args{time.Date(2020, 1, 21, 0, 0, 0, 0, time.UTC)}, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBeginningDateOfTheMonth(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBeginningDateOfTheMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBeginningDateOfLastMonth(t *testing.T) {
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"success", args{time.Date(2020, 1, 21, 0, 0, 0, 0, time.UTC)}, time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBeginningDateOfLastMonth(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBeginningDateOfLastMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLastMonthDate(t *testing.T) {
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"success", args{time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC)}, time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEquivalentDateLastMonth(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLastMonthDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roundFloat64ToTwoDecimalPlaces(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"success 1", args{1.23456789}, 1.23},
		{"success 2", args{1.23556789}, 1.24},
		{"success 3", args{4.23556789}, 4.24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundTo2DecimalPlaces(tt.args.f); got != tt.want {
				t.Errorf("roundFloat64ToTwoDecimalPlaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPercentageChange(t *testing.T) {
	type args struct {
		old float64
		new float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"success 1", args{1.0, 2.0}, 100.0},
		{"success 2", args{2.0, 1.0}, -50.0},
		{"success 3", args{0.0, 1.0}, 100.0},
		{"success 4", args{1.0, 0.0}, -100.0},
		{"success 5", args{0.0, 0.0}, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPercentageChange(tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("getPercentageChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
