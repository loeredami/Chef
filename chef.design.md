# Chef — Serialization Library Design

**Status:** Draft  
**Target:** Go  
**Primary use case:** Game save states and engine persistence  
**Registry foundation:** `github.com/loeredami/ungo`

## 1. Concept

Chef is a binary serialization library built around one assumption:

> The program decoding a value already knows, or can obtain, the recipe for its type.

A **recipe** defines how a type is encoded and decoded. A **cookbook** is the registry containing those recipes. A serialized package is the **meal** produced from a recipe and a value.

Chef is intended to encode:

- primitive values;
- strings and byte data;
- lists/slices and arrays;
- dictionaries/maps;
- optional values;
- custom types;
- custom types containing other custom types;
- containers of custom or built-in types;
- recursively nested combinations of all of the above.

The format is primarily intended for persistent game state, where compactness, predictable decoding, modding support, and schema evolution matter more than human readability.

---

## 2. Goals

### Core goals

1. Compact binary representation.
2. Stable, human-readable type identification.
3. Recursive composition of types.
4. User-defined serialization recipes.
5. Strong support for modded content.
6. Save-file longevity and migration.
7. Predictable performance.
8. Safe handling of malformed input.

### Non-goals for the first version

Chef is not intended to be:

- human-readable as its primary representation;
- a database;
- a memory dump;
- a general-purpose arbitrary-object deserializer;
- a network transport protocol;
- a system that can decode an unregistered custom type.

---

## 3. Terminology

### Recipe

A serialization definition for one logical type.

A recipe contains:

- a stable Type ID;
- a recipe/schema version;
- an encoder;
- a decoder;
- a construction/allocation mechanism.

### Cookbook / Registry

The runtime mapping from stable Type IDs to recipes.

### Meal / Package

A complete Chef-encoded document, normally representing one root value such as a save state.

### Global Type ID

The canonical string identifying a logical type.

Example:

```text
core.game.Player
```

or:

```text
com.example.inventory.Item
```

### Local Type ID

A compact integer assigned only inside one package. It references a global Type ID through the package's type table.

---

# 4. Type identity

Chef should use **string Type IDs as the persistent identity system**.

This is especially important for mods.

A mod author should be able to introduce a type without coordinating a numeric identifier with the base game or with other mods.

For example:

```text
game.Player
game.Item

mod.example.weapons.Weapon
mod.example.weapons.Ammo

mod.example.quest.Bounty
```

The string is part of the serialization contract.

## 4.1 Type IDs must be globally unique

The format should require a canonical namespace convention.

For example:

```text
game.*
mod.<author>.<mod>.*
com.<organization>.<project>.*
```

Chef itself should not require one exact namespace syntax, but the documentation should strongly recommend namespacing.

Bad:

```text
Item
Player
Weapon
```

Better:

```text
game.Item
game.Player
mod.example.weapons.Weapon
```

This avoids collisions between independently developed mods.

---

## 5. Type ID stability

A Type ID identifies the **logical serialized type**, not the Go implementation type.

Changing:

```go
type Player struct { ... }
```

to:

```go
type Character struct { ... }
```

does not require changing its serialized Type ID.

For example:

```text
game.Player
```

can remain stable even when the implementation is refactored.

Similarly, a mod should not change a Type ID merely because its Go package or source file moves.

Changing a Type ID should be treated as a serialization compatibility decision.

---

# 6. Modding model

String IDs allow Chef to treat mods as first-class serialization participants.

A save can contain:

```text
game.Player
game.Inventory
mod.example.weapons.Weapon
mod.example.weapons.Ammo
mod.othermod.quest.Bounty
```

The base game does not need to know every possible type ahead of time.

At load time, the active registry attempts to resolve each Type ID.

Conceptually:

```text
Save File
   |
   +-- game.Player
   |
   +-- mod.example.weapons.Weapon
   |
   +-- mod.othermod.quest.Bounty
          |
          v
     Current Registry
          |
     +----+----+
     |         |
   found     missing
     |         |
   decode     error /
              fallback
```

This makes the distinction between:

> "This type is unknown because the required mod is not installed"

and:

> "This save file is corrupt"

explicit.

---

# 7. Missing mod behavior

Chef itself should not silently discard an unknown custom type.

A decoder encountering:

```text
mod.example.weapons.Weapon
```

