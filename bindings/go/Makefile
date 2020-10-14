solidity_dir=$(realpath ${SOLIDITY_DIR})
solidity_files := $(foreach file,${SOLIDITY_FILES},$(file))

# Solidity filenames without .sol with subdirectories prefixes.
contract_stems := $(basename $(solidity_files))

abi_files := $(addprefix abi/,$(addsuffix .abi,$(contract_stems)))
abigen_files := $(addprefix abi/,$(addsuffix .go,$(contract_stems)))
contract_files := $(addprefix contract/,$(addsuffix .go,$(contract_stems)))

all: gen_contract_go gen_abi_go

clean:
	rm -r abi/*
	rm -r contract/*
	mkdir tmp && mv cmd/cmd*.go tmp
	rm -r cmd/*
	mv tmp/* cmd && rm -r tmp

gen_abi_go: $(abigen_files)

gen_contract_go: $(contract_files)

abi/%.abi: ${solidity_dir}/contracts/%.sol
	$(eval output=$(dir $@))
	mkdir -p $(output)
	solc openzeppelin-solidity/=${solidity_dir}/node_modules/openzeppelin-solidity/ \
	     @summa-tx/bitcoin-spv-sol/=${solidity_dir}/node_modules/@summa-tx/bitcoin-spv-sol/ \
		 @keep-network/keep-ecdsa/=${solidity_dir}/node_modules/@keep-network/keep-ecdsa/ \
		 @summa-tx/relay-sol=${solidity_dir}/node_modules/@summa-tx/relay-sol/ \
		 --allow-paths ${solidity_dir} \
		 --overwrite \
		 --abi \
		 -o $(output) $<

abi/%.go: abi/%.abi
	$(eval type=$(notdir $*))
	go run github.com/ethereum/go-ethereum/cmd/abigen --abi $< --pkg abi --type $(type) --out $@

contract/DepositLog.go cmd/DepositLog.go: abi/DepositLog.abi abi/DepositLog.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethereum $< contract/DepositLog.go cmd/DepositLog.go

contract/deposit/Deposit.go cmd/deposit/Deposit.go: abi/deposit/Deposit.abi abi/deposit/Deposit.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethereum $< contract/Deposit.go cmd/Deposit.go





