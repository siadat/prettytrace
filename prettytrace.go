package prettytrace

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/maruel/panicparse/v2/stack"
)

const packagePrefix = "github.com/siadat/prettytrace."

type CallInfo struct {
	FileName string
	FuncName string
	Line     int
}

func retrieveCallInfo() *CallInfo {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("runtime.Caller failed")
	}

	var fn = runtime.FuncForPC(pc).Name()
	var _, fileName = path.Split(file)
	return &CallInfo{
		FileName: fileName,
		FuncName: fn,
		Line:     line,
	}
}

// Print uses panicparse to print a readable panic stack.
// See https://pkg.go.dev/github.com/maruel/panicparse/v2/stack
func Print() {
	Fprint(os.Stdout)
}

func Fprint(wr io.Writer) {
	var ci = retrieveCallInfo()
	// fmt.Fprintf(wr, "Print called %s:%d inside %s\n", ci.FileName, ci.Line, ci.FuncName)

	// debug.PrintStack()
	var stream = bytes.NewReader(debug.Stack())

	var s, suffix, err = stack.ScanSnapshot(stream, wr, stack.DefaultOpts())
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Find out similar goroutine traces and group them into buckets.
	var buckets = s.Aggregate(stack.AnyValue).Buckets

	// Calculate column length.
	var colLen = 0
	for _, bucket := range buckets {
		for _, line := range filterDebugStack(ci.FuncName, bucket.Signature.Stack.Calls) {
			if l := len(formatFilename(line)); l > colLen {
				colLen = l
			}
		}
	}

	for _, bucket := range buckets {
		// Print the goroutine header.
		var extra = ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}

		if len(bucket.CreatedBy.Calls) != 0 {
			extra += fmt.Sprintf(" [Created by %s.%s @ %s:%d]",
				bucket.CreatedBy.Calls[0].Func.DirName,
				bucket.CreatedBy.Calls[0].Func.Name,
				bucket.CreatedBy.Calls[0].SrcName,
				bucket.CreatedBy.Calls[0].Line,
			)
		}
		fmt.Fprintf(wr, "%d: %s%s\n", len(bucket.IDs), bucket.State, extra)

		// Print the stack lines.
		for _, line := range filterDebugStack(ci.FuncName, bucket.Signature.Stack.Calls) {
			fmt.Fprintln(wr, formatCall(line, colLen, ci.FuncName))
		}
		if bucket.Stack.Elided {
			io.WriteString(wr, "    (...) (elided)\n")
		}
	}

	// If there was any remaining data in the pipe, dump it now.
	if len(suffix) != 0 {
		wr.Write(suffix)
	}
	if err == nil {
		io.Copy(wr, stream)
	}
}

func filterDebugStack(funcName string, lines []stack.Call) []stack.Call {
	var ret []stack.Call
	var sawDebugStack = false
	for _, line := range lines {
		// filter out debug/stack.go
		if !sawDebugStack {
			if line.Func.DirName == "debug" && line.SrcName == "stack.go" {
				sawDebugStack = true
				continue
			} else {
				continue
			}
		}

		if line.Func.DirName == "" {
			continue
		}

		// filter out this package, but not its tests
		if packagePrefix == line.Func.ImportPath+"." && !strings.HasSuffix(line.SrcName, "_test.go") {
			// ret = append(ret, line)
			continue
		}

		ret = append(ret, line)
	}
	return ret
}

func filterCallsByCurrPkg(funcName string, lines []stack.Call) []stack.Call {
	var ret []stack.Call
	for _, line := range lines {
		if strings.HasPrefix(funcName, line.ImportPath) {
			ret = append(ret, line)
			continue
		}
	}
	return ret
}

func formatCall(line stack.Call, colLen int, funcName string) string {
	var prefix = " "
	// don't mark this package, but do mark its tests
	if packagePrefix == line.Func.ImportPath+"." && !strings.HasSuffix(line.SrcName, "_test.go") {
		// noop
	} else if strings.HasPrefix(funcName, line.ImportPath) {
		prefix = ">"
	} else if line.Func.IsPkgMain {
		prefix = ">"
	}
	return fmt.Sprintf(
		"    %s %-*s %s(...)",
		prefix,
		colLen,
		formatFilename(line),
		line.Func.Name,
	)
}

func formatFilename(line stack.Call) string {
	return fmt.Sprintf("%s/%s:%d", line.Func.DirName, line.SrcName, line.Line)
}
