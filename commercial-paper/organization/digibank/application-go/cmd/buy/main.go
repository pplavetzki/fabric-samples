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
	userName := "balaji"
	wallet, err := gateway.NewFileSystemWallet(filepath.Join("..", "..", "..", "identity/user/balaji/wallet"))
	if err != nil {
		logger.Sugar().Fatalf("Failed to obtain wallet: %s", err)
	}

	// Path to the network config (CCP) file
	ccpPath := filepath.Join("..", "..", "..", "gateway/connection-org1.yaml")

	// Connect to the gateway peer(s) using the network config and identity in the wallet
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, userName),
		gateway.WithUser(userName),
	)
	if err != nil {
		logger.Sugar().Fatalf("Failed to connect to gateway: %s\n", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		logger.Sugar().Fatalf("Failed to get network: %s\n", err)
	}
	contract := network.GetContract("papercontract")

	result, err := contract.SubmitTransaction("buy", "MagnetoCorp", "00001", "MagnetoCorp", "DigiBank", "4900000", "2020-05-31")
	if err != nil {
		logger.Sugar().Fatalf("Failed to submit transaction: %s\n", err)
	}

	var paper commercialpaper.CommercialPaper
	err = paper.UnmarshalJSON(result)
	if err != nil {
		logger.Sugar().Fatalf("failed to unmarshal json result: %s", err)
	}
	logger.Sugar().Infof("%s commercial paper: %s successfully purchased by %s\n", paper.Issuer, paper.PaperNumber, paper.Owner)
	logger.Sugar().Info("Transaction complete.")
}
