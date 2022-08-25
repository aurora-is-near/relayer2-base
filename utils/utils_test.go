package utils_test

import (
	"aurora-relayer-go-common/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHexStringToAddress(t *testing.T) {
	// addressParseTests defined in types_test.go
	for _, tc := range addressParseTests {
		t.Run(tc.data, func(t *testing.T) {
			address, err := utils.ParseHexStringToAddress(tc.data)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.Nil(t, err)
				want := tc.want.Hash().Big()
				got := address.Hash().Big()
				assert.Equal(t, got, want)
			}
		})
	}
}
