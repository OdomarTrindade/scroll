package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/scroll-tech/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewBatchHeader(t *testing.T) {
	// Without L1 Msg
	templateBlockTrace, err := os.ReadFile("../testdata/blockTrace_02.json")
	assert.NoError(t, err)

	wrappedBlock := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace, wrappedBlock))
	chunk := &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock,
		},
	}
	parentBatchHeader := &BatchHeader{
		version:                1,
		batchIndex:             0,
		l1MessagePopped:        0,
		totalL1MessagePopped:   0,
		dataHash:               common.HexToHash("0x0"),
		parentBatchHash:        common.HexToHash("0x0"),
		skippedL1MessageBitmap: nil,
	}
	batchHeader, err := NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	assert.Equal(t, 0, len(batchHeader.skippedL1MessageBitmap))

	// 1 L1 Msg in 1 bitmap
	templateBlockTrace2, err := os.ReadFile("../testdata/blockTrace_04.json")
	assert.NoError(t, err)

	wrappedBlock2 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace2, wrappedBlock2))
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock2,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	assert.Equal(t, 32, len(batchHeader.skippedL1MessageBitmap))
	expectedBitmap := "00000000000000000000000000000000000000000000000000000000000003ff" // skip first 10
	assert.Equal(t, expectedBitmap, common.Bytes2Hex(batchHeader.skippedL1MessageBitmap))

	// many consecutive L1 Msgs in 1 bitmap, no leading skipped msgs
	templateBlockTrace3, err := os.ReadFile("../testdata/blockTrace_05.json")
	assert.NoError(t, err)

	wrappedBlock3 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace3, wrappedBlock3))
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock3,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 37, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	assert.Equal(t, uint64(5), batchHeader.l1MessagePopped)
	assert.Equal(t, 32, len(batchHeader.skippedL1MessageBitmap))
	expectedBitmap = "0000000000000000000000000000000000000000000000000000000000000000" // all bits are included, so none are skipped
	assert.Equal(t, expectedBitmap, common.Bytes2Hex(batchHeader.skippedL1MessageBitmap))

	// many consecutive L1 Msgs in 1 bitmap, with leading skipped msgs
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock3,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	assert.Equal(t, uint64(42), batchHeader.l1MessagePopped)
	assert.Equal(t, 32, len(batchHeader.skippedL1MessageBitmap))
	expectedBitmap = "0000000000000000000000000000000000000000000000000000001fffffffff" // skipped the first 37 messages
	assert.Equal(t, expectedBitmap, common.Bytes2Hex(batchHeader.skippedL1MessageBitmap))

	// many sparse L1 Msgs in 1 bitmap
	templateBlockTrace4, err := os.ReadFile("../testdata/blockTrace_06.json")
	assert.NoError(t, err)

	wrappedBlock4 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace4, wrappedBlock4))
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock4,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	assert.Equal(t, uint64(10), batchHeader.l1MessagePopped)
	assert.Equal(t, 32, len(batchHeader.skippedL1MessageBitmap))
	expectedBitmap = "00000000000000000000000000000000000000000000000000000000000001dd" // 0111011101
	assert.Equal(t, expectedBitmap, common.Bytes2Hex(batchHeader.skippedL1MessageBitmap))

	// many L1 Msgs in each of 2 bitmaps
	templateBlockTrace5, err := os.ReadFile("../testdata/blockTrace_07.json")
	assert.NoError(t, err)

	wrappedBlock5 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace5, wrappedBlock5))
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock5,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	assert.Equal(t, uint64(257), batchHeader.l1MessagePopped)
	assert.Equal(t, 64, len(batchHeader.skippedL1MessageBitmap))
	expectedBitmap = "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd0000000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expectedBitmap, common.Bytes2Hex(batchHeader.skippedL1MessageBitmap))
}

