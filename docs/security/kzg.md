# KZG Encoder Backend

It is important that the encoding and commitment tasks are able to be performed in seconds and that the dominating complexity of the computation is nearly linear in the degree of the polynomial. This is achieved using algorithms based on the Fast Fourier Transform (FFT) and amortized kzg multi-reveal generation.

This document describes how the KZG-FFT encoder backend implements the `Encode(data [][]byte, params EncodingParams) (BlobCommitments, []*Chunk, error)` interface, which

1. transforms the blob into a list of `params.NumChunks` `Chunks`, where each chunk is of length `params.ChunkLength`.
2. produces the associated polynomial commitments and proofs.

We will also highlight the additional constraints on the Encoding interface which arise from the KZG-FFT encoder backend.

## Deriving the polynomial coefficients and commitment

As described in the [Encoding Module Specification](encoding.md), given a blob of data, we convert the blob to a polynomial $$p(X) = \sum_{i=0}^{m-1} c_iX^i$$ by simply slicing the data into a string of symbols, and interpreting this list of symbols as the tuple $$(c_i)_{i=0}^{m-1}$$.

In the case of the KZG-FFT encoder, the polynomial lives on the field associated with the BN-254 elliptic curve, which as order \[TODO: fill in order].

Given this polynomial representation, the KZG commitment can be calculated as in [KZG polynomial commitments](https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html).

## Polynomial Evaluation with the FFT

In order to use a Discrete Fourier Transform (DFT) to evaluate a polynomial, the indices of the polynomial evaluations which will make up the Chunks must be members of a cyclic group, which we will call $$S$$. A cyclic group is the group generated by taking all of the integer powers of some generator $$v$$, i.e., $${v^k | k \in \mathbb{Z} }$$ (For this reason, the elements of a cyclic group $$S$$ of order $$|S|=m$$ will sometimes be referred to as the $$|m|$$’th roots of unity). Notice that since our polynomial lives on the BN254 field, the group $$S$$ must be a subgroup of that field (i.e. all if its elements must lie within that field).

Given a cyclic group $$S$$ of order $$m$$, we can evaluate a polynomial $$p(X)$$ of order $$n$$ at the indices contained in $$S$$ via the DFT,

$$
p_k = \sum_{i=1}^{n}c_i (v^k)^i
$$

where $$p_k$$ gives the evaluation of the polynomial at $$v^k \in S$$. Letting $$c$$ denote the vector of polynomial coefficients and $$p$$ the vector of polynomial evaluations, we can use the shorthand $$p = DFT[c]$$. The inverse relation also holds, i.e., $$c = DFT^{-1}[p]$$.

To evaluate the DFT programmatically, we want $$m = n$$. Notice that we can achieve this when $$m > n$$ by simply padding $$c$$ with zeros to be of length $$m$$.

The use of the FFT can levy an additional requirement on the size of the group $$S$$. In our implementation, we require the size of $$S$$ to be a power of $$\mathsf{2}$$. For this, we can make use of the fact that the prime field associated with BN-254 contains a subgroup of order $$2^{28}$$, which in turn contains subgroups of orders spanning every power of 2 less than $$2^{28}$$.

As the encoding interface calls for the construction of $$\mathsf{NumChunks}$$ Chunks of length $$\mathsf{ChunkLength}$$, our application requires that $$S$$ be of size $$\mathsf{NumChunks}\times \mathsf{ChunkLength}$$, which in turn must be a power of $$\mathsf{2}$$.

## Amortized Multireveal Proof Generation with the FFT

The construction of the multireveal proofs can also be performed using a DFT (as in ["Fast Amortized Kate Proofs"](https://eprint.iacr.org/2023/033.pdf)). Leaving the full details of this process to the referenced document, we describe here only 1) the index-assignment scheme used by the amortized multiproof generation approach and 2) the constraints that this creates for the overall encoder interface.

Given the group $$S$$ corresponding to the indices of the polynomial evaluations and a cyclic group $$C$$ which is a subgroup of $$S$$, the cosets of $$C$$ in $$S$$ are given by

$$
s+C = {g+c : c \in C} \text{ for } s \in S.
$$

Each coset $$s+C$$ has size $$|C|$$, and there are $$|S|/|C|$$ unique and disjoint cosets.

Given a polynomial $$p(X)$$ and the groups $$S$$ and $$C$$, the Amortized Kate Proofs approach generates $$|S|/|C|$$ different KZG multi-reveal proofs, where each proof is associated with the evaluation of $$p(X)$$ at the indices contained in a single coset $$sC$$ for $$s \in S$$. Because the Amortized Kate Proofs approach uses the FFT under the hood, $$C$$ itself must have an order which is a power of $$\mathsf{2}$$.

For the purposes of the KZG-FFT encoder, this means that we must choose $$S$$ to be of size $$\mathsf{NumChunks}\times \mathsf{ChunkLength}$$ and $$C$$ to be of size $$\mathsf{ChunkLength}$$, each of which must be powers of $$\mathsf{2}$$.

## Worked Example

As a simple illustrative example, suppose that `AssignmentCoordinator` provides the following parameters in order to meet the security requirements of given blob:

* $$\mathsf{ChunkLength = 3}$$
* $$\mathsf{NumChunks = 4}$$

Supplied with these parameters, `Encoder.GetEncodingParams` will upgrade $$\mathsf{ChunkLength}$$ to the next highest power of $$\mathsf{2}$$, i.e., $$\mathsf{ChunkLength = 4}$$, and leave $$\mathsf{NumChunks}$$ unchanged. The following figure illustrates how the indices will be assigned across the chunks in this scenario.

<figure><img src="../../../.gitbook/assets/zg-da-encoding-groups.png" alt=""><figcaption><p>Worked example of chunk indices for ChunkLength=4, NumChunks=4</p></figcaption></figure>
