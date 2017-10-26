package main

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "errors"
	// "strconv"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("example_cc0")

//SmartContract P2P Lending Smart Contract
type SmartContract struct {
}

//Account model
type Account struct {
	Name string `json:"name"`
	Risk int    `json:"risk"`
	Type string `json:"type"`
	Fund int    `json:"fund"`
	Loan int    `json:"loan"` //loan given or taken
	Auto bool   `json:"auto"`
}

//Init to intialize
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke transaction
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke is running " + function)
	logger.Info("sk0")

	if function == "initLedger" {
		return s.initLedger(stub, args)
	} else if function == "borrow" {
		return s.borrow(stub, args)
	} else if function == "updateRisk" {
		return s.updateRisk(stub, args)
	}
	if function == "query" {
		fmt.Println("sk1")
		// queries an entity state
		return s.query(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return shim.Error(fmt.Sprintf("Received unknown function invocation: " + function))
}

//initLedger
func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("initLedger called")

	stub.DelState("ACCOUNT0")
	stub.DelState("ACCOUNT1")
	stub.DelState("ACCOUNT2")

	Accounts := []Account{
		Account{Name: "Lender Harrison", Risk: 3, Type: "LENDER", Fund: 20000, Loan: 0, Auto: true},
		Account{Name: "Lender Gibson", Risk: 2, Type: "LENDER", Fund: 20000, Loan: 0, Auto: false},
		Account{Name: "Borrower Peter", Risk: 1, Type: "BORROWER", Fund: 0, Loan: 0, Auto: false},
	}

	i := 0
	for i < len(Accounts) {
		fmt.Println("i is ", i)
		accountAsBytes, _ := json.Marshal(Accounts[i])
		stub.PutState("ACCOUNT"+strconv.Itoa(i), accountAsBytes)
		fmt.Println("Added", Accounts[i])
		i = i + 1
	}

	return shim.Success(nil)
}

//borrow
func (s *SmartContract) borrow(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("borrow called")
	//step 1 : define [borrowerId, fundsNeeded, borrowerRisk]
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	borrowerID := args[0]
	fundsNeeded, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	remaining := fundsNeeded
	borrowerAsBytes, _ := stub.GetState(borrowerID)
	borrower := Account{}
	json.Unmarshal(borrowerAsBytes, &borrower)
	borrowerRisk := borrower.Risk
	borrower.Loan = borrower.Loan + fundsNeeded

	//step 2 : get [borrowerRisk,matchedLenders]
	lender1AsBytes, _ := stub.GetState("ACCOUNT0")
	lender2AsBytes, _ := stub.GetState("ACCOUNT1")
	lender1 := Account{}
	lender2 := Account{}
	json.Unmarshal(lender1AsBytes, &lender1)
	json.Unmarshal(lender2AsBytes, &lender2)
	lendersS := [2]Account{}
	lendersS[0] = lender1
	lendersS[1] = lender2

	i := 0
	for i < len(lendersS) {
		key := "ACCOUNT0"
		if i == 1 {
			key = "ACCOUNT1"
		}
		val := lendersS[i]

		if val.Risk == borrowerRisk || val.Auto == true {
			if val.Fund > 0 {
				toTransfer := remaining
				if toTransfer > val.Fund {
					toTransfer = val.Fund
				}
				remaining = remaining - toTransfer
				val.Fund = val.Fund - toTransfer
				val.Loan = val.Loan + toTransfer

				//if val.Auto != true && val.Risk != 1 {
				//    val.Risk = val.Risk - 1
				//}
				lenderAsBytes, _ := json.Marshal(val)
				e := stub.PutState(key, lenderAsBytes)
				if e != nil {
					logger.Warning("level6 err")
				}
				borrower.Fund = borrower.Fund + toTransfer
			}
		}
		if remaining == 0 {
			break
		}
		i = i + 1
	}

	if remaining == 0 {
		if borrower.Risk != 3 {
			borrower.Risk = borrower.Risk + 1
		}
	}
	borrowerAsBytes, _ = json.Marshal(borrower)
	stub.PutState(borrowerID, borrowerAsBytes)
	return shim.Success(borrowerAsBytes)
}

