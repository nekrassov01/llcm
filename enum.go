package llcm

import (
	"encoding/json"
	"fmt"
)

// OutputType represents the type of output to render.
type OutputType int

const (
	OutputTypeNone           OutputType = iota // The output type that means none.
	OutputTypeJSON                             // The output type of JSON format.
	OutputTypePrettyJSON                       // The output type of pretty JSON format.
	OutputTypeText                             // The output type of text table format.
	OutputTypeCompressedText                   // The output type of compressed text table format.
	OutputTypeMarkdown                         // The output type of markdown table format.
	OutputTypeBacklog                          // The output type of backlog table format.
	OutputTypeTSV                              // The output type of tab-separated values.
	OutputTypeChart                            // The output type that means pie chart.
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
	DesiredStateNone           DesiredState = -9999 // A value meaning none.
	DesiredStateZero           DesiredState = 0     // A value meaning delete the log group.
	DesiredStateOneDay         DesiredState = 1     // A value meaning retain the log group for one day.
	DesiredStateThreeDays      DesiredState = 3     // A value meaning retain the log group for three days.
	DesiredStateFiveDays       DesiredState = 5     // A value meaning retain the log group for five days.
	DesiredStateOneWeek        DesiredState = 7     // A value meaning retain the log group for one week.
	DesiredStateTwoWeeks       DesiredState = 14    // A value meaning retain the log group for two weeks.
	DesiredStateOneMonth       DesiredState = 30    // A value meaning retain the log group for one month.
	DesiredStateTwoMonths      DesiredState = 60    // A value meaning retain the log group for two months.
	DesiredStateThreeMonths    DesiredState = 90    // A value meaning retain the log group for three months.
	DesiredStateFourMonths     DesiredState = 120   // A value meaning retain the log group for four months.
	DesiredStateFiveMonths     DesiredState = 150   // A value meaning retain the log group for five months.
	DesiredStateSixMonths      DesiredState = 180   // A value meaning retain the log group for six months.
	DesiredStateOneYear        DesiredState = 365   // A value meaning retain the log group for one year.
	DesiredStateThirteenMonths DesiredState = 400   // A value meaning retain the log group for thirteen months.
	DesiredStateEighteenMonths DesiredState = 545   // A value meaning retain the log group for eighteen months.
	DesiredStateTwoYears       DesiredState = 731   // A value meaning retain the log group for two years.
	DesiredStateThreeYears     DesiredState = 1096  // A value meaning retain the log group for three years.
	DesiredStateFiveYears      DesiredState = 1827  // A value meaning retain the log group for five years.
	DesiredStateSixYears       DesiredState = 2192  // A value meaning retain the log group for six years.
	DesiredStateSevenYears     DesiredState = 255   // A value meaning retain the log group for seven years.
	DesiredStateEightYears     DesiredState = 2922  // A value meaning retain the log group for eight years.
	DesiredStateNineYears      DesiredState = 3288  // A value meaning retain the log group for nine years.
	DesiredStateTenYears       DesiredState = 3653  // A value meaning retain the log group for ten years.
	DesiredStateInfinite       DesiredState = 9999  // A value meaning retain the log group infinity.
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
	default:
		return DesiredStateNone, fmt.Errorf("unsupported desired state: %q", s)
	}
}
