package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	commercialpaper "github.com/hyperledger/fabric-samples/commercial-paper/organization/magnetocorp/contract-go/commercial-paper"
	event "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	provmsp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var logger *zap.Logger

func populateWallet(wallet *gateway.Wallet) error {
	credPath := viper.GetString("identityCredentials.identityCredPath")
	signCertName := viper.GetString("identityCredentials.signCertName")
	identityOrgMSP := viper.GetString("identityCredentials.identityOrgMSP")
	identityUserName := viper.GetString("identityCredentials.identityUserName")

	certPath := filepath.Join(credPath, "signcerts", signCertName)
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(identityOrgMSP, string(cert), string(key))

	err = wallet.Put(identityUserName, identity)
	if err != nil {
		return err
	}
	return nil
}

func getWallet() (*gateway.Wallet, error) {
	walletId := viper.GetString("identityCredentials.walletPath")
	return gateway.NewFileSystemWallet(walletId)
}

func getGatewayConnection() (*gateway.Gateway, error) {
	ccpPath := viper.GetString("connectionConfigPath")
	userName := viper.GetString("identityCredentials.identityUserName")

	if ccpPath == "" {
		return nil, fmt.Errorf("connectionConfigPath is required")
	}
	wallet, err := getWallet()
	if err != nil {
		return nil, err
	}

	//Connect to the gateway peer(s) using the network config and identity in the wallet
	return gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, userName),
	)
}

func startClientEventListening() {
	userName := viper.GetString("identityCredentials.identityUserName")

	contract := viper.GetString("chaincode.contract")
	logger.Sugar().Debug("chaincode.contract: ", contract)
	channelName := viper.GetString("chaincode.channel")
	logger.Sugar().Debug("chaincode.channel: ", channelName)
	eventFilter := viper.GetString("chaincode.eventFilter")
	logger.Sugar().Debug("chaincode.eventFilter: ", eventFilter)

	// This is merging relative paths
	ccpPath := viper.GetString("connectionConfigPath")
	conf := config.FromFile(filepath.Clean(ccpPath))
	// Instantiate the SDK
	sdk, err := fabsdk.New(conf)
	if err != nil {
		logger.Sugar().Fatalf("SDK init failed: %s", err)
	}
	defer sdk.Close()

	ctx := sdk.Context()
	config, err := sdk.Config()
	if err != nil {
		logger.Sugar().Fatalf("SDK config failed: %s", err)
	}
	caInstance, ok := config.Lookup("certificateAuthorities")
	if !ok {
		logger.Sugar().Fatalf("cert authorities error: %s", err)
	}
	caIns, _ := caInstance.(map[string]interface{})
	var caVal string
	for k := range caIns {
		caVal = k
		break
	}

	client, ok := config.Lookup("client")
	if !ok {
		logger.Sugar().Fatalf("failed client: %s", err)
	}
	clientBackend, _ := client.(map[string]interface{})
	orgBackend := clientBackend["organization"]

	mspClient, err := msp.New(ctx, msp.WithCAInstance(caVal))
	if err != nil {
		logger.Sugar().Fatalf("failed mspClient %s", err)
	}
	wallet, err := getWallet()
	ident, err := wallet.Get(userName)
	if err != nil {
		logger.Sugar().Fatalf("failed wallet: %s", err)
	}

	cert := []byte(ident.(*gateway.X509Identity).Certificate())
	privKey := []byte(ident.(*gateway.X509Identity).Key())

	signIdent, err := mspClient.CreateSigningIdentity(provmsp.WithCert(cert), provmsp.WithPrivateKey(privKey))
	if err != nil {
		logger.Sugar().Fatalf("failed sign identity %s", err)
	}

	chanProvider := sdk.ChannelContext(channelName, fabsdk.WithOrg(strings.ToLower(strings.ToLower(fmt.Sprintf("%s", orgBackend)))), fabsdk.WithIdentity(signIdent))
	if chanProvider == nil {
		logger.Sugar().Fatalf("failed to create channel provider: %s", err)
	}

	myclient, err := event.New(chanProvider, event.WithBlockEvents())
	if err != nil {
		logger.Sugar().Fatalf("failed to create client provider: %s", err)
	}
	reg, ccChan, err := myclient.RegisterChaincodeEvent(contract, eventFilter)
	if err != nil {
		logger.Sugar().Fatalf("error registering chaincode event: %s", err)
	}
	defer myclient.Unregister(reg)

	for {
		select {
		case event := <-ccChan:
			log.Println("************************************************************************************************")
			log.Println("From the startClientEventListening channel *****************************************")
			log.Println("************************************************************************************************")
			log.Println("transaction id:", event.TxID, "chaincode id:", event.ChaincodeID, "event name:", event.EventName)
			log.Println("Payload received:", string(event.Payload))
			log.Println("************************************************************************************************")
			var paper commercialpaper.CommercialPaper
			err = paper.UnmarshalJSON(event.Payload)
			if err != nil {
				log.Printf("failed to unmarshal: %s\n", err)
				log.Println(string(event.Payload))
			} else {
				log.Println("************************************************************************************************")
				log.Println("From the startClientEventListening channel -- from the payload  ****************************************")
				log.Println("************************************************************************************************")
				log.Printf("Commercial paper issuer: %s, paper number: %s, issued for face value %d in state %s\n", paper.Issuer, paper.PaperNumber, paper.FaceValue, paper.GetState().String())
				log.Println("************************************************************************************************")
			}
		}
	}
}

