package qparams

import (
	"net/url"
	"reflect"
	"testing"
)

func TestEncodeQuery(t *testing.T) {
	tests := []struct {
		want string
		test interface{}
	}{
		{"", struct{}{}},
		{"?b=true", struct{ B bool }{true}},
		{"?b=true&b=false&b=true", struct{ B []bool }{[]bool{true, false, true}}},
		{"?i=1&i=3", struct{ I []int }{[]int{1, 3}}},
		{"?f=1.000000&f=3.000000", struct{ F []float32 }{[]float32{1, 3}}},
		{"?s=a+string&s=str", struct{ S []string }{[]string{"a string", "str"}}},
		{"?i=27&s=str&b=true", struct {
			I int
			S string
			B bool
		}{27, "str", true}},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			got, err := Encode(tc.test)
			if err != nil {
				t.Fatalf("got error: %s", err)
			}
			if got != tc.want {
				t.Fatalf("expected %s; got %s", tc.want, got)
			}
		})
	}
}

func TestDecodeInvalidQuery(t *testing.T) {
	url, _ := url.Parse("http://example.org?a=b&c=d")
	var i int
	if err := Decode(url.Query(), &i); err == nil {
		t.Fatalf("expected an error")
	}
	var x struct{ A, C string }
	if err := Decode(url.Query(), x); err == nil {
		t.Fatalf("expected an error")
	}
}

func TestDecodeQuery(t *testing.T) {
	type data struct {
		B   bool
		I   int
		F32 float32
		F64 float64
		S   string
	}
	tests := []struct {
		url   string
		want  data
		iserr bool
	}{
		{"http://example.org", data{}, false},
		{"http://example.org?b=true", data{B: true}, false},
		{"http://example.org?b=not-a-bool", data{}, true},
		{"http://example.org?i=true", data{}, true},
		{"http://example.org?i=27", data{I: 27}, false},
		{"http://example.org?f32=2.7", data{F32: 2.7}, false},
		{"http://example.org?f32=not-a-float", data{}, true},
		{"http://example.org?f64=2.7", data{F64: 2.7}, false},
		{"http://example.org?f64=not-a-float", data{}, true},
		{"http://example.org?s=some%20string", data{S: "some string"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			url, err := url.Parse(tc.url)
			if err != nil {
				t.Fatalf("got error: %s", err)
			}
			var got data
			err = Decode(url.Query(), &got)
			if tc.iserr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("got error: %s", err)
			}
			if got != tc.want {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

func TestDecodeQueryArray(t *testing.T) {
	type data struct {
		B   []bool
		I   []int
		F32 []float32
		F64 []float64
		S   []string
	}
	// eq := func(a, b data) bool {
	// 	return false
	// }
	tests := []struct {
		url   string
		want  data
		iserr bool
	}{
		{"http://example.org", data{}, false},
		{"http://example.org?b=true&b=false", data{B: []bool{true, false}}, false},
		{"http://example.org?b=true&b=not-a-bool", data{}, true},
		{"http://example.org?i=1&i=2", data{I: []int{1, 2}}, false},
		{"http://example.org?i=1&i=not-an-int", data{}, true},
		{"http://example.org?f32=1.1&f32=2.1", data{F32: []float32{1.1, 2.1}}, false},
		{"http://example.org?f32=1.1&f32=not-a-float", data{}, true},
		{"http://example.org?f64=1.1&f64=2.1", data{F64: []float64{1.1, 2.1}}, false},
		{"http://example.org?f64=1.1&f64=not-a-float", data{}, true},
		{"http://example.org?s=first&s=second", data{S: []string{"first", "second"}}, false},
		{"http://example.org?s=first&i=1&s=second&i=2", data{S: []string{"first", "second"}, I: []int{1, 2}}, false},
	}
	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			url, err := url.Parse(tc.url)
			if err != nil {
				t.Fatalf("got error: %s", err)
			}
			var got data
			err = Decode(url.Query(), &got)
			if tc.iserr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("got error: %s", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}

		})
	}
}
