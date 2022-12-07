package engine

import (
	"aurora-relayer-go-common/log"
	error2 "aurora-relayer-go-common/types/errors"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/near/borsh-go"

	cc "aurora-relayer-go-common/types/common"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// storageSlotLength is the max length of the storage slot argument
	storageSlotLength = 32
	// ValueLength is the max length of the value argument
	ValueLength = 32
)

// ArgsForGetStorageAt is used to process GetStorageAt endpoint arguments
type ArgsForGetStorageAt struct {
	Address     cc.Address
	StorageSlot cc.Uint256
}

// NewArgsForGetStorageAt allocates and returns a new empty ArgsForGetStorageAt
func NewArgsForGetStorageAt() *ArgsForGetStorageAt {
	return &ArgsForGetStorageAt{}
}

// SetFields sets the Address and StorageSlot fields and returns a pointer of the object
func (gs *ArgsForGetStorageAt) SetFields(addr cc.Address, sSlot cc.Uint256) *ArgsForGetStorageAt {
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

// SetFields sets the Address and StorageSlot buffers and returns a pointer of the object
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
	From     *cc.Address `json:"from,omitempty"`
	To       *cc.Address `json:"to"`
	Gas      *cc.Uint256 `json:"gas,omitempty"`
	GasPrice *cc.Uint256 `json:"gasPrice,omitempty"`
	Value    *cc.Uint256 `json:"value,omitempty"`
	Data     cc.DataVec  `json:"data,omitempty"`
}

// NewTransactionForCall allocates and returns a new empty TransactionForCall
func NewTransactionForCall() *TransactionForCall {
	return new(TransactionForCall)
}

// UnmarshalJSON implements json.Unmarshaler
// This method is needed to make `.To` field of the TransactionForCall struct required/mandatory
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

	*tc = TransactionForCall(tmp)
	return nil
}