without a registered recipe should return an `ErrUnknownType`.

Higher-level game code may then decide what to do.

Possible application-level policies include:

- reject the save;
- load with degraded content;
- preserve the unknown data for later;
- substitute a placeholder object;
- invoke a mod compatibility layer.

The serialization layer should not decide which behavior is correct for the game.

---

# 8. Compactness: global strings vs local references

Using strings directly everywhere would hurt compactness.

Chef should therefore separate:

```text
Global identity
```

from:

```text
Package-local reference
```

A package contains a type table.

Example:

```text
TYPE TABLE

local 0 -> game.Player
local 1 -> game.Item
local 2 -> mod.example.weapons.Weapon
local 3 -> List<local 1>
```

The actual value payloads then use the compact local IDs:

```text
0
1
2
3
```

instead of repeatedly storing long strings.

This gives Chef both:

- mod-friendly global string identities;
- compact binary payloads.

---

# 9. Type table

Conceptually:

```text
+----------------------+
| Magic                |
+----------------------+
| Format Version       |
+----------------------+
| Flags                |
+----------------------+
| Type Table Count     |
+----------------------+
| Type Table            |
+----------------------+
| Root Type Reference  |
+----------------------+
| Root Value           |
+----------------------+
| Optional Integrity   |
+----------------------+
```

A type-table entry may contain:

```text
local ID
global Type ID
recipe/schema version
type descriptor
```

The exact representation can be optimized later.

---

# 10. Composite type identities

Built-in composite types should also have a canonical representation.

For example:

```text
List<game.Item>
Map<string,game.Item>
Optional<game.Player>
Array<3,float32>
```

The type table can either:

1. encode these structurally; or
2. represent them through reserved built-in type constructors.

Structural representation is preferable because it avoids treating every possible combination as a globally registered string.

The important distinction is:

```text
global custom identity
```

versus:

```text
composition of existing types
```

---

# 11. Recipe model

Conceptually:

```go
type Recipe interface {
    TypeID() string
    Version() uint64

    New() any
    Encode(*Encoder, any) error
    Decode(*Decoder, any) error
}
```

The public API can use typed generic adapters so normal users do not need to work with `any`.

For example:

```go
chef.Register[PlayerState](
    "game.PlayerState",
    playerRecipe,
)
```

The Type ID is part of the recipe contract.

---

# 12. Type categories

Chef distinguishes three broad categories.

### Primitive types

Built-in scalar types such as:

```text
bool
signed integers
unsigned integers
float32
float64
string
[]byte
```

### Composite types

Built-in recursive type constructors such as:

```text
List<T>
Array<N,T>
Map<K,V>
Optional<T>
Tuple<...>
```

### Custom types

Application- or mod-defined types resolved through the registry:

```text
game.Player
game.WorldState
mod.example.weapons.Weapon
```

A custom recipe may contain any other Chef-supported type.

---

# 13. Known types vs dynamic types

This is one of the main compactness optimizations.

Suppose:

```go
type Player struct {
    Name      string
    Health    int32
    Inventory []Item
}
```

The `Player` recipe already knows:

```text
Name      = string
Health    = int32
Inventory = List<Item>
```

Therefore the serialized fields do not need to repeat their type identifiers.

A statically known field contains only its value.

By contrast:

```go
Data any
```

does not have a statically known concrete type.

It therefore needs:

```text
[local type reference][value]
```

This allows Chef to support polymorphism without forcing type metadata into every ordinary field.

---

# 14. Recursive composition

Nested custom types require no special mechanism.

Example:

```go
type Weapon struct {
    ID       uint32
    Damage   int16
    Material Material
}
```

The `Weapon` recipe delegates to the `Material` recipe.

The same mechanism must support:

```text
List<Weapon>
Map<string, Weapon>
Map<PlayerID, List<Item>>
List<Map<string, CustomType>>
Optional<List<CustomType>>
```

Arbitrarily nested supported types should remain legal.

---

# 15. Containers

## Lists

A homogeneous list can encode:

```text
[element count][element payloads...]
```

Its element type comes from the containing type definition, so it does not need to be repeated.

A dynamically typed list may store a local type reference for each element.

## Maps

A map can encode:

```text
[entry count][key][value]...
```

Its key and value types are known from the map type definition when statically typed.

Map ordering requires an explicit policy. Chef should eventually provide deterministic map encoding for applications that require reproducible byte output.