//updateRisk
func (s *SmartContract) updateRisk(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("updateRisk called")

	if len(args) < 3 { //0:id, 1:risk
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	accountID := args[0]
	risk, err := strconv.Atoi(args[1])
	auto, err2 := strconv.Atoi(args[2])

	autoB := false
	if auto != 0 {
		autoB = true
	}

	if err != nil {
		return shim.Error(err.Error())
	}
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	accountAsBytes, _ := stub.GetState(accountID)
	account := Account{}
	json.Unmarshal(accountAsBytes, &account)
	account.Risk = risk
	account.Auto = autoB
	accountAsBytes, _ = json.Marshal(account)
	stub.PutState(accountID, accountAsBytes)

	return shim.Success(accountAsBytes)
}

// Query callback representing the query of a chaincode
func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("sk2")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("len(args) = " + strconv.Itoa(len(args)))

	logger.Debug("query called")
	var argA = args[0]
	//fmt.Println("query is running " + function)
	fmt.Println("query is running args[0]" + argA)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting exact 1 argument")
	}

	var queryFunctionName = args[0]
	fmt.Println("queryFunctionName is = " + queryFunctionName)

	if queryFunctionName == "read" {
		fmt.Println("sk3")
		return s.read(stub, args)
	} else if queryFunctionName == "readAll" {
		fmt.Println("sk4")
		return s.readAll(stub, args)
	}
	fmt.Println("query did not find func: " + function)
	fmt.Println("sk5")

	fmt.Println("sk2 completed")
	return shim.Error("Received unknown function = " + function)
}

//read
func (s *SmartContract) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("read called")

	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

//readAll
func (s *SmartContract) readAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("readAll - 1")
	logger.Debug("readAll called")

	var jsonResp string
	valAsbytes1, err1 := stub.GetState("ACCOUNT0")
	valAsbytes2, err2 := stub.GetState("ACCOUNT1")
	valAsbytes3, err3 := stub.GetState("ACCOUNT2")

	if err1 != nil {
		jsonResp = "{\"Error\":\"Failed to get state for ACCOUNT0\"}"
		return shim.Error(jsonResp)
	}
	if err2 != nil {
		jsonResp = "{\"Error\":\"Failed to get state for ACCOUNT1\"}"
		return shim.Error(jsonResp)
	}
	if err3 != nil {
		jsonResp = "{\"Error\":\"Failed to get state for ACCOUNT3\"}"
		return shim.Error(jsonResp)
	}

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(string(valAsbytes1))
	buffer.WriteString(",")

	buffer.WriteString(string(valAsbytes2))
	buffer.WriteString(",")
	buffer.WriteString(string(valAsbytes3))

	buffer.WriteString("]")
	//var b []byte
	//  valAsbytes := []byte("["+string(valAsbytes1)+","+ string(valAsbytes2)+"," + string(valAsbytes3)+"]")
	fmt.Println("readAll - 2")
	return shim.Success(buffer.Bytes())
}

