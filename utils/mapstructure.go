package utils

import (
	"math/big"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/aurora-is-near/relayer2-base/types/common"
)

func StringSliceToMapHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Slice || t != reflect.TypeOf(map[string]bool{}) {
			return data, nil
		}
		m := make(map[string]bool)
		for _, s := range data.([]interface{}) {
			s, ok := s.(string)
			if !ok {
				return data, nil
			}
			m[s] = true
		}
		return m, nil
	}
}

func BigIntHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(big.Int{}) {
			return data, nil
		}
		switch f.Kind() {
		case reflect.Int:
			b := big.NewInt(int64(data.(int)))
			return b, nil
		case reflect.Uint:
			b := big.NewInt(0)
			b = b.SetUint64(uint64(data.(uint)))
			return b, nil
		case reflect.String:
			b := big.NewInt(0)
			b, ok := b.SetString(data.(string), 0)
			if !ok {
				return data, nil
			}
			return b, nil
		default:
			return data, nil
		}
	}
}

func Uint256HookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Int || t != reflect.TypeOf(common.Uint256{}) {
			return data, nil
		}
		u256 := common.IntToUint256(data.(int))
		return u256, nil
	}
}
