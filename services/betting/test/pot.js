const Pot = artifacts.require("Pot");
const truffleAssert = require('truffle-assertions');

contract("Pot", function (accounts) {

    const owner = accounts[0];
    const player = accounts[1];
    var instance;

    it("Should deny owner placing bet", async function (){
        instance = await Pot.deployed();
        try {
            await instance.bet.call(owner, 5)
            throw null
        } catch(error) {
            assert(error, "Expected an error but did not get one");
        }
    });

    it("Should revert when given a zero as the betting amount", async function (){
        instance = await Pot.deployed();
        try {
            await instance.bet.call(player, 0)
            throw null
        } catch(error) {
            assert(error, "Expected an error but did not get one");
        }
    });

    it("Should allow a player to place a bet", async function (){
        instance = await Pot.deployed();
        await instance.bet(player, 10)
        assert.equal(await instance.getCurrentBet.call(player), 10)
    });

    it("Should revert if player has no rewards", async function (){
        instance = await Pot.deployed();
        try {
            await instance.collectRewards(player)
            throw null
        } catch(error) {
            assert(error, "Expected an error but did not get one");
        }
    });

    it("Should allow player to place and retrieve bet and clear their rewards from the pot", async function (){
        instance = await Pot.deployed();
        await instance.bet(player, 5);
        await instance.declareWinner(player);
        await instance.collectRewards(player)
        assert.equal(await instance.getCurrentBet.call(player), 0);
        assert.equal(await instance.getCurrentRewards.call(player), 0);
    });


    it("Should allow player to place multiple bets and retrieve the correct amount later", async function (){
        instance = await Pot.deployed();
        await instance.bet(player, 5);
        await instance.bet(player, 10);
        await instance.bet(player, 15);
        await instance.declareWinner(player);
        await instance.collectRewards(player)
        assert.equal(await instance.getCurrentRewards.call(player), 0);
    });
});