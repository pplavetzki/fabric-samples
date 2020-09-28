package main

import (
	"path/filepath"

	commercialpaper "github.com/hyperledger/fabric-samples/commercial-paper/organization/magnetocorp/contract-go/commercial-paper"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func main() {
	wallet, err := gateway.NewFileSystemWallet(filepath.Join("..", "..", "..", "identity/user/isabella/wallet"))
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// Path to the network config (CCP) file
	ccpPath := filepath.Join("..", "..", "..", "gateway/connection-org2.yaml")

	// Connect to the gateway peer(s) using the network config and identity in the wallet
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "isabella"),
		gateway.WithUser("isabella"),
	)
	if err != nil {
		logger.Sugar().Fatal("Failed to connect to gateway: %s\n", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		logger.Sugar().Fatalf("Failed to get network: %s\n", err)
	}
	contract := network.GetContract("papercontract")

	result, err := contract.SubmitTransaction("issue", "MagnetoCorp", "00001", "2020-05-31", "2020-11-30", "1000000")
	if err != nil {
		logger.Sugar().Fatalf("Failed to submit transaction: %s\n", err)
	}

	var paper commercialpaper.CommercialPaper
	err = paper.UnmarshalJSON(result)
	if err != nil {
		logger.Sugar().Fatalf("failed to unmarshal json result: %s", err)
	}
	logger.Sugar().Infof("%s commercial paper : %s successfully issued for value %d\n", paper.Issuer, paper.PaperNumber, paper.FaceValue)
	logger.Sugar().Debug("Transaction complete.")
}
