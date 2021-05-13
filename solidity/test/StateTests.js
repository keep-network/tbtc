const states = require('./states/deposit.js')
const { deployAndLinkAll } = require('./helpers/fullDeployer.js')
const { runStatePath } = require('./states/run.js')

describe('tBTC states', () => {
  describe('when running a redemption', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'awaitingWithdrawalSignature',
    )
  })
  describe('when liquidating from a courtesy call', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'courtesyCall',
      'liquidationInProgress',
    )
  })
  describe('when fraud liquidating from a courtesy call', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'courtesyCall',
      'liquidationInProgress_fraud',
    )
  })
  describe('when liquidating from severe undercollateralization', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'liquidationInProgress',
    )
  })
  describe('when liquidiating during a redemption', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'awaitingWithdrawalSignature',
      'liquidationInProgress',
    )
  })
  describe('when minting TBTC then running a redemption', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'minted',
      'awaitingWithdrawalSignature',
    )
  })
  describe('when minting TBTC then redeeming from a courtesy call', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'minted',
      'courtesyCall',
      'awaitingWithdrawalSignature',
    )
  })
  describe('when minting TBTC then liquidating from a courtesy call', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'minted',
      'courtesyCall',
      'liquidationInProgress',
    )
  })
  describe('when minting TBTC then fraud liquidating from a courtesy call', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'minted',
      'courtesyCall',
      'liquidationInProgress_fraud',
    )
  })
  describe('when minting TBTC then liquidating from severe undercollateralization', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'minted',
      'liquidationInProgress',
    )
  })
  describe('when minting TBTC then liquidiating during a redemption', () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      'start',
      'awaitingSignerSetup',
      'awaitingFundingProof',
      'active',
      'minted',
      'awaitingWithdrawalSignature',
      'liquidationInProgress',
    )
  })
})
