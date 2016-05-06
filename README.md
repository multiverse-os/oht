# oht v0.1.0
An onion distributed hash table, is a DHT that is routed through the onion network using Tor's onion services. oht is an implementation of a onion distributed hash table that is designed to be used as a framework for secure onion routed decentralized applications.  oht sets out to be a general purpose framework, design for broad set of use cases. oht can be used as the foundation decentralized web application replacement, decentralized chat, file sharing, VOIP or securely networking arm computers (IoT). 

**This software is experimental, use with caution. The code is in flux at this pre-alpha stage.** 
*The protocol specifications are still subject to major changes.*

## Develpment Progress
oht is under active development, currently only the necessary tor binaries for Linux and OSX are supplied. The code is being written to support Linux, OSX and Windows. 

P2p communication is currently onion routed using Tor onion services similar to ricochet. Currently peer communications are handled through websockets.  Basic optional account system exist using ecdsa keypairs. oht builds out configuration files in a correct structure and in appropriate locations. A basic console UI is currently the primary client. Additionaly a basic client web interface onion service exists and can be accessed through TBB.

The basic DHT functionality is still not yet implemented. Before the DHT can be started the peer communication needs to be moved to Nanomsg "scalable protocols" using protobuf. Ecdsa keys for authentication and encryption of messages. The existing websockets allow for direct connections an oht network with javascript which has been useful during development. 

## Executables

