# oht v0.1.0
An onion distributed hash table is a DHT that is routed through Tor's onion network using onion services. oht is an implementation of an onion distributed hash table that is designed to be used as a framework for secure onion routed decentralized applications.  oht sets out to be a general purpose framework, designed for a broad set of use cases. oht can be used as the foundation for decentralized web applications, chat, file sharing, VOIP or for securely networking ARM computers (IoT). 

**oht is not affiliated with or endorsed by The Tor Project. This software is experimental, use with caution. The code is in flux at this pre-alpha stage.** 
*The protocol specifications are still subject to major changes and the documentation is still mostly a patchwork of notes.*

## Development Progress
oht is under active development, and is currently only packaged with the Tor binaries necessary for Linux and OSX. The code is being written to support Linux, OSX and Windows. 

P2P communication is currently onion routed using Tor onion services similar to Ricochet. Currently peer communications are handled through WebSockets. A basic optional account system exists using ECDSA keypairs. oht creates configuration files in a correct structure and in standard locations (relative to operating system). A basic console UI is currently the primary client. Additionally, a basic client web interface exists as an onion service accessible through the Tor Browser.

The basic DHT functionality is still not yet implemented. The first step will be implementing a networking library that works with a variety of transports to support a wide number of existing protocols, such as Kademlia, Ricochet, WebRTC, and SIP. 

Unlike ricochet, instead of using the onion service keypair, oht uses a separate key pair for node Authentication. Currently ECDSA keys compatible with Ethereum are used for node authentication and encryption of messages.

Individual components that would be useful by themselves, for example the method of Tor control, will be broken out into libraries to easily implement in any program.

To demonstrate the framework, an example application will be developed alongside the framework. The application will be a encrypt decentralized geolocal and global chat that includes a set of tools to help organize groups of people.

## Executables