## Optionals

Conceptually:

```text
[present flag][payload if present]
```

## Arrays

Fixed-size arrays do not need to write their length because the type definition already contains it.

Example:

```text
Array<3,float32>
```

means exactly three `float32` values.

---

# 16. Custom type layout

Chef should not use Go struct memory layout as the wire format.

A recipe explicitly defines serialization order.

Example:

```text
Player recipe v3:

1. Name
2. Position
3. Health
4. Inventory
```

This provides:

- no field-name overhead;
- no dependency on Go field names;
- stable field ordering;
- explicit compatibility rules;
- predictable binary layout.

Reflection-based convenience APIs may be added later, but they should generate an explicit recipe rather than silently turning Go implementation details into the persistent wire contract.

---

# 17. Versioning

Game saves often outlive the program version that created them.

Chef therefore needs two different versions.

## Chef format version

Describes the serialization format itself.

Example:

```text
Chef format = 2
```

## Recipe version

Describes one logical type's data contract.

Example:

```text
game.Player recipe = 7
```

These remain independent.

---

# 18. Schema migration

A recipe should be able to evolve:

```text
game.Player v1
    Name
    Health

game.Player v2
    Name
    Health
    Position

game.Player v3
    Name
    Health
    Position
    Inventory
```

Chef should distinguish:

- direct wire compatibility;
- migration from an older representation.

A future migration API can conceptually support:

```text
game.Player v1 -> game.Player v2
game.Player v2 -> game.Player v3
```

Chained migrations are preferable to implementing every version pair independently.

New fields should be able to receive defaults when loading older saves.

---

# 19. Mod schema evolution

Mod content creates an additional versioning problem.

A mod may have:

```text
mod.example.weapons.Weapon
```

with recipe version:

```text
1
```

and later evolve to:

```text
mod.example.weapons.Weapon
version = 2
```

The Type ID remains stable while the recipe version changes.

This is preferable to embedding the version directly into the global Type ID:

```text
mod.example.weapons.Weapon.v2
```

The Type ID identifies the logical concept; the recipe version identifies the serialization schema.

An exception may be appropriate when the author intentionally creates an entirely new logical type.

---

# 20. Registry architecture

`ungo` can provide the underlying registry utility, but Chef should expose a serialization-specific abstraction over it.

Conceptually:

```text
Chef Registry
    |
    +-- Type ID string -> Recipe
    |
    +-- Runtime type -> Recipe resolver
```

The registry should support:

```text
game.Player
mod.example.weapons.Weapon
mod.example.quest.Bounty
```

as distinct namespaces.

## Registration conflicts

Two recipes attempting to register the same Type ID in the same registry should produce an error.

Silent replacement is dangerous because it could make the meaning of an existing save dependent on load order.

For modding, conflict detection should happen early during startup.

---

# 21. Registry scope

A global registry is useful for normal applications, but Chef should also support explicitly constructed registries.

That enables:

- tests with isolated schemas;
- tools working on different game versions;
- plugin-specific registries;
- multiple game configurations in one process.

Once serialization begins, the active registry should ideally be treated as immutable.

---

# 22. Dynamic values and mod types

Consider:

```go
type Component struct {
    Data any
}
```

A base game may serialize:

```text
game.Transform
```

while a mod may provide:

```text
mod.example.vehicles.VehicleComponent
```

The `Data` field therefore stores its local type reference.

Example:

```text
Component
    Data
        type = local 12
        value = ...
```

The local type table maps:

```text
local 12 -> mod.example.vehicles.VehicleComponent
```

This means the base game does not need to know the mod's concrete Go type at compile time.

It only needs an appropriate recipe to be registered when the save is loaded.

---

# 23. Unknown types and data preservation

For a mod-friendly save system, it may be useful to distinguish between:

```text
unknown type
```

and:

```text
unknown type whose raw serialized bytes can still be preserved
```

A future extension could provide an opaque value:

```go
type UnknownValue struct {
    TypeID  string
    Version uint64
    Raw     []byte
}
```

This would allow an application to load a save, modify unrelated data, and write the unknown mod data back out without understanding it.

This is particularly valuable for mod interoperability.

However, raw preservation should be considered a higher-level feature because not every encoding layout is independently self-contained.

---

# 24. Save-file integrity

The serializer should reject:

