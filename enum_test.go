package llcm

import (
	"reflect"
	"testing"
)

func TestOutputType_String(t *testing.T) {
	tests := []struct {
		name string
		tr   OutputType
		want string
	}{
		{
			name: "none",
			tr:   OutputTypeNone,
			want: "none",
		},
		{
			name: "json",
			tr:   OutputTypeJSON,
			want: "json",
		},
		{
			name: "prettyjson",
			tr:   OutputTypePrettyJSON,
			want: "prettyjson",
		},
		{
			name: "text",
			tr:   OutputTypeText,
			want: "text",
		},
		{
			name: "compressedtext",
			tr:   OutputTypeCompressedText,
			want: "compressedtext",
		},
		{
			name: "markdown",
			tr:   OutputTypeMarkdown,
			want: "markdown",
		},
		{
			name: "backlog",
			tr:   OutputTypeBacklog,
			want: "backlog",
		},
		{
			name: "tsv",
			tr:   OutputTypeTSV,
			want: "tsv",
		},
		{
			name: "chart",
			tr:   OutputTypeChart,
			want: "chart",
		},
		{
			name: "unknown",
			tr:   OutputType(12345),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("OutputType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      OutputType
		want    []byte
		wantErr bool
	}{
		{
			name: "none",
			tr:   OutputTypeNone,
			want: []byte(`"none"`),
		},
		{
			name: "json",
			tr:   OutputTypeJSON,
			want: []byte(`"json"`),
		},
		{
			name: "prettyjson",
			tr:   OutputTypePrettyJSON,
			want: []byte(`"prettyjson"`),
		},
		{
			name: "text",
			tr:   OutputTypeText,
			want: []byte(`"text"`),
		},
		{
			name: "compressedtext",
			tr:   OutputTypeCompressedText,
			want: []byte(`"compressedtext"`),
		},
		{
			name: "markdown",
			tr:   OutputTypeMarkdown,
			want: []byte(`"markdown"`),
		},
		{
			name: "backlog",
			tr:   OutputTypeBacklog,
			want: []byte(`"backlog"`),
		},
		{
			name: "tsv",
			tr:   OutputTypeTSV,
			want: []byte(`"tsv"`),
		},
		{
			name: "chart",
			tr:   OutputTypeChart,
			want: []byte(`"chart"`),
		},
		{
			name: "unknown",
			tr:   OutputType(12345),
			want: []byte(`""`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutputType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OutputType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOutputType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    OutputType
		wantErr bool
	}{
		{
			name: "json",
			args: args{
				s: "json",
			},
			want:    OutputTypeJSON,
			wantErr: false,
		},
		{
			name: "prettyjson",
			args: args{
				s: "prettyjson",
			},
			want:    OutputTypePrettyJSON,
			wantErr: false,
		},
		{
			name: "text",
			args: args{
				s: "text",
			},
			want:    OutputTypeText,
			wantErr: false,
		},
		{
			name: "compressedtext",
			args: args{
				s: "compressedtext",
			},
			want:    OutputTypeCompressedText,
			wantErr: false,
		},
		{
			name: "markdown",
			args: args{
				s: "markdown",
			},
			want:    OutputTypeMarkdown,
			wantErr: false,
		},
		{
			name: "backlog",
			args: args{
				s: "backlog",
			},
			want:    OutputTypeBacklog,
			wantErr: false,
		},
		{
			name: "tsv",
			args: args{
				s: "tsv",
			},
			want:    OutputTypeTSV,
			wantErr: false,
		},
		{
			name: "chart",
			args: args{
				s: "chart",
			},
			want:    OutputTypeChart,
			wantErr: false,
		},
		{
			name: "unknown",
			args: args{
				s: "unknown",
			},
			want:    OutputTypeNone,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOutputType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDesiredState_String(t *testing.T) {
	tests := []struct {
		name string
		tr   DesiredState
		want string
	}{
		{
			name: "none",
			tr:   DesiredStateNone,
			want: "none",
		},
		{
			name: "delete",
			tr:   DesiredStateZero,
			want: "delete",
		},
		{
			name: "1day",
			tr:   DesiredStateOneDay,
			want: "1day",
		},
		{
			name: "3days",
			tr:   DesiredStateThreeDays,
			want: "3days",
		},
		{
			name: "5days",
			tr:   DesiredStateFiveDays,
			want: "5days",
		},
		{
			name: "1week",
			tr:   DesiredStateOneWeek,
			want: "1week",
		},
		{
			name: "2weeks",
			tr:   DesiredStateTwoWeeks,
			want: "2weeks",
		},
		{
			name: "1month",
			tr:   DesiredStateOneMonth,
			want: "1month",
		},
		{
			name: "3months",
			tr:   DesiredStateThreeMonths,
			want: "3months",
		},
		{
			name: "4months",
			tr:   DesiredStateFourMonths,
			want: "4months",
		},
		{
			name: "5months",
			tr:   DesiredStateFiveMonths,
			want: "5months",
		},
		{
			name: "6months",
			tr:   DesiredStateSixMonths,
			want: "6months",
		},
		{
			name: "1year",
			tr:   DesiredStateOneYear,
			want: "1year",
		},
		{
			name: "13months",
			tr:   DesiredStateThirteenMonths,
			want: "13months",
		},
		{
			name: "18months",
			tr:   DesiredStateEighteenMonths,
			want: "18months",
		},
		{
			name: "2years",
			tr:   DesiredStateTwoYears,
			want: "2years",
		},
		{
			name: "3years",
			tr:   DesiredStateThreeYears,
			want: "3years",
		},
		{
			name: "5years",
			tr:   DesiredStateFiveYears,
			want: "5years",
		},
		{
			name: "6years",
			tr:   DesiredStateSixYears,
			want: "6years",
		},
		{
			name: "7years",
			tr:   DesiredStateSevenYears,
			want: "7years",
		},
		{
			name: "8years",
			tr:   DesiredStateEightYears,
			want: "8years",
		},
		{
			name: "9years",
			tr:   DesiredStateNineYears,
			want: "9years",
		},
		{
			name: "10years",
			tr:   DesiredStateTenYears,
			want: "10years",
		},
		{
			name: "infinite",
			tr:   DesiredStateInfinite,
			want: "infinite",
		},
		{
			name: "protect",
			tr:   DesiredStateProtected,
			want: "protect",
		},
		{
			name: "unprotect",
			tr:   DesiredStateUnprotected,
			want: "unprotect",
		},
		{
			name: "unknown",
			tr:   DesiredState(12345),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("DesiredState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDesiredState_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      DesiredState
		want    []byte
		wantErr bool
	}{
		{
			name:    "none",
			tr:      DesiredStateNone,
			want:    []byte(`"none"`),
			wantErr: false,
		},
		{
			name:    "delete",
			tr:      DesiredStateZero,
			want:    []byte(`"delete"`),
			wantErr: false,
		},
		{
			name:    "1day",
			tr:      DesiredStateOneDay,
			want:    []byte(`"1day"`),
			wantErr: false,
		},
		{
			name:    "3days",
			tr:      DesiredStateThreeDays,
			want:    []byte(`"3days"`),
			wantErr: false,
		},
		{
			name:    "5days",
			tr:      DesiredStateFiveDays,
			want:    []byte(`"5days"`),
			wantErr: false,
		},
		{
			name:    "1week",
			tr:      DesiredStateOneWeek,
			want:    []byte(`"1week"`),
			wantErr: false,
		},
		{
			name:    "2weeks",
			tr:      DesiredStateTwoWeeks,
			want:    []byte(`"2weeks"`),
			wantErr: false,
		},
		{
			name:    "1month",
			tr:      DesiredStateOneMonth,
			want:    []byte(`"1month"`),
			wantErr: false,
		},
		{
			name:    "2months",
			tr:      DesiredStateTwoMonths,
			want:    []byte(`"2months"`),
			wantErr: false,
		},
		{
			name:    "3months",
			tr:      DesiredStateThreeMonths,
			want:    []byte(`"3months"`),
			wantErr: false,
		},
		{
			name:    "4months",
			tr:      DesiredStateFourMonths,
			want:    []byte(`"4months"`),
			wantErr: false,
		},
		{
			name:    "5months",
			tr:      DesiredStateFiveMonths,
			want:    []byte(`"5months"`),
			wantErr: false,
		},
		{
			name:    "6months",
			tr:      DesiredStateSixMonths,
			want:    []byte(`"6months"`),
			wantErr: false,
		},
		{
			name:    "1year",
			tr:      DesiredStateOneYear,
			want:    []byte(`"1year"`),
			wantErr: false,
		},
		{
			name:    "13months",
			tr:      DesiredStateThirteenMonths,
			want:    []byte(`"13months"`),
			wantErr: false,
		},
		{
			name:    "18months",
			tr:      DesiredStateEighteenMonths,
			want:    []byte(`"18months"`),
			wantErr: false,
		},
		{
			name:    "2years",
			tr:      DesiredStateTwoYears,
			want:    []byte(`"2years"`),
			wantErr: false,
		},
		{
			name:    "3years",
			tr:      DesiredStateThreeYears,
			want:    []byte(`"3years"`),
			wantErr: false,
		},
		{
			name:    "5years",
			tr:      DesiredStateFiveYears,
			want:    []byte(`"5years"`),
			wantErr: false,
		},
		{
			name:    "6years",
			tr:      DesiredStateSixYears,
			want:    []byte(`"6years"`),
			wantErr: false,
		},
		{
			name:    "7years",
			tr:      DesiredStateSevenYears,
			want:    []byte(`"7years"`),
			wantErr: false,
		},
		{
			name:    "8years",
			tr:      DesiredStateEightYears,
			want:    []byte(`"8years"`),
			wantErr: false,
		},
		{
			name:    "9years",
			tr:      DesiredStateNineYears,
			want:    []byte(`"9years"`),
			wantErr: false,
		},
		{
			name:    "10years",
			tr:      DesiredStateTenYears,
			want:    []byte(`"10years"`),
			wantErr: false,
		},
		{
			name:    "infinite",
			tr:      DesiredStateInfinite,
			want:    []byte(`"infinite"`),
			wantErr: false,
		},
		{
			name:    "protect",
			tr:      DesiredStateProtected,
			want:    []byte(`"protect"`),
			wantErr: false,
		},
		{
			name:    "unprotect",
			tr:      DesiredStateUnprotected,
			want:    []byte(`"unprotect"`),
			wantErr: false,
		},
		{
			name:    "unknown",
			tr:      DesiredState(12345),
			want:    []byte(`""`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("DesiredState.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DesiredState.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDesiredState(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    DesiredState
		wantErr bool
	}{
		{
			name: "delete",
			args: args{
				s: "delete",
			},
			want:    DesiredStateZero,
			wantErr: false,
		},
		{
			name: "1day",
			args: args{
				s: "1day",
			},
			want:    DesiredStateOneDay,
			wantErr: false,
		},
		{
			name: "3days",
			args: args{
				s: "3days",
			},
			want:    DesiredStateThreeDays,
			wantErr: false,
		},
		{
			name: "5days",
			args: args{
				s: "5days",
			},
			want:    DesiredStateFiveDays,
			wantErr: false,
		},
		{
			name: "1week",
			args: args{
				s: "1week",
			},
			want:    DesiredStateOneWeek,
			wantErr: false,
		},
		{
			name: "2weeks",
			args: args{
				s: "2weeks",
			},
			want:    DesiredStateTwoWeeks,
			wantErr: false,
		},
		{
			name: "1month",
			args: args{
				s: "1month",
			},
			want:    DesiredStateOneMonth,
			wantErr: false,
		},
		{
			name: "2months",
			args: args{
				s: "2months",
			},
			want:    DesiredStateTwoMonths,
			wantErr: false,
		},
		{
			name: "3months",
			args: args{
				s: "3months",
			},
			want:    DesiredStateThreeMonths,
			wantErr: false,
		},
		{
			name: "4months",
			args: args{
				s: "4months",
			},
			want:    DesiredStateFourMonths,
			wantErr: false,
		},
		{
			name: "5months",
			args: args{
				s: "5months",
			},
			want:    DesiredStateFiveMonths,
			wantErr: false,
		},
		{
			name: "6months",
			args: args{
				s: "6months",
			},
			want:    DesiredStateSixMonths,
			wantErr: false,
		},
		{
			name: "1year",
			args: args{
				s: "1year",
			},
			want:    DesiredStateOneYear,
			wantErr: false,
		},
		{
			name: "13months",
			args: args{
				s: "13months",
			},
			want:    DesiredStateThirteenMonths,
			wantErr: false,
		},
		{
			name: "18months",
			args: args{
				s: "18months",
			},
			want:    DesiredStateEighteenMonths,
			wantErr: false,
		},
		{
			name: "2years",
			args: args{
				s: "2years",
			},
			want:    DesiredStateTwoYears,
			wantErr: false,
		},
		{
			name: "3years",
			args: args{
				s: "3years",
			},
			want:    DesiredStateThreeYears,
			wantErr: false,
		},
		{
			name: "5years",
			args: args{
				s: "5years",
			},
			want:    DesiredStateFiveYears,
			wantErr: false,
		},
		{
			name: "6years",
			args: args{
				s: "6years",
			},
			want:    DesiredStateSixYears,
			wantErr: false,
		},
		{
			name: "7years",
			args: args{
				s: "7years",
			},
			want:    DesiredStateSevenYears,
			wantErr: false,
		},
		{
			name: "8years",
			args: args{
				s: "8years",
			},
			want:    DesiredStateEightYears,
			wantErr: false,
		},
		{
			name: "9years",
			args: args{
				s: "9years",
			},
			want:    DesiredStateNineYears,
			wantErr: false,
		},
		{
			name: "10years",
			args: args{
				s: "10years",
			},
			want:    DesiredStateTenYears,
			wantErr: false,
		},
		{
			name: "infinite",
			args: args{
				s: "infinite",
			},
			want:    DesiredStateInfinite,
			wantErr: false,
		},
		{
			name: "protect",
			args: args{
				s: "protect",
			},
			want:    DesiredStateProtected,
			wantErr: false,
		},
		{
			name: "unprotect",
			args: args{
				s: "unprotect",
			},
			want:    DesiredStateUnprotected,
			wantErr: false,
		},
		{
			name: "none",
			args: args{
				s: "none",
			},
			want:    DesiredStateNone,
			wantErr: true,
		},
		{
			name: "unknown",
			args: args{
				s: "unknown",
			},
			want:    DesiredStateNone,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDesiredState(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDesiredState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDesiredState() = %v, want %v", got, tt.want)
			}
		})
	}
}