oht comes with three wrappers/executables found in
[the `ui` directory](https://github.com/onionhash/oht/tree/master/ui):

 Command  |         |
----------|---------|
`ohtd` | OHT Daemon |
`oht-cli` | OHT CLI Interface (command line interface client) |
`oht-console` | OHT Console Interface |

## APIs
oht comes with three APIs found in
[the `api` directory](https://github.com/onionhash/oht/tree/master/api):

 Proposed APIs  |         |
----------|---------|
rest | JSON Rest API |
websockets | JSON WebSocket API |
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
        /tor [start|stop]            - Start or stop Tor process
        /newtor                      - Obtain new Tor identity (Not Implemented)
        /newonions                   - Obtain new onion address
    
      NETWORK:
        /peers                       - List all connected peers (Not Implemented)
        /successor                   - Next peer in identifier ring (Not Implemented)
        /predecessor                 - Previous peer in identifier ring (Not Implemented)
        /ftable                      - List ftable peers (Not Implemented)
        /create                      - Create new ring (Not Implemented)
        /connect [onion address]     - Join ring containing peer with [onion address]
        /lookup [id]                 - Find onion address of account with [id] (Not Implemented)
        /ping [onion address]        - Ping peer (Not Implemented)
        /ringcast [message]          - Message every peer in ring (Not Implemented)
    
      DHT:
        /put [key] [value]           - Put key and value into database (Not Implemented)
        /get [key]                   - Get value of key (Not Implemented)
        /delete [key]                - Delete key and its value from database (Not Implemented)
    
      WEBUI:
        /webui [start|stop]          - Start or stop webUI server
    
      ACCOUNT:
        /accounts                    - List all local accounts (Not Implemented)
        /generate                    - Generate new account key pair (Not Implemented)
        /delete                      - Delete an account key pair (Not Implemented)
        /sign [id] [message]         - Sign with account key pair (Not Implemented)
        /verify [id] [message]       - Verify a signed message with key pair (Not Implemented)
        /encrypt [id] [message]      - Encrypt a message with key pair (Not Implemented)
        /decrypt [id] [message]      - Decrypt a message with key pair (Not Implemented)
    
      CONTACTS:
        /contacts                    - List all saved contacts (Not Implemented)
        /request [id] [message]      - Send [message] requesting account with [id] to add your id to their contacts (Not Implemented)
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
    
        /quit                        - Quit oht console

## oht Explanation and Comparison to Typical DHT

Tor is often misunderstood. People often confuse the Tor Browser Bundle (TBB) with Tor itself. Tor is a client for a decentralized p2p onion routing network, which can be simply described as adding additional proxy layers between you and your destination when accessing the internet. This "onion" of proxy layers also obscures the location of the user. TBB is a browser (based on Mozilla Firefox) bundled with the Tor binary for easy and secure access to websites through the Tor network. Tor works with any port, and is not restricted to the common http/https ports (80 and 443). The additional proxy layers provide a connection with additional security and can bypass the regional restrictions being imposed on the world wide web. An example use case for onion routing is a journalist using TBB to bypass national firewalls to report news. One aim of this project is to highlight that Tor provides more than just a solution for secure Internet browsing, Tor provides a solution for secure hosting through onion services.  Secure p2p connections can be made by having peers offer **onion services**.

**Onion services** create end-to-end encrypted onion routed connections with perfect forward secrecy. Onion services do not use Tor Exit Nodes, but instead rendezvous points outside of both peers' networks. This solves the issue of NAT transversal when connecting peers. A typical DHT when used in combination with Tor can potentially be used in correlation attacks, but when routing DHT traffic through onion services this problem is avoided. 

oht utilizes a similar onion routing design pattern to Ricochet or OnionShare, where each peer in the network establishes an onion service to communicate.  

This allows for peers to interact with a public DHT securely by limiting the amount of metadata shared. Peers do not receive the IP address and geographic location of connected peers. Instead, peers rely on ephemeral onion addresses to connect to each other and authenticate with a separate key pair that can be compatible with Bitcoin/Ethereum or telehash. 
Interoperability with Ricochet is important, this will be acheived by saving an onion address key pair in a custom meta-data field on an account or general configuration. Core functionality will include a networking library to provide compatibility with the v2 Ricochet protocol.

**Beyond providing additional security** onion routing has interesting emergent properties when combined with with the standard DHT. Because onion services use rendevous points, it allows for peers to avoid any issues with NAT transversal, which is often problematic with p2p networks. Additionaly, an onion address key pair shared across all peers (shared onion address) by either hardcoding into a client or stored in a configuration file can be used to solve problems with centralization that DHTs face with bootstrapping.

A **shared onion address** can simultaneouesly be used by mulitple peers to listen over a port. Incoming connections to the shared address are routed to one of the listeners randomly. This can be used to serve information about the DHT network to a decentralized but static location.

A shared onion address used in this manner can act both as the identifier for the p2p network and possibly more importantly as a catalyst to form the p2p network in a decentralized way. This feature allows for smaller DHT's to remain cohesive without standing up dedicated peer seeding servers.

Typically, a p2p network is reliant on predefined trusted boot strap nodes used to obtain an initial list of peers but by using a shared onion address it may be possible to solve this problem in a secure and decentralized way. 

**Decentralized bootstrapping and service discovery** can be accomplished by using a shared onion address to identify the network (oht:address). Potential peers interested in joining the p2p network can use the shared onion address either hard coded into their client (or configuration, name resolution, mnemonic) and listen for peers. Meanwhile users who are already connected to the p2p network can randomly select a known peer and submit it to the shared onion address.  If one or more is listening they will be randomly selected and receive a peer. After a potential peer collects a large enough sample of active peer addresses the potential peer connects to the network.       

Shared onion addresses also may allow for **possible decentralized DHT API and decentralized web UI**.
Active peers may optionally serve standard API or web UI defined by the client or configuration. For example, a simple REST API could be used to serve the DHT by providing basic GET, PUT, DEL commands. To verify the API results, a user could make several requests to obtain a larger sample of responses.

It may be possible to use a basic decentralized web UI. Using this a user could interact with the decentralized application without needing to run the full client, just TBB and the onion address required to connect.

There are a lot of interesting possibilities that can be explored once the software forms.

### Core Features
The core features of oht are at various stages of completion.

**Onion routed** - All connections are done using Tor's onion services, ensuring each connection is encrypted, anonymous and eliminates the need for port forwarding. 

**Decentralized Bootstrap** - Using a shared onion address, potential peers interested in joining the network can obtain peers. This also compartmentalizes a DHT and provides an ID.

**Name Resolution** - Name resolution leveraging the oht or an existing name system (GnuNet, Bitcoin/Ethereum)

**DHT** - A DHT sharded across participating peers as encrypted blocks. DHT should feature private shards and shards with expiration timers.

**Local Database** - The DBs needed to track peers, dht key/values are encrypted BoltDB databases. Alternatively, a memory only cache may be used.

**File Transfer & Streaming** - A basic system to do file transfer between peers, 1-to-1 and m-to-n. Files should be broken into blocks, tracked using merkel trees and transfered in a manner similar to torrents. Can implement using existing torrent code or leverage existing storage networks such as IPFS.

1-to-1 peer streaming of data/music/video, possibly using WebRTC. *WebRTC* can be intergrated using existing WebSocket p2p connections, the onion address bypasses the need for NAT transversal, allowing serverless p2p webRTC connections to be established without a stun/turn server. WebRTC must be modified to only offer onion service ice candidates, early research has begun on this topic.

**User Interfaces** - Provide several ways to implement a user interface for the distributed application. Web interface available through an individual onion address. Terminal command line interface and console for interacting with the DHT. A text only browser may be an effective way to rapid prototype fairly complex terminal based UIs. Basic GUI client using QT/wxWidgets or possibly a standalone browser executable for desktop clients.

*Possible Decentralized APIs and WebUI* through the use of shared public keys. Using a simple API to serve the checksums of the web UI verification may be faster. (Why are standard js libraries not already been actively checked against a published checksum on every site?)

**Localization** - Localization is important and needs to be designed to exist within the framework from the
beginning.

### Optional Modules
oht is at its core a distributed hash table built to route through Tor. oht is designed to be used by others to create more complex software. In order to satisfy and standardize the basic features of most applications oht is packaged with optional modules to extend the base features.

**Accounts** - An account system based on ECDSA keys, supports handling multiple accounts and encrypted storage of keys.

Accounts based on encrypted assymetric keys similar to telehash or Bitcoin/Ethereum. Some backup system that leverages the DHT would be interesting to experiment with.

**Contacts** - A contact system with approval, presence, relationship transversal (e.g. bestFriend.friend.aquaitance). Capable of storing Ricochet contacts. Flexible contact meta data to make it easier to extend functionality.

Easy identity sharing is important, possibly through the use of expirable human readable mnemoics to add contacts or making it easy to leverage existing name systems.

**Channels** - A broadly defined channel system, built to be compatible with Ricochet's protocol, to be used for multiuser chat rooms and any other abstraction that requires variable isolation.

Each one will include an interface.go file that matches the general structure of the oht/interface.go file. APIs and UIs will interact with these interfaces.

## Building Decentralized Applications
The best way to use the codebase in its current prealpha stage is to fork the repository and experiment.

The goal is for oht to create a framework analagous to Rails. How rails provides an intuitive framework for creating web applications rapidly, oht is planned to be a framework for creating secure decentalized applications rapidly.

The first release candidate will include tools to build a boilerplate decentralized application.

## Contribute

Everyone is encouraged to test out the software and experiment with it. Everyone is welcome to contribute to the project, report bugs and submit pull requests.

Developer communication platforms (mailing list, irc) can be established if a community grows. Until they exist everyone is welcome to create github issues to request support.
