using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("types");
$Go.import("db/types");

struct Block {
       hash             @0  :Data;
       parentHash       @1  :Data;
       height           @2  :UInt64;
       miner            @3  :Data;
       timestamp        @4  :Int64;
       gasLimit         @5  :Data;
       gasUsed          @6  :Data;
       logsBloom        @7  :Data;
       transactionsRoot @8  :Data;
       receiptsRoot     @9  :Data;
       stateRoot        @10 :Data;
       size             @11 :Data;
}

struct Filter {
	type      @0 :Text;
	createdBy @1 :Text;
	pollBlock @2 :Data;
	fromBlock @3 :Data;
	toBlock   @4 :Data;
	addresses @5 :List(Data);
	topics    @6 :List(List(Data));
}

struct Transaction {
	hash                 @0   :Data;
	blockHash            @1   :Data;
	blockHeight          @2   :UInt64;
	transactionIndex     @3   :UInt32;
	from                 @4   :Data;
	to                   @5   :Data;
	nonce                @6   :Data;
	gasPrice             @7   :Data;
	gasLimit             @8   :Data;
	gasUsed              @9   :Data;
	# maxPriorityFeePerGas @0 :Data;
	# maxFeePerGas         @0 :Data;
	value                @10  :Data;
	input                @11  :Data;
	output               @12  :Data;
	# txType               @0 :UInt8;
	status               @13  :Bool;
	contractAddress      @14  :Data;
	v                    @15  :UInt64;
	r                    @16  :Data;
	s                    @17  :Data;
        nearHash             @18  :Data;
        nearReceiptHash      @19  :Data;
}

struct Log {
	removed          @0 :Bool;
	logIndex         @1 :Data;
	transactionIndex @2 :Data;
	transactionHash  @3 :Data;
	blockHash        @4 :Data;
	blockNumber      @5 :Data;
	address          @6 :Data;
	data             @7 :Data;
	topics           @8 :List(Data);
}
