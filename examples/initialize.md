# Example: Initialize

```
import (
  	dht "github.com/archoncloud/archon-dht/archon"
  	dhtLayers "github.com/archoncloud/archon-dht/dht_permission_layers"
  
  	"github.com/pariz/gountries"

	// ...
)

	// ...
	
	seed := int64(123123) // example
	
	JP := "JP" // example Japan
	countryCode := gountries.Codes{Alpha2: JP}
  	basePort := int(9000) // example
	// note: be sure to setup router (NAT Port Forwarding) to allow 
	// ports 9000 to 9005 to be accessible from outside your LAN
	myUploadUrl := "http://myExampleSPUploadUrl.com/uploadEndpoint" 
	// example

  	archonEthAbi.SetRpcUrl([]string{"https://goerli.infura.io/v3/9ca2c17c532a09ca2c17c532a0c532a0"}) // example fake rpc url
	
	// note: the nodeID of this SP must be registered with the 
	// Archon Ethereum Smart Contract as an SP in order to connect
	// to this subnet
	eth := new(dhtLayers.Ethereum);
	ethConfigDHT := dht.DHTConnectionConfig{
	    Seed: seed,
	    Global: true,
	    IAmBootstrap: false,
	    OptInToNetworkLogging: true,
	    CountryCode: countryCode,
	    PermissionLayer: *eth,
	    Url: myUploadUrl + "/eth",
	    MyPartialMultiAddress: "/ip4/0.0.0.0/tcp/" + strconv.Itoa(basePort + 2),
	    BootstrapPeers: []string{
		"/ip4/18.220.115.81/tcp/9002/ipfs/QmNX6ASyukLch38D2Z1h4cMh39ATfqqDom1xJWv2YHc1eG"}}
    
	// note: the nodeID of this SP must be registered with the 
	// Archon Neo Smart Contract as an SP in order to connect
	// to this subnet
	neo := new(dhtLayers.Neo)
	neoConfigDHT := dht.DHTConnectionConfig{
	    Seed: seed,
	    Global: true,
	    IAmBootstrap: false,
	    OptInToNetworkLogging: true,
	    CountryCode: countryCode,
	    PermissionLayer: *neo,
	    Url: myUploadUrl + "/neo",
	    MyPartialMultiAddress: "/ip4/0.0.0.0/tcp/" + strconv.Itoa(basePort + 3),
	    BootstrapPeers: []string{
		"/ip4/18.220.115.81/tcp/9003/ipfs/QmNX6ASyukLch38D2Z1h4cMh39ATfqqDom1xJWv2YHc1eG"}}
  
	  nonPermissioned := new(dhtLayers.NonPermissioned)
	  freeConfigDHT := dht.DHTConnectionConfig{
	    Seed: seed,
	    Global: true,
	    IAmBootstrap: false,
	    OptInToNetworkLogging: true,
	    CountryCode: countryCode,
	    PermissionLayer: *nonPermissioned,
	    Url: myUploadUrl + "/non",
	    MyPartialMultiAddress: "/ip4/0.0.0.0/tcp/" + strconv.Itoa(basePort + 1),
	    BootstrapPeers: []string{
		"/ip4/18.220.115.81/tcp/9001/ipfs/QmNX6ASyukLch38D2Z1h4cMh39ATfqqDom1xJWv2YHc1eG"}}

	  var configArray []dht.DHTConnectionConfig
	  configArray = append(configArray, ethConfigDHT)
	  configArray = append(configArray, neoConfigDHT)
	  configArray = append(configArray, freeConfigDHT) 
	  // NonPermissioned config is necessary otherwise init throws

	  aDht, err := dht.Init(configArray, basePort)
	  if err != nil {
	  	// handle
	  }
```
