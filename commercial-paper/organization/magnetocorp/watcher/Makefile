
GITHUB_RUN_NUMBER ?= 1

ACR    := pplavetzki
ACR_NAME = ${ACR}.azurecr.io/hyperledger/hack
TAG    := $(shell date +"%Y%m%d.${GITHUB_RUN_NUMBER}")

ORG_NAME ?= org2

BUILD_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)

.PHONY: build-watcher-image
build-watcher-image:
	docker build -t ${ACR_NAME}/watcher:${TAG} -f watcher.dockerfile .
	docker tag ${ACR_NAME}/watcher:${TAG} ${ACR_NAME}/watcher:latest

run-watcher-image: build-watcher-image
	# docker rm -f watcher.org2.example.com
	docker run -d --name watcher.org2.example.com -h watcher.org2.example.com --network net_test -v $(BUILD_DIR)/gateway/dconnection-org2.yaml:/tmp/gateway/connection-org2.yaml \
		-v $(PWD)/config/:/tmp/config -v $(BUILD_DIR)/identity/user/isabella/wallet/:/tmp/credentials/isabella/wallet -e CONFIG_PATH=/tmp/config -e CONFIG_NAME=dconfig ${ACR_NAME}/watcher:${TAG}
