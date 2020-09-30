/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	commercialpaper "github.com/hyperledger/fabric-samples/commercial-paper/organization/magnetocorp/contract-go/commercial-paper"
	"go.uber.org/zap"
)

var (
	chaincodeID      string
	chaincodeAddress string
	logger           *zap.Logger
)

func init() {
	logger, _ = zap.NewDevelopment()
	chaincodeID = os.Getenv("CHAINCODE_CCID")
	if chaincodeID == "" {
		logger.Sugar().Fatal("CHAINCODE_CCID must not be enpty")
	}
	chaincodeAddress = os.Getenv("CHAINCODE_ADDRESS")
	if chaincodeAddress == "" {
		logger.Sugar().Fatal("CHAINCODE_ADDRESS must not be enpty")
	}
}

func main() {

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
}
