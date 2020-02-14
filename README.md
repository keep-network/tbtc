# tBTC

[![CircleCI Build Status](https://circleci.com/gh/keep-network/tbtc.svg?style=svg&circle-token=ec728f5ca814b6cb2db5ffeb7258151b752a207e)](https://circleci.com/gh/keep-network/tbtc)
[![Docs](https://img.shields.io/badge/docs-website-yellow.svg)](http://docs.keep.network/tbtc/solidity/)
[![Chat with us on Discord](https://img.shields.io/badge/chat-Discord-blueViolet.svg)](https://discord.gg/4R6RGFf)

tBTC is a trustlessly Bitcoin-backed ERC-20 token.

The goal of the project is to provide a stronger 2-way peg than federated
sidechains like [Liquid](https://blockstream.com/liquid/), expanding use cases
possible via today's Bitcoin network, while bringing superior money to other
chains.

This repo contains the Solidity [smart contracts](implementation/) and [specification](docs/).

## Getting started

 * Read the [introduction to tBTC](docs/introduction-to-tbtc.md).
 * Dig into the [specification](http://docs.keep.network/tbtc/).
 * Consult [Solidity API documentation](http://docs.keep.network/tbtc/solidity/).
 * For questions and support, join the #tbtc channel on [Discord](https://discord.gg/4R6RGFf).

## Installation

```sh
$ npm install @keep-network/tbtc
```

## Usage

**NOTE:** tBTC contracts require *solc* v0.5.10 or higher. You may have to [configure solc in your `truffle-config.js`](https://www.trufflesuite.com/docs/truffle/reference/configuration#compiler-configuration).

Once installed, you can use the contracts in the library by importing them:

```sol
pragma solidity ^0.5.0;

import "@keep-network/tbtc/contracts/deposit/Deposit.sol";

contract MySystem {
    function checkTerm(address payable _depositAddress) external {
        uint256 remainingTerm = Deposit(_depositAddress).remainingTerm();
    }
}
```

## Security

tBTC is currently under audit.

Please report any security issues you find to security@keep.network.

## Contributing

All contributions are welcome. To report bugs, please create an issue on this repository. To start a discussion, prefer Discord over GitHub issues.

**Read the [Contributing guidelines](https://github.com/keep-network/tbtc/blob/master/CONTRIBUTING.md).**

### Setup environment

You should have installed:

 * Node.js, [npm](https://docs.npmjs.com/cli/install).
 * A local Ethereum blockchain. We recommend [Ganache](https://www.trufflesuite.com/ganache).
 * [Truffle framework](https://www.trufflesuite.com/docs/truffle/overview).

### Build

Clone and install dependencies:

```sh
git clone https://github.com/keep-network/tbtc
cd tbtc/implementation
npm install
```

Deploy contracts:

```sh
truffle migrate --reset
```

### Test

Tests are written in JS using Mocha.

To run the test suite, execute `truffle test`.

To run specific tests, add [`.only`](https://jaketrent.com/post/run-single-mocha-test/) to the `contract` block:

```js
contract.only('TBTCToken', function(accounts) {
```

### Lint

We use [ESLint](https://eslint.org/) and [Ethlint](https://github.com/duaraghav8/Ethlint) for linting code. To run:

```sh
npm run sol:lint:fix
npm run js:lint:fix
```

## Documentation

The documentation includes a project overview and rationale, as well as the
on-chain specification. Docs should always be updated before or in tandem with
code.

### Prerequisites

Docs are written in [AsciiDoctor](http://asciidoctor.org/), diagrams in a LaTeX package called [Tikz](https://www.overleaf.com/learn/latex/TikZ_package).

##### macOS

 1. Install [TeX Live](https://www.tug.org/texlive/) manually, and other dependencies using CLI:

    ```sh
    gem install asciidoctor-pdf --pre
    brew install poppler
    ```

 2. Install the TikZ package to your local LaTeX environment:

    ```sh
    sudo cp docs/latex/tikz-uml.sty /usr/local/texlive/texmf-local/

    # Update TeX package tree
    sudo texhash
    ```

### Build

```sh
cd docs

# Generate diagrams
make pngs
# Generate index.pdf
asciidoctor-pdf index.adoc
```

## License

tBTC is released under the [MIT License](LICENSE).
