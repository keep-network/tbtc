## tbtc Contribution Guide

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

The following is a set of guidelines for contributing to tBTC and its packages. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

### Getting started

1.  Fork [`keep-network/tbtc`](https://github.com/keep-network/tbtc)
2.  Clone your fork
3.  Follow the [installation & build steps](https://github.com/keep-network/tbtc/tree/master/implementation#setup) in the repo's top-level README.
4.  Setup the recommended [Development Tooling](#development-tooling).
5.  Open a PR against the `master` branch and describe the change you are intending to undertake in the PR description. 

Before marking the PR as open for review, make sure:

-   It passes our linter checks (`npm run lint`)
-   It passes our continuous integration tests.
-   Your changes have sufficient test coverage (e.g regression tests have been added for bug fixes, unit tests for new features)

### Development Tooling

Commits [must be signed](https://help.github.com/en/articles/about-commit-signature-verification).

#### Continuous Integration

We use [CircleCI](https://circleci.com) for continuous integration. All CircleCI jobs (including tests, linting) must be green to merge a PR.

#### Pre-commit

Pre-commit is a tool to install hooks that check code before commits are made. It can be helpful to install this, to automatically run linter checks. Follow the [installation instructions here](https://pre-commit.com/), and then run ```pre-commit install``` to install the hooks.

#### Linting

Linters for Solidity and JavaScript code are setup and run automatically as part of pre-commit hooks.

If you want to change a rule, or add a custom rule, please make these changes to our [solium-config-keep](https://github.com/keep-network/solium-config-keep) and [eslint-config-keep](https://github.com/keep-network/eslint-config-keep) packages. All other packages have it as a dependency.