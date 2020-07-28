pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "../../../contracts/external/IMedianizer.sol";

/// @title IntermediaryMedianizer is an updatable intermediary between a real
///        medianizer and IMedianizer users.
/// @dev This is used in Keep testnets where Maker has deployed a Medianizer
///      instance that needs to authorize a single consumer, to enable multiple
///      tBTC deployments to happen in the background and be pointed at a stable
///      medianizer instance that is authorized on the Maker contract. It allows
///      the updating of the backing medianizer and therefore is NOT suitable
///      for mainnet deployment.
contract IntermediaryMedianizer is Ownable, IMedianizer {
    IMedianizer private _realMedianizer;

    constructor(IMedianizer realMedianizer) public {
        _realMedianizer = realMedianizer;
    }

    function getMedianizer() external view returns (IMedianizer) {
        return _realMedianizer;
    }

    function peek() external view returns (uint256, bool) {
        return _realMedianizer.peek();
    }

    function read() external view returns (uint256) {
        return _realMedianizer.read();
    }

    function setMedianizer(IMedianizer realMedianizer) public onlyOwner {
        _realMedianizer = realMedianizer;
    }
}
