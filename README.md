# GitUp

[![GoDoc](https://godoc.org/github.com/rvflash/gitup?status.svg)](https://godoc.org/github.com/rvflash/gitup)
[![Build Status](https://img.shields.io/travis/rvflash/gitup.svg)](https://travis-ci.org/rvflash/gitup)
[![Code Coverage](https://img.shields.io/codecov/c/github/rvflash/gitup.svg)](http://codecov.io/github/rvflash/gitup?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/gitup)](https://goreportcard.com/report/github.com/rvflash/gitup)


Automatically checks and updates if wished the latest tag of a Git repository.


## Features

You can use 3 level of strategy: Noop, Manual and Auto.
The first, does anything. The second asks a confirmation to the user on the standard input and the last,
automatically updates the repository with the latest available tag.

## Usage

See the GitUp test for an example of using.

## Use SemVer for the version tag name

The tag name must be compliant to the [Semantic Versioning 2.0](http://semver.org/spec/v2.0.0.html).
As example, v2.0.0-beta is a valid version tag but not 2.0.