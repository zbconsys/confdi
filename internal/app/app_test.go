package app

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type mockLogger struct{}

func (l *mockLogger) Info(msg string, args ...interface{})  {}
func (l *mockLogger) Error(msg string, args ...interface{}) {}

func Test_TomlMerge(t *testing.T) {
	defaultToml := `
[traces-limits]
#
# Arithmetization module limits
#
ADD                 = 524286
BIN                 = 262128
BIN_RT              = 262144
EC_DATA             = 4084
EXT                 = 131060
HUB                 = 2097150
INSTRUCTION_DECODER = 512 # Ugly hack, TODO: @franklin
MMIO                = 1048576
MMU                 = 524288
MMU_ID              = 256
MOD                 = 131064
MUL                 = 65527
MXP                 = 524284
PHONEY_RLP          = 65536 # can probably get lower
PUB_HASH            = 32768
PUB_HASH_INFO       = 8192
PUB_LOG             = 16384
PUB_LOG_INFO        = 16384
RLP                 = 504
ROM                 = 4194302
SHF                 = 65520
SHF_RT              = 4096
TX_RLP              = 110000
WCP                 = 262128
#
# Block-specific limits
#
BLOCK_TX       = 200 # max number of tx in an L2 block
BLOCK_L2L1LOGS = 16
BLOCK_KECCAK   = 8192

#
# Precompiles limits
#
PRECOMPILE_ECRECOVER = 100
PRECOMPILE_SHA2      = 100
PRECOMPILE_RIPEMD    = 100
PRECOMPILE_IDENTITY  = 10000
PRECOMPILE_MODEXP    = 1000
PRECOMPILE_ECADD     = 1000
PRECOMPILE_ECMUL     = 100
PRECOMPILE_ECPAIRING = 100
PRECOMPILE_BLAKE2F   = 512
`

	overrideToml := `
[traces-limits]

#
# Block-specific limits
#
BLOCK_TX       = 400 # max number of tx in an L2 block
BLOCK_L2L1LOGS = 116
BLOCK_KECCAK   = 81920

#
# Precompiles limits
#

PRECOMPILE_ECADD     = 10000
PRECOMPILE_ECMUL     = 1000
PRECOMPILE_ECPAIRING = 1000
PRECOMPILE_BLAKE2F   = 1512
`
	a := App{
		log: &mockLogger{},
	}

	mergeBytes, err := a.processToml(strings.NewReader(defaultToml), strings.NewReader(overrideToml))
	assert.NoError(t, err)

	var merged map[string]interface{}

	err = toml.Unmarshal(mergeBytes, &merged)
	assert.NoError(t, err)

	data, ok := merged["traces-limits"]
	assert.True(t, ok)

	dataMap, ok := data.(map[string]interface{})
	assert.True(t, ok)

	// overridden values
	assert.Equal(t, int64(400), dataMap["BLOCK_TX"])
	assert.Equal(t, int64(116), dataMap["BLOCK_L2L1LOGS"])
	assert.Equal(t, int64(81920), dataMap["BLOCK_KECCAK"])
	assert.Equal(t, int64(10000), dataMap["PRECOMPILE_ECADD"])
	assert.Equal(t, int64(1000), dataMap["PRECOMPILE_ECMUL"])
	assert.Equal(t, int64(1000), dataMap["PRECOMPILE_ECPAIRING"])
	assert.Equal(t, int64(1512), dataMap["PRECOMPILE_BLAKE2F"])

	// default values
	assert.Equal(t, int64(524286), dataMap["ADD"])
	assert.Equal(t, int64(262128), dataMap["BIN"])
	assert.Equal(t, int64(262144), dataMap["BIN_RT"])
}
