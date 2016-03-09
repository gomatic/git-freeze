# Easy Go vendoring

`git freeze` will submodule all git-package Go imports from `./...` into `vendor/`.

## Install

    go get github.com/nicerobot/git-freeze

## Usage

    git freeze

## Note

`git freeze -transitive` will traverse `./...`, so running it multiple times will continually freeze transitive dependencies. But, ideally, you do not want to use `-transitive` since transitive dependencies should be vendored/frozen by the package maintainer.
