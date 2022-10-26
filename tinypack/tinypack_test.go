package tinypack

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func checkMarshalUnmarshal[T any](t *testing.T, x *T, e *Encoder, d *Decoder) {
	data, err := e.Marshal(x)
	if !assert.NoErrorf(t, err, "marshal should be successful") {
		spew.Dump(x)
		t.FailNow()
	}
	y := new(T)
	err = d.Unmarshal(data, y)
	if !assert.NoError(t, err, "unmarshal should be successful") {
		spew.Dump(x)
		t.FailNow()
	}
	if !assert.Equal(t, x, y, "unmarshal(marshal(x)) should be equal to x") {
		spew.Dump(x)
		spew.Dump(y)
		t.FailNow()
	}
}

func ensureMarshalFails[T any](t *testing.T, x *T, e *Encoder) {
	_, err := e.Marshal(x)
	if !assert.Error(t, err, "marshal should fail") {
		spew.Dump(x)
		t.FailNow()
	}
}

var (
	boolTestValues    = []bool{false, true}
	varintTestValues  = []int64{-123451234512345, -123123123, -111, 1, 0, 1, 111, 321321321, 543215432154321}
	uvarintTestValues = []uint64{0, 1, 222, 456456456, 765765765765765}
	floatTestValues   = []float64{0, 3.1415, -3.1415, 2.7182, -2.7182, 1024, -512}
)

func TestPrimitives(t *testing.T) {
	for _, val := range boolTestValues {
		checkMarshalUnmarshal(t, &val, DefaultEncoder(), DefaultDecoder())
	}
	for _, val := range varintTestValues {
		checkMarshalUnmarshal(t, &val, DefaultEncoder(), DefaultDecoder())
	}
	for _, val := range uvarintTestValues {
		checkMarshalUnmarshal(t, &val, DefaultEncoder(), DefaultDecoder())
	}
	for _, val := range floatTestValues {
		checkMarshalUnmarshal(t, &val, DefaultEncoder(), DefaultDecoder())
	}
}

type structWithPrimitives struct {
	A bool
	B int64
	C uint64
	D float64
}

func (s *structWithPrimitives) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&s.A,
		&s.B,
		&s.C,
		&s.D,
	}, nil
}

func TestStructWithPrimitives(t *testing.T) {
	for _, boolVal := range boolTestValues {
		for _, varintVal := range varintTestValues {
			for _, uvarintVal := range uvarintTestValues {
				for _, floatVal := range floatTestValues {
					x := structWithPrimitives{
						A: boolVal,
						B: varintVal,
						C: uvarintVal,
						D: floatVal,
					}
					checkMarshalUnmarshal(t, &x, DefaultEncoder(), DefaultDecoder())
				}
			}
		}
	}
}

func TestNullable(t *testing.T) {
	var x Nullable[int64]
	val := int64(-5)
	x.Ptr = &val
	checkMarshalUnmarshal(t, &x, DefaultEncoder(), DefaultDecoder())

	x.Ptr = nil
	checkMarshalUnmarshal(t, &x, DefaultEncoder(), DefaultDecoder())
}

func TestPointer(t *testing.T) {
	var x Pointer[uint64]
	val := uint64(55555)
	x.Ptr = &val
	checkMarshalUnmarshal(t, &x, DefaultEncoder(), DefaultDecoder())

	x.Ptr = nil
	ensureMarshalFails(t, &x, DefaultEncoder())
}

type len0 struct{}
type len1 struct{}
type len3 struct{}
type len10 struct{}
type len300 struct{}

func (ld len0) GetTinyPackLength() int {
	return 0
}

func (ld len1) GetTinyPackLength() int {
	return 1
}

func (ld len3) GetTinyPackLength() int {
	return 3
}

func (ld len10) GetTinyPackLength() int {
	return 10
}

func (ld len300) GetTinyPackLength() int {
	return 300
}

func genU64List[LD LengthDescriptor](length int) *List[LD, uint64] {
	var list List[LD, uint64]
	list.Content = make([]uint64, length)
	for i := 0; i < length; i++ {
		list.Content[i] = uint64(i * i)
	}
	return &list
}

