## tbtc Contribution Guide

### Getting started

1.  Fork [`keep-network/tbtc`](https://github.com/keep-network/tbtc)
2.  Clone your fork
3.  Follow the [installation & build steps](https://github.com/keep-network/tbtc/tree/master/implementation#setup) in the repo's top-level README.
4.  Setup the recommended [Development Tooling](#development-tooling).
5.  Open a PR against the `master` branch and describe the change you are intending to undertake in the PR description. 

### Development Tooling

Commits [must be signed](https://help.github.com/en/articles/about-commit-signature-verification).

#### Continuous Integration

We use [CircleCI](https://circleci.com) for continuous integration. All CircleCI jobs (including tests, linting) must be greent to merge a PR.

#### Pre-commit

Pre-commit is a tool to install hooks that check code before commits are made. Follow the [installation instructions here](https://pre-commit.com/), and then run ```pre-commit install``` to install the hooks.

#### Linting

Linters for Solidity and JavaScript code are setup and run automatically as part of pre-commit hooks.