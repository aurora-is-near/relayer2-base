package codec

import (
	"aurora-relayer-go-common/db/types"
	"aurora-relayer-go-common/utils"
	"bytes"
	"capnproto.org/go/capnp/v3"
	"encoding/binary"
	"fmt"
)

type CapnpCodec struct {
	Encoder
	Decoder
}

type capnpEncoder struct{}
type capnpDecoder struct{}

func NewCapnpCodec() CapnpCodec {
	return CapnpCodec{
		Encoder: capnpEncoder{},
		Decoder: capnpDecoder{},
	}

}

func (e capnpEncoder) Marshal(v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case utils.Block:
		return e.encodeBlock(&val)
	case utils.LogResponse:
		return e.encodeLog(&val)
	case *utils.StoredFilter:
		return e.encodeFilter(val)
	case utils.Transaction:
		return e.encodeTransaction(&val)
	case utils.Uint256:
		return e.encodeUint256(val)
	case int64:
		return e.encodeUint64(uint64(val))
	case uint64:
		return e.encodeUint64(val)
	case []byte:
		return val, nil
	default:
		return nil, fmt.Errorf("unable to marshal unknown type: %T", v)
	}
}

func (d capnpDecoder) Unmarshal(b []byte, v interface{}) error {
	switch val := v.(type) {
	case *utils.Block:
		return d.decodeBlock(b, val)
	case *utils.LogResponse:
		return d.decodeLog(b, val)
	case *utils.StoredFilter:
		return d.decodeFilter(b, val)
	case *utils.Transaction:
		return d.decodeTransaction(b, val)
	case *utils.Uint256:
		return d.decodeUint256(b, val)
	case *int64:
		return d.decodeInt64(b, val)
	case *uint64:
		return d.decodeUint64(b, val)
	case *[]byte:
		*val = b
		return nil
	default:
		return fmt.Errorf("unable to unmarshal unknown type: %T", v)
	}
}

func (e capnpEncoder) encodeBlock(block *utils.Block) ([]byte, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	cBlock, err := types.NewRootBlock(seg)
	if err != nil {
		return nil, err
	}
	cBlock.SetHash(block.Hash.Bytes())
	cBlock.SetParentHash(block.ParentHash.Bytes())
	cBlock.SetHeight(block.Height)
	cBlock.SetMiner(block.Miner.Bytes())
	cBlock.SetTimestamp(block.Timestamp)
	cBlock.SetGasLimit(block.GasLimit.Bytes())
	cBlock.SetGasUsed(block.GasUsed.Bytes())
	cBlock.SetLogsBloom([]byte(block.LogsBloom))
	cBlock.SetTransactionsRoot(block.TransactionsRoot.Bytes())
	cBlock.SetReceiptsRoot(block.ReceiptsRoot.Bytes())
	cBlock.SetStateRoot([]byte(block.StateRoot))
	cBlock.SetSize(block.Size.Bytes())
	return msg.MarshalPacked()
}

func (e capnpEncoder) encodeTransaction(tx *utils.Transaction) ([]byte, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	cTx, err := types.NewRootTransaction(seg)
	if err != nil {
		return nil, err
	}
	cTx.SetHash(tx.Hash.Bytes())
	cTx.SetBlockHash(tx.BlockHash.Bytes())
	cTx.SetBlockHeight(tx.BlockHeight)
	cTx.SetTransactionIndex(tx.TransactionIndex)
	cTx.SetFrom(tx.From.Bytes())
	if tx.To != nil {
		cTx.SetTo(tx.To.Bytes())
	}
	cTx.SetNonce(tx.Nonce.Bytes())
	cTx.SetGasPrice(tx.GasPrice.Bytes())
	cTx.SetGasLimit(tx.GasLimit.Bytes())
	cTx.SetGasUsed(tx.GasUsed.Bytes())
	cTx.SetValue(tx.Value.Bytes())
	cTx.SetInput(tx.Input)
	cTx.SetOutput(tx.Output)
	cTx.SetStatus(tx.Status)
	cTx.SetContractAddress(tx.ContractAddress.Bytes())
	cTx.SetV(tx.V)
	cTx.SetR(tx.R.Bytes())
	cTx.SetS(tx.S.Bytes())
	cTx.SetNearHash(tx.NearTransaction.Hash.Bytes())
	cTx.SetNearReceiptHash(tx.NearTransaction.ReceiptHash.Bytes())
	return msg.MarshalPacked()
}

func (e capnpEncoder) encodeLog(l *utils.LogResponse) ([]byte, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	cLog, err := types.NewRootLog(seg)
	if err != nil {
		return nil, err
	}
	cLog.SetRemoved(l.Removed)
	cLog.SetLogIndex(l.LogIndex.Bytes())
	cLog.SetTransactionIndex(l.TransactionIndex.Bytes())
	cLog.SetTransactionHash(l.TransactionHash.Bytes())
	cLog.SetBlockHash(l.BlockHash.Bytes())
	cLog.SetBlockNumber(l.BlockNumber.Bytes())
	cLog.SetAddress(l.Address.Bytes())
	cLog.SetData(l.Data)
	if l.Topics != nil {
		list, err := cLog.NewTopics(int32(len(l.Topics)))
		if err != nil {
			return nil, err
		}
		for i, t := range l.Topics {
			list.Set(i, t)
		}
	}

	return msg.MarshalPacked()
}

