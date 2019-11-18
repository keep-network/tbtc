pragma solidity ^0.5.10;

// TODO: This is an interface holding functions which are required to be
// implemented on keep side.
interface IBondedECDSAKeep {
    // Expected behavior:
    // Error if not fraud
    // Return true if fraud
    // This means if the signature is valid, but was not approved via approveDigest
    function submitSignatureFraud(
        address _keepAddress,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes calldata _preimage
    ) external returns (bool _isFraud);

    // returns the amount of the keep's ETH bond in wei
    function checkBondAmount(address _keepAddress) external view returns (uint256);

    // seize the signer's ETH bond
    // onlyKeepOwner
    // msg.sender.transfer(bondAmount)
    function seizeSignerBonds(address _keepAddress) external returns (bool);
}
