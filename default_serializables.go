package chef

type SerializableInt32 struct {
	Value int32
}

func (s SerializableInt32) TypeID() string {
	return "int32"
}

type SerializableInt64 struct {
	Value int64
}

func (s SerializableInt64) TypeID() string {
	return "int64"
}

type SerializableInt16 struct {
	Value int16
}

func (s SerializableInt16) TypeID() string {
	return "int16"
}

type SerializableInt8 struct {
	Value int8
}

func (s SerializableInt8) TypeID() string {
	return "int8"
}

type SerializableUInt64 struct {
	Value uint64
}

func (s SerializableUInt64) TypeID() string {
	return "uint64"
}

type SerializableUInt32 struct {
	Value uint32
}

func (s SerializableUInt32) TypeID() string {
	return "uint32"
}

type SerializableUInt16 struct {
	Value uint16
}

func (s SerializableUInt16) TypeID() string {
	return "uint16"
}

type SerializableUInt8 struct {
	Value uint8
}

func (s SerializableUInt8) TypeID() string {
	return "uint8"
}

type SerializableString struct {
	Value string
}

func (s SerializableString) TypeID() string {
	return "string"
}

type SerializableBool struct {
	Value bool
}

func (s SerializableBool) TypeID() string {
	return "bool"
}

type SerializableFloat32 struct {
	Value float32
}

func (s SerializableFloat32) TypeID() string {
	return "float32"
}

type SerializableFloat64 struct {
	Value float64
}

func (s SerializableFloat64) TypeID() string {
	return "float64"
}

type SerializableList struct {
	Value []Serializable
}

func (s SerializableList) TypeID() string {
	return "list"
}

type SerializableDict struct {
	Value map[Serializable]Serializable
}

func (s SerializableDict) TypeID() string {
	return "dict"
}
