// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

contract Pot {
    address private owner;
    mapping (address => uint256) private bettingPot;
    mapping (address => uint256) private rewardsPot;

    modifier restrictAccess (address sender) {
        require(sender == owner, "This is not the owner");
        _;
    }

    constructor () public {
        owner = msg.sender;
    }

    function bet(address better, uint256 amount) public restrictAccess(msg.sender) {
        require (better != owner, "The owner cannot place a bet");
        require(amount > 0, "The betting amount must be greater than zero");
        if (bettingPot[better] == 0){
            bettingPot[better] = amount;
        } else {
            bettingPot[better] += amount;
        }
    }

    function getCurrentBet(address better) public restrictAccess(msg.sender) returns(uint256) {
        require (better != owner, "The owner has no bets");
        return bettingPot[better];
    }

    function getCurrentRewards(address player) public restrictAccess(msg.sender) returns(uint256) {
        require (player != owner, "The owner has no rewards");
        return rewardsPot[player];
    }

    function declareWinner(address winner) public restrictAccess(msg.sender) {
        require(bettingPot[winner] > 0, "This player has no bet");
        uint256 rewards = bettingPot[winner];
        delete bettingPot[winner];
        if (rewardsPot[winner] == 0){
            rewardsPot[winner] = rewards;
        } else {
            rewardsPot[winner] += rewards;
        }
    }

    function declareLoser(address loser) public restrictAccess(msg.sender) {
        require(bettingPot[loser] > 0, "This player has no bet");
        delete bettingPot[loser];
    }

    //Todo: Transfer actual funds to a players wallet
    function collectRewards(address collector) public restrictAccess(msg.sender) returns(uint256){
        require(rewardsPot[collector] > 0, "This player has no rewards");
        uint256 rewards = rewardsPot[collector];
        delete rewardsPot[collector];
        return rewards;
    }
}