oth comes with three wrappers/executables found in
[the `ui` directory](https://github.com/onionhash/oht/tree/master/ui):

 Command  |         |
----------|---------|
`othd` | OTH Daemon Client |
`oth-cli` | OTH CLI Interface (ethereum command line interface client) |
`oth-console` | OTH Console Interface |

## APIs
oth comes with three APIs found in
[the `api` directory](https://github.com/onionhash/oht/tree/master/api):

 Proposed APIs  |         |
----------|---------|
rest | JSON Rest API |
websockets | JSON websocket API |
ipc | Interprocess communication  |

### Usage
The primary client during this stage of development is the console client.

#### Console Client

    oht> /help
    COMMANDS:
      CONFIG:
        /config                      - List configuration values
        /set [config] [option]       - Change configuration options
        /unset [config]              - Unset configuration option
        /save                        - Save configuration values
    
      TOR:
        /tor [start|stop]            - Start or stop tor process (Not Implemented)
        /newtor                      - Obtain new Tor identity (Not Implemented)
        /newonion                    - Obtain new onion address (Not Implemented)
    
      NETWORK:
        /peers                       - List all connected peers (Not Implemented)
        /successor                   - Next peer in identifier ring (Not Implemented)
        /predecessor                 - Previous peer in identifier ring (Not Implemented)
        /ftable                      - List ftable peers (Not Implemented)
        /create                      - Create new ring (Not Implemented)
        /connect [onion address]     - Join to ring with peer
        /lookup [id]                 - Find onion address of account with id (Not Implemented)
        /ping [onion address]        - Ping peer (Not Implemented)
        /ringcast [message]          - Message every peer in ring (Not Implemented)
    
      DHT:
        /put [key] [value]           - Put key and value into database (Not Implemented)
        /get [key]                   - Get value of key (Not Implemented)
        /delete [key]                - Delete value of key (Not Implemented)
    
      WEBUI:
        /webui [start|stop]          - Start or stop webUI server
    
      ACCOUNT:
        /accounts                    - List all accounts (Not Implemented)
        /generate                    - Generate new account key pair (Not Implemented)
        /delete                      - Delete an account key pair (Not Implemented)
        /sign [id] [message]         - Sign with account key pair (Not Implemented)
        /verify [id] [message]       - Verify a signed message with keypair (Not Implemented)
        /encrypt [id] [message]      - Encrypt a message with keypair (Not Implemented)
        /decrypt [id] [message]      - Decrypt a message with keypair (Not Implemented)
    
      CONTACTS:
        /contacts                    - List all saved contacts (Not Implemented)
        /request [id] [message]      - Request account to add your id to their contacts (Not Implemented)
        /add [id]                    - Add account to contacts (Not Implemented)
        /rm [id]                     - Remove account from contacts (Not Implemented)
        /whisper [id] [message]      - Direct message peer (Not Implemented)
        /contactcast [message]       - Message all contacts (Not Implemented)
    
      CHANNELS:
        /channels                    - List all known channels (Not Implemented)
        /channel                     - Generates a new channel (Not Implemented)
        /join [id]                   - Join channel with id (Not Implemented)
        /leave [id]                  - Leave channel with id (Not Implemented)
        /channelcast [id] [message]  - Message all channel subscribers (Not Implemented)
    
        /quit

## oht Explanation and Comparison to Typical DHT

Tor is misunderstood, for the purpose of this documentation, we will focus on how often people will confuse Tor Browser Bundle (TBB) with Tor. Tor is a client for a decentralized p2p onion routing network, which translates to adding additional proxy layers between you and your destination when accessing the internet. Tor works with any port, and is not restricted to port 80. Additional proxy layers provide a connection with additional security and bypass the regional restrictions being imposed on the world wide web. One example use case for onion routing is a journalist using TBB to bypass national firewalls to report accurate news. The aim is highlight that Tor provides more than just a solution for secure Internet browsing, Tor provides a solution for secure hosting through onion services. 

**Onion services** create end-to-end encrypted (/w perfect forward secrecy) onion routed connections. Onion services do not use Tor Exit Nodes, but instead rendevouz points outside of both peers networks which solves the issue of NAT transversal when connecting peers. A typical DHT when used in combination with Tor can potentially be used in correlation attacks, but when routing DHT traffic through onion services this problem is avoided. 

oht utilizes a similar onion routing design pattern to ricochet or onionshare. Onion routed in this case is being defined as each peer in the network establish an onion service to communicate communication.  

This allows for peers to interact with a public DHT securely by limiting the amount of metadata. Peers do not receive the IP address and geographic location of connected peers. Instead peers rely on emphemeral onion address key pair shared with peers and authenticate with a separate key pair that can be compatible with Bitcoin/Ethereum or telehash. 
Interoperability with ricochet is important, this will be acheived by saving an onion address key pair in a custom meta-data field on an account or general configuration. Core functionality will include library in networking to provide compadibility with the v2 ricochet protocol.

**Beyond providing additional security** onion routing has interesting emergent properties when combined with with the standard DHT. Onion services allow peers to avoid any issues with NAT transversal which is often problematic with p2p networks. Additionaly, an onion address key pair shared across all peers (shared onion address) by either hardcoding into the client or using a configuration can be used to solve problems with centralization that DHTs face with bootstrapping.

A **shared onion address** can simultaneouesly be used by mulitple peers to listens over a port. Incoming connections to the shared address are randomly sent to any of the listeners. A set of information defined by the client can be sent to the shared onion address. 

A shared onion address used in this manner can act both the identifier for the p2p network and more importantly the catalyst to form the p2p network in a decentralized way. 

Typically, a p2p network is reliant on predefined trusted boot strap nodes used to obtain an initial list of peers. 

**Decentralized bootstrapping and service discovery** can be accomplished by using a shared onion address to identify the network. Potential peers can use the shared onion address either hard coded into their client or in a configuration file or a converted from a mnemoic phrase. Meanwhile peers the client of peers already connected to the p2p network optionally send a random known peer address. After a potential peer collects enough active peer addresses the potential peer connects to the network if verification succeeds.       

*Verfication may be possible by asking another random peer to give you a peer to send to the potential peer. The answer is signed by both peers and possibly include a checksum hash of the peer table. The potential peer ask for verification from both peers who generated the answer. This may allow banning based on ecdsa keys used to sign a message with incorrect data.*

Shared onion addresses also may allow for **possible decentralized DHT API and decentralized web UI**.
Active peers may optionally serve standard API defined by the protocol or configuration. For example, a simple REST API could be used to serve the DHT. Several requests can be made, checked for correctness and verified.

It may be possible to use a decentralized API to provide checksums for a Firefox plugin for javascript files served by a decentralized web UI. Using this a user could interact with the decentralized application without needing to run the full client, just TBB.

### Core Features
The core features of oht are varrying stages of completion.

**Onion routed** - All connections are done using Tor's onion services, ensuring each connection is encrypted, anonymous and eliminates the need for port forwarding. 

**Decentralized Bootstrap** - Using a known onion address keypair, randomly peers can broadcast known peers to the known kepair. Peers looking to join the network can listen using this known keypair and obtain peers to connect to. This also compartmentalizes a DHT and provides an ID.

**Name Resolution** - Name resolution leveraging the oht or an existing name system (GnuNet, Bitcoin/Ethereum)

**DHT** - A DHT sharded across participating peers as encrypted blocks. DHT should feature private shards and shards with expiration timers.

**Local Database** - The DHT for the peers and files are cached locally using encrypted BoltDB databases. Alternatively, a memory only cache may be used.

**File Transfer & Streaming** - A basic system to do file transfer between peers, 1-to-1 and m-to-n. Files should be broken into blocks, tracked using merkel trees and transfered in a manner similar to torrents. Can implement using existing torrent code or leverage existing storage networks.

1-to-1 peer streaming of data/music/video, possibly using WebRTC. *WebRTC* can be intergrated using existing websockets p2p connections, the onion address bypasses the need for NAT transversal, allowing serverless p2p webRTC connections to be established without a stun/turn server. WebRTC must be modified to only offer onion service ice candidates, early research has begun on this topic.

**User Interfaces** - Provide several ways to implement a user interface for the distributed application. Web interface available through an individual onion address. Terminal command line interface and console for interacting with the DHT. A text only browser may be an effective way to rapid prototype fairly complex terminal based UIs. Basic GUI client using QT/wxWidgets or possibly a standalone browser executable for desktop clients.

*Possible Decentralized APIs and WebUI* through the use of shared public keys and verification of checksum of all served files. (Why are standard js libraries not already been actively checked against a published checksum on every site?)

**Localization** - Localization is important and needs to be designed to exist within the framework from the
beginning.

### Optional Modules
oht is at its core a distributed hash table built to route through Tor. oht is designed to be used by others to create more complex software. In order to satisfy and standardize the basic features of most applications oht is packaged with optional modules to extend the base features.

**Accounts** - An account system based on ecdsa keys, supports handling multiple accounts and encrypted storage of keys.

Accounts based on encrypted assymetric keys similar to telehash or Bitcoin/Ethereum. Some backup system that leverages the DHT would be interesting to experiment with.

**Contacts** - A contact system with approval, presence, grandulated relationship transversal (e.g. bestFriend.friend.aquaitance). Capable of storing ricochet contacts. Flexible contact meta data to make it easier to extend functionality.

Easy identity sharing is important, possibly through the use of expirable human readable mnemoics to add contacts or making it easy to leverage existing name systems.

**Channels** - A broadly defined channel system, built to be compatible with ricochet's protocol, to be used for multiuser chat rooms and any other abstraction that requires variable isolation.

Each one will include an interface.go file that matches the general structure of the oht/interface.go file. APIs and UIs will interact with these interfaces.

## Building Decentralized Applications
The best way to use the codebase in its current prealpha stage is to fork the repository and experiment.

The goal is for oht to create a framework analagous to Rails. How rails provides an intuitive framework for creating web applications rapidly, oht is planned to be a framework for creating secure decentalized applications rapidly.

The first release candidate will include tools to build a boilerplate decentralized application.

## Contribute

Everyone is encouraged to test out the software and experiment with it. Everyone is welcome to contribute to the project, report bugs and submit pull requests.

Developer communication platforms (mailing list, irc) can be established if a community grows. Until they exist anyone is welcome to create github issues to request support.