- truncated values;
- invalid type definitions;
- unknown required types;
- impossible lengths;
- unsupported recipe versions.

Optional outer integrity metadata can provide:

```text
[Chef package]
[checksum/hash]
```

Possible future mechanisms include CRCs, cryptographic hashes, or application-level signatures.

Integrity should remain distinct from serialization semantics.

---

# 25. Decoder resource limits

Deserialization is an input boundary.

The decoder should support limits for:

- recursion depth;
- maximum container length;
- maximum string size;
- maximum byte payload;
- maximum type-table size;
- maximum total decoded allocation.

Invalid length fields should be rejected before large allocations occur.

---

# 26. Error model

Errors should preserve enough context to identify where decoding failed.

Example:

```text
chef: failed decoding
       Player.Inventory[37].Weapon
       type: mod.example.weapons.Weapon
       expected recipe version: 4
       input ended after 12 bytes
```

Useful categories include:

```text
ErrInvalidFormat
ErrUnsupportedFormat
ErrUnknownType
ErrUnsupportedRecipeVersion
ErrInvalidValue
ErrLimitExceeded
ErrUnexpectedEOF
ErrChecksumMismatch
ErrMigrationFailed
```

Errors should remain composable with Go's standard error mechanisms.

---

# 27. Determinism

Deterministic encoding is useful for:

- tests;
- binary save diffs;
- caching;
- reproducible tooling;
- content-addressed storage.

The following should be deterministic:

- header layout;
- type-table ordering;
- primitive representation;
- recipe-defined field order.

Maps need an explicit strategy because Go map iteration order is not deterministic.

Chef should eventually offer a deterministic mode, even if normal gameplay saves use a faster unordered map representation.

---

# 28. Type-table ordering

Because global Type IDs are strings, the implementation should **not** use registry insertion order as the package-local numbering scheme.

For reproducible output, Chef can assign local IDs deterministically.

For example, during package construction:

```text
game.Item
game.Player
mod.example.weapons.Weapon
```

could become:

```text
local 0 = game.Item
local 1 = game.Player
local 2 = mod.example.weapons.Weapon
```

based on a deterministic traversal or canonical ordering.

This is an implementation detail and should not affect semantic decoding.

---

# 29. Game-engine model

Chef should serialize **persistent state**, not arbitrary live engine objects.

For example, a GPU texture or runtime object should normally be represented by a stable asset identifier:

```text
textures/player/armor_01
```

rather than by its live runtime object.

The engine can reconstruct runtime resources while loading the save.

This keeps saves independent of process memory and transient engine state.

---

# 30. References and cycles

There is an unresolved design question around object identity.

Possible models:

### Tree-only

Simple, but repeated references are encoded repeatedly and cycles are unsupported.

### Fully reference-aware

Supports shared objects and cycles, but requires object/reference tables and additional complexity.

### Opt-in references

Tree serialization is the default, with reference tracking enabled for types that actually require identity preservation.

The third approach is the preferred direction for the initial design because pointer identity should not automatically imply persistence semantics.

---

# 31. Streaming API

Chef should eventually support both memory and stream APIs.

High-level:

```go
bytes, err := chef.Marshal(value)
err := chef.Unmarshal(bytes, &value)
```

Streaming:

```go
err := chef.NewEncoder(w).Encode(value)
err := chef.NewDecoder(r).Decode(&value)
```

Registry:

```go
registry.Register(recipe)
recipe, ok := registry.Lookup(typeID)
```

The exact public API should remain separate from the wire-format specification.

---

# 32. Recommended MVP

The first implementation should support:

### Scalars

```text
bool
int8/int16/int32/int64
uint8/uint16/uint32/uint64
float32/float64
string
[]byte
```

### Containers

```text
[]T
[N]T
map[K]V
optional T
```

### Custom types

```text
registered recipe with string Type ID
```

### Dynamic values

```text
any / interface-like values
```

Everything else should be composable from these.

---

# 33. Modding-oriented MVP behavior

For the first release, the following should be explicit guarantees.

### A mod can define a type independently

Example:

```text
mod.coolmod.weapon.Weapon
```

No central numeric-ID allocation is required.

### A mod's Type ID remains stable

Its Go package name, source location, or registration order does not define its persistent identity.

### Save files identify missing mods

The decoder can report:

```text
unknown type: mod.coolmod.weapon.Weapon
```

