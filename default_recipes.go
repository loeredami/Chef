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
	value := int32(0)
	for _, v := range b {
		value = (value << 8) | int32(v)
	}
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
	value := int64(0)
	for _, v := range b {
		value = (value << 8) | int64(v)
	}
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
	value := int16(0)
	for _, v := range b {
		value = (value << 8) | int16(v)
	}
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
	value := uint64(0)
	for _, v := range b {
		value = (value << 8) | uint64(v)
	}
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
	value := uint32(0)
	for _, v := range b {
		value = (value << 8) | uint32(v)
	}
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
	value := uint16(0)
	for _, v := range b {
		value = (value << 8) | uint16(v)
	}
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
	b := make([]byte, 0, 8*len(list.Value))
	for _, item := range list.Value {
		var itemBytes []byte
		var err error
		var hasRecipe bool = false
		r.Get(item.TypeID()).IfPresent(func(recipe Recipe) {
			itemBytes, err = recipe.Serialize(r, item)
			hasRecipe = true
		})
		if !hasRecipe {
			return nil, fmt.Errorf("no recipe found for type: %s", item.TypeID())
		}
		if err != nil {
			return nil, err
		}
		b = append(b, itemBytes...)
	}
	return b, nil
}

func (recipe RecipeList) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	list := make([]Serializable, 0, len(b)/8)
	for i := 0; i < len(b); i += 8 {
		itemBytes := b[i : i+8]
		var item Serializable
		var err error
		var hasRecipe bool = false
		r.Get("float64").IfPresent(func(recipe Recipe) {
			item, err = recipe.Deserialize(r, itemBytes)
			hasRecipe = true
		})
		if !hasRecipe {
			return nil, fmt.Errorf("no recipe found for type: float64")
		}
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return SerializableList{list}, nil
}

type RecipeDict struct{}

func (dict RecipeDict) TypeID() string {
	return "dict"
}

func (dict RecipeDict) Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	dict_value := s.(SerializableDict)
	b := make([]byte, 0, 8*len(dict_value.Value))
	for key, value := range dict_value.Value {
		var keyBytes, valueBytes []byte
		var err error
		var hasKeyRecipe, hasValueRecipe bool = false, false
		r.Get(key.TypeID()).IfPresent(func(recipe Recipe) {
			keyBytes, err = recipe.Serialize(r, key)
			hasKeyRecipe = true
		})
		r.Get(value.TypeID()).IfPresent(func(recipe Recipe) {
			valueBytes, err = recipe.Serialize(r, value)
			hasValueRecipe = true
		})
		if err != nil {
			return nil, err
		}
		if !hasKeyRecipe {
			return nil, fmt.Errorf("no recipe found for type: %s", key.TypeID())
		}
		if !hasValueRecipe {
			return nil, fmt.Errorf("no recipe found for type: %s", value.TypeID())
		}
		b = append(b, keyBytes...)
		b = append(b, valueBytes...)
	}
	return b, nil
}

func (dict RecipeDict) Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error) {
	dictValue := make(map[Serializable]Serializable)
	for {
		keyBytes := b[:8]
		valueBytes := b[8:]
		if len(keyBytes) == 0 {
			break
		}
		var hasKeyRecipe bool
		var hasValueRecipe bool
		var key, value Serializable
		var keyErr, valueErr error
		r.Get(string(keyBytes)).IfPresent(func(key_recipe Recipe) {
			key, keyErr = key_recipe.Deserialize(r, keyBytes)
			hasKeyRecipe = true
		})
		r.Get(string(valueBytes)).IfPresent(func(value_recipe Recipe) {
			value, valueErr = value_recipe.Deserialize(r, valueBytes)
			hasValueRecipe = true
		})
		if keyErr != nil {
			return nil, keyErr
		}
		if valueErr != nil {
			return nil, valueErr
		}
		if !hasKeyRecipe {
			return nil, fmt.Errorf("no recipe found for type: %s", string(keyBytes))
		}
		if !hasValueRecipe {
			return nil, fmt.Errorf("no recipe found for type: %s", string(valueBytes))
		}
		b = b[16:]
		dictValue[key] = value
	}
	return SerializableDict{dictValue}, nil
}
