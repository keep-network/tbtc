tBTC
====

Contains the on-chain smart contracts and tests.

### Setup

```sh
npm install
```

#### Migrating externals

tBTC interacts with external systems deployed on-chain, such as [Uniswap](https://uniswap.exchange). 

During development, these must be deployed to your local blockchain (ie. Ganache).

```sh
cd scripts/
./migrate-externals.sh
```

### Compilation

```sh
npm run compile
```

### Lint

```sh
# Show issues
npm run js:lint
npm run sol:lint

# Automatically fix issues
npm run js:lint:fix
npm run sol:lint:fix
```

Eslint errors can be disabled using a comment on the previous line. For example, 
to disable linter errors for the 'no-unused-vars' rule: 
`// eslint-disable-next-line no-unused-vars`.