func (e capnpEncoder) encodeFilter(fl *utils.StoredFilter) ([]byte, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	cFl, err := types.NewRootFilter(seg)
	if err != nil {
		return nil, err
	}
	cFl.SetType(fl.Type)
	cFl.SetCreatedBy(fl.CreatedBy)
	cFl.SetPollBlock(fl.PollBlock.Bytes())
	if fl.FromBlock != nil {
		cFl.SetFromBlock(fl.FromBlock.Bytes())
	}
	if fl.ToBlock != nil {
		cFl.SetToBlock(fl.ToBlock.Bytes())
	}
	if fl.Addresses != nil {
		list, err := cFl.NewAddresses(int32(len(fl.Addresses)))
		if err != nil {
			return nil, err
		}
		for i, addr := range fl.Addresses {
			list.Set(i, addr)
		}

	}
	if fl.Topics != nil {
		topicLists, err := cFl.NewTopics(int32(len(fl.Topics)))
		if err != nil {
			return nil, err
		}
		for i, ts := range fl.Topics {
			if len(ts) == 0 {
				topicLists.Set(i, capnp.Ptr{})
				continue
			}
			list, err := capnp.NewDataList(seg, int32(len(ts)))
			if err != nil {
				return nil, err
			}
			for j, t := range ts {
				list.Set(j, t)
			}
			topicLists.Set(i, list.ToPtr())
		}
	}
	return msg.MarshalPacked()
}

func (e capnpEncoder) encodeUint64(i uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes(), nil
}

func (e capnpEncoder) encodeUint256(u utils.Uint256) ([]byte, error) {
	return u.Bytes(), nil
}

func (d capnpDecoder) decodeBlock(b []byte, block *utils.Block) error {
	msg, err := capnp.UnmarshalPacked(b)
	if err != nil {
		return err
	}
	cBlock, err := types.ReadRootBlock(msg)
	if err != nil {
		return err
	}

	if cBlock.HasHash() {
		b, _ := cBlock.Hash()
		block.Hash.SetBytes(b)
	}
	if cBlock.HasParentHash() {
		b, _ := cBlock.ParentHash()
		block.ParentHash.SetBytes(b)
	}
	block.Height = cBlock.Height()
	if cBlock.HasMiner() {
		b, _ := cBlock.Miner()
		block.Miner.SetBytes(b)
	}
	block.Timestamp = cBlock.Timestamp()
	if cBlock.HasGasLimit() {
		b, _ := cBlock.GasLimit()
		block.GasLimit.SetBytes(b)
	}
	if cBlock.HasGasUsed() {
		b, _ := cBlock.GasUsed()
		block.GasUsed.SetBytes(b)
	}
	if cBlock.HasLogsBloom() {
		b, _ := cBlock.LogsBloom()
		block.LogsBloom = string(b)
	}
	if cBlock.HasTransactionsRoot() {
		b, _ := cBlock.TransactionsRoot()
		block.TransactionsRoot.SetBytes(b)
	}
	if cBlock.HasStateRoot() {
		b, _ := cBlock.StateRoot()
		block.StateRoot = string(b)
	}
	if cBlock.HasSize() {
		b, _ := cBlock.StateRoot()
		block.Size.SetBytes(b)
	}
	return nil
}

func (d capnpDecoder) decodeTransaction(b []byte, tx *utils.Transaction) error {
	msg, err := capnp.UnmarshalPacked(b)
	if err != nil {
		return err
	}
	cTx, err := types.ReadRootTransaction(msg)
	if err != nil {
		return err
	}

	if cTx.HasHash() {
		b, _ := cTx.Hash()
		tx.Hash.SetBytes(b)
	}
	if cTx.HasBlockHash() {
		b, _ := cTx.BlockHash()
		tx.BlockHash.SetBytes(b)
	}
	tx.BlockHeight = cTx.BlockHeight()
	tx.TransactionIndex = cTx.TransactionIndex()
	if cTx.HasFrom() {
		b, _ := cTx.From()
		tx.From.SetBytes(b)
	}
	if cTx.HasTo() {
		b, _ := cTx.To()
		tx.To = new(utils.Address)
		tx.To.SetBytes(b)
	}
	if cTx.HasNonce() {
		b, _ := cTx.Nonce()
		tx.Nonce.SetBytes(b)
	}
	if cTx.HasGasPrice() {
		b, _ := cTx.GasPrice()
		tx.GasPrice.SetBytes(b)
	}
	if cTx.HasGasLimit() {
		b, _ := cTx.GasLimit()
		tx.GasLimit.SetBytes(b)
	}
	if cTx.HasGasUsed() {
		b, _ := cTx.GasUsed()
		tx.GasUsed.SetBytes(b)
	}
	if cTx.HasValue() {
		b, _ := cTx.Value()
		tx.Value.SetBytes(b)
	}
	if cTx.HasInput() {
		b, _ := cTx.Input()
		tx.Input = b
	}
	if cTx.HasOutput() {
		b, _ := cTx.Output()
		tx.Output = b
	}
	tx.Status = cTx.Status()
	if cTx.HasContractAddress() {
		b, _ := cTx.ContractAddress()
		tx.ContractAddress.SetBytes(b)
	}
	tx.V = cTx.V()
	if cTx.HasR() {
		b, _ := cTx.R()
		tx.R.SetBytes(b)
	}
	if cTx.HasS() {
		b, _ := cTx.S()
		tx.S.SetBytes(b)
	}
	tx.NearTransaction = utils.NearTransaction{}
	if cTx.HasNearHash() {
		b, _ := cTx.NearHash()
		tx.NearTransaction.Hash.SetBytes(b)
	}
	if cTx.HasNearReceiptHash() {
		b, _ := cTx.NearReceiptHash()
		tx.NearTransaction.ReceiptHash.SetBytes(b)
	}
	return nil
}

