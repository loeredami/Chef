package chef

import (
	"encoding/json"
	"testing"

	"github.com/loeredami/ungo"
)

func TestEncodeInt64(t *testing.T) {
	intVal := SerializableInt64{int64(42)}
	encoded, err := Encode(intVal)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.TypeID() != "int64" {
		t.Fatalf("expected type ID 'int64', got '%s'", decoded.TypeID())
	}

	if decoded.(SerializableInt64).Value != intVal.Value {
		t.Fatalf("expected value %v, got %v", intVal.Value, decoded.(SerializableInt64).Value)
	}
}

func TestEncodeUInt64(t *testing.T) {
	uintVal := SerializableUInt64{uint64(42)}
	encoded, err := Encode(uintVal)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.TypeID() != "uint64" {
		t.Fatalf("expected type ID 'uint64', got '%s'", decoded.TypeID())
	}

	if decoded.(SerializableUInt64).Value != uintVal.Value {
		t.Fatalf("expected value %v, got %v", uintVal.Value, decoded.(SerializableUInt64).Value)
	}
}

func TestEncodeFloat64(t *testing.T) {
	floatVal := SerializableFloat64{42.0}
	encoded, err := Encode(floatVal)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.TypeID() != "float64" {
		t.Fatalf("expected type ID 'float64', got '%s'", decoded.TypeID())
	}

	if decoded.(SerializableFloat64).Value != floatVal.Value {
		t.Fatalf("expected value %v, got %v", floatVal.Value, decoded.(SerializableFloat64).Value)
	}
}

func TestEncodeList(t *testing.T) {
	listVal := SerializableList{[]Serializable{SerializableInt64{1}, SerializableInt64{2}}}
	encoded, err := Encode(listVal)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.TypeID() != "list" {
		t.Fatalf("expected type ID 'list', got '%s'", decoded.TypeID())
	}

	if len(decoded.(SerializableList).Value) != len(listVal.Value) {
		t.Fatalf("expected list length %d, got %d", len(listVal.Value), len(decoded.(SerializableList).Value))
	}

	for i := range decoded.(SerializableList).Value {
		if decoded.(SerializableList).Value[i] != listVal.Value[i] {
			t.Fatalf("expected value %v, got %v", listVal.Value[i], decoded.(SerializableList).Value[i])
		}
	}
}

func TestEncodeDict(t *testing.T) {
	dictVal := SerializableDict{map[Serializable]Serializable{SerializableString{"a"}: SerializableInt64{1}, SerializableString{"b"}: SerializableInt64{2}}}
	encoded, err := Encode(dictVal)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.TypeID() != "dict" {
		t.Fatalf("expected type ID 'dict', got '%s'", decoded.TypeID())
	}

	if len(decoded.(SerializableDict).Value) != len(dictVal.Value) {
		t.Fatalf("expected dict length %d, got %d", len(dictVal.Value), len(decoded.(SerializableDict).Value))
	}

	for k, v := range decoded.(SerializableDict).Value {
		if v != dictVal.Value[k] {
			t.Fatalf("expected value %v, got %v", dictVal.Value[k], v)
		}
	}
}

type customType struct {
	Name string
	Age  int
}

func (c customType) TypeID() string {
	return "customType"
}

type customTypeRecipe struct {
	Value customType
}

func (c customTypeRecipe) TypeID() string {
	return "customType"
}

func (c customTypeRecipe) Serialize(reg *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	val := s.(customType)

	encoded, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func (c customTypeRecipe) Deserialize(reg *ungo.Registry[Recipe], data []byte) (Serializable, error) {
	var val customType
	if err := json.Unmarshal(data, &val); err != nil {
		return nil, err
	}
	return val, nil
}

func TestEncodeCustomType(t *testing.T) {
	RegisterRecipe("customType", ungo.NewLazy(func() Recipe {
		return customTypeRecipe{}
	}))

	val := customType{Name: "John Doe", Age: 63}
	encoded, err := Encode(val)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.TypeID() != "customType" {
		t.Fatalf("expected type ID 'customType', got '%s'", decoded.TypeID())
	}

	if decoded.(customType).Name != val.Name {
		t.Fatalf("expected name '%s', got '%s'", val.Name, decoded.(customType).Name)
	}

	if decoded.(customType).Age != val.Age {
		t.Fatalf("expected age %d, got %d", val.Age, decoded.(customType).Age)
	}
}

type customTypeInside struct {
	customTypeData       customType
	customTypeAttachment string
}

func (c customTypeInside) TypeID() string {
	return "customTypeInside"
}

type customTypeInsideRecipe struct {
	Value customTypeInside
}

func (c customTypeInsideRecipe) TypeID() string {
	return "customTypeInside"
}

func (c customTypeInsideRecipe) Serialize(reg *ungo.Registry[Recipe], s Serializable) ([]byte, error) {
	val := s.(customTypeInside)
	encoded, err := Encode(val.customTypeData)
	if err != nil {
		return nil, err
	}
	return append(encoded, []byte(val.customTypeAttachment)...), nil
}

func TestCustomTypeInsideRecipe(t *testing.T) {
	RegisterRecipe("customTypeInside", ungo.NewLazy(func() Recipe {
		return customTypeInsideRecipe{}
	}))

	val := customTypeInside{customTypeData: customType{Name: "Joe Johnson", Age: 42}, customTypeAttachment: "attached data"}
	serialized, err := Encode(val)
	if err != nil {
		t.Fatal(err)
	}
	serialized_object, err := Decode(serialized)
	if err != nil {
		t.Fatal(err)
	}

	if serialized_object.TypeID() != "customTypeInside" {
		t.Fatal("TypeID mismatch")
	}

	if serialized_object.(customTypeInside).customTypeData.Name != "Joe Johnson" {
		t.Fatal("Name mismatch")
	}
	if serialized_object.(customTypeInside).customTypeData.Age != 42 {
		t.Fatal("Age mismatch")
	}
	if serialized_object.(customTypeInside).customTypeAttachment != "attached data" {
		t.Fatal("Attachment mismatch")
	}
}

func (c customTypeInsideRecipe) Deserialize(reg *ungo.Registry[Recipe], data []byte) (Serializable, error) {
	customDataSizeBytes := data[len("customType\000\001\000") : len("customType\000\001\000")+4]
	customDataSize := int(customDataSizeBytes[0]) | int(customDataSizeBytes[1])<<8 | int(customDataSizeBytes[2])<<16 | int(customDataSizeBytes[3])<<24
	customData := data[:len("customType\000\001\000")+4+customDataSize]
	attachment := string(data[len("customType\000\001\000")+4+customDataSize:])
	decoded, err := Decode(customData)
	if err != nil {
		return nil, err
	}
	return customTypeInside{customTypeData: decoded.(customType), customTypeAttachment: attachment}, nil
}
