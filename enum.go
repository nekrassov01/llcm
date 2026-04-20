package llcm

import (
	"encoding/json"
	"fmt"
)

// OutputType represents the type of output to render.
type OutputType int

const (
	// OutputTypeNone is the output type that means none.
	OutputTypeNone OutputType = iota

	// OutputTypeJSON is the output type that means JSON format.
	OutputTypeJSON

	// OutputTypePrettyJSON is the output type that means pretty JSON format.
	OutputTypePrettyJSON

	// OutputTypeText is the output type that means text table format.
	OutputTypeText

	// OutputTypeCompressedText is the output type that means compressed text table format.
	OutputTypeCompressedText

	// OutputTypeMarkdown is the output type that means markdown table format.
	OutputTypeMarkdown

	// OutputTypeBacklog is the output type that means backlog table format.
	OutputTypeBacklog

	// OutputTypeTSV is the output type that means tab-separated values.
	OutputTypeTSV

	// OutputTypeChart is the output type that means pie chart.
	OutputTypeChart
)

// String returns the string representation of the OutputType.
func (t OutputType) String() string {
	switch t {
	case OutputTypeNone:
		return "none"
	case OutputTypeJSON:
		return "json"
	case OutputTypePrettyJSON:
		return "prettyjson"
	case OutputTypeText:
		return "text"
	case OutputTypeCompressedText:
		return "compressedtext"
	case OutputTypeMarkdown:
		return "markdown"
	case OutputTypeBacklog:
		return "backlog"
	case OutputTypeTSV:
		return "tsv"
	case OutputTypeChart:
		return "chart"
	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the OutputType.
func (t OutputType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ParseOutputType parses a string into an OutputType.
func ParseOutputType(s string) (OutputType, error) {
	switch s {
	case OutputTypeJSON.String():
		return OutputTypeJSON, nil
	case OutputTypePrettyJSON.String():
		return OutputTypePrettyJSON, nil
	case OutputTypeText.String():
		return OutputTypeText, nil
	case OutputTypeCompressedText.String():
		return OutputTypeCompressedText, nil
	case OutputTypeMarkdown.String():
		return OutputTypeMarkdown, nil
	case OutputTypeBacklog.String():
		return OutputTypeBacklog, nil
	case OutputTypeTSV.String():
		return OutputTypeTSV, nil
	case OutputTypeChart.String():
		return OutputTypeChart, nil
	default:
		return OutputTypeNone, fmt.Errorf("unsupported output type: %q", s)
	}
}

// DesiredState represents the desired state of the log group.
type DesiredState int32

const (
	// DesiredStateNone is the desired state that means none.
	DesiredStateNone DesiredState = -1

	// DesiredStateZero is the desired state that means delete the log group.
	DesiredStateZero DesiredState = 0

	// DesiredStateOneDay is the desired state that means retain the log group for one day.
	DesiredStateOneDay DesiredState = 1

	// DesiredStateThreeDays is the desired state that means retain the log group for three days.
	DesiredStateThreeDays DesiredState = 3

	// DesiredStateFiveDays is the desired state that means retain the log group for five days.
	DesiredStateFiveDays DesiredState = 5

	// DesiredStateOneWeek is the desired state that means retain the log group for one week.
	DesiredStateOneWeek DesiredState = 7

	// DesiredStateTwoWeeks is the desired state that means retain the log group for two weeks.
	DesiredStateTwoWeeks DesiredState = 14

	// DesiredStateOneMonth is the desired state that means retain the log group for one month.
	DesiredStateOneMonth DesiredState = 30

	// DesiredStateTwoMonths is the desired state that means retain the log group for two months.
	DesiredStateTwoMonths DesiredState = 60

	// DesiredStateThreeMonths is the desired state that means retain the log group for three months.
	DesiredStateThreeMonths DesiredState = 90

	// DesiredStateFourMonths is the desired state that means retain the log group for four months.
	DesiredStateFourMonths DesiredState = 120

	// DesiredStateFiveMonths is the desired state that means retain the log group for five months.
	DesiredStateFiveMonths DesiredState = 150

	// DesiredStateSixMonths is the desired state that means retain the log group for six months.
	DesiredStateSixMonths DesiredState = 180

	// DesiredStateOneYear is the desired state that means retain the log group for one year.
	DesiredStateOneYear DesiredState = 365

	// DesiredStateThirteenMonths is the desired state that means retain the log group for thirteen months.
	DesiredStateThirteenMonths DesiredState = 400

	// DesiredStateEighteenMonths is the desired state that means retain the log group for eighteen months.
	DesiredStateEighteenMonths DesiredState = 545

	// DesiredStateTwoYears is the desired state that means retain the log group for two years.
	DesiredStateTwoYears DesiredState = 731

	// DesiredStateThreeYears is the desired state that means retain the log group for three years.
	DesiredStateThreeYears DesiredState = 1096

	// DesiredStateFiveYears is the desired state that means retain the log group for five years.
	DesiredStateFiveYears DesiredState = 1827

	// DesiredStateSixYears is the desired state that means retain the log group for six years.
	DesiredStateSixYears DesiredState = 2192

	// DesiredStateSevenYears is the desired state that means retain the log group for seven years.
	DesiredStateSevenYears DesiredState = 2557

	// DesiredStateEightYears is the desired state that means retain the log group for eight years.
	DesiredStateEightYears DesiredState = 2922

	// DesiredStateNineYears is the desired state that means retain the log group for nine years.
	DesiredStateNineYears DesiredState = 3288

	// DesiredStateTenYears is the desired state that means retain the log group for ten years.
	DesiredStateTenYears DesiredState = 3653

	// DesiredStateInfinite is the desired state that means retain the log group indefinitely.
	DesiredStateInfinite DesiredState = 9999

	// DesiredStateProtected is the desired state that means protect the log group.
	DesiredStateProtected DesiredState = 10000

	// DesiredStateUnprotected is the desired state that means unprotect the log group.
	DesiredStateUnprotected DesiredState = 10001
)

// String returns the string representation of the DesiredState.
func (t DesiredState) String() string {
	switch t {
	case DesiredStateNone:
		return "none"
	case DesiredStateZero:
		return "delete"
	case DesiredStateOneDay:
		return "1day"
	case DesiredStateThreeDays:
		return "3days"
	case DesiredStateFiveDays:
		return "5days"
	case DesiredStateOneWeek:
		return "1week"
	case DesiredStateTwoWeeks:
		return "2weeks"
	case DesiredStateOneMonth:
		return "1month"
	case DesiredStateTwoMonths:
		return "2months"
	case DesiredStateThreeMonths:
		return "3months"
	case DesiredStateFourMonths:
		return "4months"
	case DesiredStateFiveMonths:
		return "5months"
	case DesiredStateSixMonths:
		return "6months"
	case DesiredStateOneYear:
		return "1year"
	case DesiredStateThirteenMonths:
		return "13months"
	case DesiredStateEighteenMonths:
		return "18months"
	case DesiredStateTwoYears:
		return "2years"
	case DesiredStateThreeYears:
		return "3years"
	case DesiredStateFiveYears:
		return "5years"
	case DesiredStateSixYears:
		return "6years"
	case DesiredStateSevenYears:
		return "7years"
	case DesiredStateEightYears:
		return "8years"
	case DesiredStateNineYears:
		return "9years"
	case DesiredStateTenYears:
		return "10years"
	case DesiredStateInfinite:
		return "infinite"
	case DesiredStateProtected:
		return "protect"
	case DesiredStateUnprotected:
		return "unprotect"
	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the DesiredState.
func (t DesiredState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// ParseDesiredState parses a string into a DesiredState.
func ParseDesiredState(s string) (DesiredState, error) {
	switch s {
	case DesiredStateZero.String():
		return DesiredStateZero, nil
	case DesiredStateOneDay.String():
		return DesiredStateOneDay, nil
	case DesiredStateThreeDays.String():
		return DesiredStateThreeDays, nil
	case DesiredStateFiveDays.String():
		return DesiredStateFiveDays, nil
	case DesiredStateOneWeek.String():
		return DesiredStateOneWeek, nil
	case DesiredStateTwoWeeks.String():
		return DesiredStateTwoWeeks, nil
	case DesiredStateOneMonth.String():
		return DesiredStateOneMonth, nil
	case DesiredStateTwoMonths.String():
		return DesiredStateTwoMonths, nil
	case DesiredStateThreeMonths.String():
		return DesiredStateThreeMonths, nil
	case DesiredStateFourMonths.String():
		return DesiredStateFourMonths, nil
	case DesiredStateFiveMonths.String():
		return DesiredStateFiveMonths, nil
	case DesiredStateSixMonths.String():
		return DesiredStateSixMonths, nil
	case DesiredStateOneYear.String():
		return DesiredStateOneYear, nil
	case DesiredStateThirteenMonths.String():
		return DesiredStateThirteenMonths, nil
	case DesiredStateEighteenMonths.String():
		return DesiredStateEighteenMonths, nil
	case DesiredStateTwoYears.String():
		return DesiredStateTwoYears, nil
	case DesiredStateThreeYears.String():
		return DesiredStateThreeYears, nil
	case DesiredStateFiveYears.String():
		return DesiredStateFiveYears, nil
	case DesiredStateSixYears.String():
		return DesiredStateSixYears, nil
	case DesiredStateSevenYears.String():
		return DesiredStateSevenYears, nil
	case DesiredStateEightYears.String():
		return DesiredStateEightYears, nil
	case DesiredStateNineYears.String():
		return DesiredStateNineYears, nil
	case DesiredStateTenYears.String():
		return DesiredStateTenYears, nil
	case DesiredStateInfinite.String():
		return DesiredStateInfinite, nil
	case DesiredStateProtected.String():
		return DesiredStateProtected, nil
	case DesiredStateUnprotected.String():
		return DesiredStateUnprotected, nil
	default:
		return DesiredStateNone, fmt.Errorf("unsupported desired state: %q", s)
	}
}