func (d capnpDecoder) decodeLog(b []byte, l *utils.LogResponse) error {
	msg, err := capnp.UnmarshalPacked(b)
	if err != nil {
		return err
	}
	cLog, err := types.ReadRootLog(msg)
	if err != nil {
		return err
	}

	l.Removed = cLog.Removed()
	if cLog.HasLogIndex() {
		b, _ := cLog.LogIndex()
		l.LogIndex.SetBytes(b)
	}
	if cLog.HasTransactionIndex() {
		b, _ := cLog.TransactionIndex()
		l.TransactionIndex.SetBytes(b)
	}
	if cLog.HasTransactionHash() {
		b, _ := cLog.TransactionHash()
		l.TransactionHash.SetBytes(b)
	}
	if cLog.HasBlockHash() {
		b, _ := cLog.BlockHash()
		l.BlockHash.SetBytes(b)
	}
	if cLog.HasBlockNumber() {
		b, _ := cLog.BlockNumber()
		l.BlockNumber.SetBytes(b)
	}
	if cLog.HasAddress() {
		b, _ := cLog.Address()
		l.Address.SetBytes(b)
	}
	if cLog.HasData() {
		b, _ := cLog.Data()
		l.Data = b
	}
	if cLog.HasTopics() {
		ts, _ := cLog.Topics()
		l.Topics = make([]utils.Bytea, ts.Len())
		for i := range l.Topics {
			t, _ := ts.At(i)
			l.Topics[i] = t
		}
	}
	return nil
}

func (d capnpDecoder) decodeFilter(b []byte, fl *utils.StoredFilter) error {
	msg, err := capnp.UnmarshalPacked(b)
	if err != nil {
		return err
	}
	cFl, err := types.ReadRootFilter(msg)
	if err != nil {
		return err
	}
	if cFl.HasType() {
		s, _ := cFl.Type()
		fl.Type = s
	}
	if cFl.HasCreatedBy() {
		s, _ := cFl.CreatedBy()
		fl.CreatedBy = s
	}
	if cFl.HasPollBlock() {
		b, _ := cFl.PollBlock()
		fl.PollBlock.SetBytes(b)
	}
	if cFl.HasFromBlock() {
		b, _ := cFl.FromBlock()
		from := utils.IntToUint256(0)
		from.SetBytes(b)
		fl.FromBlock = &from
	}
	if cFl.HasToBlock() {
		b, _ := cFl.ToBlock()
		to := utils.IntToUint256(0)
		to.SetBytes(b)
		fl.ToBlock = &to
	}
	if cFl.HasAddresses() {
		addrs, _ := cFl.Addresses()
		fl.Addresses = make([][]byte, addrs.Len())
		for i := range fl.Topics {
			a, _ := addrs.At(i)
			fl.Addresses[i] = a
		}
	}
	if cFl.HasTopics() {
		topicLists, _ := cFl.Topics()
		fl.Topics = make([][][]byte, topicLists.Len())
		for i := range fl.Topics {
			topics, err := topicLists.At(i)
			if err != nil {
				continue
			}
			list, _ := capnp.NewDataList(cFl.Segment(), int32(topics.List().Len()))
			list = list.DecodeFromPtr(topics)
			fl.Topics[i] = make([][]byte, list.Len())

			for j := range fl.Topics[i] {
				topic, _ := list.At(j)
				fl.Topics[i][j] = topic
			}
		}
	}
	return nil
}

func (d capnpDecoder) decodeInt64(b []byte, i *int64) error {
	rd := bytes.NewReader(b)
	return binary.Read(rd, binary.BigEndian, i)
}

func (d capnpDecoder) decodeUint64(b []byte, i *uint64) error {
	rd := bytes.NewReader(b)
	return binary.Read(rd, binary.BigEndian, i)
}

func (d capnpDecoder) decodeUint256(b []byte, u *utils.Uint256) error {
	u.SetBytes(b)
	return nil
}
