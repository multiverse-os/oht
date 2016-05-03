# oht v0.1.0
oht is an onion distributed hash table used to create decentralized application framework. **This software is experimental, use with caution.** *The code is in flux at this prealpha stage. The protocol specifications are still subject to major changes.*

Utilizing a similar design pattern as ricochet or onionshare, oht creates an onion routed distributed hash table by having each peer establish an onion service to communicate with peers. This allows for peers to interact with a public DHT in a secure and anonymous manner. However, instead of relying on the onion address keys to authenticate oht opts to use a separate esdca keypair and uses emphermal onion addresses. This provides better security and makes it easier to migrate an account to other terminals.

Onion routing does not only improve security, oht gains other interesting features by leveraging the onion network to route peers. By using onion addresses over ip addresses, peers can directly without port forwarding or worrying about NAT transversal.

The goal of oht is to function as a decentralized application framework for rapid design of onion routed decentralized applications. oht is a general purpose framework that will hopefully be useful for a wide variety of usecases. oht can be used as decentralized web application replacement, decentralized chat, file sharing or connecting "IoT" arm computers for secure interconnections or decentralized VOIP.


## Progress
oht is under active development, currently only the necessary tor binaries for Linux and OSX are supplied. Peer to peer communication is currently handled through websockets and routed through Tor in a manner similar to ricochet. A basic web interface exists served over a onion service. Basic account system exist using ecdsa keypairs. oht builds out configuration files in a correct structure and in appropriate locations.

### Features
Below is a list of the functional requirements to meet the goal of the project. The features below are at varrying stages of completion. Features will be designed to modular and optional so decentralized applications can be scaffolded quickly using only the features needed.

**Accounts** - Accounts based on encrypted assymetric keys. Elliptic curve similar to telehash or Bitcoin/Ethereum. Some backup system that leverages the DHT would be interesting to experiment with. GNUnet naming system or flexibile name system using a variety of existing solutions.

**Name Resolution** - Name resolution leveraging the oht or an existing name system.

**Contact List** - Contact list, optional transversal through contacts network by creating relationship chains (bestFriend.friend.aquaitance). Flexible contact meta data to make it easier to extend functionality.

Easy identity sharing is important, possibly through the use of expirable human readable mnemoics to add contacts or making it easy to leverage existing name systems.

**Direct/Multiuser Chat** - A global DHT chat, hash defined multiuser chat and direct messaging managed with the account system keys. Presence, offline messages and message encryption/validation. Multiuser chats can be encrypted, limited to a set of keys or public. API for bots/plugins to encourage extending the functionality.

Easy public chat sharing is important, possibly through the use of expirable human readable mnemoics to join chat.

**Decentralized Database** - A DHT sharded across participating peers as encrypted blocks. DHT should feature private shards and shards with expiration timers.

**Local Database** - The DHT for the peers and files are cached locally using encrypted BoltDB databases. Alternatively, a memory only cache may be used.

**File Transfer** - A basic system to do file transfer between peers, 1-to-1 and m-to-n. Files should be broken into blocks, tracked using merkel trees and transfered in a manner similar to torrents. Can implement using existing torrent code or leverage existing storage networks.

**Streaming** - 1-to-1 peer streaming of data/music/video, possibly using WebRTC. *WebRTC* can be intergrated using existing websockets p2p connections, the onion address bypasses the need for NAT transversal, allowing serverless p2p webRTC connections to be established without a stun/turn server. WebRTC must be modified to only offer onion service ice candidates, early research has begun on this topic.

**User Interface** - Provide several ways to implement a user interface for the distributed application. Web interface available through an individual onion address. Terminal command line interface and console for interacting with the DHT. A text only browser may be an effective way to rapid prototype fairly complex terminal based UIs. Basic GUI client using QT/wxWidgets or possibly a standalone browser executable for desktop clients.

**Localization** - Localization is important and needs to be designed to exist within the framework from the
beginning.

### Usage

    oht> /help
    COMMANDS:
    
        /config               - List configuration values (Not Implemented)
    
      DHT NETWORK:
    
        /peers                - List all connected peers (Not Implemented)
        /successor            - Next peer in identifier ring (Not Implemented)
        /predecessor          - Previous peer in identifier ring (Not Implemented)
        /ftable               - List ftable peers (Not Implemented)
        /connect [ip address] - Direct connect to peer (Not Implemented)
    
      ACCOUNTS:
    
        /contacts             - List all saved contacts (Not Implemented)
        /add     [user]       - Add account to contacts (Not Implemented)
        /rm      [user]       - Remove account from contacts (Not Implemented)
        /whisper [user]       - Direct message peer (Not Implemented)
    
        /quit

### Roadmap
Currently the basic DHT functionality is still not yet implemented, this core functionality is the first major milestone. Before this can be completed the p2p communication needs to be lower level then the current websockets connections, utilize protobuf, authenticate peers and use ephermeral keys for communication. The existing websockets connections will remain and it is useful for the existing webui.

Building decentralized networks still has a few major roadblocks when trying to create truly trustless decentralized applications:

* **Service Discovery** - Currently here is no simple way for anonymous peers to connect without leveraging a a trusted seed for peers. One possible solution is leveraging a separate larger established DHT to maintain a list of potential peers.

The most interesting idea so far is using a known hard coded onion address key pair to temporarily listen for peers. Requests from active peers will routinely be made to this known onion address. If a reliable way to validate real peers is workable this could be a excellent way to solve the problem.

* **Consensus** - While not as critical as service discovery, it would be preferable to to have some optional and basic methods of consensus for portions of the DHT. Optionally we can make it easy to leverage existing larger file storage networks or cryptocurrencies.

*Notes* for name resolution Bitcoin OP_RETURN based name resolution (https://github.com/telehash/blockname) appears to be a possibility. This could also be used for seeding peers. Research Chord4S as it promises to provide decentralized service discovery but it may not work in an anonymous environment.

## Building Decentralized Applications
The best way to use the codebase in its current prealpha stage is to fork the repository and experiment.

The goal is for oht to create a framework analagous to Rails. How rails provides an intuitive framework for creating web applications rapidly, oht is planned to be a framework for creating secure decentalized applications rapidly.

the first release candidate would include tools to build a boilerplate decentralized application.

## Contribute

Everyone is encouraged to test out the software and experiment with it. Everyone is welcome to contribute to the project, report bugs and submit pull requests.

Developer communication platforms (mailing list, irc) can be established if a community grows. Until they exist anyone is welcome to create github issues to request support.
