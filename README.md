# Easy Go vendoring

`git freeze` will submodule all git-package Go imports from `./...` into `vendor/`.

## Install

    go get github.com/nicerobot/git-freeze

## Usage

    git freeze

## Note

Since `git freeze` will traverse `./...`, running `git freeze` multiple times will freeze transitive dependencies. This may not be what you want. Ideally, all transitive dependencies will be vendored/frozen, providing build consistency.
