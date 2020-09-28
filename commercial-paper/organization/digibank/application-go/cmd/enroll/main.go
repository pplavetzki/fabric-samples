package main

import (
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func main() {
	conf := config.FromFile(filepath.Join("..", "..", "..", "gateway/connection-org1.yaml"))
	// Instantiate the SDK
	sdk, err := fabsdk.New(conf)
	if err != nil {
		logger.Sugar().Fatalf("SDK init failed: %s", err)
	}
	defer sdk.Close()

	// gets the sdk from a contrived connection path
	ctx := sdk.Context()
	if err != nil {
		logger.Sugar().Fatalf("SDK Context failed: %s", err)
	}

	mspClient, err := msp.New(ctx, msp.WithCAInstance("ca.org1.example.com"))
	if err != nil {
		logger.Sugar().Fatalf("Could not create MSP Client: %s", err)
	}

	wallet, err := gateway.NewFileSystemWallet(filepath.Join("..", "..", "..", "identity/user/balaji/wallet"))
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	if err := mspClient.Enroll("user1", msp.WithSecret("user1pw")); err != nil {
		logger.Sugar().Fatal(err)
	}

	signIdent, err := mspClient.GetSigningIdentity("user1")
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	pk, err := signIdent.PrivateKey().Bytes()
	if err != nil {
		logger.Sugar().Fatal(err)
	}
	identity := gateway.NewX509Identity("Org1MSP", string(signIdent.EnrollmentCertificate()), string(pk))

	err = wallet.Put("balaji", identity)
	if err != nil {
		logger.Sugar().Error(err)
	}
}
