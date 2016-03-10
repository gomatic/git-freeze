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
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

//
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	branch := ""
	{
		cmd := exec.Command(_git, "rev-parse", "--abbrev-ref", "HEAD")
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\noutput:\n%s", cmd.Args, output)
		} else {
			branch = strings.TrimSpace(string(output))
		}
	}

	force, verbose, transitive, subtree, list, dryrun := false, false, false, false, false, false
	flag.BoolVar(&transitive, "transitive", transitive, "Traverse transitive imports, i.e. vendor/")
	flag.BoolVar(&subtree, "subtree", subtree, "Use a subtree instead of a submodule.")
	flag.BoolVar(&dryrun, "dry-run", dryrun, "Just print the command but do not run it.")
	flag.BoolVar(&list, "list", list, "Only list the imports that can be frozen.")
	flag.BoolVar(&verbose, "verbose", verbose, "More output.")
	flag.BoolVar(&force, "force", force, "Force.")
	flag.StringVar(&branch, "branch", branch, "Git branch/commit to submodule/subtree.")
	flag.Usage = func() {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	flag.Parse()
	patterns := make([]*regexp.Regexp, len(flag.Args()))
	for i, a := range flag.Args() {
		patterns[i] = regexp.MustCompile(a)
	}

	imports := exec.Command(_go, "list", "-f", `{{$p := .ImportPath}}{{range $imp := .Imports}}{{printf "%s\t%s\n" $p $imp}}{{end}}`, "./...")

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
		reponly := map[string]int{}
		for err == nil {
			if line, err = reader.ReadString('\n'); err == nil {
				pi := strings.Split(strings.TrimSpace(string(line)), "\t")
				p, i := pi[0], pi[1]
				parts := strings.Split(i, "/")
				if len(parts) >= 3 && parts[2] != self {
					if !transitive && strings.Contains(p, "/vendor") {
						continue
					}
					matched := 0
					for _, p := range patterns {
						if p.MatchString(i) {
							matched += 1
						}
					}
					if matched != len(patterns) {
						continue
					}
					repo := strings.Join(parts[:3], "/")
					reponly[repo] += 1
				}
			}
		}
		repos := make([]string, len(reponly))
		ri := 0
		for repo, _ := range reponly {
			repos[ri] = repo
			ri += 1
		}
		sort.Strings(repos)
		for _, repo := range repos {
			if list {
				fmt.Printf("%s\n", repo)
				continue
			}
			var cmd *exec.Cmd
			b := path.Join("vendor", repo)
			if _, err := os.Stat(b); err == nil {
				if verbose {
					fmt.Printf("%s exists\n", repo)
				}
				continue
			} else if subtree {
				cmd = exec.Command(_git, "subtree", "add", "--prefix", fmt.Sprintf("vendor/%s", repo), fmt.Sprintf("https://%s", repo), branch, "--squash")
			} else {
				fullSubmodule := []string{"-f", "-b", branch, fmt.Sprintf("https://%s", repo), fmt.Sprintf("vendor/%s", repo)}
				var submoduler []string
				if force {
					submoduler = fullSubmodule
				} else {
					submoduler = fullSubmodule[1:]
				}
				cmd = exec.Command(_git, append([]string{"submodule", "add"}, submoduler...)...)
			}

			if dryrun {
				fmt.Println(strings.Join(cmd.Args, " "))
			} else {
				if output, err := cmd.CombinedOutput(); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\noutput:\n%s", cmd.Args, output)
				} else if verbose {
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
