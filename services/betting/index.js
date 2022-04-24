const fs = require("fs-extra");
const path = require("path");
const solc = require("solc");
const { stringify } = require('flatted');
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.WebsocketProvider('ws://blockchain:8545'));
//const server = require('./server/server.js');

const compile = function () {
    const buildFolder = path.resolve(__dirname,"./build");
    fs.removeSync(buildFolder);

    const contract = path.resolve(__dirname,"./contracts","Pot.sol");
    const source = fs.readFileSync(contract, "utf8");

    var input = {
        language: 'Solidity',
        sources: {
            'Pot.sol': {
                content: source
            }
        },
        settings: {
            outputSelection: {
                '*': {
                    '*': ['*']
                }
            }
        }
    };

    const output = JSON.parse(solc.compile(JSON.stringify(input)));
    fs.ensureDirSync(buildFolder);
    fs.outputJSONSync(path.resolve(buildFolder, "Pot"+".json"), output.contracts['Pot.sol']['Pot']);
};

const deploy = async (owner) => {
    try {
        const json = fs.readFileSync(path.resolve(__dirname,"./build","Pot.json"), "utf8");
        let obj = JSON.parse(json);
        //console.log(JSON.stringify(obj.evm.bytecode.object))
        /*
        console.log(JSON.stringify(obj))
        console.log("============================================================================ evm ==============================================================================================")
        console.log(JSON.stringify(obj.evm))
        console.log("============================================================================ bytcode ==============================================================================================")
        console.log(JSON.stringify(obj.evm.bytecode))
        console.log("============================================================================ object ==============================================================================================")
        console.log(JSON.stringify(obj.evm.bytecode.object))
        */
        const result = await new web3.eth.Contract(obj.abi)
            .deploy({data: obj.evm.bytecode.object})
            .send({gas: 6721975, from: owner});
        const serialised = stringify(result.options);
    
        return await serialised;
    } catch (error) {
        console.error(error);
        return error;
    }
}

var getOwnerAddress = async function () {
    var address;
    const getAccount = function () {
        return new Promise((resolve, reject) => {
            web3.eth.getAccounts((err, accounts) => {
                if (err === null) {
                    resolve(accounts[0]);
                } else {
                    reject(err);
                }
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
    console.log(`found owner => ${address}`)
    return address;
};

(async function () {
    console.log('starting owner')
    const owner = await getOwnerAddress();
    console.log('got owner');
    compile();
    console.log('finished compilation')
    await deploy(owner);
    console.log('finished deploy')
    console.log('starting server')
    //server(owner, web3)
})();