package chef

import (
	"encoding/binary"
	"fmt"

	"github.com/loeredami/ungo"
)

var recipeRegistry *ungo.Registry[Recipe] = ungo.NewRegistry[Recipe](1024)

func init() {
	recipeRegistry.Register("int32", ungo.NewLazy(func() Recipe {
		return RecipeInt32{}
	}))
	recipeRegistry.Register("int64", ungo.NewLazy(func() Recipe {
		return RecipeInt64{}
	}))
	recipeRegistry.Register("int16", ungo.NewLazy(func() Recipe {
		return RecipeInt16{}
	}))
	recipeRegistry.Register("int8", ungo.NewLazy(func() Recipe {
		return RecipeInt8{}
	}))
	recipeRegistry.Register("uint64", ungo.NewLazy(func() Recipe {
		return RecipeUInt64{}
	}))
	recipeRegistry.Register("uint32", ungo.NewLazy(func() Recipe {
		return RecipeUInt32{}
	}))
	recipeRegistry.Register("uint16", ungo.NewLazy(func() Recipe {
		return RecipeUInt16{}
	}))
	recipeRegistry.Register("uint8", ungo.NewLazy(func() Recipe {
		return RecipeUInt8{}
	}))
	recipeRegistry.Register("float64", ungo.NewLazy(func() Recipe {
		return RecipeFloat64{}
	}))
	recipeRegistry.Register("float32", ungo.NewLazy(func() Recipe {
		return RecipeFloat32{}
	}))
	recipeRegistry.Register("bool", ungo.NewLazy(func() Recipe {
		return RecipeBool{}
	}))
	recipeRegistry.Register("string", ungo.NewLazy(func() Recipe {
		return RecipeString{}
	}))
	recipeRegistry.Register("dict", ungo.NewLazy(func() Recipe {
		return RecipeDict{}
	}))
	recipeRegistry.Register("list", ungo.NewLazy(func() Recipe {
		return RecipeList{}
	}))
}

func RegisterRecipe(name string, recipe ungo.Lazy[Recipe]) {
	recipeRegistry.Register(name, recipe)
}

func GetRecipe(name string) ungo.Optional[Recipe] {
	return recipeRegistry.Get(name)
}

func Encode(s Serializable) ([]byte, error) {
	var result []byte
	var err error
	var wasPresent bool = false
	TYPE_ID := s.TypeID() + "\000\001\000"
	type_ID_HEADER := []byte(TYPE_ID)
	encodedSize := make([]byte, 4)

	recipeRegistry.Get(s.TypeID()).IfPresent(func(r Recipe) {
		result, err = r.Serialize(recipeRegistry, s)
		wasPresent = true
	})
	if !wasPresent {
		return nil, fmt.Errorf("no recipe found for type %s", s.TypeID())
	}
	if err != nil {
		return nil, err
	}
	SIZE := len(result)
	binary.LittleEndian.PutUint32(encodedSize, uint32(SIZE))

	FULL_HEADER := append(type_ID_HEADER, encodedSize...)
	return append(FULL_HEADER, result...), nil
}

func DecodeNext(data []byte) (Serializable, int, error) {
	var typeIDB []byte
	var delimiterIdx = -1
	for i := 0; i <= len(data)-3; i++ {
		if data[i] == 0 && data[i+1] == 1 && data[i+2] == 0 {
			typeIDB = data[:i]
			delimiterIdx = i
			break
		}
	}
	if delimiterIdx == -1 {
		return nil, 0, fmt.Errorf("no type ID header found")
	}
	typeID := string(typeIDB)

	sizeOffset := delimiterIdx + 3
	if len(data) < sizeOffset+4 {
		return nil, 0, fmt.Errorf("insufficient data for size header")
	}
	size := binary.LittleEndian.Uint32(data[sizeOffset : sizeOffset+4])

	payloadStart := sizeOffset + 4
	payloadEnd := payloadStart + int(size)
	if len(data) < payloadEnd {
		return nil, 0, fmt.Errorf("insufficient data for payload")
	}

	var result Serializable
	var err error
	var wasPresent bool = false

	recipeRegistry.Get(typeID).IfPresent(func(r Recipe) {
		result, err = r.Deserialize(recipeRegistry, data[payloadStart:payloadEnd])
		wasPresent = true
	})
	if !wasPresent {
		return nil, 0, fmt.Errorf("no recipe found for type %s", typeID)
	}
	if err != nil {
		return nil, 0, err
	}
	return result, payloadEnd, nil
}

func Decode(data []byte) (Serializable, error) {
	val, _, err := DecodeNext(data)
	return val, err
}
