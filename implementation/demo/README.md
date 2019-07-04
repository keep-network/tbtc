## Set Up

### Bitcoin SPV

1. Clone [bitcoin-spv](https://github.com/summa-tx/bitcoin-spv) repository
   ```sh
   git clone https://github.com/summa-tx/bitcoin-spv.git
   ```
   
2. Configure environment variable with a path to the cloned repository, e.g.:
    ```sh
    export BITCOIN_SPV_DIR="/Users/jakub/workspace/bitcoin-spv/"
    ```

3. Install python 3.6.6:
    ```sh
    brew install pyenv
    pyenv install 3.6.6
    ```

4. Set up python environment:
    ```sh
    brew install pipenv
    cd $BITCOIN_SPV_DIR && pipenv install --python ~/.pyenv/versions/3.6.6/bin/python
    ```

5. Compile contracts
   ```sh
   cd $BITCOIN_SPV_DIR
   npm install
   npm run compile
   ```
   Requires `node` version 12.3.1 or greater. To install latest node version run
   `npm install -g node`
