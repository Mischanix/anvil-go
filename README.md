# anvil-go

[GoDoc](http://godoc.org/github.com/Mischanix/anvil-go)

- Read chunks `Chunk(x,z)`
- Write chunks `SetChunk(x, z, c)`
- Supply only a region path (e.g. `world/region/`, `world/DIM_7/region/`)

`New(path string) *Level`

`*Level`
- `Chunk(x,z) *Chunk`
- `SetChunk(x, z, *Chunk) err`
- `ResetChunk(x, z) err // free the sectors in the corresponding header`

`*Chunk`
- `X, Z int`
- `Data nbt.TagCompound`
