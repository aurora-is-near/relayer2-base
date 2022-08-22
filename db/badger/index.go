package badger

import "bytes"

type Index struct {
	numFields        uint
	indexEmptyFields bool
	maxBitmask       uint8
}

/*
	In case if there will be no searches by empty values - set indexEmptyFields to false,
	it can dramatically save space.
*/
func NewIndex(numFields uint, indexEmptyFields bool) *Index {
	if numFields > 8 {
		panic("NewIndex: numFields can't be greater than 8")
	}
	return &Index{
		numFields:        numFields,
		indexEmptyFields: indexEmptyFields,
		maxBitmask:       (1 << numFields) - 1,
	}
}

/*
	Missing and nil fields are considered empty.
	Extra fields are dropped.
	Fields longer than 255-characters are chopped.
*/
func (index *Index) Insert(
	indexTablePrefix []byte,
	fields [][]byte,
	primaryKey []byte,
	receiver func(key, value []byte) error,
) error {

	emptyFieldsBitmask := uint8(0)
	if !index.indexEmptyFields {
		for pos := 0; pos < int(index.numFields); pos++ {
			if len(fields) <= pos || len(fields[pos]) == 0 {
				emptyFieldsBitmask = emptyFieldsBitmask | (1 << pos)
			}
		}
	}

	for keyBitmask := uint8(0); keyBitmask <= index.maxBitmask; keyBitmask++ {
		if !index.indexEmptyFields && (keyBitmask&emptyFieldsBitmask > 0) {
			continue
		}

		var keyBuff, valueBuff bytes.Buffer
		keyBuff.Write(indexTablePrefix)
		keyBuff.WriteByte(keyBitmask)

		for pos := 0; pos < int(index.numFields); pos++ {
			var field []byte
			if pos < len(fields) {
				field = fields[pos]
			}
			if field == nil {
				field = []byte{}
			}
			if len(field) > 255 {
				field = field[:255]
			}

			if keyBitmask&(1<<pos) > 0 {
				keyBuff.WriteByte(uint8(len(field)))
				keyBuff.Write(field)
			} else {
				valueBuff.WriteByte(uint8(len(field)))
				valueBuff.Write(field)
			}
		}

		keyBuff.Write(primaryKey)

		if err := receiver(keyBuff.Bytes(), valueBuff.Bytes()); err != nil {
			return err
		}
	}

	return nil
}
