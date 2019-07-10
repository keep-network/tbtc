pragma solidity 0.4.25;

/**
 * @title Keep interface
 */

interface IKeep {

    // request a new m-of-n group
    // should return a 256 unique keep id
    function requestKeepGroup(uint256 _m, uint256 _n) external payable returns (uint256 _keepID);

    // get the result of a keep formation
    // should return a 64 byte packed pubkey (x and y)
    // error if not ready yet
    function getKeepPubkey(uint256 _keepID) external view returns (bytes);

    // onlyKeepOwner
    // record a digest as approved for signing
    function approveDigest(uint256 _keepID, bytes32 _digest) external returns (bool _success);

    // returns the timestamp when it was approved
    function wasDigestApprovedForSigning(uint256 _keepID, bytes32 _digest) external view returns (uint256);

    // Expected behavior:
    // Error if not fraud
    // Return true if fraud
    //     This means if the signature is valid, but was not approved via approveDigest
    function submitSignatureFraud(
        uint256 _keepID,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) external returns (bool _isFraud);

    // Allow sending funds to a keep group
    // Expected: increment their existing ETH bond
    function distributeEthToKeepGroup(uint256 _keepID) external payable returns (bool);

    // Allow sending tokens to a keep group
    // Useful for sending signers their TBTC
    // The Keep contract should call transferFrom on the token contrcact
    function distributeERC20ToKeepGroup(uint256 _keepID, address _asset, uint256 _value) external returns (bool);

    // returns the amount of the keep's ETH bond in wei
    function checkBondAmount(uint256 _keepID) external view returns (uint256);

    // seize the signer's ETH bond
    // onlyKeepOwner
    // msg.sender.transfer(bondAmount)
    function seizeSignerBonds(uint256 _keepID) external returns (bool);
}
