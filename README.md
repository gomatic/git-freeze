# Easy Go vendoring

`git freeze` will submodule or subtree all git-package Go imports from `./...` into `vendor/`.

# Status

[![Build Status](https://travis-ci.org/nicerobot/git-freeze.png?branch=master)](https://travis-ci.org/nicerobot/git-freeze)

## Install

		go get github.com/nicerobot/git-freeze

## Usage

If your `GOBIN` is also in your `PATH`, `git-freeze` will be accessible as:

		git freeze

## Help

		Usage:
			-branch string
						Git branch/commit to submodule/subtree. (defaults to the parent's branch)
			-dry-run
						Just print the command but do not run it.
			-force
						Force.
			-list
						Only list the imports that can be frozen.
			-notests
						Do not freeze test-imports.
			-subtree
						Use a subtree instead of a submodule.
			-transitive
						Traverse transitive imports, i.e. vendor/
			-verbose
						More output.

## Note

`git freeze -transitive` will traverse `./...`, so running it multiple times will continually freeze transitive dependencies. But, ideally, you do not want to use `-transitive` since transitive dependencies should be vendored/frozen by the package maintainer.
