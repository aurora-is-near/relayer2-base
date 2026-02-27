package codec

import (
	"fmt"

	"github.com/aurora-is-near/relayer2-base/db/codec/dbpb"
	"github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"google.golang.org/protobuf/proto"
)

type ProtobufCodec struct{}

func NewProtobufCodec() Codec {
	return ProtobufCodec{}
}

func (c ProtobufCodec) Marshal(v any) ([]byte, error) {
	msg, err := toProto(v)
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("protobuf marshal: %w", err)
	}
	return data, nil
}

func (c ProtobufCodec) Unmarshal(data []byte, v any) error {
	switch val := v.(type) {
	case *db.Block:
		pb := &dbpb.Block{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		blockFromProto(pb, val)
		return nil
	case *db.Transaction:
		pb := &dbpb.Transaction{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		transactionFromProto(pb, val)
		return nil
	case *db.Log:
		pb := &dbpb.Log{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		logFromProto(pb, val)
		return nil
	case *db.BlockKey:
		pb := &dbpb.BlockKey{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		val.Height = pb.Height
		return nil
	case *db.TransactionKey:
		pb := &dbpb.TransactionKey{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		val.BlockHeight = pb.BlockHeight
		val.TransactionIndex = pb.TransactionIndex
		return nil
	case *db.LogKey:
		pb := &dbpb.LogKey{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		val.BlockHeight = pb.BlockHeight
		val.TransactionIndex = pb.TransactionIndex
		val.LogIndex = pb.LogIndex
		return nil
	case *db.BlockFilter:
		pb := &dbpb.BlockFilter{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		blockFilterFromProto(pb, val)
		return nil
	case *db.TransactionFilter:
		pb := &dbpb.TransactionFilter{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		transactionFilterFromProto(pb, val)
		return nil
	case *db.LogFilter:
		pb := &dbpb.LogFilter{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		logFilterFromProto(pb, val)
		return nil
	case *primitives.Data32:
		pb := &dbpb.Data{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		*val = data32(pb.Content)
		return nil
	case *primitives.Data20:
		pb := &dbpb.Data{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		*val = data20(pb.Content)
		return nil
	case *primitives.VarData:
		pb := &dbpb.Data{}
		if err := proto.Unmarshal(data, pb); err != nil {
			return fmt.Errorf("protobuf unmarshal: %w", err)
		}
		*val = primitives.VarDataFromBytes(pb.Content)
		return nil
	default:
		return fmt.Errorf("protobuf codec: unsupported type for unmarshal: %T", v)
	}
}

func toProto(v any) (proto.Message, error) {
	switch val := v.(type) {
	case *db.Block:
		return blockToProto(val), nil
	case *db.Transaction:
		return transactionToProto(val), nil
	case *db.Log:
		return logToProto(val), nil
	case *db.BlockKey:
		return &dbpb.BlockKey{Height: val.Height}, nil
	case *db.TransactionKey:
		return &dbpb.TransactionKey{
			BlockHeight:      val.BlockHeight,
			TransactionIndex: val.TransactionIndex,
		}, nil
	case *db.LogKey:
		return &dbpb.LogKey{
			BlockHeight:      val.BlockHeight,
			TransactionIndex: val.TransactionIndex,
			LogIndex:         val.LogIndex,
		}, nil
	case *db.BlockFilter:
		return blockFilterToProto(val), nil
	case *db.TransactionFilter:
		return transactionFilterToProto(val), nil
	case *db.LogFilter:
		return logFilterToProto(val), nil
	case *primitives.Data32:
		return &dbpb.Data{Content: val.Content}, nil
	case *primitives.Data20:
		return &dbpb.Data{Content: val.Content}, nil
	case *primitives.VarData:
		return &dbpb.Data{Content: val.Content}, nil
	default:
		return nil, fmt.Errorf("protobuf codec: unsupported type for marshal: %T", v)
	}
}

// Nil-preserving helpers: proto3 cannot distinguish nil from empty bytes,
// so we preserve nil to avoid padding zero-value fields.

func quantity(b []byte) primitives.Quantity {
	if b == nil {
		return primitives.Quantity{}
	}
	return primitives.QuantityFromBytes(b)
}

func data20(b []byte) primitives.Data20 {
	if b == nil {
		return primitives.Data20{}
	}
	return primitives.Data20FromBytes(b)
}

func data32(b []byte) primitives.Data32 {
	if b == nil {
		return primitives.Data32{}
	}
	return primitives.Data32FromBytes(b)
}

func data256(b []byte) primitives.Data256 {
	if b == nil {
		return primitives.Data256{}
	}
	return primitives.Data256FromBytes(b)
}

func blockToProto(b *db.Block) *dbpb.Block {
	return &dbpb.Block{
		ParentHash:       b.ParentHash.Content,
		Miner:            b.Miner.Content,
		Timestamp:        b.Timestamp,
		GasLimit:         b.GasLimit.Content,
		GasUsed:          b.GasUsed.Content,
		LogsBloom:        b.LogsBloom.Content,
		TransactionsRoot: b.TransactionsRoot.Content,
		StateRoot:        b.StateRoot.Content,
		ReceiptsRoot:     b.ReceiptsRoot.Content,
		Size:             b.Size,
	}
}

func blockFromProto(pb *dbpb.Block, b *db.Block) {
	b.ParentHash = data32(pb.ParentHash)
	b.Miner = data20(pb.Miner)
	b.Timestamp = pb.Timestamp
	b.GasLimit = quantity(pb.GasLimit)
	b.GasUsed = quantity(pb.GasUsed)
	b.LogsBloom = data256(pb.LogsBloom)
	b.TransactionsRoot = data32(pb.TransactionsRoot)
	b.StateRoot = data32(pb.StateRoot)
	b.ReceiptsRoot = data32(pb.ReceiptsRoot)
	b.Size = pb.Size
}

func transactionToProto(tx *db.Transaction) *dbpb.Transaction {
	pb := &dbpb.Transaction{
		Type:                 tx.Type,
		From:                 tx.From.Content,
		IsContractDeployment: tx.IsContractDeployment,
		Nonce:                tx.Nonce.Content,
		GasPrice:             tx.GasPrice.Content,
		GasLimit:             tx.GasLimit.Content,
		GasUsed:              tx.GasUsed,
		CumulativeGasUsed:    tx.CumulativeGasUsed.Content,
		Value:                tx.Value.Content,
		Input:                tx.Input.Content,
		NearReceiptHash:      tx.NearReceiptHash.Content,
		Status:               tx.Status,
		V:                    tx.V,
		R:                    tx.R.Content,
		S:                    tx.S.Content,
		LogsBloom:            tx.LogsBloom.Content,
		MaxPriorityFeePerGas: tx.MaxPriorityFeePerGas.Content,
		MaxFeePerGas:         tx.MaxFeePerGas.Content,
	}
	if tx.ToOrContract != nil {
		pb.ToOrContract = tx.ToOrContract.Content
	}
	if tx.NearHash != nil {
		pb.NearHash = tx.NearHash.Content
	}
	pb.AccessList = make([]*dbpb.AccessListEntry, len(tx.AccessList))
	for i, ale := range tx.AccessList {
		pbAle := &dbpb.AccessListEntry{
			Address: ale.Address.Content,
		}
		pbAle.StorageKeys = make([][]byte, len(ale.StorageKeys))
		for j, sk := range ale.StorageKeys {
			pbAle.StorageKeys[j] = sk.Content
		}
		pb.AccessList[i] = pbAle
	}
	pb.AuthorizationList = make([]*dbpb.Authorization, len(tx.AuthorizationList))
	for i, auth := range tx.AuthorizationList {
		pb.AuthorizationList[i] = &dbpb.Authorization{
			ChainId: auth.ChainId,
			Address: auth.Address.Content,
			Nonce:   auth.Nonce,
			YParity: uint32(auth.YParity),
			R:       auth.R.Content,
			S:       auth.S.Content,
		}
	}
	if tx.YParity != nil {
		v := uint32(*tx.YParity)
		pb.YParity = &v
	}
	return pb
}

func transactionFromProto(pb *dbpb.Transaction, tx *db.Transaction) {
	tx.Type = pb.Type
	tx.From = data20(pb.From)
	tx.IsContractDeployment = pb.IsContractDeployment
	tx.Nonce = quantity(pb.Nonce)
	tx.GasPrice = quantity(pb.GasPrice)
	tx.GasLimit = quantity(pb.GasLimit)
	tx.GasUsed = pb.GasUsed
	tx.CumulativeGasUsed = quantity(pb.CumulativeGasUsed)
	tx.Value = quantity(pb.Value)
	tx.Input = primitives.VarDataFromBytes(pb.Input)
	tx.NearReceiptHash = data32(pb.NearReceiptHash)
	tx.Status = pb.Status
	tx.V = pb.V
	tx.R = quantity(pb.R)
	tx.S = quantity(pb.S)
	tx.LogsBloom = data256(pb.LogsBloom)
	tx.MaxPriorityFeePerGas = quantity(pb.MaxPriorityFeePerGas)
	tx.MaxFeePerGas = quantity(pb.MaxFeePerGas)

	if pb.ToOrContract != nil {
		addr := data20(pb.ToOrContract)
		tx.ToOrContract = &addr
	}
	if pb.NearHash != nil {
		hash := data32(pb.NearHash)
		tx.NearHash = &hash
	}

	if pb.AccessList != nil {
		tx.AccessList = make([]db.AccessListEntry, len(pb.AccessList))
		for i, pbAle := range pb.AccessList {
			ale := db.AccessListEntry{
				Address: data20(pbAle.Address),
			}
			ale.StorageKeys = make([]primitives.Data32, len(pbAle.StorageKeys))
			for j, sk := range pbAle.StorageKeys {
				ale.StorageKeys[j] = data32(sk)
			}
			tx.AccessList[i] = ale
		}
	}

	if pb.AuthorizationList != nil {
		tx.AuthorizationList = make([]db.Authorization, len(pb.AuthorizationList))
		for i, pbAuth := range pb.AuthorizationList {
			tx.AuthorizationList[i] = db.Authorization{
				ChainId: pbAuth.ChainId,
				Address: data20(pbAuth.Address),
				Nonce:   pbAuth.Nonce,
				YParity: uint8(pbAuth.YParity),
				R:       quantity(pbAuth.R),
				S:       quantity(pbAuth.S),
			}
		}
	}
	if pb.YParity != nil {
		v := uint8(*pb.YParity)
		tx.YParity = &v
	}
}

func logToProto(l *db.Log) *dbpb.Log {
	pb := &dbpb.Log{
		Address: l.Address.Content,
		Data:    l.Data.Content,
	}
	pb.Topics = make([][]byte, len(l.Topics))
	for i, t := range l.Topics {
		pb.Topics[i] = t.Content
	}
	return pb
}

func logFromProto(pb *dbpb.Log, l *db.Log) {
	l.Address = data20(pb.Address)
	l.Data = primitives.VarDataFromBytes(pb.Data)
	l.Topics = make([]primitives.Data32, len(pb.Topics))
	for i, t := range pb.Topics {
		l.Topics[i] = data32(t)
	}
}

func blockFilterToProto(f *db.BlockFilter) *dbpb.BlockFilter {
	return &dbpb.BlockFilter{
		CreatedAt: f.CreatedAt,
		Metadata:  f.Metadata.Content,
		From:      &dbpb.BlockKey{Height: f.From.Height},
		Next:      &dbpb.BlockKey{Height: f.Next.Height},
		To:        &dbpb.BlockKey{Height: f.To.Height},
	}
}

func blockFilterFromProto(pb *dbpb.BlockFilter, f *db.BlockFilter) {
	f.CreatedAt = pb.CreatedAt
	f.Metadata = primitives.VarDataFromBytes(pb.Metadata)
	if pb.From != nil {
		f.From = db.BlockKey{Height: pb.From.Height}
	}
	if pb.Next != nil {
		f.Next = db.BlockKey{Height: pb.Next.Height}
	}
	if pb.To != nil {
		f.To = db.BlockKey{Height: pb.To.Height}
	}
}

func transactionFilterToProto(f *db.TransactionFilter) *dbpb.TransactionFilter {
	return &dbpb.TransactionFilter{
		CreatedAt: f.CreatedAt,
		Metadata:  f.Metadata.Content,
		From: &dbpb.TransactionKey{
			BlockHeight:      f.From.BlockHeight,
			TransactionIndex: f.From.TransactionIndex,
		},
		Next: &dbpb.TransactionKey{
			BlockHeight:      f.Next.BlockHeight,
			TransactionIndex: f.Next.TransactionIndex,
		},
		To: &dbpb.TransactionKey{
			BlockHeight:      f.To.BlockHeight,
			TransactionIndex: f.To.TransactionIndex,
		},
	}
}

func transactionFilterFromProto(pb *dbpb.TransactionFilter, f *db.TransactionFilter) {
	f.CreatedAt = pb.CreatedAt
	f.Metadata = primitives.VarDataFromBytes(pb.Metadata)
	if pb.From != nil {
		f.From = db.TransactionKey{
			BlockHeight:      pb.From.BlockHeight,
			TransactionIndex: pb.From.TransactionIndex,
		}
	}
	if pb.Next != nil {
		f.Next = db.TransactionKey{
			BlockHeight:      pb.Next.BlockHeight,
			TransactionIndex: pb.Next.TransactionIndex,
		}
	}
	if pb.To != nil {
		f.To = db.TransactionKey{
			BlockHeight:      pb.To.BlockHeight,
			TransactionIndex: pb.To.TransactionIndex,
		}
	}
}

func logFilterToProto(f *db.LogFilter) *dbpb.LogFilter {
	pb := &dbpb.LogFilter{
		CreatedAt: f.CreatedAt,
		Metadata:  f.Metadata.Content,
		From: &dbpb.LogKey{
			BlockHeight:      f.From.BlockHeight,
			TransactionIndex: f.From.TransactionIndex,
			LogIndex:         f.From.LogIndex,
		},
		Next: &dbpb.LogKey{
			BlockHeight:      f.Next.BlockHeight,
			TransactionIndex: f.Next.TransactionIndex,
			LogIndex:         f.Next.LogIndex,
		},
		To: &dbpb.LogKey{
			BlockHeight:      f.To.BlockHeight,
			TransactionIndex: f.To.TransactionIndex,
			LogIndex:         f.To.LogIndex,
		},
	}
	pb.Addresses = make([][]byte, len(f.Addresses))
	for i, addr := range f.Addresses {
		pb.Addresses[i] = addr.Content
	}
	pb.Topics = make([]*dbpb.TopicList, len(f.Topics))
	for i, topicGroup := range f.Topics {
		tl := &dbpb.TopicList{}
		tl.Topics = make([][]byte, len(topicGroup))
		for j, t := range topicGroup {
			tl.Topics[j] = t.Content
		}
		pb.Topics[i] = tl
	}
	return pb
}

func logFilterFromProto(pb *dbpb.LogFilter, f *db.LogFilter) {
	f.CreatedAt = pb.CreatedAt
	f.Metadata = primitives.VarDataFromBytes(pb.Metadata)
	if pb.From != nil {
		f.From = db.LogKey{
			BlockHeight:      pb.From.BlockHeight,
			TransactionIndex: pb.From.TransactionIndex,
			LogIndex:         pb.From.LogIndex,
		}
	}
	if pb.Next != nil {
		f.Next = db.LogKey{
			BlockHeight:      pb.Next.BlockHeight,
			TransactionIndex: pb.Next.TransactionIndex,
			LogIndex:         pb.Next.LogIndex,
		}
	}
	if pb.To != nil {
		f.To = db.LogKey{
			BlockHeight:      pb.To.BlockHeight,
			TransactionIndex: pb.To.TransactionIndex,
			LogIndex:         pb.To.LogIndex,
		}
	}
	if pb.Addresses != nil {
		f.Addresses = make([]primitives.Data20, len(pb.Addresses))
		for i, addr := range pb.Addresses {
			f.Addresses[i] = data20(addr)
		}
	}
	if pb.Topics != nil {
		f.Topics = make([][]primitives.Data32, len(pb.Topics))
		for i, tl := range pb.Topics {
			f.Topics[i] = make([]primitives.Data32, len(tl.Topics))
			for j, t := range tl.Topics {
				f.Topics[i][j] = data32(t)
			}
		}
	}
}
