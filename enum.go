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

// FilterKey represents the key to filter by.
type FilterKey int

const (
	FilterKeyNone      FilterKey = iota // The key meaning none.
	FilterKeyName                       // The key meaning "name" corresponding LogGroupName.
	FilterKeySource                     // The key meaning "source" corresponding Source.
	FilterKeyClass                      // The key meaning "class" corresponding Class.
	FilterKeyElapsed                    // The key meaning "elapsed" corresponding ElapsedDays.
	FilterKeyRetention                  // The key meaning "retention" corresponding RetentionInDays.
	FilterKeyBytes                      // The key meaning "bytes" corresponding StoredBytes.
)

// String returns the string representation of the FilterKey.
func (t FilterKey) String() string {
	switch t {
	case FilterKeyNone:
		return "none"
	case FilterKeyName:
		return "name"
	case FilterKeySource:
		return "source"
	case FilterKeyClass:
		return "class"
	case FilterKeyElapsed:
		return "elapsed"
	case FilterKeyRetention:
		return "retention"
	case FilterKeyBytes:
		return "bytes"
	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the FilterKey.
func (t FilterKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// parseFilterKey parses a string into a FilterKey.
func parseFilterKey(s string) (FilterKey, error) {
	switch s {
	case FilterKeyName.String():
		return FilterKeyName, nil
	case FilterKeySource.String():
		return FilterKeySource, nil
	case FilterKeyClass.String():
		return FilterKeyClass, nil
	case FilterKeyElapsed.String():
		return FilterKeyElapsed, nil
	case FilterKeyRetention.String():
		return FilterKeyRetention, nil
	case FilterKeyBytes.String():
		return FilterKeyBytes, nil
	default:
		return FilterKeyNone, fmt.Errorf("unsupported filter key: %q", s)
	}
}

// FilterOperator represents the operator to filter by.
type FilterOperator int

const (
	FilterOperatorNone  FilterOperator = iota // The operator meaning none.
	FilterOperatorGT                          // The operator meaning "greater than".
	FilterOperatorGTE                         // The operator meaning "greater than or equal to".
	FilterOperatorLT                          // The operator meaning "less than".
	FilterOperatorLTE                         // The operator meaning "less than or equal to".
	FilterOperatorEQ                          // The operator meaning "equal".
	FilterOperatorEQI                         // The operator meaning "equal" (case-insensitive).
	FilterOperatorNEQ                         // The operator meaning "not equal".
	FilterOperatorNEQI                        // The operator meaning "not equal" (case-insensitive).
	FilterOperatorREQ                         // The operator meaning regular expression match.
	FilterOperatorREQI                        // The operator meaning regular expression match (case-insensitive).
	FilterOperatorNREQ                        // The operator meaning regular expression unmatch.
	FilterOperatorNREQI                       // The operator meaning regular expression unmatch (case-insensitive).

)

// String returns the string representation of the FilterOperator.
func (t FilterOperator) String() string {
	switch t {
	case FilterOperatorNone:
		return "none"
	case FilterOperatorGT:
		return ">"
	case FilterOperatorGTE:
		return ">="
	case FilterOperatorLT:
		return "<"
	case FilterOperatorLTE:
		return "<="
	case FilterOperatorEQ:
		return "=="
	case FilterOperatorEQI:
		return "==*"
	case FilterOperatorNEQ:
		return "!="
	case FilterOperatorNEQI:
		return "!=*"
	case FilterOperatorREQ:
		return "=~"
	case FilterOperatorREQI:
		return "=~*"
	case FilterOperatorNREQ:
		return "!~"
	case FilterOperatorNREQI:
		return "!~*"
	default:
		return ""
	}
}

// MarshalJSON returns the JSON representation of the FilterOperator.
func (t FilterOperator) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// parseFilterOperator parses a string into a FilterOperator.
func parseFilterOperator(s string) (FilterOperator, error) {
	switch s {
	case FilterOperatorGT.String():
		return FilterOperatorGT, nil
	case FilterOperatorGTE.String():
		return FilterOperatorGTE, nil
	case FilterOperatorLT.String():
		return FilterOperatorLT, nil
	case FilterOperatorLTE.String():
		return FilterOperatorLTE, nil
	case FilterOperatorEQ.String():
		return FilterOperatorEQ, nil
	case FilterOperatorEQI.String():
		return FilterOperatorEQI, nil
	case FilterOperatorNEQ.String():
		return FilterOperatorNEQ, nil
	case FilterOperatorNEQI.String():
		return FilterOperatorNEQI, nil
	case FilterOperatorREQ.String():
		return FilterOperatorREQ, nil
	case FilterOperatorREQI.String():
		return FilterOperatorREQI, nil
	case FilterOperatorNREQ.String():
		return FilterOperatorNREQ, nil
	case FilterOperatorNREQI.String():
		return FilterOperatorNREQI, nil
	default:
		return FilterOperatorNone, fmt.Errorf("unsupported filter operator: %q", s)
	}
}
