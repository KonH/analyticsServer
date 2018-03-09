package server

import (
	"reflect"
	"testing"
)

func Test_toFlatMap(t *testing.T) {
	type args struct {
		m map[string][]string
	}

	mapOrd := make(map[string][]string)
	arrOrd := make([]string, 0)
	arrOrd = append(arrOrd, "value")
	mapOrd["key"] = arrOrd
	flatMapOrd := make(map[string]string)
	flatMapOrd["key"] = "value"

	mapEmpty := make(map[string][]string)
	arrEmpty := make([]string, 0)
	mapEmpty["key"] = arrEmpty
	flatMapEmpty := make(map[string]string)

	mapMulti := make(map[string][]string)
	arrMulti := make([]string, 0)
	arrMulti = append(arrMulti, "value0")
	arrMulti = append(arrMulti, "value1")
	mapMulti["key"] = arrMulti
	flatMapMulti := make(map[string]string)
	flatMapMulti["key"] = "value0"

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"Empty", args{make(map[string][]string)}, make(map[string]string)},
		{"Ordinal", args{mapOrd}, flatMapOrd},
		{"EmptyArr", args{mapEmpty}, flatMapEmpty},
		{"MultipleValues", args{mapMulti}, flatMapMulti},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toFlatMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toFlatMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
