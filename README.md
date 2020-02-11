# tBTC

[![CircleCI Build Status](https://circleci.com/gh/keep-network/tbtc.svg?style=svg&circle-token=ec728f5ca814b6cb2db5ffeb7258151b752a207e)](https://circleci.com/gh/keep-network/tbtc)
[![Docs](https://img.shields.io/badge/docs-website-yellow.svg)](http://docs.keep.network/tbtc/solidity/)
[![Chat with us on Discord](https://img.shields.io/badge/chat-Discord-blueViolet.svg)](https://discord.gg/4R6RGFf)

tBTC is a trustlessly Bitcoin-backed ERC-20 token.

The goal of the project is to provide a stronger 2-way peg than federated
sidechains like [Liquid](https://blockstream.com/liquid/), expanding use cases
possible via today's Bitcoin network, while bringing superior money to other
chains.

This repo contains the Solidity [smart contracts](implementation/) and [whitepaper](documentation/).

## Getting started

 * Read the [Whitepaper](http://docs.keep.network/tbtc/)
 * Consult [Solidity API documentation](http://docs.keep.network/tbtc/solidity/)
 * For questions and support, join us on [Discord](https://discord.gg/4R6RGFf).

## Installation

```sh
$ npm install @tbtc/contracts
```

## Usage

Once installed, you can use the contracts in the library by importing them:

```sol
pragma solidity ^0.5.0;

import "@tbtc/contracts/deposit/Deposit.sol";

contract MySystem {
    function checkTerm() external {
        uint256 remainingTerm = Deposit(_depositAddress).remainingTerm();
    }
}
```

## Contributing

We strongly recommend that the community help us make improvements and determine the future direction of the protocol. To report bugs within this package, please create an issue in this repository.

**Read our [Contributing guidelines](https://github.com/keep-network/tbtc/blob/master/CONTRIBUTING.md).**

### Prerequisites

 * Node.js, [npm](https://docs.npmjs.com/cli/install)
 * A local Ethereum blockchain. We recommend [Ganache](https://www.trufflesuite.com/ganache).
 * Truffle framework.

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

We use ESLint and Solium for linting code. To run:

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