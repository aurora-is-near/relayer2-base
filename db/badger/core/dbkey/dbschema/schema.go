package dbschema

import (
	"aurora-relayer-go-common/log"
)

type SchemaPath struct {
	consts []SchemaConst
	vars   []SchemaVar
	length int
}

type SchemaConst struct {
	value  byte
	offset int
}

type SchemaVar struct {
	size   int
	offset int
}

func Const(value byte) SchemaConst {
	return SchemaConst{value: value}
}

func Var(size int) SchemaVar {
	return SchemaVar{size: size}
}

func Path(tokens ...any) *SchemaPath {
	p := &SchemaPath{}
	for _, t := range tokens {
		switch tt := t.(type) {
		case SchemaConst:
			tt.offset = p.length
			p.consts = append(p.consts, tt)
			p.length += 1
		case SchemaVar:
			tt.offset = p.length
			p.vars = append(p.vars, tt)
			p.length += tt.size
		default:
			log.Log().Fatal().Msgf("object of type %T can't be token of dbschema.Path", t)
		}
	}
	return p
}

func (p *SchemaPath) Get(vars ...any) []byte {
	if len(vars) != len(p.vars) {
		log.Log().Fatal().Msg("wrong vars count")
	}

	result := make([]byte, p.length)

	for _, c := range p.consts {
		result[c.offset] = c.value
	}

	for i, varDesc := range p.vars {
		switch vt := vars[i].(type) {
		case []byte:
			copy(result[varDesc.offset:], vt[:varDesc.size])
		case string:
			copy(result[varDesc.offset:], vt[:varDesc.size])
		case uint64:
			putBigEndian(result[varDesc.offset:][:varDesc.size], vt)
		default:
			log.Log().Fatal().Msgf("%T can't be path var", vars[i])
		}
	}
	return result
}

func (p *SchemaPath) Matches(key []byte) bool {
	if len(key) < p.length {
		return false
	}
	for _, c := range p.consts {
		if key[c.offset] != c.value {
			return false
		}
	}
	return true
}

func (p *SchemaPath) ReadVar(key []byte, varIndex int) []byte {
	varDesc := p.vars[varIndex]
	return key[varDesc.offset:][:varDesc.size]
}

func (p *SchemaPath) ReadUintVar(key []byte, varIndex int) uint64 {
	return readBigEndian(p.ReadVar(key, varIndex))
}
