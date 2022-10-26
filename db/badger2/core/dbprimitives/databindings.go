package dbprimitives

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
	VarDataFromHex = DataFromHex[VarLen]
	Data8FromHex   = DataFromHex[Len8]
	Data20FromHex  = DataFromHex[Len20]
	Data32FromHex  = DataFromHex[Len32]
	Data256FromHex = DataFromHex[Len256]
)
