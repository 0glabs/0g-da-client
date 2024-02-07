# Rollups

## Deciding Trust Assumption

Rollups need to define the quorums they want to sign off on the availability of their data and their trust assumptions on each of those quorums. The rollups need to define:

1. `AdversaryThresholdBPs`. This is the maximum ratio (in basis points) that can be adversarial in the quorum.
2. `QuorumThresholdBPs`. This is the minimum ratio (in basis points) that can need to sign on the availability of the rollups data for the rollup's contracts to consider the data available.

## Requests to Store Data

When the rollup has data that they want to make available, they construct a request to the disperser of the form [`BlobStoreRequest`](broken-reference) and receive a response of the form [`BlobStoreResponse`](broken-reference).

This flow is detailed [here](disperser.md#requests-to-store-data).

## Optimistic Rollups

### State Claims

When making state claims, validators of the optimistic rollup should resubmit the header of the dataStore (or load dataStore's quorum params from storage \[TODO: @gpsanant @Sidu28 fill this out after on chain verification]) and validate their trust assumptions from [Deciding Trust Assumptions](rollups.md#deciding-trust-assumption). In addition, the start and end degree of the rollup's data taken from the `BlobStoreResponse` received from the disperser.

### Revealing Data Onchain during Fraud Proofs (Optimistic rollups)

To keep the interface between EIP-4844 and ZGDA the same, optimistic rollups need to reveal data onchain against the commitment to their own data instead of the concatenated data of all of the `BlobStoreRequests` in the batch. To prove the rollups own data commitment against the batched (concatenated) commitment that was posted onchain in the dataStore header, rollups generate the following proof.

A challenger retrieves the data corresponding to the KZG commitment pointed to by validators of their rollup and parses their rollup's data from the claimed start and end degree. They then prove to a smart contract the commitment to the rollups data, along with their fraud proof, via the [subcommitment proofs](rollups.md#subcommitment-proof). Note that the rollup's smart contract will implement the verifier described in the proof.

\[TODO: Explain how the powers of tau are put on chain (use logarithmic adding)]

## ZK Rollups

## Subcommitment Proof

Call $C \in \mathbb{G}\_1$ is the commitment to the polynomial $c(x)$ batching the data among many `BlobStoreRequest`s. The retriever then finds the data of a rollup by taking the coefficients between degrees $m$ and $n$. Call the polynomial that corresponds to the rollup's data $b(x)$. In math:

$$c(x) = f(x) + x^n b(x) + x^m g(x)$$

Where $\text{degree}(f(x) < n)$, $\text{degree}(b(x)) < m-n$, and $\text{degree}(g(x)) < \text{degree}(c(x)) - m + 1$

Suppose a prover wants to prove the commitment to $b(x)$ to a lightweight verifier given the verifier has knowledge of $C$ and $m, n$.

The challenger can generate a proof of the commitment to $b(x)$, $B \in \mathbb{G}\_1$, by

* Generate a commitment to $f(x)$ and $x^{\text{max degree} - n}f(x)$. Call these $F$ and $L\_F$.
* Generate a commitment to $b(x)$ and $x^{\text{max degree} - (m - n)}b(x)$. Call these $B$ and $L\_B$.
* Generate a commitment to $g(x)$ and $x^{\text{max degree} - (\text{degree}(c(x)) - m + 1)}g(x)$. Call these $G$ and $L\_G$.
* Calculate $\gamma = keccak256(C, F, B, G)$.
* Calculate $\pi\_F$, the commitment to $\pi\_F(x) = \frac{f(x) - f(\gamma)}{x - \gamma}$
* Calculate $\pi\_B$, the commitment to $\pi\_B(x) = \frac{b(x) - b(\gamma)}{x - \gamma}$
* Calculate $\pi\_G$, the commitment to $\pi\_G(x) = \frac{g(x) - g(\gamma)}{x - \gamma}$
* Calculate $\pi\_C$, the commitment to $\pi\_C(x) = \frac{c(x) - c(\gamma)}{x - \gamma}$
* Calculate $\beta = keccak256(\gamma, \pi\_F, \pi\_B, \pi\_G, \pi\_C)$
* Calculate $\pi = \pi\_F + \beta \pi\_B + \beta^2 \pi\_G + \beta^3 \pi\_C$

The prover then submits to the verifier $F, B, G, L\_F, L\_B, L\_G, \pi, f(\gamma), b(\gamma), g(\gamma), c(\gamma)$ along with $C$ from the dataStore header of the blob in question. The verifier then verifies:

* $e(F, \[x^{\text{max degree} - n}]\_2) = e(L\_F, \[1]\_2)$. This verifies the low degreeness of $F$.
* $e(B, \[x^{\text{max degree} - (m - n)}]\_2) = e(L\_B, \[1]\_2)$. This verifies the low degreeness of $B$.
* $e(G, \[x^{\text{max degree} - (\text{degree} - m)}]\_2) = e(L\_G, \[1]\_2)$. This verifies the low degreeness of $G$.
* Calculate $\gamma = keccak256(C, F, B, G)$.
* Calculate $\beta = keccak256(\gamma, \pi\_F, \pi\_B, \pi\_G, \pi\_C)$
* $e(F - \[f(\gamma)]\_1 + \beta(B - \[b(\gamma)]\_1) + \beta^2(G - \[g(\gamma)]\_1) + \beta^3(C - \[c(\gamma)]\_1), \[1]\_2) = e(\pi, \[x-\gamma]\_2)$. This verifies a random opening of all of the claimed polynomials at the same x-coordinate.
* $c(\gamma) = f(\gamma) + \gamma^nb(\gamma) + \gamma^mg(\gamma)$. This verifies that the polynomials have the claimed shifted relationship with $c(x)$.
