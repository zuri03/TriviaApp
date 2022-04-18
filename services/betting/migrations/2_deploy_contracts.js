var Pot = artifacts.require("../contracts/Pot.sol")

module.exports = function(deployer) {
    deployer.deploy(Pot).then(function () {
        console.log(Pot.address)
        console.log("Pot Deployed")
    });
}