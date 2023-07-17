package engine

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/types/common"
	error2 "github.com/aurora-is-near/relayer2-base/types/errors"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/near/borsh-go"
)

const (
	addrLength   = 20
	raw256Length = 32

	// Engine TxsStatus errors
	errStackOverflow = "ERR_STACK_OVERFLOW"
	errUnreachable   = "FunctionCallError(WasmTrap(Unreachable))"
)

// ArgsForGetStorageAt is used to process GetStorageAt endpoint arguments
type ArgsForGetStorageAt struct {
	address     primitives.Data20
	storageSlot primitives.Quantity
}

// NewArgsForGetStorageAt allocates and returns a new empty ArgsForGetStorageAt
func NewArgsForGetStorageAt(addr common.Address, sSlot common.Uint256) *ArgsForGetStorageAt {
	return &ArgsForGetStorageAt{
		address:     primitives.Data20FromBytes(addr.Bytes()),
		storageSlot: primitives.QuantityFromBytes(sSlot.Bytes()),
	}
}

// Serialize transforms ArgsForGetStorageAt to argsForGetStorageAtEngine, calls its Serialize method
// and returns the received buffer
func (gs *ArgsForGetStorageAt) Serialize() ([]byte, error) {
	args := argsForGetStorageAtEngine{}
	copy(args.Address[:], gs.address.Bytes())
	copy(args.Key[:], gs.storageSlot.Bytes())
	buff, err := args.serialize()
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// argsForGetStorageAtEngine is the data format accepted by engine for GetStorageAt endpoint
type argsForGetStorageAtEngine struct {
	Address [addrLength]byte
	Key     [raw256Length]byte
}

// serialize argsForGetStorageAtEngine to a buffer using borsh so to communicate with engine
func (gse argsForGetStorageAtEngine) serialize() ([]byte, error) {
	buff, err := borsh.Serialize(gse)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// TransactionForCall is the type used to serialize eth_call input
type TransactionForCall struct {
	From     *primitives.Data20   `json:"from,omitempty"`
	To       *primitives.Data20   `json:"to,omitempty"`
	Gas      *primitives.Quantity `json:"gas,omitempty"`
	GasPrice *primitives.Quantity `json:"gasPrice,omitempty"`
	Value    *primitives.Quantity `json:"value,omitempty"`
	Data     *primitives.VarData  `json:"data,omitempty"`
}

// UnmarshalJSON implements jsoniter.Unmarshaler
// This method is needed to make `.To` field of the TransactionForCall struct required/mandatory
func (tc *TransactionForCall) UnmarshalJSON(data []byte) error {
	type tmpType TransactionForCall
	tmp := tmpType{}
	err := jsoniter.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	*tc = TransactionForCall(tmp)
	return nil
}

// Serialize transforms TransactionForCall to transactionForCallEngine, calls its Serialize method
// and returns the received buffer
func (tc *TransactionForCall) Serialize() ([]byte, error) {

	tce := transactionForCallEngine{}

	if tc.To == nil {
		copy(tce.To[:], primitives.Data20FromHex("0x0").Bytes())
	} else {
		copy(tce.To[:], tc.To.Bytes())
	}

	if tc.From == nil {
		copy(tce.From[:], primitives.Data20FromHex("0x0").Bytes())
	} else {
		copy(tce.From[:], tc.From.Bytes())
	}

	if tc.Value != nil {
		copy(tce.Value[:], tc.Value.Bytes())
	}
	if tc.Data != nil {
		tce.Data = make([]byte, len(tc.Data.Bytes()))
		copy(tce.Data[:], tc.Data.Bytes())
	}

	buff, err := tce.Serialize()
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// Validate checks the incoming object and returns error for incorrect fields
func (tc *TransactionForCall) Validate() error {
	if tc.Gas == nil || len(tc.Gas.Bytes()) == 0 {
		return nil
	}
	// return OutOfGas error if Gas field is "0x0"
	if tc.Gas.IsZero() {
		return &error2.TxsStatusError{Message: "Ok(OutOfGas)"}
	}

	return nil
}

// transactionForCallEngine is the type send to engine for eth_call endpoint
type transactionForCallEngine struct {
	From  [addrLength]byte
	To    [addrLength]byte
	Value [raw256Length]byte
	Data  []byte
}

// Serialize transactionForCallEngine to a buffer using borsh so to communicate with engine
func (tce transactionForCallEngine) Serialize() ([]byte, error) {
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
	result, ok := resp.(map[string]interface{})["result"].([]interface{})
	if !ok {
		log.Log().Error().Msgf("query response is not in correct format: %s", result)
		return nil, errors.New("query response is not in correct format")
	}
	return &QueryResult{Result: result}, nil
}

// ToUint256Response processes the engine query response, retrieves the `result` map and converts it to Uint256 response
func (r *QueryResult) ToUint256Response() (*common.Uint256, error) {
	buf := r.resultToByteBuffer()
	ui256 := common.Uint256FromBytes(buf)
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
	Address [addrLength]uint8
	Topics  []RawU256
	Data    []uint8
}

// RawU256 is the type used to handle engine's LogEventWithAddress response
type RawU256 struct {
	Value [raw256Length]uint8
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
	resp, ok := respArg.(map[string]interface{})["result"].([]interface{})
	if !ok {
		err, ok := respArg.(map[string]interface{})["error"].(string)
		if !ok {
			return nil, errors.New("call response is not in correct format")
		}
		log.Log().Error().Msgf("errors returned to eth_call: %v", err)
		// Check for specific TxsStatus errors
		if strings.Contains(err, errStackOverflow) {
			return nil, &error2.TxsStatusError{Message: "EvmError(StackOverflow)"}
		} else if strings.Contains(err, errUnreachable) {
			return nil, &error2.TxsStatusError{Message: "WasmTrap(Unreachable)"}
		}
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
			rReason, err := utils.ParseEVMRevertReason(ts.Revert.Output)
			rOutputStr := "0x" + hex.EncodeToString(ts.Revert.Output)
			if err != nil {
				log.Log().Error().Msgf("execution reverted with data 0x%s, got err %s", hex.EncodeToString(ts.Revert.Output), err.Error())
				return nil, &error2.TxsStatusError{Message: fmt.Sprintf("execution reverted. error thrown while parsing revert msg %s", err.Error())}
			}
			return nil, &error2.TxsRevertError{
				Code:    3,
				Message: "execution reverted: " + rReason,
				Data:    rOutputStr,
			}
		} else {
			return nil, &error2.TxsRevertError{
				Code:    3,
				Message: "execution reverted",
			}
		}
	case 2: // OutOfGas case
		return nil, &error2.TxsStatusError{Message: "Ok(OutOfGas)"}
	case 3: // OutOfFund case
		return nil, &error2.TxsStatusError{Message: "Ok(OutOfFund)"}
	case 4: // OutOfOffset case
		return nil, &error2.TxsStatusError{Message: "Ok(OutOfOffset)"}
	case 5: // CallTooDeep case
		return nil, &error2.InvalidParamsError{Message: "Call Too Deep)"}
	}
	log.Log().Debug().Msgf("unhandled transaction status: %d", ts.Enum)
	return nil, errors.New("unhandled transaction status")
}

// ToResponse processes the engine query response (`TransactionStatus`) and returns output buffer or errors
func (ts *TransactionStatus) ToResponse() (*string, error) {
	buf, err := ts.Validate()
	if err != nil {
		return nil, err
	}
	str := "0x" + hex.EncodeToString(buf)
	return &str, nil
}

// SubmitStatus is the type received from engine for submit (eg: sendRawTransactionSync) calls
type SubmitStatus struct {
	StatusMap    map[string]interface{}
	SubmitResult *SubmitResultV2
	ResponseHash string
}

// NewSubmitStatus allocates and returns a new SubmitStatus object
func NewSubmitStatus(respArg interface{}, txsHash string) (*SubmitStatus, error) {
	resp, ok := respArg.(map[string]interface{})
	if !ok {
		log.Log().Error().Msgf("submit response is not in correct format: %s", respArg)
		return nil, errors.New("submit response is not in correct format")
	}

	status, ok := resp["status"].(map[string]interface{})
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
		failTypeMap, ok := fail.(map[string]interface{})["ActionError"].(map[string]interface{})["kind"].(map[string]interface{})
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
			execErr, ok := failTypeMap["FunctionCallError"].(map[string]interface{})["ExecutionError"].(string)
			if ok {
				errMsg := strings.Replace(execErr, "Smart contract panicked: ", "", 1)
				logger.Debug().Msgf("submit request failed with ExecutionError: %s, txs hash: %s", errMsg, ss.ResponseHash)
				return errors.New(errMsg)
			} else {
				jsn, err := jsoniter.Marshal(failTypeMap["FunctionCallError"])
				if err != nil {
					logger.Error().Msgf("submit request failed while marshalling FunctionCallError: %s, txs hash: %s", err.Error(), ss.ResponseHash)
					return err
				}
				jsnStr := reg.ReplaceAllString(string(jsn), "")
				logger.Error().Msgf("submit request failed with errors: %s, txs hash: %s", jsnStr, ss.ResponseHash)
				return fmt.Errorf("failure:%s", jsnStr)
			}
		case "MethodNotFound":
			jsn, err := jsoniter.Marshal(failTypeMap["MethodNotFound"])
			if err != nil {
				logger.Error().Msgf("submit request failed while marshalling MethodNotFound: %s, txs hash: %s", err.Error(), ss.ResponseHash)
				return err
			}
			jsnStr := reg.ReplaceAllString(string(jsn), "")
			logger.Error().Msgf("submit request failed with MethodNotFound: %s, txs hash: %s", jsnStr, ss.ResponseHash)
			return fmt.Errorf("failure: %s", jsnStr)
		default:
			jsn, err := jsoniter.Marshal(fail.(map[string]interface{})["ActionError"].(map[string]interface{})["kind"])
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
			sucBuf, err := base64.StdEncoding.DecodeString(sucStrB64)
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
	// Validate can generate either `errors.InvalidParamError` or `errors.GenericError` errors
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
