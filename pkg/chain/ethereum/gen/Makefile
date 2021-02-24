# Environment `SOLIDITY_DIR` provides the solidity directory as a
# potentially-relative path, which we resolve. However, because
# Solidity files are grouped into separate subdirectories, specific
# files must be provided explicitly via `SOLIDITY_FILES` variable
# to handle them in a correct way.
solidity_dir=$(realpath ${SOLIDITY_DIR})
solidity_files := $(foreach file,${SOLIDITY_FILES},$(file))

# Solidity filenames without .sol with subdirectories prefixes.
contract_stems := $(basename $(solidity_files))
# Go abigen bindings in abi/ subdirectory with .go suffix, alongside solc ABI
# files with .abi suffix.
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

contract/deposit/Deposit.go cmd/deposit/Deposit.go: abi/deposit/Deposit.abi abi/deposit/Deposit.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethlike $< contract/Deposit.go

contract/system/TBTCSystem.go cmd/system/TBTCSystem.go: abi/system/TBTCSystem.abi abi/system/TBTCSystem.go *.go
	go run github.com/keep-network/keep-common/tools/generators/ethlike $< contract/TBTCSystem.go




