# Example: Stored

```
import (
	"fmt"

	// ...
)

	// ...
	
	// using initialized "aDht". see initialization
	// An uploader "example_joe" uploads their file "cat.jpg" to the
	// Archon Cloud. This file is sharded into <example-number> of shards
	// and "example_joe" makes this upload using "ETH" permission layer.
	// Below is example of logic for Storage Provider who receives a
	// shard from "example_joe", where they prepare and make the "Stored"
	// function call.
	
	shardPath := ShardPathFromShard(shard) // pseudocode
	// can derive shard path from shard container
	// shardPath has format 
	// permissionLayerID/username/path/shardIdx.shardFileSuffix 
	
	versionData := VersionDataFromShard(shard) // pseudocode
	// shard container has VersionData
	
	err = aDht.Stored(shardPath, versionData)
	if err != nil {
		// handle
	}

```
