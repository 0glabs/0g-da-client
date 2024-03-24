# 0G DA Spec

## Organization

The 0G DA repo is organized as a monorepo, with each project adhering to the "Ben Johnson" project structure style. Within the core project directories (e.g., `core`, `disperser`, `retriever`), the main interfaces and data types are defined at the root of the project, while implementations are organized by dependency.

In general, the `core` project contains implementation of all the important business logic responsible for the security guarantees of the 0G DA protocol, while the other projects add the networking layers needed to run the distributed system.

## Directory structure

```
┌── : api: Protobuf definitions
├── : common: contract bindings and other basic components
┌── : core: Core logic of the 0G DA protocol
├── : disperser: Disperser service
├── : docs: Documentation and specification
├── : pkg: Dependent pkg
|   ├── : encoding: Core encoding/decoding functionality and multiproof generation
|   └── : kzg: kzg libraries
├── : retriever: Retriever service
├── : tests: Tools for running integration tests
```
