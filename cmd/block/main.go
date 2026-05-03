package main

import (
	"fmt"
	"github.com/Yeah114/sc2-world-operator/content"
)

func main() {
	cell := content.EncodeBlockState(
		&content.Block{
			BlockID:    content.IDWoodenTrapdoor,
			LightLevel: 3,
		},
		&content.Trapdoor{
			Facing:     content.FacingNorth,
			Open:       true,
			UpsideDown: false,
		},
	)
	fmt.Printf("cell = 0x%08x\n", cell)

	block, state := content.DecodeBlockState(cell)
	fmt.Printf("BlockID: %d\n", block.BlockID)
	fmt.Printf("LightLevel: %d\n", block.LightLevel)

	trapdoor := state.(*content.Trapdoor)
	fmt.Printf("Facing: %d\n", trapdoor.Facing)
	fmt.Printf("Open: %t\n", trapdoor.Open)
	fmt.Printf("UpsideDown: %t\n", trapdoor.UpsideDown)

	_, properties := content.DecodeBlockProperties(cell)
	fmt.Printf("%+v\n", properties)
}
