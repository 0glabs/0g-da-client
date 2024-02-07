package core

func AllocateChunks(encodedBlobs []*BlobLocation) uint {
	n := len(encodedBlobs)
	allocated := make([]int, n)
	segments := uint(0)
	for i := 0; i < n; {
		offset := uint(0)
		// allocate by column
		for j := i; i < n; {
			if allocated[j] == int(encodedBlobs[j].ChunkNum) {
				// encoded blob is fully allocated
				if j == i {
					i++
				}
			} else {
				// try to fill one chunk + proof
				l := encodedBlobs[j].ChunkLength*CoeffSize + ProofSize
				if offset+l <= SegmentSize {
					encodedBlobs[j].SegmentIndexes[allocated[j]] = segments
					encodedBlobs[j].Offsets[allocated[j]] = offset
					allocated[j]++
					offset += l
				} else {
					break
				}
			}
			// move to next blob
			j++
			if j >= n {
				j = i
			}
		}
		if offset > 0 {
			segments++
		}
	}
	return segments
}
