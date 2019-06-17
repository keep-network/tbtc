## tbtc Contribution Guide

We welcome contributions from anyone on the internet and are grateful for even the smallest contributions. This document will help get you setup to start contributing to tBTC.

### Getting started

1.  Fork `keep-network/tbtc`
2.  Clone your fork
3.  Follow the [installation & build steps](https://github.com/keep-network/tbtc/tree/master/implementation#setup) in the repo's top-level README.
4.  Setup the recommended [Development Tooling](#development-tooling).
5.  Open a PR with the `[WIP]` flag against the `development` branch and describe the change you are intending to undertake in the PR description. 

### Development Tooling

#### Pre-commit

Pre-commit is a tool to install hooks that check code before commits are made. Follow the [installation instructions here](https://pre-commit.com/), and then run ```pre-commit install``` to install the hooks.

#### Linting
Linters for Solidity and JavaScript code are setup and run automatically as part of pre-commit hooks.