func TestList(t *testing.T) {
	checkMarshalUnmarshal(t, genU64List[len0](0), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genU64List[len0](1), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len0](10), DefaultEncoder())

	ensureMarshalFails(t, genU64List[len1](0), DefaultEncoder())
	checkMarshalUnmarshal(t, genU64List[len1](1), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genU64List[len1](10), DefaultEncoder())

	ensureMarshalFails(t, genU64List[len10](0), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len10](1), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len10](9), DefaultEncoder())
	checkMarshalUnmarshal(t, genU64List[len10](10), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genU64List[len10](11), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len10](30), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len10](150), DefaultEncoder())

	ensureMarshalFails(t, genU64List[len300](0), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len300](1), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len300](299), DefaultEncoder())
	checkMarshalUnmarshal(t, genU64List[len300](300), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genU64List[len300](301), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len300](350), DefaultEncoder())
	ensureMarshalFails(t, genU64List[len300](999), DefaultEncoder())
}

func genU64VarList(length int) *VarList[uint64] {
	var list VarList[uint64]
	list.Content = make([]uint64, length)
	for i := 0; i < length; i++ {
		list.Content[i] = uint64(i * i)
	}
	return &list
}

func TestVarList(t *testing.T) {
	checkMarshalUnmarshal(t, genU64VarList(0), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genU64VarList(1), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genU64VarList(12), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genU64VarList(123), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genU64VarList(1234), DefaultEncoder(), DefaultDecoder())
}

func genData[LD LengthDescriptor](length int) *Data[LD] {
	var data Data[LD]
	data.Content = make([]byte, length)
	for i := 0; i < length; i++ {
		data.Content[i] = byte((i * i) % 256)
	}
	return &data
}

func TestData(t *testing.T) {
	checkMarshalUnmarshal(t, genData[len0](0), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genData[len0](1), DefaultEncoder())
	ensureMarshalFails(t, genData[len0](10), DefaultEncoder())

	ensureMarshalFails(t, genData[len1](0), DefaultEncoder())
	checkMarshalUnmarshal(t, genData[len1](1), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genData[len1](10), DefaultEncoder())

	ensureMarshalFails(t, genData[len10](0), DefaultEncoder())
	ensureMarshalFails(t, genData[len10](1), DefaultEncoder())
	ensureMarshalFails(t, genData[len10](9), DefaultEncoder())
	checkMarshalUnmarshal(t, genData[len10](10), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genData[len10](11), DefaultEncoder())
	ensureMarshalFails(t, genData[len10](30), DefaultEncoder())
	ensureMarshalFails(t, genData[len10](150), DefaultEncoder())

	ensureMarshalFails(t, genData[len300](0), DefaultEncoder())
	ensureMarshalFails(t, genData[len300](1), DefaultEncoder())
	ensureMarshalFails(t, genData[len300](299), DefaultEncoder())
	checkMarshalUnmarshal(t, genData[len300](300), DefaultEncoder(), DefaultDecoder())
	ensureMarshalFails(t, genData[len300](301), DefaultEncoder())
	ensureMarshalFails(t, genData[len300](350), DefaultEncoder())
	ensureMarshalFails(t, genData[len300](999), DefaultEncoder())
}

func genVarData(length int) *VarData {
	var data VarData
	data.Content = make([]byte, length)
	for i := 0; i < length; i++ {
		data.Content[i] = byte((i * i) % 256)
	}
	return &data
}

func TestVarData(t *testing.T) {
	checkMarshalUnmarshal(t, genVarData(0), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genVarData(1), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genVarData(12), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genVarData(123), DefaultEncoder(), DefaultDecoder())
	checkMarshalUnmarshal(t, genVarData(1234), DefaultEncoder(), DefaultDecoder())
}

type struct1 struct {
	A Pointer[int64]
	B List[len0, float64]
}

func (s *struct1) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{&s.A, &s.B}, nil
}

type struct2 struct {
	A VarList[struct1]
	B Pointer[VarList[Nullable[List[len3, struct1]]]]
	C Data[len1]
	D Nullable[bool]
}

