# Random idea on templates (rather than generics) in Go 2.0
Why not to have specialized template packages, where all declared public items (structs, functions, variables) will have to be explicitly qualified with their types before usage?

### Example
##### Declaration
```go
// template package header, KeyType and ValueType can be used as type names below
package dicts[KeyType, ValueType]

// ...
```
##### Usage: 
```go
// ...

import (
    "..../dicts"
)

// ...
    a := dicts.NewDenseHash[string, int]()
// ...
```

And there should be no way to avoid type qualification. In other words, something like
```go
c := misc.Min(a, b)
```
will not work and
```go
c := misc.Min[float32](a, b)
```
will.

## Benefits of this approach:
There will be less intention to overuse templates due to more verbose templated code.

# PS
I need algebraic types with their decomposition via pattern matching much more than generics.
