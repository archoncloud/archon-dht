# Example:  GetUrlsOfNodesHoldingKeysFromAllLayers

```
import (
	"fmt"
	"time"

	permLayer "github.com/archoncloud/archon-dht/permission_layer"
	
	// ...
)

	// using initialized "aDht". see initialization

	
	// A download client wants to download "cat.jpg"
	// Recall: An uploader "example_joe" uploads their file "cat.jpg" to the
	// Archon Cloud. This file is sharded into <example-number> of shards
	// and "example_joe" makes this upload using "ETH" permission layer.

	

	permissionLayerID := permLayer.PermissionLayerID("ETH") // example
	username := "example_joe"
	path := "path/to/my/example/file/cat"
	shardFileSuffix := ".jpg"

	var shards []string
	exampleNumber := <example-number>
	
	for i := 0; i < exampleNumber; i++ {
		shardIdx := i 
		shardPath := fmt.Sprintf("%s/%s/%s/%d.%s", permissionLayerID, 
							   username, 
							   path, 
							   shardIdx, 
							   shardFileSuffix) 
		shards = append(shards, shardPath)
	}
	timeout := 8*time.Second	
	urls, err := aDht.GetUrlsOfNodesHoldingKeysFromAllLayers(shards, timeout)
	if err != nil {
		// handle
	}

```
