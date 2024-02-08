# Introduction

The ZeroGDA (ZGDA) system is a scalable Data Availability (DA) service layer which is directly built on top of a decentralized storage system and addresses the scalability issue by minimizing the data transfer volume required for broadcast. Its general decentralized storage design further enables it to support a variety of availability data types from diversified scenarios not limited to Layer 2 networks but also inclusion of decentralized AI infrastructures.

&#x20;In informal terms, DA is a guarantee that a given piece of data is available to anyone who wishes to retrieve it. ZGDA is focused on providing DA with both high security and throughput.

At a high level, a DA system is one which accepts blobs of data via some interface and then makes them available to retrievers through another interface.

Two important aspects of a DA system are:

1. [Security](overview-1/): The security of a DA system constitutes the set of conditions which are sufficient to ensure that all data blobs certified by the system as available are indeed available for retrievers to download.
2. Throughput: The throughput of a DA system is the rate at which the system is able to accept blobs of data, typically measured in bytes/second.
