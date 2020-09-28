#!/bin/bash

export CORE_PEER_ADDRESS="localhost:9051"
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_MSPCONFIGPATH="/home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp"
export CORE_PEER_TLS_ENABLED="true"
export CORE_PEER_TLS_ROOTCERT_FILE="/home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
export FABRIC_CFG_PATH="/home/pplavetzki/development/hack/fabric-samples/config"
export ORDERER_CA="/home/pplavetzki/development/hack/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
export PATH="/home/pplavetzki/development/hack/fabric-samples/bin:/home/pplavetzki/development/hack/fabric-samples/test-network:/home/pplavetzki/.nvm/versions/node/v13.9.0/bin:/home/pplavetzki/.cargo/bin:/home/pplavetzki/.local/bin:/home/pplavetzki/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin:/usr/local/go/bin:/home/pplavetzki/go/bin:/home/pplavetzki/.cargo/bin:/home/pplavetzki/bin:/home/pplavetzki/development/hack/fabric-samples/bin"
export PEER0_ORG1_CA="/home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
export PEER0_ORG2_CA="/home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
export PEER0_ORG3_CA="/home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt"
export PEER_PARMS="--peerAddresses localhost:9051 --tlsRootCertFiles /home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/pplavetzki/development/hack/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"