instead of only reporting an opaque numeric identifier.

### Multiple mods can coexist

Namespaced IDs prevent ordinary collisions.

### Load order does not determine meaning

Registration order must never change the meaning of a save.

---

# 34. Example

Given:

```go
type Player struct {
    Name      string
    Health    int32
    Inventory []Item
}

type Item struct {
    ID       uint32
    Quantity uint16
}
```

with a mod-defined type:

```go
type Weapon struct {
    ID     uint32
    Damage int16
}
```

a package might conceptually contain:

```text
CHEF
format = 1

TYPE TABLE
    local 0 = game.Player
    local 1 = game.Item
    local 2 = mod.example.weapons.Weapon
    local 3 = List<local 1>
    local 4 = string
    local 5 = int32
    local 6 = uint32
    local 7 = uint16
    local 8 = int16

ROOT
    local 0

PLAYER
    Name       -> string payload
    Health     -> int32 payload
    Inventory  -> list payload
                    count = N
                    Item payload
                    Item payload
                    ...

SOME DYNAMIC COMPONENT
    type = local 2
    value = Weapon payload
```

The global identities remain readable:

```text
game.Player
game.Item
mod.example.weapons.Weapon
```

while actual payloads can use compact local references.

---

# 35. Open design questions

These should remain open until prototype benchmarks and tests exist.

### Type ID namespace rules

Should Chef define a recommended canonical syntax such as:

```text
<namespace>.<domain>.<type>
```

or leave naming entirely to applications?

### Type ID normalization

Should IDs be case-sensitive?

Should Unicode be permitted?

Should whitespace be forbidden?

A conservative initial rule would be:

```text
UTF-8
case-sensitive
no leading/trailing whitespace
canonical ASCII-recommended namespaces
```

### Type ID length

String IDs improve modding, but extremely long IDs create unnecessary package overhead.

The type table naturally limits this cost because the string is normally stored only once per package.

A maximum encoded Type ID length should still be enforced.

### Alias support

A future registry could support:

```text
old.mod.Weapon -> mod.example.weapons.Weapon
```

for migration or mod renaming.

Aliases must never silently create two independent identities for the same type.

### Unknown-value preservation

Should Chef v1 expose opaque values so missing mods can be preserved across save cycles?

This is highly useful for mod-heavy environments and deserves further consideration.

### Cross-language compatibility

String IDs make the format more attractive for cross-language tooling, but the initial target remains Go.

Cross-language support would need explicit rules for:

- integer representation;
- floating-point semantics;
- Unicode/string encoding;
- Type ID canonicalization;
- schema representation.

---

# 36. Design invariants

The following should become hard requirements of Chef:

1. A serialized custom type has a stable **string** identity.
2. Type IDs are independent of registration order.
3. Type IDs are independent of Go implementation names.
4. Registration order never changes the meaning of a save.
5. An unknown custom type is never guessed.
6. Known-type fields do not carry redundant type metadata.
7. Dynamic values contain enough type information to select a recipe.
8. Custom recipes can recursively contain any Chef-supported type.
9. Chef format version and recipe version are separate.
10. Malformed input returns errors rather than panicking.
11. Decoding has resource limits.
12. Go memory layout is never itself the persistent wire contract.
13. Package-local numeric type references are merely compression and have no global meaning.
14. Mod-defined types can be introduced without coordination with the base game.
15. Type ID collisions are detected rather than silently resolved.

---

# 37. Design direction

Chef should be implemented as a **typed, schema-assisted binary serialization system with string-based global type identities**.

The core model is:

```text
Canonical String Type ID
          |
          v
        Recipe
          |
          v
Package-Local Type Table
          |
          v
Compact Value Payloads
```

The crucial separation is:

```text
STRING ID
= persistent identity
```

while:

```text
LOCAL INTEGER
= compact reference within one package
```

This gives Chef the properties needed for a mod-friendly save system:

- mod authors can create types independently;
- save files can identify the exact logical type involved;
- registration order does not matter;
- type names remain understandable during debugging;
- repeated type identifiers do not bloat the payload;
- custom types can recursively contain arbitrary Chef-supported types;
- future migration and unknown-type preservation remain possible.

The cookbook therefore becomes more than a runtime registry. It is the authoritative collection of recipes for the ecosystem of types that a game and its mods can understand.
