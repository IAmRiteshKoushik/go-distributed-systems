# Distributed Systems from Scratch

This repository tracks the exercises from the Build Distributed Systems from Scratch course.

## Course Index

### Foundations
Core message passing, IDs, and broadcast primitives.

- [Messenger](./messenger)
- [Identifier](./identifier)
- [Gossipper](./gossipper)

### Consensus
State agreement, leader election, and replicated storage.

- [Counter](./counter)
- [Elector](./elector)
- [Consensus](./consensus)
- [Store](./store)

### Scalability
Building blocks for scaling reads, routing, and work distribution.

- [Caches](./caches)
- [Proxies](./proxies)
- [Indexes](./indexes)
- [Load Balancers](./load_balancers)
- [Queues](./queues)

### Advanced
Larger distributed workflows and coordination-heavy systems.

- [Sharder](./sharder)
- [Coordinator](./coordinator)
- [Advanced](./advanced)

## Notes

Each directory is a self-contained Go module with a `main.go` entrypoint and a short module-specific README.