func contractListening() {
	contract := viper.GetString("chaincode.contract")
	channel := viper.GetString("chaincode.channel")
	eventFilter := viper.GetString("chaincode.eventFilter")

	//Connect to the gateway peer(s) using the network config and identity in the wallet
	gw, err := getGatewayConnection()
	if err != nil {
		logger.Sugar().Fatalf("gateway error: %s\n", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channel)
	if err != nil {
		logger.Sugar().Fatalf("Failed to get network: %s\n", err)
	}
	paperContract := network.GetContract(contract)
	reg, ccChan, err := paperContract.RegisterEvent(eventFilter)
	if err != nil {
		logger.Sugar().Fatalf("Failed to get register event: %s\n", err)
	}
	log.Printf("listening for events from contract %s with filtered events: %s", contract, eventFilter)
	defer paperContract.Unregister(reg)
	for {
		select {
		case event := <-ccChan:
			log.Println("************************************************************************************************")
			log.Println("From the contractListening channel *************************************************************")
			log.Println("************************************************************************************************")
			log.Println("transaction id:", event.TxID, "chaincode id:", event.ChaincodeID, "event name:", event.EventName)
			log.Println("************************************************************************************************")
			var paper commercialpaper.CommercialPaper
			err = paper.UnmarshalJSON(event.Payload)
			if err != nil {
				log.Printf("failed to unmarshal: %s\n", err)
				log.Println(string(event.Payload))
			} else {
				log.Println("************************************************************************************************")
				log.Println("From the contractListening channel -- from the payload  ****************************************")
				log.Println("************************************************************************************************")
				log.Printf("%s commercial paper : %s issued for value %d in state %s\n", paper.Issuer, paper.PaperNumber, paper.FaceValue, paper.GetState().String())
				log.Println("************************************************************************************************")
			}
		}
	}
}

func networkListening() {
	// Path to the network config (CCP) file
	channel := viper.GetString("chaincode.channel")
	//Connect to the gateway peer(s) using the network config and identity in the wallet
	gw, err := getGatewayConnection()
	if err != nil {
		logger.Sugar().Fatalf("Failed to connect to gateway: %s\n", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channel)
	if err != nil {
		logger.Sugar().Fatalf("Failed to get network: %s\n", err)
	}
	reg, blkChan, err := network.RegisterBlockEvent()
	if err != nil {
		logger.Sugar().Fatalf("registering block event error")
	}
	defer network.Unregister(reg)
	for {
		select {
		case event := <-blkChan:
			blk := event.Block.GetData()
			log.Println("************************************************************************************************")
			log.Println("From the RegisterBlockEvent from the network *****************************************************")
			log.Println("************************************************************************************************")
			log.Println(blk)
			log.Println("************************************************************************************************")
		}
	}
}

func init() {
	logger, _ = zap.NewDevelopment()
	logger.Sugar().Info("inside the init function")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}
	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "config"
	}
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logger.Sugar().Fatalf("Fatal error config file: %s \n", err)
		return
	}
}

func main() {
	go startClientEventListening()
	go networkListening()
	go contractListening()

	logger.Sugar().Info("listening on all cylinders")

	select {}
}
