package llcm

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func TestEvaluateFilter(t *testing.T) {
	type args struct {
		expressions []string
	}
	tests := []struct {
		name    string
		args    args
		want    []Filter
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				expressions: nil,
			},
			want:    []Filter{},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				expressions: []string{""},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "space",
			args: args{
				expressions: []string{" "},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "spaces",
			args: args{
				expressions: []string{"    "},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "name == foo",
			args: args{
				expressions: []string{"name == foo"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "foo",
				},
			},
			wantErr: false,
		},
		{
			name: "name ==* Foo",
			args: args{
				expressions: []string{"name ==* Foo"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQI,
					Value:    "Foo",
				},
			},
			wantErr: false,
		},
		{
			name: "name != foo",
			args: args{
				expressions: []string{"name != foo"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorNEQ,
					Value:    "foo",
				},
			},
			wantErr: false,
		},
		{
			name: "name !=* Foo",
			args: args{
				expressions: []string{"name !=* Foo"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorNEQI,
					Value:    "Foo",
				},
			},
			wantErr: false,
		},
		{
			name: "name =~ ^.*foo.*$",
			args: args{
				expressions: []string{"name =~ ^.*foo.*$"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorREQ,
					Value:    "^.*foo.*$",
				},
			},
			wantErr: false,
		},
		{
			name: "name =~* ^.*Foo.*$",
			args: args{
				expressions: []string{"name =~* ^.*Foo.*$"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorREQI,
					Value:    "^.*Foo.*$",
				},
			},
			wantErr: false,
		},
		{
			name: "name !~ ^.*foo.*$",
			args: args{
				expressions: []string{"name !~ ^.*foo.*$"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorNREQ,
					Value:    "^.*foo.*$",
				},
			},
			wantErr: false,
		},
		{
			name: "name !~* ^.*Foo.*$",
			args: args{
				expressions: []string{"name !~* ^.*Foo.*$"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorNREQI,
					Value:    "^.*Foo.*$",
				},
			},
			wantErr: false,
		},
		{
			name: "bytes > 100",
			args: args{
				expressions: []string{"bytes > 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorGT,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "bytes >= 100",
			args: args{
				expressions: []string{"bytes >= 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorGTE,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "bytes < 100",
			args: args{
				expressions: []string{"bytes < 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorLT,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "bytes <= 100",
			args: args{
				expressions: []string{"bytes <= 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorLTE,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "bytes == 100",
			args: args{
				expressions: []string{"bytes == 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorEQ,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "bytes != 100",
			args: args{
				expressions: []string{"bytes != 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorNEQ,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "value with spaces",
			args: args{
				expressions: []string{"name == \"a b c\""},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "\"a b c\"",
				},
			},
			wantErr: false,
		},
		{
			name: "no value",
			args: args{
				expressions: []string{"name =="},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no operator",
			args: args{
				expressions: []string{"name foo"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid key",
			args: args{
				expressions: []string{"invalid == value"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid operator",
			args: args{
				expressions: []string{"name invalid value"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no value with spaces",
			args: args{
				expressions: []string{"name == "},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no value with double quotes",
			args: args{
				expressions: []string{"name == \"\""},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "\"\"",
				},
			},
			wantErr: false,
		},
		{
			name: "no value with single quotes",
			args: args{
				expressions: []string{"name == ''"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "''",
				},
			},
			wantErr: false,
		},
		{
			name: "no value with back quotes",
			args: args{
				expressions: []string{"name == ``"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "``",
				},
			},
			wantErr: false,
		},
		{
			name: "no value with many spaces",
			args: args{
				expressions: []string{"name == '    '"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "' '",
				},
			},
			wantErr: false,
		},
		{
			name: "expr with many spaces",
			args: args{
				expressions: []string{"   name    ==    foo   "},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "foo",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple expressions string+number",
			args: args{
				expressions: []string{"name == foo", "bytes > 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "foo",
				},
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorGT,
					Value:    "100",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple expressions string+string",
			args: args{
				expressions: []string{"name == foo", "class != STANDARD"},
			},
			want: []Filter{
				{
					Key:      FilterKeyName,
					Operator: FilterOperatorEQ,
					Value:    "foo",
				},
				{
					Key:      FilterKeyClass,
					Operator: FilterOperatorNEQ,
					Value:    "STANDARD",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple expressions number+number",
			args: args{
				expressions: []string{"elapsed <= 90", "bytes > 100"},
			},
			want: []Filter{
				{
					Key:      FilterKeyElapsed,
					Operator: FilterOperatorLTE,
					Value:    "90",
				},
				{
					Key:      FilterKeyBytes,
					Operator: FilterOperatorGT,
					Value:    "100",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateFilter(tt.args.expressions)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_setFilter(t *testing.T) {
	type args struct {
		filters []Filter
	}
	type want struct {
		entry *entry
		fns   []bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "name == foo",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "foo",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "foo"},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "class == STANDARD",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyClass,
						Operator: FilterOperatorEQI,
						Value:    "standard",
					},
				},
			},
			want: want{
				entry: &entry{Class: types.LogGroupClassStandard},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "class != STANDARD",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyClass,
						Operator: FilterOperatorNEQI,
						Value:    "standard",
					},
				},
			},
			want: want{
				entry: &entry{Class: types.LogGroupClassStandard},
				fns:   []bool{false},
			},
			wantErr: false,
		},
		{
			name: "name =~ ^test-.*",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorREQ,
						Value:    "^test-.*",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "test-log-group"},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "name =~* ^Test-.*",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorREQI,
						Value:    "^Test-.*",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "test-log-group"},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "name !~ ^test-.*",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorNREQ,
						Value:    "^test-.*",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "test-log-group"},
				fns:   []bool{false},
			},
			wantErr: false,
		},
		{
			name: "name !~* ^Test-.*",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorNREQI,
						Value:    "^Test-.*",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "test-log-group"},
				fns:   []bool{false},
			},
			wantErr: false,
		},
		{
			name: "elapsed == 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyElapsed,
						Operator: FilterOperatorEQ,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{ElapsedDays: 100},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "retention != 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorNEQ,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{RetentionInDays: 100},
				fns:   []bool{false},
			},
			wantErr: false,
		},
		{
			name: "bytes < 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorLT,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{StoredBytes: 50},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "bytes =< 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorLTE,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{StoredBytes: 100},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "bytes > 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorGT,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{StoredBytes: 150},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "bytes >= 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorGTE,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{StoredBytes: 150},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "elapsed >= 1year",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyElapsed,
						Operator: FilterOperatorGTE,
						Value:    "1year",
					},
				},
			},
			want: want{
				entry: &entry{ElapsedDays: 731},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "retention >= 1year",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorGTE,
						Value:    "1year",
					},
				},
			},
			want: want{
				entry: &entry{RetentionInDays: 731},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "elapsed >= delete",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyElapsed,
						Operator: FilterOperatorGTE,
						Value:    "delete",
					},
				},
			},
			want:    want{},
			wantErr: true,
		},
		{
			name: "retention >= delete",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorGTE,
						Value:    "delete",
					},
				},
			},
			want:    want{},
			wantErr: true,
		},
		{
			name: "elapsed >= none",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyElapsed,
						Operator: FilterOperatorGTE,
						Value:    "none",
					},
				},
			},
			want:    want{},
			wantErr: true,
		},
		{
			name: "retention >= none",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyRetention,
						Operator: FilterOperatorGTE,
						Value:    "none",
					},
				},
			},
			want:    want{},
			wantErr: true,
		},
		{
			name: "name == foo && bytes > 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "foo",
					},
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorGT,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{
					LogGroupName: "foo",
					StoredBytes:  150,
				},
				fns: []bool{true, true},
			},
			wantErr: false,
		},
		{
			name: "elapsed <= 90 && bytes > 100",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyElapsed,
						Operator: FilterOperatorLTE,
						Value:    "90",
					},
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorGT,
						Value:    "100",
					},
				},
			},
			want: want{
				entry: &entry{
					ElapsedDays: 90,
					StoredBytes: 150,
				},
				fns: []bool{true, true},
			},
			wantErr: false,
		},
		{
			name: "name == aaa && name != bbb",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "aaa",
					},
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorNEQ,
						Value:    "bbb",
					},
				},
			},
			want: want{
				entry: &entry{
					LogGroupName: "aaa",
				},
				fns: []bool{true, true},
			},
			wantErr: false,
		},
		{
			name: "empty filters",
			args: args{
				filters: []Filter{},
			},
			want: want{
				entry: &entry{
					LogGroupName:    "foo",
					RetentionInDays: 45,
				},
				fns: []bool{},
			},
			wantErr: false,
		},
		{
			name: "empty value",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: ""},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "back quotes",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "`invalid`",
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "`invalid`"},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "long value",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorEQ,
						Value:    "a" + strings.Repeat("b", 1000),
					},
				},
			},
			want: want{
				entry: &entry{LogGroupName: "a" + strings.Repeat("b", 1000)},
				fns:   []bool{true},
			},
			wantErr: false,
		},
		{
			name: "error parsing number",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyBytes,
						Operator: FilterOperatorGT,
						Value:    "invalid_number",
					},
				},
			},
			want:    want{},
			wantErr: true,
		},
		{
			name: "invalid operator",
			args: args{
				filters: []Filter{
					{
						Key:      FilterKeyName,
						Operator: FilterOperatorNone,
						Value:    "foo",
					},
				},
			},
			want:    want{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := &Manager{}
			err := man.setFilter(tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Fatalf("setFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			for i, fn := range man.filterFns {
				if i >= len(tt.want.fns) {
					t.Errorf("Unexpected additional filter function at index %d", i)
					continue
				}
				passed := fn(tt.want.entry)
				if passed != tt.want.fns[i] {
					t.Errorf("Filter function %d: expected %v, got %v", i, tt.want.fns[i], passed)
				}
			}
			if len(man.filterFns) != len(tt.want.fns) {
				t.Errorf("Number of filter functions mismatch: expected %d, got %d", len(tt.want.fns), len(man.filterFns))
			}
		})
	}
}
