package utils

import (
	"aurora-relayer-go-common/log"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/near/borsh-go"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	// storageSlotLength is the max length of the storage slot argument
	storageSlotLength = 32
	// ValueLength is the max length of the value argument
	ValueLength = 32
)

// BlockNum represents the blocks either in hex number format or in block tags
type BlockNum struct{ rpc.BlockNumber }

// ArgsForGetStorageAt is used to process GetStorageAt endpoint arguments
type ArgsForGetStorageAt struct {
	Address     Address
	StorageSlot Uint256
}

// NewArgsForGetStorageAt allocates and returns a new empty ArgsForGetStorageAt
func NewArgsForGetStorageAt() *ArgsForGetStorageAt {
	return &ArgsForGetStorageAt{}
}

// SetFields sets the Address and StorageSlot fields and returns a pointer of the object
func (gs *ArgsForGetStorageAt) SetFields(addr Address, sSlot Uint256) *ArgsForGetStorageAt {
	gs.Address = addr
	gs.StorageSlot = sSlot
	return gs
}

// Serialize transforms ArgsForGetStorageAt to ArgsForGetStorageAtEngine, calls its Serialize method
// and returns the received buffer
func (gs ArgsForGetStorageAt) Serialize() ([]byte, error) {
	tmpObj := NewArgsForGetStorageAtEngine().SetFields(gs.Address.Address[:], gs.StorageSlot.Bytes())
	buff, err := tmpObj.Serialize()
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// ArgsForGetStorageAtEngine is the data format accepted by engine for GetStorageAt endpoint
type ArgsForGetStorageAtEngine struct {
	Address [common.AddressLength]uint8
	Key     [storageSlotLength]uint8
}

// NewArgsForGetStorageAtEngine allocates and returns a new empty ArgsForGetStorageAtEngine
func NewArgsForGetStorageAtEngine() *ArgsForGetStorageAtEngine {
	return &ArgsForGetStorageAtEngine{}
}

// SetFields sets the Address and SotrageSlot buffers and returns a pointer of the object
func (gse *ArgsForGetStorageAtEngine) SetFields(addrBuf, keyBuf []byte) *ArgsForGetStorageAtEngine {
	startIndex := storageSlotLength - len(keyBuf)
	copy(gse.Address[:], addrBuf)
	copy(gse.Key[startIndex:], keyBuf)
	return gse
}

// Serialize ArgsForGetStorageAtEngine to a buffer using borsh so to communicate with engine
func (gse ArgsForGetStorageAtEngine) Serialize() ([]byte, error) {
	buff, err := borsh.Serialize(gse)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// TransactionForCall is the type used to serialize eth_call input
type TransactionForCall struct {
	From     *Address
	To       *Address // mandatory field
	Gas      *Uint256
	GasPrice *Uint256
	Value    *Uint256
	Data     *Uint256
}

// NewTransactionForCall allocates and returns a new empty TransactionForCall
func NewTransactionForCall() *TransactionForCall {
	return new(TransactionForCall)
}

// UnmarshalJSON implements json.Unmarshaler.
func (tc *TransactionForCall) UnmarshalJSON(data []byte) error {
	type tmpType TransactionForCall
	tmp := tmpType{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	if tmp.To == nil {
		return errors.New("missing value for `To` address")
	}

	tc.From = tmp.From
	tc.To = tmp.To
	tc.Gas = tmp.Gas
	tc.GasPrice = tmp.GasPrice
	tc.Value = tmp.Value
	tc.Data = tmp.Data
	return nil
}

// Serialize transforms TransactionForCall to TransactionForCallEngine, calls its Serialize method
// and returns the received buffer
func (tc TransactionForCall) Serialize() ([]byte, error) {

	to := *tc.To
	from := HexStringToAddress("0x0")
	if tc.From != nil {
		from = *tc.From
	}
	var value [ValueLength]uint8
	if tc.Value != nil {
		startIndexForValue := ValueLength - len(tc.Value.Bytes())
		copy(value[startIndexForValue:], tc.Value.Bytes())
	}
	var data []uint8
	if tc.Data != nil {
		copy(data[:], tc.Data.Bytes())
	}

	tmpObj := NewTransactionForCallEngine().SetFields(to.Address, from.Address, value, data)
	buff, err := tmpObj.Serialize()
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// TransactionForCallEngine is the type send to engine for eth_call endpoint
type TransactionForCallEngine struct {
	From  [common.AddressLength]uint8
	To    [common.AddressLength]uint8
	Value [32]uint8
	Data  []uint8
}

// NewTransactionForCall allocates and returns a new empty TransactionForCall
func NewTransactionForCallEngine() *TransactionForCallEngine {
	return new(TransactionForCallEngine)
}

// SetFields sets the Address and SotrageSlot buffers and returns a pointer of the object
func (tce *TransactionForCallEngine) SetFields(to, from [common.AddressLength]uint8, value [32]uint8, data []uint8) *TransactionForCallEngine {
	copy(tce.To[:], to[:])
	copy(tce.From[:], from[:])
	copy(tce.Value[:], value[:])
	copy(tce.Data[:], data[:])
	return tce
}

// Serialize TransactionForCallEngine to a buffer using borsh so to communicate with engine
func (tce TransactionForCallEngine) Serialize() ([]byte, error) {
	buff, err := borsh.Serialize(tce)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// QueryResult is the type received from engine for query (readonly) calls
type QueryResult struct {
	Result []interface{}
}

// ToResponse processes the engine query response, retrieves the `result` map and converts it to Uint256 response
func NewQueryResult(resp interface{}) (*QueryResult, error) {
	result, ok := resp.((map[string]interface{}))["result"].([]interface{})
	if !ok {
		log.Log().Error().Msgf("query response is not in correct format: %s", result)
		return nil, errors.New("query response is not in correct format")
	}
	return &QueryResult{Result: result}, nil
}

// ToResponse processes the engine query response, retrieves the `result` map and converts it to Uint256 response
func (r *QueryResult) ToResponse() *Uint256 {
	len := len(r.Result)
	buf := make([]byte, len)
	for i, b := range r.Result {
		if b, ok := b.(json.Number); ok {
			t, _ := b.Int64()
			buf[i] = byte(t)
		}
	}
	ui256 := IntToUint256(0)
	ui256.SetBytes(buf)
	return &ui256
}

// SubmitResultV2 is the type used to handle engine response for sendRawTransactionSync endpoint
type SubmitResultV2 struct {
	Version uint8
	Status  SendRawTxsStatus
	GasUsed uint64                `borsh_skip:"true"`
	logs    []LogEventWithAddress `borsh_skip:"true"`
}

// Deserialize uses borsh to initialize the SubmitResultV2 from the provided buffer
func (sr *SubmitResultV2) Deserialize(buf []byte) error {
	return borsh.Deserialize(sr, buf)
}

// Validate checks `SubmitResultV2.Status` to return the success or error
func (sr *SubmitResultV2) Validate() error {
	switch sr.Status.Enum {
	case 0: //SuccessStatus case
		return nil
	case 1: //RevertStatus case
		if len(sr.Status.Revert.Output) > 0 {
			return &InvalidParamsError{Message: fmt.Sprintf("execution error:TXS Revert with status %v", sr.Status.Revert.Output)}
		} else {
			return &InvalidParamsError{Message: "execution error: TXS Revert without any status"}
		}
	case 2: //OutOfGas case
		return &InvalidParamsError{Message: "execution error: Out Of Gas"}
	case 3: //OutOfFund case
		return &InvalidParamsError{Message: "execution error: Out Of Fund"}
	case 4: //OutOfOffset case
		return &InvalidParamsError{Message: "execution error: Out Of Offset"}
	case 5: //CallTooDeep case
		return &InvalidParamsError{Message: "execution error: Call Too Deep"}
	}
	log.Log().Debug().Msgf("unhandled SubmitResultV2 status: %d", sr.Status.Enum)
	return errors.New("execution error: unhandled SubmitResultV2 status")
}

// LogEventWithAddress is the type used to handle engine's SubmitResultV2 response
type LogEventWithAddress struct {
	Address [common.AddressLength]uint8
	Topics  []RawU256
	Data    []uint8
}

// RawU256 is the type used to handle engine's LogEventWithAddress response
type RawU256 struct {
	Value [32]uint8
}

// SendRawTxsStatus is the type used to handle engine's SubmitResultV2 response
type SendRawTxsStatus struct {
	Enum        borsh.Enum `borsh_enum:"true"` // treat struct as complex enum when serializing/deserializing
	Success     SendRawTxsSuccessStatus
	Revert      SendRawTxsRevertStatus
	OutOfGas    borsh.Enum
	OutOfFund   borsh.Enum
	OutOfOffset borsh.Enum
	CallTooDeep borsh.Enum
}

// SendRawTxsSuccessStatus is the type used to handle engine's SendRawTxsStatus response
type SendRawTxsSuccessStatus struct {
	Output []uint8
}

// SendRawTxsRevertStatus is the type used to handle engine's SendRawTxsStatus response
type SendRawTxsRevertStatus struct {
	Output []uint8
}

// SubmitStatus is the type received from engine for submit (eg: sendRawTransactionSync) calls
type SubmitStatus struct {
	StatusMap    map[string]interface{}
	SubmitResult *SubmitResultV2
	ResponseHash string
}

// NewSubmitStatus allocates and returns a new NewSubmitStatus object
func NewSubmitStatus(respArg interface{}, txsHash string) (*SubmitStatus, error) {
	resp, ok := respArg.((map[string]interface{}))
	if !ok {
		log.Log().Error().Msgf("submit response is not in correct format: %s", respArg)
		return nil, errors.New("submit response is not in correct format")
	}

	status, ok := resp["status"].((map[string]interface{}))
	if !ok {
		log.Log().Error().Msgf("submit status is not in correct format: %s", status)
		return nil, errors.New("submit status is not in correct format")
	}
	return &SubmitStatus{
		StatusMap:    status,
		SubmitResult: &SubmitResultV2{},
		ResponseHash: txsHash,
	}, nil
}

//Validate checks `SubmitStatus.StatusMap` to return the success or error
func (ss *SubmitStatus) Validate() error {
	logger := log.Log()
	// Check if any error returned
	fail, ok := ss.StatusMap["Failure"]
	if ok {
		failTypeMap, ok := fail.((map[string]interface{}))["ActionError"].((map[string]interface{}))["kind"].((map[string]interface{}))
		if !ok {
			logger.Error().Msgf("submit request failure while parsing `Failure` object, txs hash: %s", ss.ResponseHash)
			return errors.New("submit request failure while parsing `Failure` object")
		}

		reg := regexp.MustCompile("[\n\r\t\"]")
		failType := ""
		// Access the first element of the map
		for k := range failTypeMap {
			failType = k
			break
		}
		switch failType {
		case "FunctionCallError":
			execErr, ok := failTypeMap["FunctionCallError"].((map[string]interface{}))["ExecutionError"].(string)
			if ok {
				errMsg := strings.Replace(execErr, "Smart contract panicked: ", "", 1)
				logger.Debug().Msgf("submit request failed with ExecutionError: %s, txs hash: %s", errMsg, ss.ResponseHash)
				return errors.New(errMsg)
			} else {
				jsn, err := json.Marshal(failTypeMap["FunctionCallError"])
				if err != nil {
					logger.Error().Msgf("submit request failed while marshalling FunctionCallError: %s, txs hash: %s", err.Error(), ss.ResponseHash)
					return err
				}
				jsnStr := reg.ReplaceAllString(string(jsn), "")
				logger.Error().Msgf("submit request failed with error: %s, txs hash: %s", jsnStr, ss.ResponseHash)
				return fmt.Errorf("failure:%s", jsnStr)
			}
		case "MethodNotFound":
			jsn, err := json.Marshal(failTypeMap["MethodNotFound"])
			if err != nil {
				logger.Error().Msgf("submit request failed while marshalling MethodNotFound: %s, txs hash: %s", err.Error(), ss.ResponseHash)
				return err
			}
			jsnStr := reg.ReplaceAllString(string(jsn), "")
			logger.Error().Msgf("submit request failed with MethodNotFound: %s, txs hash: %s", jsnStr, ss.ResponseHash)
			return fmt.Errorf("failure: %s", jsnStr)
		default:
			jsn, err := json.Marshal(fail.((map[string]interface{}))["ActionError"].((map[string]interface{}))["kind"])
			if err != nil {
				logger.Error().Msgf("submit request failed while marshalling Default case: %s, txs hash: %s", err.Error(), ss.ResponseHash)
				return err
			}
			jsnStr := reg.ReplaceAllString(string(jsn), "")
			logger.Error().Msgf("submit request failed with Default case error: %s, txs hash: %s", jsnStr, ss.ResponseHash)
			return fmt.Errorf("failure:%s", jsnStr)
		}
	} else if ss.StatusMap["SuccessValue"] != nil {
		sucStrB64, ok := ss.StatusMap["SuccessValue"].(string)
		if ok {
			sucBuf, err := base64.URLEncoding.DecodeString(sucStrB64)
			if err != nil {
				return err
			}
			if len(sucBuf) == 0 {
				return nil
			}
			if err = ss.SubmitResult.Deserialize(sucBuf); err != nil {
				return err
			}
			return ss.SubmitResult.Validate()
		}
	}
	logger.Error().Msgf("submit request returned with unhandled error, txs hash: %s", ss.ResponseHash)
	return errors.New("submit request returned with unhandled error")

}

// ToResponse processes the engine query response, retrieves the `result` map and converts it to Uint256 response
func (ss *SubmitStatus) ToResponse() *H256 {
	h256Hash := HexStringToHash(ss.ResponseHash)
	return &h256Hash
}
