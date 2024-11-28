package primitives

type (
	VarData = Data[VarLen]
	Data8   = Data[Len8]
	Data20  = Data[Len20]
	Data32  = Data[Len32]
	Data256 = Data[Len256]
)

var (
	VarDataFromBytes = DataFromBytes[VarLen]
	Data8FromBytes   = DataFromBytes[Len8]
	Data20FromBytes  = DataFromBytes[Len20]
	Data32FromBytes  = DataFromBytes[Len32]
	Data256FromBytes = DataFromBytes[Len256]
)

var (
	MustVarDataFromHex = MustDataFromHex[VarLen]
	MustData8FromHex   = MustDataFromHex[Len8]
	MustData20FromHex  = MustDataFromHex[Len20]
	MustData32FromHex  = MustDataFromHex[Len32]
	MustData256FromHex = MustDataFromHex[Len256]
)
