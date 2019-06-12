import sys
import json
import asyncio

from connectrum.svr_info import ServerInfo
from connectrum.client import StratumClient

from riemann import tx
from riemann import utils as rutils

from typing import Tuple

# # # # # # # # # # # # # # #
# Use this script Sparingly #
# # # # # # # # # # # # # # #

# try:
#     with open('../build/contracts/ValidateSPV.json', 'r') as jsonfile:
#         j = json.loads(jsonfile.read())
#         ABI = json.loads(j['interface'])
# except Exception:
#     raise ValueError('Couldn\'t open ABI file. Hint: Did you npm run compile?')

CLIENT = None

async def get_client():
    global CLIENT
    if CLIENT is not None:
        return CLIENT
    else:
        CLIENT = await setup_client()
        return CLIENT

async def setup_client():
    if CLIENT is not None:
        return CLIENT

    server = ServerInfo({
        "nickname": None,
        "hostname": "testnet.hsmiths.com",
        "ip_addr": None,
        "ports": [
            "s53012"
        ],
        "version": "1.4",
        "pruning_limit": 0,
        "seen_at": 1533670768.8676858
    })

    client = StratumClient()

    await asyncio.wait_for(
        client.connect(
            server_info=server,
            proto_code='s',
            use_tor=False,
            disable_cert_verify=True),
        timeout=5)

    # await asyncio.wait_for(
    #     client.RPC(
    #         'server.version',
    #         'bitcoin-spv-merkle',
    #         '1.2'),
    #     timeout=5)

    return client

# def make_ether_data(t: tx.Tx, proof: bytes, index: bytes, header: bytes):
#     '''Creates a data blob for a transaction calling validateTransaction'''
#     calldata.call(
#         'validateTransaction',
#         [t.to_bytes(), proof, index, header],
#         ABI)

async def get_latest_blockheight() -> int:
    '''Gets the electrum server's latest known blockheight'''
    client = await get_client()
    fut, _ = client.subscribe('blockchain.headers.subscribe')
    block_dict = await fut
    height = block_dict['height'] \
        if 'height' in block_dict \
        else block_dict['block_height']
    return height


async def get_block_merkle_root(height: int) -> bytes:
    '''Gets the merkle root of the block at a specified height'''
    client = await get_client()

    header_dict = await client.RPC('blockchain.block.headers', height, 1)
    merkle_root = bytes.fromhex(header_dict['hex'])[36:68]

    return merkle_root


async def get_tx(tx_id: str) -> Tuple[dict, tx.Tx]:
    client = await get_client()
    tx_dict = await client.RPC('blockchain.transaction.get', tx_id, True)
    t = tx.Tx.from_hex(tx_dict['hex'])

    latest_blockheight = await get_latest_blockheight()

    # NB: This has a small probability of causing failure
    #     If the server updates its highest between when we get the TX dict
    #     when we get the latest height

    tx_dict['block_height'] = latest_blockheight - tx_dict['confirmations'] + 1

    return tx_dict, t


async def get_header_chain(start_height: int,
                           count: int) -> str:
    client = await get_client()

    res = await client.RPC(
        'blockchain.block.headers', start_height, count)

    return res['hex']
    

async def get_merkle_proof_from_api(tx_id: str, hght: int) -> Tuple[str, int]:
    client = await get_client()

    res = await client.RPC('blockchain.transaction.get_merkle', tx_id, hght)

    pos = res['pos']

    proof = bytearray()
    proof.extend(bytes.fromhex(tx_id)[::-1])
    for tx_id in res['merkle']:
        proof.extend(bytes.fromhex(tx_id)[::-1])

    block_root = await get_block_merkle_root(hght)

    # print(block_root)
    proof.extend(block_root)

    # NB: add 1 because our proof uses 1-indexed position
    return proof.hex(), pos + 1


def verify_proof(proof: bytes, index: int):
    index = index  # This is 1 indexed
    # TODO: making creating and verifying indexes the same
    root = proof[-32:]
    current = proof[0:32]

    # For all hashes between first and last
    for i in range(1, len(proof) // 32 - 1):
        # If the current index is even,
        # The next hash goes before the current one
        if index % 2 == 0:
            current = rutils.hash256(
                proof[i * 32: (i + 1) * 32]
                + current
            )
            # Halve and floor the index
            index = index // 2
        else:
            # The next hash goes after the current one
            current = rutils.hash256(
                current
                + proof[i * 32: (i + 1) * 32]
            )
            # Halve and ceil the index
            index = index // 2 + 1
    # At the end we should have made the root
    if current != root:
        return False
    return True


async def do_it_all(tx_id: str, num_headers: int):
    (tx_json, t) = await get_tx(tx_id)

    proof, index = await get_merkle_proof_from_api(
        t.tx_id.hex(), tx_json['block_height'])

    # Create a header chain
    chain = await get_header_chain(
        tx_json['block_height'],
        num_headers + 1)

    # submission = make_ether_data(t, proof, index, header)

    # Error if the proof isn't valid
    assert(verify_proof(bytes.fromhex(proof), index))

    datadict = {
        "tx": t.hex(),
        "proof": proof,
        "index": index,
        "chain": chain,
        "txid": sys.argv[1],
        'chainLen': sys.argv[2]
    }

    with open('spv_out.json', 'w', encoding='utf-8') as outfile:
        json.dump(datadict, outfile, ensure_ascii=False, indent=2)


def main(tx_id: str, num_headers: int):
    asyncio.get_event_loop().run_until_complete(do_it_all(tx_id, num_headers))

if __name__ == '__main__':
    # Read tx_id from args, and then get it and its block from explorers
    tx_id = str(sys.argv[1])
    num_headers = int(sys.argv[2]) if len(sys.argv) > 2 else 6
    main(tx_id, num_headers)