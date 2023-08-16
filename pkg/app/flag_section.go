package app

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

type FlagSets struct {
	Order []string
	Sets map[string]*pflag.FlagSet
}

func (fs *FlagSets) FlagSet(name string) *pflag.FlagSet{
	if fs.Sets == nil {
		fs.Sets = map[string]*pflag.FlagSet{}
	}
	if _,ok := fs.Sets[name];!ok {
		fs.Sets[name] = pflag.NewFlagSet(name,pflag.ExitOnError)
		fs.Order = append(fs.Order, name)
	}
	return fs.Sets[name]
}

func PrintSections(w io.Writer, fss FlagSets, cols int) {
	for _, name := range fss.Order {
		fs := fss.Sets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFS.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, buf.String())
		}
	}
}