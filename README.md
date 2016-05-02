# oht
oht is an onion distributed hash table used to create secure decentralized application framework. **This software is experimental, use with caution.** *Everything is in flux at this stage and the protocol specifications are still subject to major changes.* 

Utilizing a similar design pattern as ricochet or onionshare, oht creates an onion routed distributed hash table to connect peers to create a DHT in a secure and anonymous way. 

The goal of oht is to function as a application framework for designing onion routed decentralized applications. The design is focused on creating a general purpose framework that can create decentralized web application replacements, multiuser chat or even VOIP.

Leveraging the onion network peers can connect directly without port forwarding or puncturing a NAT while providing additional security, less meta data and increased anonymity. 

## Progress
oht is under active development, the software only provides the necessary tor binaries for Linux and OSX. Peer to peer communication is currently handled through websockets and a basic web interface. 

### Features
There are some basic but optional features, some completed and others to be completed. The goal of the features is to be modular, optional so decentralized applications can be scaffolded quickly.

*Accounts* - Accounts based on assymetric keys. Early design used the onion service key but this will be upgraded soon because of possible security issues. Elliptic curve similar to telehash or Bitcoin seems like a good choice since it could intergrate into existing services easier. Some backup system that leverages the DHT would be interesting to experiment with. GNUnet naming system or flexibile name resolution using a variety of existing solutions. 

*Contact List* - Contact list with optional availability, optional transversal through contact network. Include variable metadata. 

*Direct/Multiuser Chat* - A global chat, custom multiuser chat and direct messaging encrypted with the account asymmetric key. Presence, offline messages and several levels of encryption. 

*User Interface* - Web interface available through an individual onion address. Terminal command line interface and console for interacting with the DHT. Basic GUI client using QT/wxWidgets or possibly a standalone browser executable for desktop clients. 

*Localization* - Localization built into the framework

*Database* - Encrypted database using BoltDB or optional in memory database. 

*WebRTC* - Using the websockets peer to peer connections serverless p2p web RTC connections can be established without a stun/turn server. WebRTC must be modified to only offer onion service ice candidates.

### Roadmap 
Currently the basic DHT functionality is still not yet implemented, this core functionality is the first major milestone. The peer to peer communication should be lower level and use protobuf. The websockets is useful for the webui but it is not ideal for peer to peer communication.

There are a few major roadblocks in creating a trustless decentralized application:

* *Service Discovery* - currently there is no practical way for peers to connect without leveraging a a trusted seed for peers. One possible solution is leveraging a separate larger established DHT to maintain a list of potential peers.
* *Consensus* - It would be preferable to to have some optional method of consensus for DHT settings over hard coding the settings. One possible solution is leveraging a seperate larger established DHT to maintain a cryptographically verified configuration.

## Building Decentralized Applications
The best way to use the codebase in its current stage is to fork the repository. 

The goal is for oht to be analagous to Rails for creating web applications: oht would be a framework for creating decentalized applications. This would include tools to build a boilerplate project, or possibly a package of configuration files and templates that can be loaded into the client. 

## Contribute

Everyone is encouraged to test out the software, report bugs and submit pull requests. 
