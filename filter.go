package llcm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Filter represents a filter expression for filtering log group entries.
type Filter struct {
	Key      FilterKey      // The key of the filter expression
	Operator FilterOperator // The operator of the filter expression
	Value    string         // The value of the filter expression
}

// EvaluateFilter evaluates the filter expressions.
func EvaluateFilter(expressions []string) ([]Filter, error) {
	if len(expressions) == 0 {
		return []Filter{}, nil
	}
	filters := make([]Filter, 0, len(expressions))
	for _, expression := range expressions {
		if strings.TrimSpace(expression) == "" {
			return nil, fmt.Errorf("invalid syntax: empty expression")
		}
		tokens := strings.Fields(expression)
		if len(tokens) < 3 {
			return nil, fmt.Errorf("invalid syntax: %q", expression)
		}
		if len(tokens) > 3 {
			tokens[2] = strings.Join(tokens[2:], " ")
		}
		key, err := parseFilterKey(tokens[0])
		if err != nil {
			return nil, err
		}
		operator, err := parseFilterOperator(tokens[1])
		if err != nil {
			return nil, err
		}
		filter := Filter{
			Key:      key,
			Operator: operator,
			Value:    tokens[2],
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func (man *Manager) setFilter(filters []Filter) error {
	man.Filters = filters
	man.filterFns = make([]func(*entry) bool, 0, len(filters))
	for _, filter := range filters {
		switch filter.Key {
		case FilterKeyName, FilterKeySource, FilterKeyClass:
			fn, err := stringFilterFunc(filter)
			if err != nil {
				return err
			}
			man.filterFns = append(man.filterFns, func(e *entry) bool {
				switch filter.Key {
				case FilterKeyName:
					return fn(e.LogGroupName)
				case FilterKeySource:
					return fn(e.Source)
				case FilterKeyClass:
					return fn(string(e.Class))
				default:
					return true
				}
			})
		case FilterKeyElapsed, FilterKeyRetention, FilterKeyBytes:
			fn, err := numberFilterFunc(filter)
			if err != nil {
				return err
			}
			man.filterFns = append(man.filterFns, func(e *entry) bool {
				switch filter.Key {
				case FilterKeyElapsed:
					return fn(e.ElapsedDays)
				case FilterKeyRetention:
					return fn(e.RetentionInDays)
				case FilterKeyBytes:
					return fn(e.StoredBytes)
				default:
					return true
				}
			})
		default:
			return fmt.Errorf("invalid key: %q", filter.Key)
		}
	}
	return nil
}

func stringFilterFunc(filter Filter) (func(string) bool, error) {
	var (
		re       *regexp.Regexp
		fn       func(string) bool
		err      error
		operator = filter.Operator
		value    = filter.Value
	)
	if operator == FilterOperatorREQ || operator == FilterOperatorREQI || operator == FilterOperatorNREQ || operator == FilterOperatorNREQI {
		if operator == FilterOperatorREQI || operator == FilterOperatorNREQI {
			value = "(?i)" + value
		}
		re, err = regexp.Compile(value)
		if err != nil {
			return nil, err
		}
	}
	switch operator {
	case FilterOperatorEQ:
		fn = func(v string) bool {
			return v == value
		}
	case FilterOperatorEQI:
		fn = func(v string) bool {
			return strings.EqualFold(v, value)
		}
	case FilterOperatorNEQ:
		fn = func(v string) bool {
			return v != value
		}
	case FilterOperatorNEQI:
		fn = func(v string) bool {
			return !strings.EqualFold(v, value)
		}
	case FilterOperatorREQ, FilterOperatorREQI:
		fn = func(v string) bool {
			return re.MatchString(v)
		}
	case FilterOperatorNREQ, FilterOperatorNREQI:
		fn = func(v string) bool {
			return !re.MatchString(v)
		}
	default:
		return nil, fmt.Errorf("invalid operator: %q", operator)
	}
	return fn, nil
}

func numberFilterFunc(filter Filter) (func(int64) bool, error) {
	var (
		operator = filter.Operator
		value    = filter.Value
	)
	n, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		if filter.Key == FilterKeyRetention {
			d, err := ParseDesiredState(value)
			if err != nil || d <= 0 {
				return nil, fmt.Errorf("invalid desired state: %q", value)
			}
			n = int64(d)
		} else {
			return nil, err
		}
	}
	var fn func(int64) bool
	switch operator {
	case FilterOperatorGT:
		fn = func(v int64) bool {
			return v > n
		}
	case FilterOperatorGTE:
		fn = func(v int64) bool {
			return v >= n
		}
	case FilterOperatorLT:
		fn = func(v int64) bool {
			return v < n
		}
	case FilterOperatorLTE:
		fn = func(v int64) bool {
			return v <= n
		}
	case FilterOperatorEQ:
		fn = func(v int64) bool {
			return v == n
		}
	case FilterOperatorNEQ:
		fn = func(v int64) bool {
			return v != n
		}
	default:
		return nil, fmt.Errorf("invalid operator: %q", operator)
	}
	return fn, nil
}
