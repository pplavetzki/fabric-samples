/*
 * SPDX-License-Identifier: Apache-2.0
 */

package commercialpaper

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing commercial
// paper
type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

// Query returns the paper for the given issuer and paper number
func (c *Contract) Query(ctx TransactionContextInterface, issuer string, paperNumber string) (*CommercialPaper, error) {
	return ctx.GetPaperList().GetPaper(issuer, paperNumber)
}

// QueryByOwner gets all records by owner
func (c *Contract) QueryByOwner(ctx TransactionContextInterface, owner string) ([]*CommercialPaper, error) {
	query := fmt.Sprintf("{\"selector\":{\"owner\":\"%s\"}}", owner)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var papers []*CommercialPaper
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var paper CommercialPaper
		err = json.Unmarshal(queryResponse.Value, &paper)
		if err != nil {
			return nil, err
		}
		papers = append(papers, &paper)
	}

	return papers, nil
}

// Issue creates a new commercial paper and stores it in the world state
func (c *Contract) Issue(ctx TransactionContextInterface, issuer string, paperNumber string, issueDateTime string, maturityDateTime string, faceValue int) (*CommercialPaper, error) {
	paper := CommercialPaper{PaperNumber: paperNumber, Issuer: issuer, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: issuer}
	paper.SetIssued()

	err := ctx.GetPaperList().AddPaper(&paper)
	if err != nil {
		return nil, err
	}
	payload, err := paper.MarshalJSON()
	if err != nil {
		fmt.Println("Failed to marshal payload of paper", err)
	} else {
		stub := ctx.GetStub()
		stub.SetEvent("issue", payload)
	}

	return &paper, nil
}

// Buy updates a commercial paper to be in trading status and sets the new owner
func (c *Contract) Buy(ctx TransactionContextInterface, issuer string, paperNumber string, currentOwner string, newOwner string, price int, purchaseDateTime string) (*CommercialPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(issuer, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.Owner != currentOwner {
		return nil, fmt.Errorf("Paper %s:%s is not owned by %s", issuer, paperNumber, currentOwner)
	}

	if paper.IsIssued() {
		paper.SetTrading()
	}

	if !paper.IsTrading() {
		return nil, fmt.Errorf("Paper %s:%s is not trading. Current state = %s", issuer, paperNumber, paper.GetState())
	}

	paper.Owner = newOwner

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	payload, err := paper.MarshalJSON()
	if err != nil {
		fmt.Println("Failed to marshal payload of paper", err)
	} else {
		stub := ctx.GetStub()
		stub.SetEvent("buy", payload)
	}

	return paper, nil
}

// Redeem updates a commercial paper status to be redeemed
func (c *Contract) Redeem(ctx TransactionContextInterface, issuer string, paperNumber string, redeemingOwner string, redeenDateTime string) (*CommercialPaper, error) {
	paper, err := ctx.GetPaperList().GetPaper(issuer, paperNumber)

	if err != nil {
		return nil, err
	}

	if paper.Owner != redeemingOwner {
		return nil, fmt.Errorf("Paper %s:%s is not owned by %s", issuer, paperNumber, redeemingOwner)
	}

	if paper.IsRedeemed() {
		return nil, fmt.Errorf("Paper %s:%s is already redeemed", issuer, paperNumber)
	}

	paper.Owner = paper.Issuer
	paper.SetRedeemed()

	err = ctx.GetPaperList().UpdatePaper(paper)

	if err != nil {
		return nil, err
	}
	payload, err := paper.MarshalJSON()
	if err != nil {
		fmt.Println("Failed to marshal payload of paper", err)
	} else {
		stub := ctx.GetStub()
		stub.SetEvent("redeem", payload)
	}

	return paper, nil
}
