# Example: GetArchonSPProfilesForMarketplace

```
import (
  	permLayer "github.com/archoncloud/archon-dht/permission_layer"
	
	// ...
)
	
	// ...

	// using initialized "aDht". see initialization

	permissionLayerID := permLayer.PermissionLayerID("ETH") // example

	res, err := aDht.GetArchonSPProfilesForMarketplace(permissionLayerID)
	if err != nil {
		// handle
	}
```
