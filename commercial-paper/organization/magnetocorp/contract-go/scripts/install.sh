#!/bin/bash

ORG_NAME=$1
ORG_ENV=$2
BUILD_DIR=$3
OUTPUT_DIR=$4
CHAINCODE_CONTAINER_ORG1_PUBLISH_PORT=$5
CHAINCODE_CONTAINER_ORG2_PUBLISH_PORT=$6
CC_VER=$7
CC_SEQ=$8

GITHUB_RUN_NUMBER=1

CONTRACT_NAME=paperchain
PACKAGE_NAME=papercontract
PACKAGE_LABEL=$PACKAGE_NAME"_"$CC_VER"_"$CC_SEQ
CONTAINER_NAME=$CONTRACT_NAME-$ORG_NAME.example.com
CHAINCODE_SERVER_ADDRESS=$CONTRACT_NAME-$ORG_NAME.example.com:7052
CHAINCODE_CONTAINER_ADDRESS=0.0.0.0:7052

ACR=pplavetzki
ACR_NAME=$ACR.azurecr.io/hyperledger/hack
TAG=$(date +"%Y%m%d.$GITHUB_RUN_NUMBER")

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done

SCRIPTS_DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
ROOT_DIR="$(dirname "$SCRIPTS_DIR")"
PACKAGE_DIR=$ROOT_DIR/packaging

source $ORG_ENV

if [ $ORG_NAME == "org1" ]
then
    CHAINCODE_CONTAINER_PORT=$5
fi

if [ $ORG_NAME == "org2" ]
then
    CHAINCODE_CONTAINER_PORT=$6
fi

if [ ! -d $BUILD_DIR ]
then
    echo 'Creating build directory:' $BUILD_DIR
    mkdir -p $BUILD_DIR
fi

if [ ! -d $OUTPUT_DIR ]
then
    echo 'Creating output directory:' $OUTPUT_DIR
    mkdir -p $OUTPUT_DIR
fi

query_cc() {
    COMMAND="peer lifecycle chaincode queryinstalled -O json | jq -r '.installed_chaincodes[] | select(.label == \""$PACKAGE_LABEL"\") | .package_id'"
    eval $COMMAND
}

echo "Package the chaincode"
echo "1. Copy Packaging to Build Directory"
cp -R $PACKAGE_DIR/* $BUILD_DIR
cp $BUILD_DIR/metadata.tpl $BUILD_DIR/metadata.json
cp $BUILD_DIR/connection.tpl $BUILD_DIR/connection.json
echo "Create the connection.json file from the template and replace the CHAINCODE_SERVER_ADDRESS with:" $CHAINCODE_SERVER_ADDRESS
sed -i "s/REPLACE_CHAINCODE_SERVER_ADDRESS/$CHAINCODE_SERVER_ADDRESS/g" $BUILD_DIR/connection.json
echo "Create the metadata.json file from the template and replace the PACKAGE_LABEL with:" $PACKAGE_LABEL
sed -i "s/REPLACE_LABEL/$PACKAGE_LABEL/g" $BUILD_DIR/metadata.json
echo "Step 2: tar connection.json"
cd $BUILD_DIR  
tar cfz code.tar.gz connection.json
echo "Step 3: tar meta.json and step 1 together"
cd $BUILD_DIR 
tar cfz paperchain.tgz metadata.json code.tar.gz
cp $BUILD_DIR/paperchain.tgz $OUTPUT_DIR
echo "Packaged chaincode paperchain at" $OUTPUT_DIR/paperchain.tgz
# rm -rf $BUILD_DIR/*
echo "--------------------------------------------------------------------------------------------------------------------------------"
echo "Build and Deploy image"
echo "--------------------------------------------------------------------------------------------------------------------------------"
peer lifecycle chaincode install $OUTPUT_DIR/paperchain.tgz
echo "--------------------------------------------------------------------------------------------------------------------------------"
PACKAGE_ID=$(query_cc)
echo "package id:" $PACKAGE_ID
echo "--------------------------------------------------------------------------------------------------------------------------------"
echo "Building the docker image"
docker build -t $ACR_NAME/contract-go:$TAG -f $ROOT_DIR/Dockerfile $ROOT_DIR/.
docker tag $ACR_NAME/contract-go:$TAG $ACR_NAME/contract-go:latest
echo "Removing old container if necessary"
docker rm -f $CONTAINER_NAME
echo "--------------------------------------------------------------------------------------------------------------------------------"
echo "Running the paperchain chaincode container"
docker run -d --name $CONTAINER_NAME --network net_test -p $CHAINCODE_CONTAINER_PORT:$CHAINCODE_CONTAINER_PORT -e CHAINCODE_ADDRESS=$CHAINCODE_CONTAINER_ADDRESS \
	-e CHAINCODE_CCID=$PACKAGE_ID -h $CONTAINER_NAME $ACR_NAME/contract-go:$TAG
echo "--------------------------------------------------------------------------------------------------------------------------------"
echo "Approving the chaincode"
peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name papercontract -v $CC_VER --package-id $PACKAGE_ID --sequence $CC_SEQ --tls --cafile $ORDERER_CA
echo "--------------------------------------------------------------------------------------------------------------------------------"
if [ $ORG_NAME == "org1" ]
then
    echo "Committing the chaincode"
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} --channelID mychannel --name papercontract -v $CC_VER --sequence $CC_SEQ --tls --cafile $ORDERER_CA --waitForEvent
fi