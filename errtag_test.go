package errtag_test

import (
	"github.com/pi9min/errtag"
	"testing"
)

type T1 int

func (t T1) CallErrorTag() string {
	return errtag.ErrorTag()
}

type T2 struct{}

func (t *T2) CallErrorTag() string {
	return errtag.ErrorTag()
}

func TestErrorTag(t *testing.T) {
	t1 := T1(1)
	t2 := &T2{}
	tests := []struct {
		name   string
		actual string
		expect string
	}{
		{
			name:   "Normal Functionからの呼び出し",
			actual: errtag.ErrorTag(),
			expect: "errtag_test.TestErrorTag",
		},
		{
			name:   "Primitive型のReceiver Function内からの呼び出し",
			actual: t1.CallErrorTag(),
			expect: "errtag_test/T1.CallErrorTag",
		},
		{
			name:   "構造体のReceiver Function内からの呼び出し",
			actual: t2.CallErrorTag(),
			expect: "errtag_test/T2.CallErrorTag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expect {
				t.Errorf("Failed:\nactual: %s\nexpect: %s", tt.actual, tt.expect)
			}
		})
	}
}

func BenchmarkErrorTag(b *testing.B) {
	t1 := T1(1)
	t2 := &T2{}
	tests := []struct {
		name string
		call func() string
	}{
		{
			name: "normal",
			call: func() string { return errtag.ErrorTag() },
		},
		{
			name: "primitive",
			call: func() string { return t1.CallErrorTag() },
		},
		{
			name: "struct",
			call: func() string { return t2.CallErrorTag() },
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tt.call()
			}
		})
	}
}
