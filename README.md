# GitUp

Automatically checks and updates if wished the latest tag of a Git repository.


## Features

You can use 3 level of strategy: Noop, Manual and Auto.
The first, does anything. The second asks a confirmation to the user on the standard input and the last,
automatically updates the repository with the latest available tag.

## Usage

See the GitUp test for an example of using.

## Use SemVer for the version tag name

The tag name must be compliant to the Semantic Versioning 2.0, see semver.org/spec/v2.0.0.html for more information.
As example, v2.0.0-beta is a valid version tag but not 2.0.