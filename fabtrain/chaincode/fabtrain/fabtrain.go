/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the train structure, with 4 properties.  Structure tags are used by encoding/json library
type Train struct {
	Fname  string `json:"fname"`
	Gender  string `json:"gender"`
	Place string `json:"place"`
	Class  string `json:"class"`
	Status  string `json:"status"`
}

/*
 * The Init method is called when the Smart Contract "fabtrain" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabtrain"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryTrain" {
		return s.queryTrain(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createTrain" {
		return s.createTrain(APIstub, args)
	} else if function == "queryAllTrains" {
		return s.queryAllTrains(APIstub)
	} else if function == "changeTrainStatus" {
		return s.changeTrainStatus(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryTrain(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	trainAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(trainAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	trains := []Train{
		Train{Fname: "Lokesh", Gender: "Male", Place: "Bangalore-Mysore", Class: "AC", Status: "Waiting list"},
		Train{Fname: "Subhiksha", Gender: "Female", Place: "Bangalore-Chennai", Class: "AC", Status: "RAC"},
		Train{Fname: "Rakshita", Gender: "Female", Place: "Chennai-Mysore", Class: "Sleeper", Status: "Confirmed"},
		Train{Fname: "Sachin", Gender: "Male", Place: "Bangalore-Mumbai", Class: "General", Status: "Waiting list"},
		Train{Fname: "Kiran", Gender: "Male", Place: "Mysore-Mangalore", Class: "AC", Status: "Waiting list"},
	}

	i := 0
	for i < len(trains) {
		fmt.Println("i is ", i)
		trainAsBytes, _ := json.Marshal(trains[i])
		APIstub.PutState("TRAIN"+strconv.Itoa(i), trainAsBytes)
		fmt.Println("Added", trains[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createTrain(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var train = Train{Fname: args[1], Gender: args[2], Place: args[3], Class: args[4], Status: args[5]}

	trainAsBytes, _ := json.Marshal(train)
	APIstub.PutState(args[0], trainAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllTrains(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "TRAIN0"
	endKey := "TRAIN999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTrains:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeTrainStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	trainAsBytes, _ := APIstub.GetState(args[0])
	train := Train{}

	json.Unmarshal(trainAsBytes, &train)
	train.Status = args[1]

	trainAsBytes, _ = json.Marshal(train)
	APIstub.PutState(args[0], trainAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
