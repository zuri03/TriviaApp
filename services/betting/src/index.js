const fs = require("fs-extra");
const path = require("path");
const solc = require("solc");
//const { stringify } = require('flatted');
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.WebsocketProvider('ws://blockchain:8545'));
const server = require('./server/server.js');

const compile = function () {
    const contract = path.resolve(__dirname,"../contracts","Pot.sol");
    const source = fs.readFileSync(contract, "utf8");

    var input = {
        language: 'Solidity',
        sources: {
            'Pot.sol': { content: source }
        },
        settings: {
            outputSelection: {
                '*': { '*': ['*'] }
            }
        }
    };

    const output = JSON.parse(solc.compile(JSON.stringify(input)));
    return output.contracts['Pot.sol']['Pot'];
};

const deploy = async (owner, obj) => {
    let contract = await new web3.eth.Contract(obj.abi);
    console.log(`initial address => ${contract.options.address}`)

    let result = await contract.deploy({data: obj.evm.bytecode.object})
        .send({gas: 6721975, from: owner});
    contract.options.address = result['_address'];
    return contract;
}

const getOwnerAddress = async function () {
    var address;
    const getAccount = function () {
        return new Promise((resolve, reject) => {
            web3.eth.getAccounts((err, accounts) => {
                err === null ? resolve(accounts[0]) : reject(err)
            });
        });
    };

    //We have to wait for the blockchain service to start so if it does not work the first time try again
    let attempts = 0;
    while (attempts < 100) {
        try {
            address = await getAccount(); 
            break;
        } catch(error) {
            attempts++;
            await new Promise(r => setTimeout(r, 2000));
        } 
    }
    return address;
};

(async function () {
    const owner = await getOwnerAddress();
    console.log('got owner');
    const compiledContract = compile();
    console.log('finished compilation')
    const contract = await deploy(owner, compiledContract);
    console.log('finished deploy')
    console.log('starting server')
    server(owner, web3, contract)
})();