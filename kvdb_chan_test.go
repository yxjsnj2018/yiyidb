package yiyidb

import (
	"fmt"
	"testing"
	"path/filepath"
	"os"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestKvdb_AllByKVChan(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, true, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	kv.PutChan("jac", []byte("yudfuud dekjrker"), 0)
	kv.PutChan("jac", []byte("dfdfseeee ee value2"), 0)

	kv.PutChan("yum", []byte("test value1"), 0)
	kv.PutChan("yum", []byte("test value2"), 0)

	kv.Del(idToKey("jac", 1))

	all := kv.AllByKVChan("jac")
	assert.Equal(t, all[0].Value, []byte("dfdfseeee ee value2"))

	kv.Drop()
}
