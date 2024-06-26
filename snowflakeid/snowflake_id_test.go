package snowflakeid

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowflakeId(t *testing.T) {
	worker, err := NewWorker(1)
	if err != nil {
		t.Error(err)
	}

	wg, lock := &sync.WaitGroup{}, &sync.Mutex{}
	set := make(map[int64]bool, 1000)

	f := func(wg *sync.WaitGroup) {
		defer wg.Done()

		id := worker.NextSnowflakeID()
		lock.Lock()
		defer lock.Unlock()
		set[id] = true
	}

	cnt := 100000
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go f(wg)
	}

	wg.Wait()
	t.Log("set len:", len(set))
	assert.Equal(t, cnt, len(set))
}

func TestInvalidWorkId(t *testing.T) {
	_, err := NewWorker(-1)
	assert.Error(t, err, errors.New("Worker ID 不合法"))
}
