# Pow challenge

+ [Go 1.18](https://go.dev/dl/)

In this application, I designed and implemented a "Word of Wisdom" tcp server, 
which is protected from DDOS attacks using the "Proof of work" algorithm.

## Choice of “Proof of work” algorithm
In the process of choosing an algorithm, I considered the following:

### Hashcash
https://ru.wikipedia.org/wiki/Hashcash \
Is a proof-of-work system used to reduce spam and DoS attacks

#### Pros
+ Well documented
+ Easy enough to implement
+ Easy header validation by the server
+ Dynamic throttling.

#### Cons
+ Hashcash requires significant computational for a client  and may not be suitable for systems with low performance

### Merkle tree
https://en.wikipedia.org/wiki/Merkle_tree \
A hash tree allows efficient and secure verification of the contents of a large data structure.

#### Pros
+ Used in many cryptocurrencies such as bitcoin and ethereum
+ Used in the Rate Limiting Nullifier algoritms (Cloudflare) for protecting from spammers

#### Cons
+ Requires complex calculations from the server


### Guided tour puzzle protocol
https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol \
Guided tour puzzle (GTP) protocol is a cryptographic protocol for mitigating application layer denial of service attacks.

#### Pros
+ Dynamically enables protection if the server suspects it is currently under denial of service attack
+ Clients are not burdened with heavy computations

#### Cons
+ More complex implementation than Hashcash algorithm
+ Guided tour puzzle protocol enforces delay on the clients through round trip delays,

### Conclusion
After a detailed analysis of the algorithms, 
I decided to use Hashcash because it was used in similar systems to protect against spam and it is easy to implement.

## Clone the project

```
$ git clone git@github.com:mirogindev/pow-challenge.git
$ cd pow-challenge
```

## Start the project
### If you have docker-compose installed
```
make compose-start
```
#### Or
```
docker-compose up
```

### If you have Go 1.18 installed
```
make server-start
make client-start
```

### Run tests
```
make test
```


