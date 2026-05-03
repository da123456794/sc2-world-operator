package define

// CellValue bit layout (int32):
//
//	bits  0- 9 : Contents (block ID), 0-1023
//	bits 10-13 : Light, 0-15
//	bits 14-31 : Data (block metadata), 0-262143

// ExtractContents returns the block ID (bits 0-9).
func ExtractContents(v int32) int32 { return v & 0x3FF }

// ExtractLight returns the light value (bits 10-13).
func ExtractLight(v int32) int32 { return (v & 0x3C00) >> 10 }

// ExtractData returns the block data (bits 14-31).
func ExtractData(v int32) int32 { return (v & -16384) >> 14 }

// ReplaceContents returns v with the block ID replaced.
func ReplaceContents(v, contents int32) int32 {
	return v ^ ((v ^ contents) & 0x3FF)
}

// ReplaceLight returns v with the light field replaced.
func ReplaceLight(v, light int32) int32 {
	return v ^ ((v ^ (light << 10)) & 0x3C00)
}

// ReplaceData returns v with the data field replaced.
func ReplaceData(v, data int32) int32 {
	return v ^ ((v ^ (data << 14)) & -16384)
}

// ShaftValue bit layout (int32):
//
//	bits  8-11 : Temperature, 0-15
//	bits 12-15 : Humidity, 0-15

// ExtractTemperature returns the temperature from a shaft value (bits 8-11).
func ExtractTemperature(shaft int32) int32 { return (shaft & 0xF00) >> 8 }

// ExtractHumidity returns the humidity from a shaft value (bits 12-15).
func ExtractHumidity(shaft int32) int32 { return (shaft & 0xF000) >> 12 }

// ReplaceTemperature returns shaft with the temperature field replaced.
func ReplaceTemperature(shaft, temp int32) int32 {
	return shaft ^ ((shaft ^ (temp << 8)) & 0xF00)
}

// ReplaceHumidity returns shaft with the humidity field replaced.
func ReplaceHumidity(shaft, humidity int32) int32 {
	return shaft ^ ((shaft ^ (humidity << 12)) & 0xF000)
}
