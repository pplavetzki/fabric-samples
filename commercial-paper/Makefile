
GITHUB_RUN_NUMBER ?= 1

ACR    ?= pplavetzki
ACR_NAME ?= ${ACR}.azurecr.io/hyperledger/hack
TAG    := $(shell date +"%Y%m%d.${GITHUB_RUN_NUMBER}")

ORG_NAME ?= org2

BUILD_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)
PAPER_DIR = $(BUILD_DIR)/commercial-paper
ORG2_DIR = $(PAPER_DIR)/organization/magnetocorp

.PHONY: build-watcher-image
build-watcher-image:
	docker build -t ${ACR_NAME}/watcher:${TAG} -f $(ORG2_DIR)/watcher/watcher.dockerfile $(ORG2_DIR)/watcher/.
	docker tag ${ACR_NAME}/watcher:${TAG} ${ACR_NAME}/watcher:latest

.PHONY: run-watcher
run-watcher: build-watcher-image
	# docker rm -f watcher.org2.example.com
	docker run -d --name watcher.org2.example.com -h watcher.org2.example.com --network net_test -v $(ORG2_DIR)/gateway/dconnection-org2.yaml:/tmp/gateway/connection-org2.yaml \
		-v $(ORG2_DIR)/watcher/config/:/tmp/config -v $(ORG2_DIR)/identity/user/isabella/wallet/:/tmp/credentials/isabella/wallet -e CONFIG_PATH=/tmp/config -e CONFIG_NAME=dconfig ${ACR_NAME}/watcher:${TAG}

.PHONY: install-go-cc-org2
install-go-cc-org2:
	./install-go-cc.sh org2 $(PAPER_DIR)/organization/magnetocorp/magnetocorp.sh $(PAPER_DIR)/build $(PAPER_DIR)/bin 8052 7052 0 1

.PHONY: install-go-cc-org1
install-go-cc-org1:
	./install-go-cc.sh org1 $(PAPER_DIR)/organization/digibank/digibank.sh $(PAPER_DIR)/build $(PAPER_DIR)/bin 8052 7052 0 1

.PHONY: install-go-cc
install-go-cc: install-go-cc-org2 install-go-cc-org1

.PHONY: clean-all
clean-all:
	docker container rm $$(docker container list -f network=net_test -aq) -f
	./network-clean.sh

.PHONY: network-start
network-start:
	./network-starter.sh

.PHONY: deploy
deploy: network-start install-go-cc

.PHONY: enroll-org1
enroll-org1:
	cd $(PAPER_DIR)/organization/digibank/application-go/cmd/enroll && go run .

.PHONY: enroll-org2
enroll-org2:
	cd $(PAPER_DIR)/organization/magnetocorp/application-go/cmd/enroll && go run .

.PHONY: enroll
enroll: enroll-org2 enroll-org1

.PHONY: issue
issue:
	cd $(PAPER_DIR)/organization/magnetocorp/application-go/cmd/issue && go run .

.PHONY: buy
buy:
	cd $(PAPER_DIR)/organization/digibank/application-go/cmd/buy && go run .

.PHONY: redeem
redeem:
	cd $(PAPER_DIR)/organization/digibank/application-go/cmd/redeem && go run .