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

    // Allow sending funds to a keep group
    // Expected: increment their existing ETH bond
    function distributeEthToKeepGroup(address _keepAddress) external payable returns (bool);

    // Allow sending tokens to a keep group
    // Useful for sending signers their TBTC
    // The Keep contract should call transferFrom on the token contract
    function distributeERC20ToKeepGroup(address _keepAddress, address _asset, uint256 _value) external returns (bool);

    // returns the amount of the keep's ETH bond in wei
    function checkBondAmount(address _keepAddress) external view returns (uint256);

    // seize the signer's ETH bond
    // onlyKeepOwner
    // msg.sender.transfer(bondAmount)
    function seizeSignerBonds(address _keepAddress) external returns (bool);
}
