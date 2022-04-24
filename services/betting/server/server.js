//Endpoints
/*
Endpoints:
    - Submit: Provides clients the ability to declar a player as correct 
        - this triggers a function on the contract to move the players funds out of the betting map and into the rewards map
        - add measures to protect this endpoint and ensure the only internal services can call it 
        - this only moves funds from the betting map to the rewards map it is up to the player to pull their rewards from the contract
    
    - Bet: Provides clients the ability to add funds to a player's current bet or create a new entry in the betting map
        - accepts player account address and an amount (for now this amount will simply be a number)

    - Get rewards: A post request that allows clients to call a function on the contract to send user rewards to their accounts
        - accepts player account address 


import dotenv from "dotenv";
dotenv.config();
import Web3 from 'web3';
import fs from 'fs';
import express from 'express';
*/

const dotenv = require("dotenv");
dotenv.config();
const express = require('express');
const cors = require('cors');

module.exports = async function (owner, web3) {

    const { abi } = require('../build/Pot.json');
    const abi = "abi";
    const address = process.env.POT_CONTRACT_ADDRESS;
    
    const port = process.env.PORT;
    const contract = new web3.eth.Contract(abi, address, { from: owner });
    const app = express();

    
    app.use(express.urlencoded());
    app.use(cors());

    app.get("*/:address", (req, res, next) => {
        console.log("got get dynamic url request");
        var address = req.params.address;
        console.log(`testing address => ${address}`);
        if (!web3.utils.isAddress(address)) {
            res.status(400).json({data: "", error: `${address} is not a valid address`})
        } else {
            next();
        } 
    });

    app.post("*", (req, res, next) => {
        const address = req.body.address;
        console.log(`post request address => ${address}`)
        if (!web3.utils.isAddress(address)) {
            res.status(400).json({data: "", error: `${address} is not a valid address`})
        } else {
            next();
        } 
    })
    
    app.get('/rewards/:address', (req, res) => {
        
        const address = req.params.address;

        contract.methods.getCurrentRewards(address).call({from: owner}, function (error, result) {
            if (error){
                console.log(error);
                res.status(500).json({data: "", error: "Internal Server Error"})
            } else {
                res.status(200).json({data: result, error: ""});
            }   
        });
    });

    app.post('/collect', (req, res) => {

        const address = req.body.address;

        contract.methods.collectRewards(address).send({from: owner}, function (error, result) {
            if (error) {
                const transactionKey = Object.keys(error["data"])[0];
                if (!transactionKey && error["data"][transactionKey] !== "revert") {
                    res.status(500).json({data: "", error: "Internal Server Error"})
                } else {
                    res.status(200).json({data: "", error: "Player has no rewards"})
                }
            } else {
                res.status(200).json({data: result, error: ""});
            }   
        });
    });

    app.get('/bet/:address', (req, res) => {

        const address = req.params.address;

        contract.methods.getCurrentBet(address).call({from: owner}, function (error, result) {
            if (error) {
                console.log(error)
                res.status(500).json({data: "", error: "Internal server error"})
            } else {
                res.status(200).json({data: result, error: ""});
            }   
        });
    });


    //Add some authentication to this method
    app.post('/winner', (req, res) => {

        const address = req.body.address;
        console.log(`address => ${address} is a winner`)

        contract.methods.declareWinner(address).send({from: owner}, function (error, result) {
            if (error) {
                const transactionKey = Object.keys(error["data"])[0];
                if (!transactionKey && error["data"][transactionKey]["error"] !== "revert") {
                    res.status(500).json({data: "", error: "Internal Server Error"})
                } else {
                    let reason = error["data"][transactionKey]["reason"]
                    res.status(200).json({data: "", error: reason})
                }
            } else {
                res.status(200).json({data: result, error: ""});
            }   
        });
    });


    app.post('/loser', (req, res) => {

        const address = req.body.address;

        contract.methods.declareLoser(address).send({from: owner}, function (error, result) {
            if (error) {
                const transactionKey = Object.keys(error["data"])[0];
                if (!transactionKey && error["data"][transactionKey]["error"] !== "revert") {
                    res.status(500).json({data: "", error: "Internal Server Error"})
                } else {
                    let reason = error["data"][transactionKey]["reason"]
                    res.status(200).json({data: "", error: reason})
                }
            } else {
                res.status(200).json({data: result, error: ""});
            }   
        });
    });
    
    app.post('/bet', (req, res) => {

        const address = req.body.address;
        const amount = req.body.amount;

        contract.methods.bet(address, amount).send({from: owner}, function (error, result) {
            if (error) {
                res.status(500).json({data: "", error: "Internal Server Error"})
            } else {
                contract.methods.getCurrentBet(address).call({from: owner}, function (error, result){
                    if (error) {
                        res.status(500).json({data: "", error: "Internal Server Error"})
                    } else {
                        res.status(200).json({data: result, error: ""});
                    }
                });
            }   
        });
    });
    
    app.listen(port, () => {
      return console.log(`Express is listening at http://localhost:${port}`);
    });
    
}
