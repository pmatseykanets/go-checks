package commentspace

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const Doc = `checks if a body of a comment is separated by a space`

var Analyzer = &analysis.Analyzer{
	Name: "commentspace",
	Doc:  Doc,
	Run:  run,
}

var directives = []string{
	// https://golang.org/cmd/compile/#hdr-Compiler_Directives
	"//line ",
	"/*line ",
	"//go:",
	// https://staticcheck.io/docs/#line-based-linter-directives
	"//lint:",
}

func run(pass *analysis.Pass) (interface{}, error) {
	var nlen int
	for _, file := range pass.Files {
		for _, group := range file.Comments {
			// Skip the package comments.
			if group == file.Doc {
				continue
			}
			for _, comment := range group.List {
				nlen = len(comment.Text)
				if nlen < 3 {
					continue
				}
				// Skip known directives.
				for _, directive := range directives {
					//lint:ignore SA4017 false-positive.
					if strings.HasPrefix(comment.Text, directive) {
						continue
					}
				}

				if strings.HasPrefix(comment.Text, "//") && nlen > 2 && comment.Text[2:3] != " " {
					fmt.Println(comment.Pos(), comment.End())
					pass.Report(analysis.Diagnostic{
						Pos: comment.Pos(), Message: "A comment should have a leading space",
						SuggestedFixes: []analysis.SuggestedFix{
							{Message: "Add space", TextEdits: []analysis.TextEdit{
								{Pos: comment.Pos(), End: comment.End(), NewText: []byte("// " + comment.Text[2:])},
							}},
						},
					})

					continue
				}

				if strings.HasPrefix(comment.Text, "/*") && nlen > 4 && !strings.Contains(comment.Text, "\n") {
					// An inline comment /* foo */.
					var (
						failed  bool
						newText = comment.Text
					)
					if comment.Text[2:3] != " " {
						failed = true
						newText = "/* " + newText[2:]
					}
					if comment.Text[nlen-3:nlen-2] != " " {
						failed = true
						newText = newText[0:len(newText)-2] + " */"
					}
					if failed {
						pass.Report(analysis.Diagnostic{
							Pos: comment.Pos(), Message: "A an in-line comment should have a leading and trailing spaces",
							SuggestedFixes: []analysis.SuggestedFix{
								{Message: "Add spaces", TextEdits: []analysis.TextEdit{
									{Pos: comment.Pos(), End: comment.End(), NewText: []byte(newText)},
								}},
							},
						})
					}
					continue
				}
			}
		}
	}

	return nil, nil
}
