pragma solidity ^0.5.10;

/**
 * @title Keep interface
 */

interface IKeep {

    // get the result of a keep formation
    // should return a 64 byte packed pubkey (x and y)
    // error if not ready yet
    function getKeepPubkey(address _keepAddress) external view returns (bytes memory);

    /// @notice Approves digest for signing.
    /// @param _keepAddress Keep address
    /// @param _digest Digest to sign
    /// @return True if successful.
    function approveDigest(address _keepAddress, bytes32 _digest) external returns (bool _success);

    /// @notice Gets timestamp of digest approval for signing.
    /// @param _keepAddress Keep address
    /// @param _digest Digest to sign
    /// @return Timestamp from the moment of recording the digest for signing.
    /// Returns 0 if the digest was not recorded for signing for the given keep.
    function wasDigestApprovedForSigning(address _keepAddress, bytes32 _digest) external view returns (uint256);

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
    // The Keep contract should call transferFrom on the token contrcact
    function distributeERC20ToKeepGroup(address _keepAddress, address _asset, uint256 _value) external returns (bool);

    // returns the amount of the keep's ETH bond in wei
    function checkBondAmount(address _keepAddress) external view returns (uint256);

    // seize the signer's ETH bond
    // onlyKeepOwner
    // msg.sender.transfer(bondAmount)
    function seizeSignerBonds(address _keepAddress) external returns (bool);
}
