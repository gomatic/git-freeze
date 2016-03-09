package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	transitive := func() bool {
		transitive := flag.Bool("transitive", false, "Traverse transitive imports, i.e. vendor/")
		flag.Parse()
		return *transitive
	}()

	self := "."
	{
		if path, err := filepath.Abs("."); err != nil {
			fmt.Fprintf(os.Stderr, "abs . failed: %s\n", err)
			os.Exit(1)
		} else {
			self = filepath.Base(path)
		}
	}

	if _, err := os.Stat(".git"); err != nil {
		fmt.Fprintf(os.Stderr, ".git does not exist\n")
		os.Exit(1)
	}

	_git, _go := "", ""
	if git_, err := exec.LookPath("git"); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		_git = git_
	}

	if go_, err := exec.LookPath("go"); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		_go = go_
	}

	imports := exec.Command(_go, `list`, `-f`, `{{$p := .ImportPath}}{{range $imp := .Imports}}{{printf "%s\t%s\n" $p $imp}}{{end}}`, `./...`)

	r, w := io.Pipe()
	imports.Stdout = w

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(r io.Reader) {
		var (
			err  error = nil
			line string
		)
		reader := bufio.NewReader(r)
		repos := map[string]int{}
		for err == nil {
			if line, err = reader.ReadString('\n'); err == nil {
				pi := strings.Split(strings.TrimSpace(string(line)), "\t")
				p, i := pi[0], pi[1]
				parts := strings.Split(i, "/")
				if len(parts) >= 3 && parts[2] != self {
					if !transitive && strings.Contains(p, "/vendor") {
						continue
					}
					repo := strings.Join(parts[:3], "/")
					repos[repo] += 1
				}
			}
		}
		for repo, _ := range repos {
			b := path.Join("vendor", repo)
			if _, err := os.Stat(b); err == nil {
				fmt.Printf("%s exists\n", repo)
			} else {
				submodule := exec.Command(_git, `submodule`, `add`, `-f`, fmt.Sprintf(`https://%s`, repo), fmt.Sprintf(`vendor/%s`, repo))
				if output, err := submodule.CombinedOutput(); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\noutput:\n%s", submodule.Args, output)
				} else {
					fmt.Printf("%s", output)
				}
			}
		}
		wg.Done()
	}(r)

	imports.Start()
	imports.Wait()
	w.Close()
	wg.Wait()

}