// Serialize transforms TransactionForCall to TransactionForCallEngine, calls its Serialize method
// and returns the received buffer
func (tc TransactionForCall) Serialize() ([]byte, error) {
	to := *tc.To
	from := cc.HexStringToAddress("0x0")
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
		data = make([]uint8, len(tc.Data))
		copy(data, tc.Data)
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

// NewTransactionForCallEngine allocates and returns a new empty TransactionForCall
func NewTransactionForCallEngine() *TransactionForCallEngine {
	return new(TransactionForCallEngine)
}

// SetFields sets the Address and StorageSlot buffers and returns a pointer of the object
func (tce *TransactionForCallEngine) SetFields(to, from [common.AddressLength]uint8, value [32]uint8, data []uint8) *TransactionForCallEngine {
	copy(tce.To[:], to[:])
	copy(tce.From[:], from[:])
	copy(tce.Value[:], value[:])
	tce.Data = make([]uint8, len(data))
	copy(tce.Data, data)
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

// NewQueryResult allocates and returns a new QueryResult object
func NewQueryResult(resp interface{}) (*QueryResult, error) {
	result, ok := resp.((map[string]interface{}))["result"].([]interface{})
	if !ok {
		log.Log().Error().Msgf("query response is not in correct format: %s", result)
		return nil, errors.New("query response is not in correct format")
	}
	return &QueryResult{Result: result}, nil
}

// ToUint256Response processes the engine query response, retrieves the `result` map and converts it to Uint256 response
func (r *QueryResult) ToUint256Response() (*cc.Uint256, error) {
	buf := r.resultToByteBuffer()
	ui256 := cc.Uint256FromBytes(buf)
	return &ui256, nil
}

// ToStringResponse processes the engine query response, retrieves the `result` map and converts it to string response
func (r *QueryResult) ToStringResponse() (*string, error) {
	buf := r.resultToByteBuffer()
	strHex := "0x" + hex.EncodeToString(buf)
	return &strHex, nil
}

// resultToByteBuffer processes the engine query response and creates a byte buffer
func (r *QueryResult) resultToByteBuffer() []byte {
	length := len(r.Result)
	buf := make([]byte, length)
	for i, b := range r.Result {
		if b, ok := b.(json.Number); ok {
			t, _ := b.Int64()
			buf[i] = byte(t)
		}
	}
	return buf
}

// SubmitResultV2 is the type used to handle engine response for sendRawTransactionSync endpoint
type SubmitResultV2 struct {
	Version uint8
	Status  TransactionStatus
	GasUsed uint64                `borsh_skip:"true"`
	logs    []LogEventWithAddress `borsh_skip:"true"`
}

// Deserialize uses borsh to initialize the SubmitResultV2 from the provided buffer
func (sr *SubmitResultV2) Deserialize(buf []byte) error {
	return borsh.Deserialize(sr, buf)
}

// Validate checks `SubmitResultV2.Status` to return the success or errors
func (sr *SubmitResultV2) Validate() error {
	_, err := sr.Status.Validate()
	return err
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

// TransactionStatus is the type used to handle engine's SubmitResultV2 response
type TransactionStatus struct {
	Enum        borsh.Enum `borsh_enum:"true"` // treat struct as complex enum when serializing/deserializing
	Success     TransactionSuccessStatus
	Revert      TransactionRevertStatus
	OutOfGas    borsh.Enum
	OutOfFund   borsh.Enum
	OutOfOffset borsh.Enum
	CallTooDeep borsh.Enum
}

// TransactionSuccessStatus is the type used to handle engine's TransactionStatus response
type TransactionSuccessStatus struct {
	Output []uint8
}

// TransactionRevertStatus is the type used to handle engine's TransactionStatus response
type TransactionRevertStatus struct {
	Output []uint8
}

// NewTransactionStatus allocates and returns a new TransactionStatus object
func NewTransactionStatus(respArg interface{}) (*TransactionStatus, error) {
	resp, ok := respArg.((map[string]interface{}))["result"].([]interface{})
	if !ok {
		err, ok := respArg.((map[string]interface{}))["error"].(string)
		if !ok {
			return nil, errors.New("call response is not in correct format")
		}
		log.Log().Error().Msgf("errors returned to eth_call: %v", err)
		return nil, fmt.Errorf("%v", err)
	}
	lenResp := len(resp)
	buf := make([]byte, lenResp)
	for i, b := range resp {
		if b, ok := b.(json.Number); ok {
			t, _ := b.Int64()
			buf[i] = byte(t)
		}
	}

	ts := &TransactionStatus{}
	// TODO -- An interim solution to handle "OutOfGas, OutOfFund, OutOfOffset, CallTooDeep" txs statuses
	if len(buf) == 1 && buf[0] > 1 && buf[0] < 8 {
		tmp := make([]byte, 1)
		tmp[0] = 1 << buf[0]
		buf = append(buf, tmp[0])
	}

	err := borsh.Deserialize(ts, buf)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

// Validate checks `TransactionStatus` to return the success or errors
func (ts *TransactionStatus) Validate() ([]uint8, error) {
	switch ts.Enum {
	case 0: // SuccessStatus case
		if len(ts.Success.Output) > 0 {
			return ts.Success.Output, nil
		} else {
			return []uint8{}, nil
		}
	case 1: // RevertStatus case
		if len(ts.Revert.Output) > 0 {
			return ts.Revert.Output, &error2.InvalidParamsError{Message: fmt.Sprintf("execution errors: transaction revert with status %v", ts.Revert.Output)}
		} else {
			return []uint8{}, &error2.InvalidParamsError{Message: "execution errors: transaction revert without any status"}
		}
	case 2: // OutOfGas case
		return nil, &error2.InvalidParamsError{Message: "execution errors: Out Of Gas"}
	case 3: // OutOfFund case
		return nil, &error2.InvalidParamsError{Message: "execution errors: Out Of Fund"}
	case 4: // OutOfOffset case
		return nil, &error2.InvalidParamsError{Message: "execution errors: Out Of Offset"}
	case 5: // CallTooDeep case
		return nil, &error2.InvalidParamsError{Message: "execution errors: Call Too Deep"}
	}
	log.Log().Debug().Msgf("unhandled transaction status: %d", ts.Enum)
	return nil, errors.New("execution errors: unhandled transaction status")
}

// ToResponse processes the engine query response (`TransactionStatus`) and returns output buffer or errors
func (ts *TransactionStatus) ToResponse() (*string, error) {
	buf, err := ts.Validate()
	if buf != nil {
		str := "0x"
		for _, b := range buf {
			tmp := fmt.Sprint(int(b))
			str = str + tmp
		}
		return &str, nil
	}
	return nil, err
}

// SubmitStatus is the type received from engine for submit (eg: sendRawTransactionSync) calls
type SubmitStatus struct {
	StatusMap    map[string]interface{}
	SubmitResult *SubmitResultV2
	ResponseHash string
}

// NewSubmitStatus allocates and returns a new SubmitStatus object
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

// Validate checks `SubmitStatus.StatusMap` to return the success or errors
func (ss *SubmitStatus) Validate() error {
	logger := log.Log()
	// Check if any errors returned
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
				logger.Error().Msgf("submit request failed with errors: %s, txs hash: %s", jsnStr, ss.ResponseHash)
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
			logger.Error().Msgf("submit request failed with Default case errors: %s, txs hash: %s", jsnStr, ss.ResponseHash)
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
	logger.Error().Msgf("submit request returned with unhandled errors, txs hash: %s", ss.ResponseHash)
	return errors.New("submit request returned with unhandled errors")

}

// ToResponse processes the engine query response, retrieves the `result` map and returns the hash
func (ss *SubmitStatus) ToResponse() (*string, error) {
	err := ss.Validate()
	// Validate can generate either `utils.InvalidParams` or `errors.New` errors
	if err != nil {
		_, ok := err.(*error2.InvalidParamsError)
		if ok {
			return nil, err
		} else {
			return nil, &error2.GenericError{Err: err}
		}
	}
	return &(ss.ResponseHash), nil
}
