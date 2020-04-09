# Example: Stored

```
import (
	permLayer "github.com/archoncloud/archon-dht/permission_layer"
	permLayers "github.com/archoncloud/archon-dht/dht_permission_layers"

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

	permissionLayerID := permLayer.PermissionLayerID("ETH") // example
	username := "example_joe"
	path := "path/to/my/example/file/cat"
	shardIdx := int(3) // example (assume <example-number> > 3)
	shardFileSuffix := ".jpg"

	shardPath := fmt.Sprintf("%s/%s/%s/%d.%s", permissionLayerID, 
						   username, 
						   path, 
						   shardIdx, 
						   shardFileSuffix) 

	layer := permLayers.NewPermissionLayer(permissionLayerID)
	if layer == nil {
		// handle
	}
	vd, err := layer.NewVersionData()
	if err != nil {
		// handle
	}
	err = aDht.Stored(shardPath, vd)
	if err != nil {
		// handle
	}
```
