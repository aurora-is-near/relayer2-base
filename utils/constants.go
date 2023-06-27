package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/carlmjohnson/versioninfo"
	"gopkg.in/yaml.v2"
)

const (
	buildInfoFileNameStandalone = "version"
	buildInfoFileNameNative     = "build.info"
)

type constants struct {
	clientVersion  string
	zeroStrUint256 string
	zeroStrUint160 string
	zeroStrUint128 string
	response0x     string
	syncing        bool
	mining         bool
	full           bool
	gasLimit       uint64
	emptyArray     []string
	zeroUint256    common.Uint256
	relayerVersion string
}

type buildInfo struct {
	Branch string `yaml:"branch"`
	Commit string `yaml:"commit"`
	Tag    string `yaml:"tag"`
}

func (c *constants) ClientVersion() *string {
	return &c.clientVersion
}

func (c *constants) ZeroStrUint256() *string {
	return &c.zeroStrUint256
}

func (c *constants) ZeroStrUint160() *string {
	return &c.zeroStrUint160
}

func (c constants) ZeroStrUint128() *string {
	return &c.zeroStrUint128
}

func (c *constants) Response0x() *string {
	return &c.response0x
}

func (c *constants) Syncing() *bool {
	return &c.syncing
}

func (c *constants) Mining() *bool {
	return &c.mining
}

func (c *constants) Full() *bool {
	return &c.full
}

func (c *constants) GasLimit() *uint64 {
	return &c.gasLimit
}

func (c *constants) EmptyArray() *[]string {
	return &c.emptyArray
}

func (c *constants) ZeroUint256() *common.Uint256 {
	return &c.zeroUint256
}

func (c *constants) RelayerVersion() *string {
	return &c.relayerVersion
}

var Constants constants

func init() {
	Constants.clientVersion = "Aurora"
	Constants.zeroStrUint256 = "0x0000000000000000000000000000000000000000000000000000000000000000"
	Constants.zeroStrUint160 = "0x0000000000000000000000000000000000000000"
	Constants.zeroStrUint128 = "0x00000000000000000000000000000000"
	Constants.response0x = "0x"
	Constants.syncing = false
	Constants.mining = false
	Constants.full = false
	Constants.gasLimit = 9007199254740991 // hex value 0x1fffffffffffff
	Constants.emptyArray = []string{}
	Constants.zeroUint256 = common.IntToUint256(0)

	bi, err := getBuildInfo()
	// If err is not nil then it means version info can't be read, error message is logged
	// still the commit hash is	returned in `bi.Tag` field.
	if err != nil {
		log.Log().Err(err).Msg("failed to get version tag")
	}

	if len(strings.TrimSpace(bi.Tag)) == 0 && len(strings.TrimSpace(bi.Commit)) > 0 {
		Constants.relayerVersion = bi.Commit
	} else {
		Constants.relayerVersion = bi.Tag
	}
}

// getBuildInfo checks the files containing build information and returns
// generated buildInfo object.
// If no file found or version tag is empty, then the latest commit hash is returned in tag field
func getBuildInfo() (*buildInfo, error) {
	var err error
	// first, check the StandaloneRPC app build file which is located at root
	file2open := "/" + buildInfoFileNameStandalone
	if _, err = os.Stat(file2open); err != nil {
		// Second, check native apps build file
		// executable path is extracted to be able to locate the build info file
		ex, _ := os.Executable()
		exPath := filepath.Dir(ex)
		file2open = exPath + "/" + buildInfoFileNameNative
		if _, err = os.Stat(file2open); err != nil {
			// if both failed, then the latest commit hash is used as tag
			re := buildInfo{
				Tag: versioninfo.Revision,
			}
			return &re, errors.New("neither " + buildInfoFileNameStandalone + " nor " + buildInfoFileNameNative + " exist")
		}
	}
	file, err := ioutil.ReadFile(file2open)
	if err != nil {
		re := buildInfo{
			Tag: versioninfo.Revision,
		}
		return &re, errors.New(file2open + " couldn't be opened")
	}

	var bi buildInfo
	err = yaml.Unmarshal(file, &bi)
	if err != nil {
		re := buildInfo{
			Tag: versioninfo.Revision,
		}
		return &re, errors.New("File reading error for " + file2open)
	}
	return &bi, nil
}
