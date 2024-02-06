# docs

## Organization

The ZGDA repo is organized as a monorepo, with each project adhering to the "Ben Johnson" project structure style. Within the core project directories (e.g., `core`, `disperser`, `retriever`), the main interfaces and data types are defined at the root of the project, while implementations are organized by dependency.

In general, the `core` project contains implementation of all the important business logic responsible for the security guarantees of the ZGDA protocol, while the other projects add the networking layers needed to run the distributed system.

## Directory structure

```
┌──  Protobuf definitions and contract bindings
┌── : Core logic of the ZGDA protocol
├── : Disperser service
├── : Documentation and specification
├── : Dependent pkg
|   ├── : Core encoding/decoding functionality and multiproof generation
|   └── : kzg libraries
├── : Retriever service
├── : Tools for running integration tests
```
