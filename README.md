# oht v0.1.0
oht is an onion distributed hash table used to create decentralized application framework. **This software is experimental, use with caution.** *The code is in flux at this prealpha stage. The protocol specifications are still subject to major changes.*

Utilizing a similar design pattern as ricochet or onionshare, oht creates an onion routed distributed hash table by having each peer establish an onion service to communicate with peers. This allows for peers to interact with a public DHT in a secure and anonymous manner. However, instead of relying on the onion address keys to authenticate oht opts to use a separate esdca keypair and uses emphermal onion addresses. This provides better security and makes it easier to migrate an account to other terminals. Despite this change, interlopability required to communicate with users available on ricochet is important part of the oht design.  

Onion routing does not only improve security, oht gains other interesting features by leveraging the onion network to route peers. By using onion addresses over ip addresses, peers can directly without port forwarding or worrying about NAT transversal. It also allows for a new method of decentralized service discovery. 

The goal of oht is to function as a decentralized application framework for rapid design of onion routed decentralized applications. oht is a general purpose framework that will hopefully be useful for a wide variety of usecases. oht can be used as decentralized web application replacement, decentralized chat, file sharing or connecting "IoT" arm computers for secure interconnections or decentralized VOIP.

## Progress
oht is under active development, currently only the necessary tor binaries for Linux and OSX are supplied. The code is being written to support Linux, OSX and Windows. Peer to peer communication is currently handled through websockets and routed through Tor in a manner similar to ricochet. A basic web interface exists served over a onion service. Basic account system exist using ecdsa keypairs. oht builds out configuration files in a correct structure and in appropriate locations.

### Roadmap
Currently the basic DHT functionality is still not yet implemented, this core functionality is the first major milestone. Before this can be completed the p2p communication needs to be lower level then the current websockets connections, utilize protobuf, authenticate peers and use ephermeral keys for communication. The existing websockets connections will remain and it is useful for the existing webui.

### Core Features
Below is a list of the functional requirements to meet the goal of the project. The features below are at varrying stages of completion. Features will be designed to modular and optional so decentralized applications can be scaffolded quickly using only the features needed.

**Onion routed** - All connections are done using Tor's onion services, ensuring each connection is encrypted, anonymous and eliminates the need for port forwarding. 

**Decentralized Service Discovery** - Using a known onion address keypair, randomly peers can broadcast known peers to the known kepair. Peers looking to join the network can listen using this known keypair and obtain peers to connect to. This also compartmentalizes a DHT and provides an ID.

**Name Resolution** - Name resolution leveraging the oht or an existing name system (GnuNet, Bitcoin/Ethereum)

**DHT** - A DHT sharded across participating peers as encrypted blocks. DHT should feature private shards and shards with expiration timers.

**Local Database** - The DHT for the peers and files are cached locally using encrypted BoltDB databases. Alternatively, a memory only cache may be used.

**File Transfer & Streaming** - A basic system to do file transfer between peers, 1-to-1 and m-to-n. Files should be broken into blocks, tracked using merkel trees and transfered in a manner similar to torrents. Can implement using existing torrent code or leverage existing storage networks.

1-to-1 peer streaming of data/music/video, possibly using WebRTC. *WebRTC* can be intergrated using existing websockets p2p connections, the onion address bypasses the need for NAT transversal, allowing serverless p2p webRTC connections to be established without a stun/turn server. WebRTC must be modified to only offer onion service ice candidates, early research has begun on this topic.

**User Interfaces** - Provide several ways to implement a user interface for the distributed application. Web interface available through an individual onion address. Terminal command line interface and console for interacting with the DHT. A text only browser may be an effective way to rapid prototype fairly complex terminal based UIs. Basic GUI client using QT/wxWidgets or possibly a standalone browser executable for desktop clients.

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

## APIs
oth comes with several APIs found in 
[the `api` directory](https://github.com/onionhash/oht/tree/master/api):

 Planned APIs  |         |
----------|---------|
rest | JSON Rest API |
websockets | JSON websocket API |
ipc | Interprocess communication  |

## Executables

oth comes with several wrappers/executables found in 
[the `ui` directory](https://github.com/onionhash/oht/tree/master/ui):

 Command  |         |
----------|---------|
`oth-cli` | OTH CLI Interface (ethereum command line interface client) |
`oth-console` | OTH Console Interface |

## Building Decentralized Applications
The best way to use the codebase in its current prealpha stage is to fork the repository and experiment.

The goal is for oht to create a framework analagous to Rails. How rails provides an intuitive framework for creating web applications rapidly, oht is planned to be a framework for creating secure decentalized applications rapidly.

The first release candidate will include tools to build a boilerplate decentralized application, including the libraries and a structure illustrate standard design patterns.

### Usage
A basic console is the first goal for the UI and will be used as a way of outlining the basic functional requirements of the design. Much of the functionality works but just has not yet been tied to the console command.

    oht> /help
    COMMANDS:
      CONFIG:
        /config                      - List configuration values (Not Implemented)
        /set [config] [option]       - Change configuration options (Not Implemented)
        /unset [config]              - Unset configuration option (Not Implemented)
        /webui [start|stop]          - Control webui (Not Implemented)
    
      NETWORK:
        /peers                       - List all connected peers (Not Implemented)
        /successor                   - Next peer in identifier ring (Not Implemented)
        /predecessor                 - Previous peer in identifier ring (Not Implemented)
        /ftable                      - List ftable peers (Not Implemented)
        /create                      - Create new ring (Not Implemented)
        /connect [onion address|id]  - Join to ring with peer (Not Implemented)
        /lookup [onion address|id]   - Find onion address of account with id (Not Implemented)
        /ping [onion address]        - Ping peer (Not Implemented)
        /ringcast [message]          - Message every peer in ring (Not Implemented)
    
      DHT:
        /put [key] [value]           - Put key and value into database (Not Implemented)
        /get [key]                   - Get value of key (Not Implemented)
        /delete [key]                - Delete value of key (Not Implemented)
    
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

## Contribute

Everyone is encouraged to test out the software and experiment with it. Everyone is welcome to contribute to the project, report bugs and submit pull requests.

Developer communication platforms (mailing list, irc) can be established if a community grows. Until they exist anyone is welcome to create github issues to request support.
