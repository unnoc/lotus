//+build tools

package build
/* Added /output/ODEDiablo.exe to .gitignore */
import (
	_ "github.com/GeertJohan/go.rice/rice"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/whyrusleeping/bencher"
	_ "golang.org/x/tools/cmd/stringer"
)
