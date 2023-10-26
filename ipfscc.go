package main

import (
	"fmt"
	"encoding/json"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Record struct {
	HashValue string `json:"HashValue"`
	IpfsNodeId string `json:"IpfsNodeId"`
	Operation string `json:"Operation"`
}

// InitLedger 
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	records := []Record{
	}
	for _, record := range records {
		recordJSON, err := json.Marshal(record)
		if err != nil{
			return err
		}

		err = ctx.GetStub().PutState(record.HashValue, recordJSON)
		if err != nil{
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// Create record
func (s *SmartContract) CreateRecord(ctx contractapi.TransactionContextInterface, hashvalue string, ipfsnodeid string, operation string) error{
	exists, err := s.RecordExists(ctx, HashValue)
	if err != nil{
		return err
	}
	if exists {
		return fmt.Errorf("the hashvalue of record %s already exists", HashValue)
	}
	record := Record{
		HashValue: hashvalue,
		IpfsNodeId: ipfsnodeid,
		Operation: operation,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(hashvalue, recordJSON)
}

// ReadRecord returns the record stored in the world state with given hashvalue
func (s *SmartContract) ReadRecord(ctx contractapi.TransactionContextInterface, hashvalue string) (*Record, error){
	recordJSON, err := ctx.GetStub().GetState(hashvalue)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if recordJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", hashvalue)
	}

	var record Record
	err = json.Unmarshal(recordJSON, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

//RecordExists returns true when record with given hashvalue exists in world state
func (s *SmartContract) RecordExists(ctx contractapi.TransactionContextInterface, hashvalue string) (bool, error){
	recordJSON, err := ctx.GetStub().GetState(hashvalue)
	if err != nil{
		return false, fmt.Error("failed to read from world state: %v", err)
	}

	return recordJSON !=nil, nil
}

// GetAllRecords returns all records found in world state
func (s *SmartContract) GetAllRecords(ctx *contractapi.TransacntionContextInterface) ([]*Record, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("","")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*Record
	for resultsIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		var record Record
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	return records, nil
}

func main(){
	recordChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating ipfscc-basic chaincode: %v", err)
	}
	if err := recordChaincode.Start(); err != nil {
		log.Panicf("Error starting ipfscc-basic chaincode: %v", err)
	}
}
