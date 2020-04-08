# archon-dht

Documentation in progress.

### Contents:

  1. Overview

  2. High-Level protocol description

  3. Initialization

  4. APIs

  5. Further reading 

--------------------------------------------------------------------


### 1. Overview

This module provides the peer-to-peer networking stack for the Archon Cloud (AC). This networking stack consists of an extension of libp2p's implementation of [Kademlia](http://www.scs.stanford.edu/~dm/home/papers/kpos.pdf) to support Blockchain permissioned subnets enabling more efficient and secure decentralized file storage. While currently this module directly supports [Ethereum](https://github.com/archoncloud/archoncloud-ethereum) and [Neo](https://github.com/archoncloud/archoncloud-go/blockchainAPI/neo) permission layers, the extension is generic so that permission layers hosted by other blockchains can be easily integrated.

--------------------------------------------------------------------

### 2. High-Level protocol description

This is a very high-level description. Many details are glossed over in order to keep this brief. For a more detailed description, see the Archon White Paper, the links from further reading section, or read the code from our official repositories (including this one).

For this simple protocol description, we define the players in the Archon Cloud to be storage providers S, uploaders U, and downloaders D. The intent of these players are what you think they would be: the U want to make their versioned content available, the D want to obtain the content of U, and S wants to earn cryptocurrency by acting as a conduit serving the needs of U and D.

To bootstrap this protocol, we start with the S. Any storage provider S must be registered with a permission layer (Archon Smart Contract SC) as storage providers. This registration includes providing information about their storage capabilities, marketplace ask, routing information, as well as staking token (cryptocurrency in the respective blockchain). Once an S is registered with a permission layer, the S is able to establish a permissioned archon-dht connection with other storage providers in each respective permissioned subnet. All S are able to establish a non-permissioned p2p connection with AC. But the permission layer is the mechanism that enforces optimal outcomes for all players in AC. 

An uploader U must also be registered with the SC. This registration includes establishing a namespace, and publishing the public key corresponding to their pseudo-identity. The U do not need to be connected to archon-dht, but we see below how U calls on proxy S to interact with AC.

We follow an upload u from U to its final target, the downloader D.

The U prepares u using some encoding and cryptographically signs u to get {u, sig(u), versionData, {other-metadata}}. Either now, or in the past, U has accumulated a subset of storage providers S_ = {S_1,S_2,...,S_n} (a subset of all storage providers in AC) from one or a few S. This is where some S act as proxy to U to retrieve data from the permission layer and the archon-dht to build and return S_. Locally, U runs the AC marketplace to determine the best S from S_ to accomodate u.

U concurrently makes a proposeUpload transaction pu_tx to SC with metadata of u and S that was determined by the marketplace and sends {u, sig(u), versionData, {other-metadata}}. S caches this upload and listens to SC for pu_tx to be confirmed by the blockchain. The proposeUpload transaction includes a payment to S for storing u, documents metadata of u including sig(u) and versionData, and also validates the result of the marketplace. Assuming pu_tx is confirmed, S announces to the networking overlay (archon-dht) of AC that it is storing {u, ...} and stores {u,...} for the period paid for by U in pu_tx. 

The downloader D knows of u from some other channel. Perhaps U advertised on, say, reddit that U uploaded u. D contacts some storage provider S' asking for the AC download url of u. Storage provider S' queries its networking overlay (archon-dht) for the url(s) of any S holding the latest version of u and returns these values to D. Downloader D downloads {u, sig(u), versionData, {other-metadata}} from S and retrieves the public key of U from the SC. D validates sig(u) and versionData with this public key and accepts u in the ideal case.

We will see below which API's each of the players call in order to participate in this protocol. Please keep in mind, this description glossed over some very important implementation details in order to be brief. For a more detailed protocol description, refer to the Archon Whitepaper, or inspect the source of the official repositories (including this one).


--------------------------------------------------------------------

### 3. Initialization

Set the logging level

```
  common.InitLogging(common.DefaultToExecutable("testingLoggingFolder/logging.log"))
  common.SetLoggingLevelFromName("debug")
```

Initialize ArchonDHT 

```
  aDht, err := dht.Init(configArray, basePort);
  // see example for explicit configArray construction
  if err != nil {
    fmt.Println("debug aDht err ", err)
  }
```
[example](https://github.com/archoncloud/archon-dht/examples/initialize.md)

--------------------------------------------------------------------

### 4. APIs 

--------------------------------------------------------------------

### 5. Further reading 

 - [Kademlia](http://www.scs.stanford.edu/~dm/home/papers/kpos.pdf)

 - [s/Kademlia](https://www.researchgate.net/publication/4319659_SKademlia_A_practicable_approach_towards_secure_key-based_routing)

 - [libp2p](https://github.com/libp2p/go-libp2p)

 - [archoncloud-ethereum](https://github.com/archoncloud/archoncloud-ethereum)

 - [archoncloud-neo](https://github.com/archoncloud/archoncloud-go/blockchainAPI/neo)

 - [archoncloud-contracts](https://github.com/archoncloud/archoncloud-contracts)


--------------------------------------------------------------------

