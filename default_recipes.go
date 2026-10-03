package chef

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/loeredami/ungo"
)

type RecipeInt32 struct{}

func (recipe RecipeInt32) TypeID() string {
	return "int32"
}

func (recipe RecipeInt32) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableInt32).Value
	b := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b[i] = byte(value >> (8 * i))
	}
	return b, nil
}

func (recipe RecipeInt32) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 4 {
		return nil, fmt.Errorf("invalid int32 byte length: %d", len(b))
	}
	value := int32(binary.LittleEndian.Uint32(b))
	return SerializableInt32{value}, nil
}

type RecipeInt64 struct{}

func (recipe RecipeInt64) TypeID() string {
	return "int64"
}

func (recipe RecipeInt64) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableInt64).Value
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(value >> (8 * i))
	}
	return b, nil
}

func (recipe RecipeInt64) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 8 {
		return nil, fmt.Errorf("invalid int64 byte length: %d", len(b))
	}
	value := int64(binary.LittleEndian.Uint64(b))
	return SerializableInt64{value}, nil
}

type RecipeInt16 struct{}

func (recipe RecipeInt16) TypeID() string {
	return "int16"
}

func (recipe RecipeInt16) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableInt16).Value
	b := make([]byte, 2)
	for i := 0; i < 2; i++ {
		b[i] = byte(value >> (8 * i))
	}
	return b, nil
}

func (recipe RecipeInt16) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 2 {
		return nil, fmt.Errorf("invalid int16 byte length: %d", len(b))
	}
	value := int16(binary.LittleEndian.Uint16(b))
	return SerializableInt16{value}, nil
}

type RecipeInt8 struct{}

func (recipe RecipeInt8) TypeID() string {
	return "int8"
}

func (recipe RecipeInt8) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableInt8).Value
	b := make([]byte, 1)
	b[0] = byte(value)
	return b, nil
}

func (recipe RecipeInt8) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 1 {
		return nil, fmt.Errorf("invalid int8 byte length: %d", len(b))
	}
	return SerializableInt8{int8(b[0])}, nil
}

type RecipeUInt64 struct{}

func (recipe RecipeUInt64) TypeID() string {
	return "uint64"
}

func (recipe RecipeUInt64) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableUInt64).Value
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(value >> (8 * i))
	}
	return b, nil
}

func (recipe RecipeUInt64) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 8 {
		return nil, fmt.Errorf("invalid uint64 byte length: %d", len(b))
	}
	value := binary.LittleEndian.Uint64(b)
	return SerializableUInt64{value}, nil
}

type RecipeUInt32 struct{}

func (recipe RecipeUInt32) TypeID() string {
	return "uint32"
}

func (recipe RecipeUInt32) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableUInt32).Value
	b := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b[i] = byte(value >> (8 * i))
	}
	return b, nil
}

func (recipe RecipeUInt32) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 4 {
		return nil, fmt.Errorf("invalid uint32 byte length: %d", len(b))
	}
	value := binary.LittleEndian.Uint32(b)
	return SerializableUInt32{value}, nil
}

type RecipeUInt16 struct{}

func (recipe RecipeUInt16) TypeID() string {
	return "uint16"
}

func (recipe RecipeUInt16) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableUInt16).Value
	b := make([]byte, 2)
	for i := 0; i < 2; i++ {
		b[i] = byte(value >> (8 * i))
	}
	return b, nil
}

func (recipe RecipeUInt16) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 2 {
		return nil, fmt.Errorf("invalid uint16 byte length: %d", len(b))
	}
	value := binary.LittleEndian.Uint16(b)
	return SerializableUInt16{value}, nil
}

type RecipeUInt8 struct{}

func (recipe RecipeUInt8) TypeID() string {
	return "uint8"
}

func (recipe RecipeUInt8) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableUInt8).Value
	b := make([]byte, 1)
	b[0] = byte(value)
	return b, nil
}

func (recipe RecipeUInt8) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 1 {
		return nil, fmt.Errorf("invalid uint8 byte length: %d", len(b))
	}
	value := uint8(b[0])
	return SerializableUInt8{value}, nil
}

type RecipeString struct{}

func (recipe RecipeString) TypeID() string {
	return "string"
}

func (recipe RecipeString) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableString).Value
	return []byte(value), nil
}

func (recipe RecipeString) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	value := string(b)
	return SerializableString{value}, nil
}

type RecipeBool struct{}

func (recipe RecipeBool) TypeID() string {
	return "bool"
}

func (recipe RecipeBool) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableBool).Value
	b := make([]byte, 1)
	b[0] = ungo.If(value, byte(1), byte(0))
	return b, nil
}

func (recipe RecipeBool) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 1 {
		return nil, fmt.Errorf("invalid bool byte length: %d", len(b))
	}
	value := b[0] != 0
	return SerializableBool{value}, nil
}

type RecipeFloat32 struct{}

func (recipe RecipeFloat32) TypeID() string {
	return "float32"
}

func (recipe RecipeFloat32) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableFloat32).Value
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(value))
	return b, nil
}

func (recipe RecipeFloat32) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 4 {
		return nil, fmt.Errorf("invalid float32 byte length: %d", len(b))
	}
	value := math.Float32frombits(binary.LittleEndian.Uint32(b))
	return SerializableFloat32{value}, nil
}

type RecipeFloat64 struct{}

func (recipe RecipeFloat64) TypeID() string {
	return "float64"
}

func (recipe RecipeFloat64) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	value := s.(SerializableFloat64).Value
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(value))
	return b, nil
}

func (recipe RecipeFloat64) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	if len(b) != 8 {
		return nil, fmt.Errorf("invalid float64 byte length: %d", len(b))
	}
	value := math.Float64frombits(binary.LittleEndian.Uint64(b))
	return SerializableFloat64{value}, nil
}

type RecipeList struct{}

func (recipe RecipeList) TypeID() string {
	return "list"
}

func (recipe RecipeList) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	list := s.(SerializableList)
	var b []byte
	for _, item := range list.Value {
		itemBytes, err := Encode(item)
		if err != nil {
			return nil, err
		}
		b = append(b, itemBytes...)
	}
	return b, nil
}

func (recipe RecipeList) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	var list []Serializable
	offset := 0
	for offset < len(b) {
		item, n, err := DecodeNext(b[offset:])
		if err != nil {
			return nil, err
		}
		list = append(list, item)
		offset += n
	}
	return SerializableList{list}, nil
}

type RecipeDict struct{}

func (dict RecipeDict) TypeID() string {
	return "dict"
}

func (dict RecipeDict) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	dictValue := s.(SerializableDict)
	var b []byte
	for key, value := range dictValue.Value {
		keyBytes, err := Encode(key)
		if err != nil {
			return nil, err
		}
		valBytes, err := Encode(value)
		if err != nil {
			return nil, err
		}
		b = append(b, keyBytes...)
		b = append(b, valBytes...)
	}
	return b, nil
}

func (dict RecipeDict) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	dictValue := make(map[Serializable]Serializable)
	offset := 0
	for offset < len(b) {
		key, nKey, err := DecodeNext(b[offset:])
		if err != nil {
			return nil, err
		}
		offset += nKey

		val, nVal, err := DecodeNext(b[offset:])
		if err != nil {
			return nil, err
		}
		offset += nVal

		dictValue[key] = val
	}
	return SerializableDict{dictValue}, nil
}
