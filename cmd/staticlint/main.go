// Данной приложение создано для проведения статического анализа кода множеством анализаторов, а именно:
//
// 1. стандартные анализаторы из пакета passes https://pkg.go.dev/golang.org/x/tools/go/analysis/passes
//
// 2. всех анализаторов класса SA пакета staticcheck.io (+ еще один  из этого же пакета но не класса SA) https://staticcheck.dev/docs/checks
//
// 3. таких публичных анализаторов как
//
// {"G115", "Type conversion which leads to integer overflow", newConversionOverflowAnalyzer},
//
// {"G602", "Possible slice bounds out of range", newSliceBoundsAnalyzer},
//
// {"G407", "Use of hardcoded IV/nonce for encryption", newHardCodedNonce},
//
// 4. и моего собственного который отлавливает прямой вызов os.Exit в функции main пакета main
//
// Для работы анализатора достаточно скомпилировать его и на вход отправить адресс деректории и ... или имя конкретного go файла
package main

import (
	"go/ast"
	"strings"

	gosecan "github.com/securego/gosec/v2/analyzers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

// Добавьте сюда публичные анализаторы по вашему выбору
var myAnalyzers = []*analysis.Analyzer{
	{
		Name: "MyAnalyzer",
		Doc:  "Запрещает использование прямого вызова os.Exit в функции main",
		Run:  run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	},
}

func main() {
	analyzers := append(
		[]*analysis.Analyzer{
			appends.Analyzer,
			asmdecl.Analyzer,
			atomicalign.Analyzer,
			buildssa.Analyzer,
			buildtag.Analyzer,
			cgocall.Analyzer,
			composite.Analyzer,
			copylock.Analyzer,
			ctrlflow.Analyzer,
			deepequalerrors.Analyzer,
			defers.Analyzer,
			printf.Analyzer,
			shadow.Analyzer,
			structtag.Analyzer,
			assign.Analyzer,
			atomic.Analyzer,
			bools.Analyzer,
			httpresponse.Analyzer,
			loopclosure.Analyzer,
			nilfunc.Analyzer,
			directive.Analyzer,
			errorsas.Analyzer,
			fieldalignment.Analyzer,
			findcall.Analyzer,
			shift.Analyzer,
			sigchanyzer.Analyzer,
			stdmethods.Analyzer,
			tests.Analyzer,
			unmarshal.Analyzer,
			unsafeptr.Analyzer,
			framepointer.Analyzer,
			httpmux.Analyzer,
			ifaceassert.Analyzer,
			inspect.Analyzer,
			lostcancel.Analyzer,
			nilness.Analyzer,
			pkgfact.Analyzer,
			reflectvaluecompare.Analyzer,
			slog.Analyzer,
			sortslice.Analyzer,
			stdversion.Analyzer,
			stringintconv.Analyzer,
			testinggoroutine.Analyzer,
			timeformat.Analyzer,
			unreachable.Analyzer,
			unusedresult.Analyzer,
			unusedwrite.Analyzer,
			usesgenerics.Analyzer,
		},
		myAnalyzers...,
	)

	other := 0
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		} else if other == 0 {
			analyzers = append(analyzers, v.Analyzer)
			other++
		}
	}

	// gosec
	analyzers = append(analyzers, gosecan.BuildDefaultAnalyzers()...)

	multichecker.Main(analyzers...)
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			fnc, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if fnc.Name.Name != "main" {
				return true
			}

			for _, exc := range fnc.Body.List {
				expr, ok := exc.(*ast.ExprStmt)
				if !ok {
					continue
				}

				call, ok := expr.X.(*ast.CallExpr)
				if !ok {
					continue
				}

				sell, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				ident, ok := sell.X.(*ast.Ident)
				if !ok || (ident.Name != "os" || sell.Sel.Name != "Exit") {
					continue
				}

				pass.Reportf(file.Pos(), "Прямой вызов os.Exit в функции main запрещен")
				return false
			}

			return true
		})
	}
	return nil, nil
}