func TestBatchHeaderEncode(t *testing.T) {
	// Without L1 Msg
	templateBlockTrace, err := os.ReadFile("../testdata/blockTrace_02.json")
	assert.NoError(t, err)

	wrappedBlock := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace, wrappedBlock))
	chunk := &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock,
		},
	}
	parentBatchHeader := &BatchHeader{
		version:                1,
		batchIndex:             0,
		l1MessagePopped:        0,
		totalL1MessagePopped:   0,
		dataHash:               common.HexToHash("0x0"),
		parentBatchHash:        common.HexToHash("0x0"),
		skippedL1MessageBitmap: nil,
	}
	batchHeader, err := NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	bytes := batchHeader.Encode()
	assert.Equal(t, 89, len(bytes))
	assert.Equal(t, "0100000000000000010000000000000000000000000000000010a64c9bd905f8caf5d668fbda622d6558c5a42cdb4b3895709743d159c22e534136709aabc8a23aa17fbcc833da2f7857d3c2884feec9aae73429c135f94985", common.Bytes2Hex(bytes))

	// With L1 Msg
	templateBlockTrace2, err := os.ReadFile("../testdata/blockTrace_04.json")
	assert.NoError(t, err)

	wrappedBlock2 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace2, wrappedBlock2))
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock2,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	bytes = batchHeader.Encode()
	assert.Equal(t, 121, len(bytes))
	assert.Equal(t, "010000000000000001000000000000000b000000000000000b34f419ce7e882295bdb5aec6cce56ffa788a5fed4744d7fbd77e4acbf409f1ca4136709aabc8a23aa17fbcc833da2f7857d3c2884feec9aae73429c135f9498500000000000000000000000000000000000000000000000000000000000003ff", common.Bytes2Hex(bytes))
}

func TestBatchHeaderHash(t *testing.T) {
	// Without L1 Msg
	templateBlockTrace, err := os.ReadFile("../testdata/blockTrace_02.json")
	assert.NoError(t, err)

	wrappedBlock := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace, wrappedBlock))
	chunk := &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock,
		},
	}
	parentBatchHeader := &BatchHeader{
		version:                1,
		batchIndex:             0,
		l1MessagePopped:        0,
		totalL1MessagePopped:   0,
		dataHash:               common.HexToHash("0x0"),
		parentBatchHash:        common.HexToHash("0x0"),
		skippedL1MessageBitmap: nil,
	}
	batchHeader, err := NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	hash := batchHeader.Hash()
	assert.Equal(t, "d69da4357da0073f4093c76e49f077e21bb52f48f57ee3e1fbd9c38a2881af81", common.Bytes2Hex(hash.Bytes()))

	templateBlockTrace, err = os.ReadFile("../testdata/blockTrace_03.json")
	assert.NoError(t, err)

	wrappedBlock2 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace, wrappedBlock2))
	chunk2 := &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock2,
		},
	}
	batchHeader2, err := NewBatchHeader(1, 2, 0, batchHeader.Hash(), []*Chunk{chunk2})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader2)
	hash2 := batchHeader2.Hash()
	assert.Equal(t, "34de600163aa745d4513113137a5b54960d13f0d3f2849e490c4b875028bf930", common.Bytes2Hex(hash2.Bytes()))

	// With L1 Msg
	templateBlockTrace3, err := os.ReadFile("../testdata/blockTrace_04.json")
	assert.NoError(t, err)

	wrappedBlock3 := &WrappedBlock{}
	assert.NoError(t, json.Unmarshal(templateBlockTrace3, wrappedBlock3))
	chunk = &Chunk{
		Blocks: []*WrappedBlock{
			wrappedBlock3,
		},
	}
	batchHeader, err = NewBatchHeader(1, 1, 0, parentBatchHeader.Hash(), []*Chunk{chunk})
	assert.NoError(t, err)
	assert.NotNil(t, batchHeader)
	hash = batchHeader.Hash()
	assert.Equal(t, "1c3007880f0eafe74572ede7d164ff1ee5376e9ac9bff6f7fb837b2630cddc9a", common.Bytes2Hex(hash.Bytes()))
}

func TestBatchHeaderDecode(t *testing.T) {
	header := &BatchHeader{
		version:                1,
		batchIndex:             10,
		l1MessagePopped:        20,
		totalL1MessagePopped:   30,
		dataHash:               common.HexToHash("0x01"),
		parentBatchHash:        common.HexToHash("0x02"),
		skippedL1MessageBitmap: []byte{0x01, 0x02, 0x03},
	}

	encoded := header.Encode()
	decoded, err := DecodeBatchHeader(encoded)
	assert.NoError(t, err)
	assert.Equal(t, header, decoded)
}
