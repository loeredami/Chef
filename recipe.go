package chef

import "github.com/loeredami/ungo"

type Recipe interface {
	TypeID() string

	Serialize(*ungo.Registry[Recipe], Serializable) ([]byte, error)
	Deserialize(*ungo.Registry[Recipe], []byte) (Serializable, error)
}
