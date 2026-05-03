package world

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// InventorySlot represents one slot in any inventory (chest, furnace, player, etc.).
// Contents is a full CellValue: use define.ExtractContents to get block ID,
// define.ExtractData to get block metadata.
type InventorySlot struct {
	Index    int   // slot index (0-based)
	Contents int32 // CellValue (blockID in bits 0-9, data in bits 14-31)
	Count    int32 // stack count
}

// Inventory holds the parsed slots for any inventory component.
type Inventory struct {
	ID    int32
	Slots []InventorySlot // only non-empty slots are populated
}

// BlockEntityData is a block entity placed in the world
// (chest, furnace, dispenser, crafting table, etc.).
type BlockEntityData struct {
	EntityID  string     // key in "Entities" dict
	Name      string     // entity template name (e.g. "Chest", "Furnace")
	X, Y, Z   int        // world block coordinates
	OwnerGUID string     // player GUID who placed it
	Inventory *Inventory // nil if no inventory component
	// Raw component overrides for advanced access
	Raw map[string]json.RawMessage
}

// LoadBlockEntities reads all block entities from Project.json.
// Block entities are any entities that have a "BlockEntity" component
// (chest, furnace, dispenser, crafting table, etc.).
func (w *World) LoadBlockEntities() ([]*BlockEntityData, error) {
	proj, err := w.LoadProject()
	if err != nil {
		return nil, err
	}

	rawEntities, ok := proj.Raw["Entities"]
	if !ok {
		return nil, nil
	}
	var entities map[string]json.RawMessage
	if err := json.Unmarshal(rawEntities, &entities); err != nil {
		return nil, fmt.Errorf("world.LoadBlockEntities: unmarshal entities: %w", err)
	}

	var result []*BlockEntityData
	for eid, rawEnt := range entities {
		var ent struct {
			Name      json.RawMessage            `json:"Name"`
			Overrides map[string]json.RawMessage `json:"Overrides"`
		}
		if err := json.Unmarshal(rawEnt, &ent); err != nil {
			continue
		}
		// Only block entities have a "BlockEntity" component.
		beRaw, hasBE := ent.Overrides["BlockEntity"]
		if !hasBE {
			continue
		}

		bed := &BlockEntityData{
			EntityID: eid,
			Raw:      ent.Overrides,
		}

		// Name field: ["string", "Chest"]
		if ent.Name != nil {
			var pair [2]json.RawMessage
			if json.Unmarshal(ent.Name, &pair) == nil {
				json.Unmarshal(pair[1], &bed.Name) //nolint:errcheck
			}
		}

		// BlockEntity component: {"Coordinates": ..., "Owner": ...}
		parseBlockEntityComp(beRaw, bed)

		// Inventory component (Chest, Furnace, Dispenser, CraftingTable, etc.)
		for _, compName := range []string{"Inventory", "Chest", "Furnace", "Dispenser"} {
			if invRaw, ok := ent.Overrides[compName]; ok {
				if inv := parseInventory(invRaw); inv != nil {
					bed.Inventory = inv
					break
				}
			}
		}

		result = append(result, bed)
	}
	return result, nil
}

// parseBlockEntityComp fills X,Y,Z and OwnerGUID from the BlockEntity component JSON.
// The component stores values using the ["typeName", value] pair format.
func parseBlockEntityComp(raw json.RawMessage, bed *BlockEntityData) {
	var comp map[string]json.RawMessage
	if json.Unmarshal(raw, &comp) != nil {
		return
	}
	// Coordinates: ["Game.Point3", "X,Y,Z"] or {"X":["int",n], "Y":["int",n], "Z":["int",n]}
	if coordRaw, ok := comp["Coordinates"]; ok {
		// Try type-value pair first: ["Game.Point3", "10, 64, 20"]
		var pair [2]json.RawMessage
		if json.Unmarshal(coordRaw, &pair) == nil {
			var s string
			if json.Unmarshal(pair[1], &s) == nil {
				fmt.Sscanf(s, "%d, %d, %d", &bed.X, &bed.Y, &bed.Z) //nolint:errcheck
			}
		}
	}
	// Owner: ["System.Guid", "xxxxxxxx-xxxx-..."]
	if ownerRaw, ok := comp["Owner"]; ok {
		var pair [2]json.RawMessage
		if json.Unmarshal(ownerRaw, &pair) == nil {
			json.Unmarshal(pair[1], &bed.OwnerGUID) //nolint:errcheck
		}
	}
}

// parseInventory parses an inventory component from its raw JSON.
//
// Expected format (from ComponentInventoryBase.Save):
//
//	{
//	  "Id":         ["int", 5],
//	  "SlotsCount": ["int", 16],
//	  "Slots": {
//	    "Slot3": { "Contents": ["int", 40960], "Count": ["int", 5] },
//	    ...
//	  }
//	}
func parseInventory(raw json.RawMessage) *Inventory {
	var comp map[string]json.RawMessage
	if json.Unmarshal(raw, &comp) != nil {
		return nil
	}

	inv := &Inventory{}

	// Id
	if idRaw, ok := comp["Id"]; ok {
		inv.ID = intFromPair(idRaw)
	}

	// Slots
	slotsRaw, ok := comp["Slots"]
	if !ok {
		return inv
	}
	var slotsMap map[string]json.RawMessage
	if json.Unmarshal(slotsRaw, &slotsMap) != nil {
		return inv
	}
	for key, slotRaw := range slotsMap {
		// Key format: "Slot0", "Slot1", ...
		if len(key) <= 4 || key[:4] != "Slot" {
			continue
		}
		idx, err := strconv.Atoi(key[4:])
		if err != nil {
			continue
		}
		var slotFields map[string]json.RawMessage
		if json.Unmarshal(slotRaw, &slotFields) != nil {
			continue
		}
		slot := InventorySlot{
			Index:    idx,
			Contents: intFromPair(slotFields["Contents"]),
			Count:    intFromPair(slotFields["Count"]),
		}
		if slot.Count > 0 {
			inv.Slots = append(inv.Slots, slot)
		}
	}
	return inv
}

// intFromPair extracts an int32 from a ["typeName", number] JSON pair.
func intFromPair(raw json.RawMessage) int32 {
	if raw == nil {
		return 0
	}
	var pair [2]json.RawMessage
	if json.Unmarshal(raw, &pair) != nil {
		return 0
	}
	var n int32
	json.Unmarshal(pair[1], &n) //nolint:errcheck
	return n
}
