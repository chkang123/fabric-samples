/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Commodity struct {
	ID              	string 		`json:"id"`
	Trust_score     	float32 	`json:"trust_score"`
	Owners    			[]string 	`json:"owners"`
	Readings        	float32 	`json:"readings"`
	Ideal_Temp      	float32 	`json:"ideal_temp"`
	Parents_IDs			[]string	`json:"parents_ids"`
 }

type Participant struct {
	ID              	string 		`json:"id"`
	Reputation_score    float32 	`json:"reputation_score"`
	Device_IDs    		[]string 	`json:"device_ids"`
}

type Device struct {
	ID              	string 		`json:"id"`
	Type     			string 		`json:"type"`
}

type QueryResult_Com struct {
	Key    string `json:"Key"`
	Record *Commodity
}

type QueryResult_Part struct {
	Key    string `json:"Key"`
	Record *Participant
}

type QueryResult_Dev struct {
	Key    string `json:"Key"`
	Record *Device
}

type QueryResult_All struct {
	Record_Com []QueryResult_Com `json:"Com"`
	Record_Part []QueryResult_Part `json:"Part"`
	Record_Dev []QueryResult_Dev `json:"Dev"`
}

// InitLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	commodities := []Commodity{
		Commodity{ID: "COM1", Trust_score: 5, Owners: []string{"PART1"}, Readings: 25, Ideal_Temp: 10, Parents_IDs: []string{}},
		Commodity{ID: "COM2", Trust_score: 5, Owners: []string{"PART2"}, Readings: 25, Ideal_Temp: 20, Parents_IDs: []string{}},
		Commodity{ID: "COM3", Trust_score: 5, Owners: []string{"PART3"}, Readings: 25, Ideal_Temp: 30, Parents_IDs: []string{}},
		Commodity{ID: "COM4", Trust_score: 5, Owners: []string{"PART4"}, Readings: 25, Ideal_Temp: 40, Parents_IDs: []string{}},
	}

	for i, com := range commodities {
		comAsBytes, _ := json.Marshal(com)
		err := ctx.GetStub().PutState("COM"+strconv.Itoa(i), comAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	participants := []Participant{
		Participant{ID: "PART1", Reputation_score: 1, Device_IDs: []string{"DEV1"}},
		Participant{ID: "PART2", Reputation_score: 1, Device_IDs: []string{"DEV2"}},
		Participant{ID: "PART3", Reputation_score: 2, Device_IDs: []string{"DEV3"}},
		Participant{ID: "PART4", Reputation_score: 3, Device_IDs: []string{"DEV4"}},
	}

	for i, part := range participants {
		partAsBytes, _ := json.Marshal(part)
		err := ctx.GetStub().PutState("PART"+strconv.Itoa(i), partAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	devices := []Device{
		Device{ID: "DEV1", Type: "Temp"},
		Device{ID: "DEV2", Type: "Temp"},
		Device{ID: "DEV3", Type: "Temp"},
		Device{ID: "DEV4", Type: "Temp"},
	}

	for i, dev := range devices {
		devAsBytes, _ := json.Marshal(dev)
		err := ctx.GetStub().PutState("DEV"+strconv.Itoa(i), devAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) CreateCom(ctx contractapi.TransactionContextInterface, id string, score float32, owner string, readings float32, ideal float32) error {
	com := Commodity{
		ID: id,
		Trust_score: score,
		Owners: []string{owner},
		Readings: readings,
		Ideal_Temp: ideal,
		Parents_IDs: []string{},
	}

	comAsBytes, _ := json.Marshal(com)

	return ctx.GetStub().PutState(id, comAsBytes)
}

func (s *SmartContract) CreatePart(ctx contractapi.TransactionContextInterface, id string, score float32) error {
	part := Participant{
		ID: id,
		Reputation_score: score,
		Device_IDs: []string{},
	}

	partAsBytes, _ := json.Marshal(part)

	return ctx.GetStub().PutState(id, partAsBytes)
}

func (s *SmartContract) CreateDev(ctx contractapi.TransactionContextInterface, id string, dev_type string) error {
	dev := Device{
		ID: id,
		Type: dev_type,
	}

	devAsBytes, _ := json.Marshal(dev)

	return ctx.GetStub().PutState(id, devAsBytes)
}

func (s *SmartContract) DeleteData(ctx contractapi.TransactionContextInterface, id string) error {
	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) QueryCom(ctx contractapi.TransactionContextInterface, comID string) (*Commodity, error) {
	comAsBytes, err := ctx.GetStub().GetState(comID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if comAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", comID)
	}

	com := new(Commodity)
	_ = json.Unmarshal(comAsBytes, com)

	return com, nil
}

func (s *SmartContract) QueryPart(ctx contractapi.TransactionContextInterface, partID string) (*Participant, error) {
	partAsBytes, err := ctx.GetStub().GetState(partID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if partAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", partID)
	}

	part := new(Participant)
	_ = json.Unmarshal(partAsBytes, part)

	return part, nil
}

func (s *SmartContract) QueryDev(ctx contractapi.TransactionContextInterface, devID string) (*Device, error) {
	devAsBytes, err := ctx.GetStub().GetState(devID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if devAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", devID)
	}

	dev := new(Device)
	_ = json.Unmarshal(devAsBytes, dev)

	return dev, nil
}

func (s *SmartContract) QueryAll(ctx contractapi.TransactionContextInterface) (QueryResult_All, error) {
	/*
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results_Com := []QueryResult_Com{}
	results_Part := []QueryResult_Part{}
	results_Dev := []QueryResult_Dev{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		switch queryResponse.(type) {
		case Commodity:
			com := new(Commodity)
			_ = json.Unmarshal(queryResponse.Value, com)
			queryResult_com := QueryResult_Com{Key: queryResponse.Key, Record: com}
			results_Com = append(results_Com, queryResult_com)
			break
		case Participant:
			part := new(Participant)
			_ = json.Unmarshal(queryResponse.Value, part)
			queryResult_part := QueryResult_Part{Key: queryResponse.Key, Record: part}
			results_Part = append(results_Part, queryResult_part)
			break
		case Device:
			dev := new(Device)
			_ = json.Unmarshal(queryResponse.Value, dev)
			queryResult_dev := QueryResult_Dev{Key: queryResponse.Key, Record: dev}
			results_Dev = append(results_Dev, queryResult_dev)
			break
		default:
			return nil, nil
		}

	}

	queryResult := QueryResult_All{Record_Com: results_Com, Record_Part: results_Part, Record_Dev: results_Dev}

	return queryResult, nil
	*/

	results_Com := []QueryResult_Com{}
	results_Part := []QueryResult_Part{}
	results_Dev := []QueryResult_Dev{}

	com1, _ := s.QueryCom(ctx, "COM1")
	results_Com = append(results_Com, QueryResult_Com{Key:"COM1", Record:com1})
	com2, _ := s.QueryCom(ctx, "COM2")
	results_Com = append(results_Com, QueryResult_Com{Key:"COM2", Record:com2})
	com3, _ := s.QueryCom(ctx, "COM3")
	results_Com = append(results_Com, QueryResult_Com{Key:"COM3", Record:com3})
	com4, _ := s.QueryCom(ctx, "COM4")
	results_Com = append(results_Com, QueryResult_Com{Key:"COM4", Record:com4})

	part1, _ := s.QueryPart(ctx, "PART1")
	results_Part = append(results_Part, QueryResult_Part{Key:"PART1", Record:part1})
	part2, _ := s.QueryPart(ctx, "PART2")
	results_Part = append(results_Part, QueryResult_Part{Key:"PART2", Record:part2})
	part3, _ := s.QueryPart(ctx, "PART3")
	results_Part = append(results_Part, QueryResult_Part{Key:"PART3", Record:part3})
	part4, _ := s.QueryPart(ctx, "PART4")
	results_Part = append(results_Part, QueryResult_Part{Key:"PART4", Record:part4})

	dev1, _ := s.QueryDev(ctx, "DEV1")
	results_Dev = append(results_Dev, QueryResult_Dev{Key:"DEV1", Record:dev1})
	dev2, _ := s.QueryDev(ctx, "DEV2")
	results_Dev = append(results_Dev, QueryResult_Dev{Key:"DEV2", Record:dev2})
	dev3, _ := s.QueryDev(ctx, "DEV3")
	results_Dev = append(results_Dev, QueryResult_Dev{Key:"DEV3", Record:dev3})
	dev4, _ := s.QueryDev(ctx, "DEV4")
	results_Dev = append(results_Dev, QueryResult_Dev{Key:"DEV4", Record:dev4})

	queryResult := QueryResult_All{Record_Com: results_Com, Record_Part: results_Part, Record_Dev: results_Dev}

	return queryResult, nil
}

// Trade
func (s *SmartContract) TradeCom(ctx contractapi.TransactionContextInterface, comNumber string, newOwner string) error {
	com, err := s.QueryCom(ctx, comNumber)

	if err != nil {
		return err
	}

	com.Owners = append(com.Owners, newOwner)

	comAsBytes, _ := json.Marshal(com)

	return ctx.GetStub().PutState(comNumber, comAsBytes)
}

// Produce
func (s *SmartContract) ProduceCom(ctx contractapi.TransactionContextInterface, comNumbers []string, id string, ideal float32) error {
	coms := []*Commodity{}
	var trust_score float32 = 0
	parents_ids := []string{}

	for _, comNumber := range comNumbers {
		com, _ := s.QueryCom(ctx, comNumber)
		coms = append(coms, com)
		trust_score = trust_score + com.Trust_score
		parents_ids = append(parents_ids, com.ID)
	}
	trust_score = trust_score / float32(len(coms))

	newCom := Commodity{
		ID: id,
		Trust_score: trust_score,
		Owners: []string{coms[0].Owners[len(coms[0].Owners)-1]},
		Readings: float32(0),
		Ideal_Temp: ideal,
		Parents_IDs: parents_ids,
	}

	comAsBytes, _ := json.Marshal(newCom)

	return ctx.GetStub().PutState(id, comAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create dairy chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting dairy chaincode: %s", err.Error())
	}
}
