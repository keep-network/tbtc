# tBTC

tBTC is a trustlessly Bitcoin-backed ERC-20 token.

The goal of the project is to provide a stronger 2-way peg than federated
sidechains like [Liquid](https://blockstream.com/liquid/), expanding use cases
possible via today's Bitcoin network, while bringing superior money to other
chains.

This is the primary repo of the project, containing the
[project documentation](docs/index.adoc) and majority of the [on-chain code](implementation/).

## Building

[![CircleCI](https://circleci.com/gh/keep-network/tbtc.svg?style=svg&circle-token=ec728f5ca814b6cb2db5ffeb7258151b752a207e)](https://circleci.com/gh/keep-network/tbtc)

The docs and Solidity are built [via CircleCI](.circleci/).

## Documentation

The documentation includes a project overview and rationale, as well as the
on-chain specification. Docs should always be updated before or in tandem with
code. 

Latest build from master is available at [http://docs.keep.network/tbtc/](http://docs.keep.network/tbtc/) ([pdf](http://docs.keep.network/tbtc/index.pdf)).

### Building

Docs are written in [AsciiDoctor](http://asciidoctor.org/), diagrams in a LaTeX package called [Tikz](https://www.overleaf.com/learn/latex/TikZ_package). To build:

```sh
cd docs

# Generate diagrams
make pngs
# Generate index.pdf
asciidoctor-pdf index.adoc
```

#### macOS

Install [TeX Live](https://www.tug.org/texlive/) manually, and other dependencies using CLI:

```sh
gem install asciidoctor-pdf --pre
brew install poppler
```

Install the TikZ package to your local LaTeX environment:

```sh
sudo cp docs/latex/tikz-uml.sty /usr/local/texlive/texmf-local/

# Update TeX package tree
sudo texhash
```
