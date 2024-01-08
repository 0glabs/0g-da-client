package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zero-gravity-labs/zgda/common"
	"github.com/zero-gravity-labs/zgda/common/store"
)

var (
	inmemBucketStoreSize = 1000
)

func TestLocalStore(t *testing.T) {

	localStore, err := store.NewLocalParamStore[common.RateBucketParams](inmemBucketStoreSize)
	assert.NoError(t, err)

	ctx := context.Background()

	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now(),
	}

	p2, err := localStore.GetItem(ctx, "testRetriever")
	assert.Error(t, err)
	assert.Nil(t, p2)

	err = localStore.UpdateItem(ctx, "testRetriever", p)
	assert.NoError(t, err)

	p2, err = localStore.GetItem(ctx, "testRetriever")

	assert.NoError(t, err)
	assert.Equal(t, p, p2)

}
