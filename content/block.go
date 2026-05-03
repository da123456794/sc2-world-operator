package content

import (
	"github.com/Yeah114/sc2-world-operator/bit"
	"github.com/mitchellh/mapstructure"
)

// Block holds the basic data of a block in the world.
type Block struct {
	// BlockID is the content ID of the block.
	BlockID int32
	// LightLevel is the light level emitted or stored by the block.
	LightLevel int32
}

// Marshal encodes or decodes the basic block data.
func (b *Block) Marshal(io bit.IO) {
	io.Int32(&b.BlockID, bit.BitsBlockID)
	io.Int32(&b.LightLevel, bit.BitsLightLevel)
}

// BlockState represents the additional state of a specific block type.
type BlockState interface {
	Marshal(io bit.IO)
}

// DecodeBlockState decodes a packed 32-bit cell value into a Block and its corresponding BlockState.
func DecodeBlockState(cell int32) (*Block, BlockState) {
	block := &Block{}
	reader := bit.NewReader(&cell)
	block.Marshal(reader)
	blockStateFunc, ok := blockStatePool[block.BlockID]
	if !ok {
		return block, nil
	}
	blockState := blockStateFunc()
	blockState.Marshal(reader)
	return block, blockState
}

// EncodeBlockState encodes a Block and its BlockState into a packed 32-bit cell value.
func EncodeBlockState(block *Block, blockState BlockState) int32 {
	var cell int32
	writer := bit.NewWriter(&cell)
	block.Marshal(writer)
	if blockState != nil {
		blockState.Marshal(writer)
	}
	return cell
}

// DecodeBlockProperties decodes a packed 32-bit cell value into a Block and a property map.
// It converts the BlockState into a map for easy serialization and manipulation.
func DecodeBlockProperties(cell int32) (*Block, map[string]any) {
	block, blockState := DecodeBlockState(cell)
	if blockState == nil {
		return block, nil
	}
	var properties map[string]any
	_ = mapstructure.WeakDecode(blockState, &properties)
	return block, properties
}

// EncodeBlockProperties encodes a Block and a property map back into a packed 32-bit cell value.
// It reconstructs the corresponding BlockState from the property map before encoding.
func EncodeBlockProperties(block *Block, properties map[string]any) int32 {
	var state BlockState
	if properties != nil {
		blockStateFunc, ok := blockStatePool[block.BlockID]
		if !ok {
			return EncodeBlockState(block, nil)
		}
		state = blockStateFunc()
		_ = mapstructure.WeakDecode(properties, state)
	}
	return EncodeBlockState(block, state)
}