/*
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if function == "initLedger" {
        return s.initLedger(APIstub, args)
    } else if function == "borrow" {
        return s.borrow(APIstub, args)
    } else if function == "updateRisk" {
        return s.updateRisk(APIstub, args)
    }

    fmt.Println("invoke did not find func: " + function)
    return nil, errors.New("Received unknown function invocation: " + function)
}
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

    APIstub.DelState("ACCOUNT0")
    APIstub.DelState("ACCOUNT1")
    APIstub.DelState("ACCOUNT2")

    Accounts := []Account{
        Account{Name:"Lender Harrison", Risk:3, Type:"LENDER", Fund:20000, Loan:0,Auto:true},
        Account{Name:"Lender Gibson", Risk:2, Type:"LENDER", Fund:20000, Loan:0,Auto:false},
        Account{Name:"Borrower Peter", Risk:1, Type:"BORROWER", Fund:0, Loan:0,Auto:false},
    }

    i := 0
    for i < len(Accounts) {
        fmt.Println("i is ", i)
        accountAsBytes, _ := json.Marshal(Accounts[i])
        APIstub.PutState("ACCOUNT"+strconv.Itoa(i), accountAsBytes)
        fmt.Println("Added", Accounts[i])
        i = i + 1
    }

    return nil, nil
}
func (s *SmartContract) borrow(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    //step 1 : define [borrowerId, fundsNeeded, borrowerRisk]
    if len(args) < 2 {
        //return shim.Error("Incorrect number of arguments. Expecting 2")
        return nil, errors.New("Incorrect number of arguments. Expecting 2")
    }
    borrowerId := args[0]
    fundsNeeded, err := strconv.Atoi(args[1]);
    if err != nil {
        return nil, errors.New(err.Error())
    }

    remaining := fundsNeeded
    borrowerAsBytes, _ := APIstub.GetState(borrowerId)
    borrower := Account{}
    json.Unmarshal(borrowerAsBytes, &borrower)
    borrowerRisk := borrower.Risk
    borrower.Loan = borrower.Loan + fundsNeeded

    //step 2 : get [borrowerRisk,matchedLenders]
    lender1AsBytes, _ := APIstub.GetState("ACCOUNT0")
    lender2AsBytes, _ := APIstub.GetState("ACCOUNT1")
    lender1 := Account{}
    lender2 := Account{}
    json.Unmarshal(lender1AsBytes, &lender1)
    json.Unmarshal(lender2AsBytes, &lender2)
    lendersS := [2]Account{}
    lendersS[0] = lender1
    lendersS[1] = lender2

    i := 0
    for i < len(lendersS) {
        key := "ACCOUNT0"
        if i == 1 {
            key = "ACCOUNT1"
        }
        val := lendersS[i]

        if val.Risk == borrowerRisk || val.Auto == true {
            if val.Fund > 0 {
                toTransfer := remaining
                if toTransfer > val.Fund {
                    toTransfer = val.Fund
                }
                remaining = remaining - toTransfer
                val.Fund = val.Fund - toTransfer
                val.Loan = val.Loan + toTransfer

                //if val.Auto != true && val.Risk != 1 {
                //    val.Risk = val.Risk - 1
                //}
                lenderAsBytes, _ := json.Marshal(val)
                e := APIstub.PutState(key, lenderAsBytes)
                if e != nil {
                    logger.Warning("level6 err")
                }
                borrower.Fund = borrower.Fund + toTransfer
            }
        }
        if remaining == 0 {
            break
        }
        i = i + 1
    }

    if remaining == 0 {
        if borrower.Risk != 3 {
            borrower.Risk = borrower.Risk + 1
        }
    }
    borrowerAsBytes, _ = json.Marshal(borrower)
    APIstub.PutState(borrowerId, borrowerAsBytes)
    return borrowerAsBytes, nil
}
func (s *SmartContract) updateRisk(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    if len(args) < 3 { //0:id, 1:risk
        return nil, errors.New("Incorrect number of arguments. Expecting 3")
    }
    accountId := args[0]
    risk, err := strconv.Atoi(args[1])
    auto, err2 := strconv.Atoi(args[2])

    autoB := false
    if auto != 0 {
        autoB = true
    }

    if err != nil {
        return nil, errors.New(err.Error())
    }
    if err2 != nil {
        return nil, errors.New(err2.Error())
    }

    accountAsBytes, _ := APIstub.GetState(accountId)
    account := Account{}
    json.Unmarshal(accountAsBytes, &account)
    account.Risk = risk
    account.Auto = autoB
    accountAsBytes, _ = json.Marshal(account)
    APIstub.PutState(accountId, accountAsBytes)

    return accountAsBytes, nil
}

func (s *SmartContract) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    if function == "read" {
        return s.read(stub, args)
    } else if function == "readAll" {
        return s.readAll(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}
func (s *SmartContract) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}
func (s *SmartContract) readAll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var jsonResp string
    valAsbytes1, err1 := stub.GetState("ACCOUNT0")
    valAsbytes2, err2 := stub.GetState("ACCOUNT1")
    valAsbytes3, err3 := stub.GetState("ACCOUNT2")

    if err1 != nil {
        jsonResp = "{\"Error\":\"Failed to get state for ACCOUNT0\"}"
        return nil, errors.New(jsonResp)
    }
    if err2 != nil {
        jsonResp = "{\"Error\":\"Failed to get state for ACCOUNT1\"}"
        return nil, errors.New(jsonResp)
    }
    if err3 != nil {
        jsonResp = "{\"Error\":\"Failed to get state for ACCOUNT3\"}"
        return nil, errors.New(jsonResp)
    }

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(string(valAsbytes1))
	buffer.WriteString(",")

	buffer.WriteString(string(valAsbytes2))
	buffer.WriteString(",")
	buffer.WriteString(string(valAsbytes3))

	buffer.WriteString("]")
    //var b []byte
  //  valAsbytes := []byte("["+string(valAsbytes1)+","+ string(valAsbytes2)+"," + string(valAsbytes3)+"]")

    return buffer.Bytes(), nil
}
*/

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
