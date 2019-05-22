# End-to-End Hacky Test With Single Signer

Below are described steps required to run the test.

Prepare 3 terminals in which you will run each section.

## Ganache

Execute command to run Ganache:
```sh
ganache-cli
```

## Keep TECDSA

**Repository**: [keep-network/keep-tecdsa](https://github.com/keep-network/keep-tecdsa)

**Branch**: `git checkout client-cmd`

1. Deploy contracts
    ```sh
    cd solidity && truffle migrate -—reset
    ```

    You will need two contract's addresses later:
    - `ECDSAKeepFactory`
    - `KeepRegistry`

1. Update `config.toml` with `ECDSAKeepFactory` contract address.

1. Build keep client:
   ```sh
   cd ../ && go build -a -o keep-tecdsa .
   ```

1. Run keep client:
   ```sh
   ./keep-tecdsa start
   ```
   You should see `Client started.` output message.

## tBTC

**Repository**: [keep-network/tbtc](https://github.com/keep-network/tbtc)

**Branch**: `git checkout e2e-signer-test`

1. Deploy contracts
    ```sh
    cd implementation/
    truffle migrate -—reset
    ```

1. Prepare [demo/create_new_deposit.js](create_new_deposit.js) script:
   - Provide `KeepRegistry` contract address (copied from Keep TECDSA) in constant:
   ```js
   const keepRegistry = "0x41649ff7E9E4512fbe8b42A51d73f33377D012c5"; // KeepRegistry contract address
   ```

1. Execute the `create_new_deposit.js` script with truffle:
   ```sh
   truffle exec demo/create_new_deposit.js
   ```

## Expected results
After setting everything up and running the script you should see following results in terminal:
- From running `create_new_deposit.js` script:
    ```
    Using network 'development'.

    Call createNewDeposit
    Creation tx: 0x503e1b5bb19a47285dcaee5c331a71bb3caf47103673ed1ee4067d66a6a35121
    ```
- From `keep-tecdsa` client:
    ```
    Client started.
    New ECDSA Keep created [&{KeepAddress:[144 33 160 172 142 172 215 7 122 214 228 75 19 211 63 143 213 106 86 79]}]
    Signer for keep [0x9021a0AC8eAcD7077AD6E44b13d33f8fd56a564f] initialized with Bitcoin P2WPKH address [tb1qenqac9xua2fe2vxxhwkvk5w5ql2uhgadx7j2nw]
    ```
    
    `tb1qenqac9xua2fe2vxxhwkvk5w5ql2uhgadx7j2nw` is an address which needs to be
    funded with bitcoin transaction.
