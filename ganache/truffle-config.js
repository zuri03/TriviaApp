module.exports = {
    // See <http://truffleframework.com/docs/advanced/configuration>
    // for more about customizing your Truffle configuration!
    networks: {
      development: {
        host: "blockchain",
        port: 8545,
        network_id: "*",// Match any network id,
        webSockets: true
      },
      develop: {
        port: 8545
      },
      compilers: {
        solc: {
          version: "^0.8.11",    // Fetch exact version from solc-bin (default: truffle's version)
          // docker: true,        // Use "0.5.1" you've installed locally with docker (default: false)
          // settings: {          // See the solidity docs for advice about optimization and evmVersion
          //  optimizer: {
          //    enabled: false,
          //    runs: 200
          //  },
          //  evmVersion: "byzantium"
          // }
        }
      }
    }
  }
  
  