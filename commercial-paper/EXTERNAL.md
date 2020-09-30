# Commercial Paper Contract Using External Builders

This is a fork of the fabric-samples.

## External Builders

This is to show an example of using HLF's new > v2.0 features of external [builders and launchers](https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_launcher.html) to [build, deploy and launch](https://hyperledger-fabric.readthedocs.io/en/release-2.2/cc_service.html) chaincode.

Currently, the chaincode as an external service model is only supported by GO chaincode shim. In Fabric v2.0, the GO shim API adds a ChaincodeServer type that developers should use to create a chaincode server.

```
contract := new(commercialpaper.Contract)
contract.TransactionContextHandler = new(commercialpaper.TransactionContext)
contract.Name = "papercontract"
contract.Info.Version = "0.0.1"

chaincode, err := contractapi.NewChaincode(contract)
if err != nil {
    logger.Panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
}
chaincode.Info.Title = "CommercialPaperChaincode"
chaincode.Info.Version = "0.0.1"

server := &shim.ChaincodeServer{
    CCID:    chaincodeID,
    Address: chaincodeAddress,
    CC:      chaincode,
    TLSProps: shim.TLSProperties{
        Disabled: true,
    },
}

err = server.Start()

if err := server.Start(); err != nil {
    logger.Panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
}
```

## Watcher

There is also a ```Watcher``` element to the example.  You can deploy a container to view the transactions that are associated to the network/channel/contract.

## Setup

You will need a machine with the following

- Docker and docker-compose installed
- GO installed https://golang.org/doc/install

You will need to install the peer cli binaries and this fabric-samples repository available. For more information
[Install the Samples, Binaries and Docker Images](https://hyperledger-fabric.readthedocs.io/en/latest/install.html) in the Hyperledger Fabric documentation.

It is advised to have 3 console windows open; one to monitor the infrastructure and one each for MagnetoCorp and DigiBank. Once you've cloned the fabric-samples - change to the commercial-paper directory in each window.

```
cd fabric-samples/commercial-paper
```

## Running the Infrastructure

In one console window, run the `./network-starter.sh` script; this will start the basic infrastructure.

You can re-use this console window if you wish, but it is recommended to run a docker container monitoring script. This will let you view what Fabric is doing and help diagnose any failures.

### Deploy the smart contract to the channel

This is where we diverge from the standard deployment of chaincode.  Although you can still deploy the compiled/packaged chaincode to the peer.  The purpose of this sample is to deploy using external builders and launchers.

I'll only summarize the differences:

1. I added a builders directory to the peer containers that's responsible for building and executing the external chaincode.  This code is in the `fabric-samples/builders` directory.  For more info on what this looks like...[builders and launchers](https://hyperledger-fabric.readthedocs).
2. 