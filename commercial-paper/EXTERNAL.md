# Commercial Paper Contract Using External Builders

This is a fork of the fabric-samples.

## TL;DR

1. `cd fabric-samples/commercial-paper`
2. `make deploy`
3. `make enroll`
4. `make run-watcher`
5. `make issue`
6. `make buy`
7. `make redeem`

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

It is advised to have 2 console windows open; one to monitor the infrastructure and one run the make recipes. Once you've cloned the fabric-samples - change to the commercial-paper directory in each window.

```
cd fabric-samples/commercial-paper
```

## Running the Infrastructure

### Deploy the network and smart contract to the channel

This is where we diverge from the standard deployment of chaincode.  Although you can still deploy the compiled/packaged chaincode to the peer.  The purpose of this sample is to deploy using external builders and launchers.

I'll only summarize the differences:

1. I added a builders directory to the peer containers that's responsible for building and executing the external chaincode.  This code is in the `fabric-samples/builders` directory.  For more info on what this looks like...[builders and launchers](https://hyperledger-fabric.readthedocs).
2. External chaincode can run as containers (services), so that's what we're going to do.  There is a make recipe to install the chaincode using the previous step.  This will package/deploy/approve/commit the chaincode to both MagnetoCorp and Digibank.  The recipe also includes the building and running of the chaincode as services in our deployed network.  It creates two services one each for MagnetoCorp and Digibank.

From the `fabric-samples/commercial-paper` directory run the make recipe for creating the network, installing the chaincode and running the services:

```make deploy```

You will see a lot of commands flow through the terminal window. It should finish with the following couple of lines something similar to this (txid will be different):

```
Committing the chaincode
2020-09-30 16:28:46.701 CDT [chaincodeCmd] ClientWait -> INFO 001 txid [ff82213c065f4573834f1be645f6f006ed6a09a7c68007f9b67e3d8c67a6e34a] committed with status (VALID) at localhost:7051
2020-09-30 16:28:46.719 CDT [chaincodeCmd] ClientWait -> INFO 002 txid [ff82213c065f4573834f1be645f6f006ed6a09a7c68007f9b67e3d8c67a6e34a] committed with status (VALID) at localhost:9051
```

After this is finished you can verified that the chaincode services are running: ```docker ps -f name=paperchain``` you should see something similar to the following:

```
CONTAINER ID        IMAGE                                                           COMMAND             CREATED             STATUS              PORTS                    NAMES
97fd5e9421fe        contract-go:20200930.1   "/contract-cc"      2 minutes ago       Up 2 minutes        0.0.0.0:8052->8052/tcp   paperchain-org1.example.com
1322b5eaa81b        contract-go:20200930.1   "/contract-cc"      2 minutes ago       Up 2 minutes        0.0.0.0:7052->7052/tcp   paperchain-org2.example.com
```

We are now ready to perform some administration for user access by creating some wallets.  Keep in mind that we are imitating that we have two organizations: MagnetoCorp (org2) and Digibank (org1).  We need to enroll our org users.  We will create the user `isabella` in MagnetoCorp and the user `balaji` in Digibank.  Run the following make recipe to enroll both users in their respective organization:

```
make enroll
```

This example makes use of the file system wallet so after running the enroll recipe, you will see that you have an id file at `identity/user/[USER_NAME]/wallet/[USER_NAME].id` in the `magnetocorp` and in the `digibank` organization.

Now that the users are enrolled, we can **optionally** deploy the watcher application to view transactions on the ledger.  This application needs an identity to access the ledger, so we will use the `isabella` identity created in MagnetoCorp in the previous step.  Keep in mind you would deploy a 'watcher' app in each organization, but since we are using a test docker network we'll only deploy this once.

Run this recipe to build and deploy the watcher application to the network:

```
make run-watcher
```

You can verify that the application is running by using this command: `docker ps -f name=watcher` you should see something similar to the following:

```
1682        watcher:20200930.1   "/watcher"          16 seconds ago      Up 14 seconds                           watcher.org2.example.com
```

Now we can use our application to follow the commercial paper [workflow](https://hyperledger-fabric.readthedocs.io/en/release-2.2/tutorial/commercial_paper.html#issue-application)

```
make issue   # this is run as MagnetoCorp
make buy     # this is run as Digibank
make redeem  # this is run as Digibank
```

If you started the watcher application you can view the logs to see the various transactions that hit the magnetocorp ledger:

```
docker logs watcher.org2.example.com
```

## Cleaning Up

Run the following recipe from the `fabric-samples/commercial-paper` directory:

```
make clean-all
```