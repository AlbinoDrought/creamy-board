package markup

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []Node
	}{
		{
			"empty",
			"",
			[]Node{},
		},

		{
			"default",
			"foo\nbar\n\nbaz",
			[]Node{
				{
					Content: "foo\nbar\n\nbaz",
				},
			},
		},

		{
			"quote",
			">foo\ntrue\n>bar\nfalse",
			[]Node{
				{
					Quoted:  true,
					Content: ">foo",
				},
				{
					Content: "\ntrue\n",
				},
				{
					Quoted:  true,
					Content: ">bar",
				},
				{
					Content: "\nfalse",
				},
			},
		},

		{
			"not quote",
			"foo >bar\nbaz",
			[]Node{
				{
					Quoted:  false,
					Content: "foo >bar\nbaz",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
