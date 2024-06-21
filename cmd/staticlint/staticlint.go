// Package main contains project multicheck linter:
// - golang.org/x/tools/go/analysis
// - honnef.co/go/tools/staticcheck : SA*,S1006,ST1000,QF1009
// - github.com/alexkohler/nakedret
// - github.com/orijtech/structslop
// - github.com/YaNeAndrey/ya-metrics/internal/osexit
//

package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/osexit"
	"github.com/alexkohler/nakedret"
	"github.com/orijtech/structslop"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	analyzers := []*analysis.Analyzer{
		appends.Analyzer,          // detects if there is only one variable in append
		asmdecl.Analyzer,          // reports mismatches between assembly files and Go declarations
		assign.Analyzer,           // detects useless assignments
		atomic.Analyzer,           // checks for common mistakes using the sync/atomic package
		atomicalign.Analyzer,      // checks for non-64-bit-aligned arguments to sync/atomic functions
		bools.Analyzer,            // detects common mistakes involving boolean operators
		buildtag.Analyzer,         // constructs the SSA representation of an error-free package and returns the set of all functions within it
		cgocall.Analyzer,          // detects some violations of the cgo pointer passing rules
		composite.Analyzer,        // checks for unkeyed composite literals
		copylock.Analyzer,         // checks for locks erroneously passed by value
		deepequalerrors.Analyzer,  // checks for the use of reflect.DeepEqual with error values
		defers.Analyzer,           // checks for common mistakes in defer statements
		directive.Analyzer,        // checks known Go toolchain directives
		errorsas.Analyzer,         // checks that the second argument to errors.As is a pointer to a type implementing error
		framepointer.Analyzer,     // reports assembly code that clobbers the frame pointer before saving it
		httpresponse.Analyzer,     // checks for mistakes using HTTP responses
		ifaceassert.Analyzer,      // flags impossible interface-interface type assertions
		loopclosure.Analyzer,      // checks for references to enclosing loop variables from within nested functions
		lostcancel.Analyzer,       // checks for failure to call a context cancellation function
		nilfunc.Analyzer,          // checks for useless comparisons against nil
		printf.Analyzer,           // checks consistency of Printf format strings and arguments
		shadow.Analyzer,           // checks for shadowed variables
		shift.Analyzer,            // checks for shifts that exceed the width of an integer
		sigchanyzer.Analyzer,      // detects misuse of unbuffered signal as argument to signal.Notify
		slog.Analyzer,             // checks for mismatched key-value pairs in log/slog calls
		stdmethods.Analyzer,       // checks for misspellings in the signatures of methods similar to well-known interfaces
		stringintconv.Analyzer,    // flags type conversions from integers to strings
		structtag.Analyzer,        // checks struct field tags are well formed
		testinggoroutine.Analyzer, // detect calls to Fatal from a test goroutine
		tests.Analyzer,            // checks for common mistaken usages of tests and examples
		timeformat.Analyzer,       // checks for the use of time.Format or time.Parse calls with a bad format
		unmarshal.Analyzer,        // checks for passing non-pointer or non-interface types to unmarshal and decode functions
		unreachable.Analyzer,      // checks for unreachable code
		unsafeptr.Analyzer,        // checks for invalid conversions of uintptr to unsafe.Pointer.
		unusedresult.Analyzer,     // checks for unused results of calls to certain pure functions
		unusedwrite.Analyzer,      // checks for unused writes to the elements of a struct or array object

		structslop.Analyzer,              // checks if struct fields can be re-arranged to optimize size
		nakedret.NakedReturnAnalyzer(50), //  find naked returns in functions greater than 50 lines.
		osexit.Analyzer,
	}

	/*	for _, v := range defaultAnalyzers {
			analyzers = append(analyzers, v)
		}
	*/
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") ||
			v.Analyzer.Name == "S1006" || // Use for { ... } for infinite loops
			v.Analyzer.Name == "ST1000" || // Incorrect or missing package comment
			v.Analyzer.Name == "QF1009" { // Use time.Time.Equal instead of == operator
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	multichecker.Main(
		analyzers...,
	)
}
