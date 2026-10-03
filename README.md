# Chef

`chef` is a binary serialization and deserialization library for Go built on top of `ungo` registry types[cite: 5, 6].

## Core Interfaces

To make a type serializable, implement `Serializable`[cite: 4]:

```go
type Serializable interface {
    TypeID() string
}
```

To define custom serialization logic, implement `Recipe`[cite: 5]:

```go
type Recipe interface {
    TypeID() string
    Serialize(r *ungo.Registry[Recipe], s Serializable) ([]byte, error)
    Deserialize(r *ungo.Registry[Recipe], b []byte) (Serializable, error)
}
```

## Built-In Types

`chef` provides default serializables and recipes for standard Go types out of the box[cite: 1, 2, 6]:
- **Signed Integers**: `SerializableInt8`, `SerializableInt16`, `SerializableInt32`, `SerializableInt64`[cite: 1]
- **Unsigned Integers**: `SerializableUInt8`, `SerializableUInt16`, `SerializableUInt32`, `SerializableUInt64`[cite: 1]
- **Floats & Bools**: `SerializableFloat32`, `SerializableFloat64`, `SerializableBool`[cite: 1]
- **Strings**: `SerializableString`[cite: 1]
- **Containers**: `SerializableList` (`[]Serializable`), `SerializableDict` (`map[Serializable]Serializable`)[cite: 1]

## Usage Example

### Basic Encoding & Decoding

```go
package main

import (
	"fmt"
	"log"

	"[github.com/loeredami/chef](https://github.com/loeredami/chef)"
)

func main() {
	list := chef.SerializableList{
		Value: []chef.Serializable{
			chef.SerializableString{Value: "hello"},
			chef.SerializableInt64{Value: 42},
		},
	}

	encoded, err := chef.Encode(list)
	if err != nil {
		log.Fatalf("Encode failed: %v", err)
	}

	decoded, err := chef.Decode(encoded)
	if err != nil {
		log.Fatalf("Decode failed: %v", err)
	}

	result := decoded.(chef.SerializableList)
	fmt.Println("Decoded list items:", len(result.Value))
}
```

### Registering Custom Recipes

Register custom recipes using `chef.RegisterRecipe` with `ungo.NewLazy`[cite: 6]:

```go
package main

import (
	"encoding/json"
	"fmt"

	"[github.com/loeredami/chef](https://github.com/loeredami/chef)"
	"[github.com/loeredami/ungo](https://github.com/loeredami/ungo)"
)

type CustomUser struct {
	Name string `json:"name"`
}

func (u CustomUser) TypeID() string {
	return "customUser"
}

type CustomUserRecipe struct{}

func (r CustomUserRecipe) TypeID() string {
	return "customUser"
}

func (r CustomUserRecipe) Serialize(reg *ungo.Registry[chef.Recipe], s chef.Serializable) ([]byte, error) {
	return json.Marshal(s.(CustomUser))
}

func (r CustomUserRecipe) Deserialize(reg *ungo.Registry[chef.Recipe], b []byte) (chef.Serializable, error) {
	var user CustomUser
	err := json.Unmarshal(b, &user)
	return user, err
}

func main() {
	chef.RegisterRecipe("customUser", ungo.NewLazy(func() chef.Recipe {
		return CustomUserRecipe{}
	}))

	encoded, _ := chef.Encode(CustomUser{Name: "Alice"})
	decoded, _ := chef.Decode(encoded)

	fmt.Printf("User: %+v\n", decoded.(CustomUser))
}
```
