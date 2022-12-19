package prettytrace_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/siadat/prettytrace"
	"github.com/siadat/prettytrace/testpkg"
)

func TestOutput(t *testing.T) {
	var testCases = []struct {
		f    func()
		want []lineSpec
	}{
		{
			testpkg.DivideByZero,
			[]lineSpec{
				{"1: running [Created by testing.(*T).Run @ testing.go:", ""},
				{"    >", "TestOutput.func1.1(...)"},
				{"     ", "divideByZero(...)"},
				{"     ", "DivideByZero(...)"},
				{"    >", "TestOutput.func1(...)"},
				{"    >", "TestOutput(...)"},
				{"     ", "tRunner(...)"},
			},
		},
		{
			testpkg.F1,
			[]lineSpec{
				{"1: running [Created by testing.(*T).Run @ testing.go:", ""},
				{"    >", "TestOutput.func1.1(...)"},
				{"     ", "S.f1(...)"},
				{"     ", "F1(...)"},
				{"    >", "TestOutput.func1(...)"},
				{"    >", "TestOutput(...)"},
				{"     ", "tRunner(...)"},
			},
		},
		{
			testpkg.F2,
			[]lineSpec{
				{"1: running [Created by testing.(*T).Run @ testing.go:", ""},
				{"    >", "TestOutput.func1.1(...)"},
				{"     ", "(*S).f2(...)"},
				{"     ", "F2(...)"},
				{"    >", "TestOutput.func1(...)"},
				{"    >", "TestOutput(...)"},
				{"     ", "tRunner(...)"},
			},
		},
	}

	var successfulCases = 0
	for _, tc := range testCases {
		func() {
			defer func() {
				var _ = recover()
				var buf bytes.Buffer
				prettytrace.Fprint(&buf)
				checkTrace(t, buf.String(), tc.want)
				successfulCases++
			}()
			tc.f()

			t.Fatalf("should not be here")
		}()
	}
	if successfulCases != len(testCases) {
		t.Fatalf("want %d cases to succeed, got %d", len(testCases), successfulCases)
	}
}

type lineSpec struct {
	prefix, suffix string
}

func checkTrace(t *testing.T, got string, want []lineSpec) {
	var lines = strings.Split(strings.TrimSpace(got), "\n")
	var printWantAndGot = func() {
		fmt.Println("Want:")
		for i, line := range want {
			fmt.Printf("  line[%d]: %q\n", i, line)
		}
		fmt.Println("Got:")
		for i, line := range lines {
			fmt.Printf("  line[%d]: %s\n", i, line)
		}
	}
	if len(lines) != len(want) {
		printWantAndGot()
		t.Fatalf("want %d lines, got %d lines", len(want), len(lines))
	}
	for i := range lines {
		if !strings.HasPrefix(lines[i], want[i].prefix) {
			printWantAndGot()
			t.Fatalf("in line %d\nwant prefix:\n    %q\ngot:\n    %q", i, want[i].prefix, lines[i])
		}
		if !strings.HasSuffix(lines[i], want[i].suffix) {
			printWantAndGot()
			t.Fatalf("in line %d\nwant suffix:\n    %s%q\ngot:\n    %q", i, strings.Repeat(".", len(lines[i])-len(want[i].suffix)), want[i].suffix, lines[i])
		}
	}
}
