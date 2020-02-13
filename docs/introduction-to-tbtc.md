# Introduction to tBTC

This guide serves as an introduction to the tBTC system and its core concepts. It is written for developers who want to understand tBTC and how it works at a high-level.

If you are looking for a reference guide, consult the [official specification](http://docs.keep.network/tbtc/).

*You will learn about*: the Deposit, TBTC vs. TDT's, funding, redemption and liquidation.

## What is tBTC?

tBTC is a trustless bridge of Bitcoin onto Ethereum. It produces the TBTC token, an ERC20 backed 1:1 by bitcoin, redeemable at any time. 

Using the Keep protocol, a Bitcoin wallet is created and secured by a decentralized set of signers. A user deposits BTC into this wallet, and proves the deposit transaction on the Ethereum chain. From this deposit, TBTC can be minted and transacted. TBTC can be redeemed for the underlying BTC at any time, and the signers will sign a transaction sending it to a user-provided address.

All Bitcoin transactions are proven using SPV proofs, which rely on an [on-chain relay](https://github.com/summa-tx/relays), operated by Summa.

## Why tBTC?

Bitcoin was the first ever digital money controlled by no government or corporation. Since its inception over a decade ago, it has remained robust against attack, and simple in its mission - digital, peer-to-peer money.

If Bitcoin is the spreadsheet of accounts and balances, Ethereum is the entire computer. Using Ethereum's smart contracts, more complex applications can be built with digital money. For the first time ever, users can do more than just transact - they can get loans without fear of discrimination, start digital organisations without financial borders, and *build new ways to collaborate with money*. 

A whole new world of digital money is out there, and at its core is openness. tBTC bridges the world's two biggest networks, Bitcoin and Ethereum, enabling value to flow freely between them. Using tBTC, bitcoin can be used as collateral for loans, lent to generate interest [on Compound](https://compound.finance/), owned by decentralised autonomous organisations...you name it!

Now that we've covered a bit of why tBTC exists, let's continue to understand how it works.

## Deposits 

The **`Deposit` contract** is the core of tBTC; Deposits are to tBTC as Vaults are to MakerDAO. Deposits are created by the `DepositFactory`, and are parameterised by:

* **lot size**: the BTC size of the deposit. *eg. 1.0 BTC*
* **signer fee**: the fees paid to signers. *eg. 0.5%*

These are governed by the `TBTCSystem` contract. Once a deposit is created, these parameters are set in stone.

Deposits have a fixed term of 6 months. While they can be redeemed at anytime, after 6 months they must be closed. This ensures signers are paid.

Lastly, Deposits can be locked or unlocked. A locked deposit can only be redeemed by the deposit owner, who holds a token called the tBTC Deposit Token. *Hold up! What's this new token?*

## TDT's and TBTC's

The **tBTC Deposit Token (TDT)** is a non-fungible counterpart to TBTC, representing a claim to the underlying UTXO on the Bitcoin blockchain. It is useful for more *advanced use cases* of Bitcoin on Ethereum. Some of these include:

* *Assigning value to higher-risk deposits*: Different deposits have different risk profiles. For example, a 1 BTC deposit is more valuable than 0.001 BTC to steal, and thus might be more susceptible to attacks like re-orgs etc. An [NFT](https://en.wikipedia.org/wiki/Non-fungible_token) facilitates this risk being '*priced in*', which is highly relevant for applications that use bitcoin as collateral.
* *Private/anonymous currency:* TDTs are an ideal target for secret fixed-size "notes" or other financial privacy improvements.

To recap:

1. TBTC - the fungible ERC20, intended for the average user.
2. TDT - the non-fungible ERC721, intended for advanced usage.

The TDT and TBTC are interchangeable for each other using a contract called the *vending machine*, which we'll get to in a moment. 

## Funding a Deposit

The funding flow is as follows:

1. Create deposit.
   1. Request a Keep signing group.
   2. Depositor receives a TDT.
2. Signers generate a Bitcoin address.
3. Depositor sends BTC and awaits confirmation on the Bitcoin chain.
4. Depositor proves funding transaction to Deposit.
   1. Deposit is marked as *qualified* for minting TBTC.

## Minting TBTC

After a Deposit is funded and *qualified*, it is ready for minting TBTC. The TDT represents the now-confirmed UXTO, and we are ready to swap it using the vending machine.

The **vending machine** manages the changing of TDT into TBTC and vice-versa. 

1. Given a TDT, it will mint TBTC.
2. Given TBTC, it will burn it and return the TDT.

## Redeeming a Deposit

The redemption flow is as follows:

1. Select a deposit to redeem and request redemption.
   1. Relinquish TDT and pay signer fees.
2. Signers receive redemption request, and produce a signature.
3. Redeemer builds the redemption transaction, broadcasts it, and awaits confirmation.
4. Redeemer receives BTC.
5. Signers prove the transaction, and receive their fees and bond back.

## Liquidation

The "happy paths" have been covered, but there hasn't been discussion of when *things fall apart*. 

To disincentivise signers from stealing bitcoin, deposits are overcollateralised with ETH. The collateral is priced using an on-chain [ETH:BTC price feed](https://github.com/keep-network/tbtc/blob/master/implementation/contracts/price-feed/BTCETHPriceFeed.sol), which is operated by MakerDAO. This collateral also ensures  guarantees of signer availability within the protocol.

---

*If you have any questions, suggestions, improvements, don't hesitate to reach out via GitHub and/or [chat on Discord](https://discord.gg/4R6RGFf).*