func (s *struct2) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{&s.A, &s.B, &s.C, &s.D}, nil
}

type struct3 struct {
	A List[len3, Nullable[Data[len10]]]
	B VarList[Pointer[struct2]]
	C VarData
}

func (s *struct3) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{&s.A, &s.B, &s.C}, nil
}

func toPtr[T any](val T) *T {
	return &val
}

func TestVeryComplicatedHierarchy(t *testing.T) {
	s := &struct3{
		A: CreateList[len3](
			CreateNullable[Data[len10]](nil),
			CreateNullable(toPtr(CreateData[len10](0, 8, 5, 7, 1, 179, 250, 118, 0, 0))),
			CreateNullable[Data[len10]](nil),
		),
		B: CreateVarList(
			CreatePointer(&struct2{
				A: CreateVarList(
					struct1{
						A: CreatePointer(toPtr(int64(0))),
						B: CreateList[len0, float64](),
					},
					struct1{
						A: CreatePointer(toPtr(int64(123451234512345))),
						B: CreateList[len0, float64](),
					},
					struct1{
						A: CreatePointer(toPtr(int64(0))),
						B: CreateList[len0, float64](),
					},
					struct1{
						A: CreatePointer(toPtr(int64(0))),
						B: CreateList[len0, float64](),
					},
				),
				B: CreatePointer(toPtr(CreateVarList(
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable(toPtr(CreateList[len3](
						struct1{
							A: CreatePointer(toPtr(int64(-543543543))),
							B: CreateList[len0, float64](),
						},
						struct1{
							A: CreatePointer(toPtr(int64(0))),
							B: CreateList[len0, float64](),
						},
						struct1{
							A: CreatePointer(toPtr(int64(0))),
							B: CreateList[len0, float64](),
						},
					))),
				))),
				C: CreateData[len1](255),
				D: CreateNullable(toPtr(false)),
			}),
			CreatePointer(&struct2{
				A: CreateVarList[struct1](),
				B: CreatePointer(toPtr(CreateVarList(
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable(toPtr(CreateList[len3](
						struct1{
							A: CreatePointer(toPtr(int64(1))),
							B: CreateList[len0, float64](),
						},
						struct1{
							A: CreatePointer(toPtr(int64(-1))),
							B: CreateList[len0, float64](),
						},
						struct1{
							A: CreatePointer(toPtr(int64(0))),
							B: CreateList[len0, float64](),
						},
					))),
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable(toPtr(CreateList[len3](
						struct1{
							A: CreatePointer(toPtr(int64(0))),
							B: CreateList[len0, float64](),
						},
						struct1{
							A: CreatePointer(toPtr(int64(0))),
							B: CreateList[len0, float64](),
						},
						struct1{
							A: CreatePointer(toPtr(int64(0))),
							B: CreateList[len0, float64](),
						},
					))),
				))),
				C: CreateData[len1](0),
				D: CreateNullable[bool](nil),
			}),
			CreatePointer(&struct2{
				A: CreateVarList(
					struct1{
						A: CreatePointer(toPtr(int64(-5435435))),
						B: CreateList[len0, float64](),
					},
					struct1{
						A: CreatePointer(toPtr(int64(5435435))),
						B: CreateList[len0, float64](),
					},
				),
				B: CreatePointer(toPtr(CreateVarList(
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable[List[len3, struct1]](nil),
					CreateNullable[List[len3, struct1]](nil),
				))),
				C: CreateData[len1](128),
				D: CreateNullable(toPtr(true)),
			}),
		),
		C: CreateVarData(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
	}
	checkMarshalUnmarshal(t, s, DefaultEncoder(), DefaultDecoder())

	// var buf bytes.Buffer
	// Write(&buf, s, DefaultEncoder())
	// fmt.Println(len(buf.Bytes()))
	// fmt.Println(buf.Bytes())
	// Just 97 bytes!

	// data, _ := DefaultEncoder().Marshal(s)
	// fmt.Println(len(data))
	// fmt.Println(data)
	// Just 60 bytes!
}
