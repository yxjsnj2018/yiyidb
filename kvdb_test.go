package yiyidb

import (
	"testing"
	"path/filepath"
	"os"
	"fmt"
	"time"
	"strconv"
	"github.com/stretchr/testify/assert"
)

func TestKvdb_AllByKVMix(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	//参数说明
	//1数据库路径,2是否开启ttl自动删除记录,3数据碰测优化，输入可能出现key的最大长度
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	kv.PutMix( "jac", "testkey1", []byte("testkey1"), 0)
	kv.PutMix( "jac", "testkey22", []byte("testkey22"), 0)
	kv.PutMix( "jac", "testke", []byte("testk3563453e"), 0)

	one, err := kv.GetMix("jac", "testke")
	assert.NoError(t, err)
	assert.Equal(t, one, []byte("testk3563453e"))

	kv.DelColMix("jac", "testkey1")
	all := kv.AllByKVMix("jac")
	assert.Equal(t, all[0].Value,[]byte("testk3563453e"))

	kv.DelMix("jac")
	all = kv.AllByKVMix("jac")
	assert.Equal(t, len(all), 0)

	type object struct {
		Value int
	}

	kv.PutObjectMix("jjj","test",&object{111}, 0)
    var o object
    err = kv.GetObjectMix("jjj","test", &o)
    assert.NoError(t, err)
	assert.Equal(t, o.Value, 111)

	all = kv.AllByObjectMix("jjj", o)
	assert.Equal(t, all[0].Object, &object{111})

	kv.Drop()
}

func TestKvdb_IterStartWith(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	//参数说明
	//1数据库路径,2是否开启ttl自动删除记录,3数据碰测优化，输入可能出现key的最大长度
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	kv.Put([]byte("testkey1"), []byte("testkey1"), 0)
	kv.Put([]byte("testkey22"), []byte("testkey22"), 0)
	kv.Put([]byte("testke"), []byte("testke"), 0)

	iter, err := kv.IterStartWith([]byte("testkey"))
	assert.NoError(t, err)
	iter.Next()
	assert.Equal(t, iter.Key(), []byte("testkey1"))

	kv.Drop()
}

func TestKvdb_KeyRangeByObject(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	//参数说明
	//1数据库路径,2是否开启ttl自动删除记录,3数据碰测优化，输入可能出现key的最大长度
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	type object struct {
		Value int
	}
	kv.PutObject([]byte("testkey1"), object{1}, 0)
	kv.PutObject([]byte("testkey22"), object{2}, 0)
	kv.PutObject([]byte("testke"), object{3}, 0)

	var o object
	all, err := kv.KeyRangeByObject([]byte("testkey"), []byte("testkey25"), o)
	assert.NoError(t, err)
	assert.Equal(t, all[0].Object, &object{1})

	kv.Drop()
}

func TestKvdb_KeyStartByObject(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	type object struct {
		Value int
	}

	kv.PutObject([]byte("testkey1"), object{1}, 0)
	kv.PutObject([]byte("testkey2"), object{2}, 0)
	kv.PutObject([]byte("testke"), object{3}, 0)

	var o object
	all, err := kv.KeyStartByObject([]byte("testkey"), o)
	assert.NoError(t, err)
	assert.Equal(t, all[0].Object, &object{1})

	kv.Drop()
}

func TestKvdb_KeysByRegexp(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()
	topics := []string{"foo/bar/baz",
		"foo/jhg/baz",
		"foo/uyt/kkk",
		"foo/bar/+",
		"foo/#",
		"foo/bar/#",
		"foo/ytr/#"}
	for _, v := range topics {
		kv.Put([]byte(v), []byte(v), 0)
	}
	all, err := kv.RegexpKeys(`foo/((?:[^/#+]+/)*)`)
	assert.NoError(t, err)
	assert.Equal(t, all[0], "foo/#")
	kv.Drop()
}

func TestKvdb_AllByKV(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	kv.Put([]byte("testkey"), []byte("test value1"), 0)
	kv.Put([]byte("testkey1"), []byte("test value2"), 0)

	kv.Put([]byte("testkey4"), []byte("test value1"), 0)
	kv.Put([]byte("testkey5"), []byte("test value2"), 0)

	all := kv.AllByKV()
	assert.Equal(t, all[0].Value, []byte("test value1"))

	kv.Drop()
}

func TestKvdb_AllByObject(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	type object struct {
		Value int
	}

	kv.PutObject([]byte("testkey"), object{1}, 0)
	kv.PutObject([]byte("testkey1"), object{2}, 0)

	var o object
	all := kv.AllByObject(o)
	assert.Equal(t, all[0].Object, &object{1})

	kv.Drop()
}

func TestKvdb_Drop(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, false, false, 10)
	assert.NoError(t, err)
	defer kv.Close()

	kv.Drop()
}

//func TestKvdb_Put(t *testing.T) {
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		panic(err)
//	}
//	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
//	kv, err := OpenKvdb(dir, false, false, 10)
//	if err != nil {
//		panic(err)
//	}
//	defer kv.Close()
//
//	start := time.Now()
//	//for i := 0; i < 1000000; i++ {
//	ks, err := kv.KeyRange([]byte("key0"), []byte("key1000000"))
//	if err != nil {
//		panic(err)
//	}
//	//}
//	exp := time.Now().Sub(start)
//	fmt.Println(exp, ks)
//
//	kv.Drop()
//}

func TestKvdb_BatPutOrDel(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, false, false, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	items := make([]BatItem, 0)
	//add 50000 items to database
	for i := 1; i < 50000; i++ {
		item := BatItem{
			Op:    "put",
			Ttl:   0,
			Key:   []byte("test" + strconv.Itoa(i)),
			Value: []byte("bat values"),
		}
		items = append(items, item)
	}
	kv.BatPutOrDel(&items)
	assert.NoError(t, err)
	last, err := kv.Get([]byte("test9999"))
	assert.NoError(t, err)
	assert.Equal(t, last, []byte("bat values"))

	//remove 50000 items from database
	for i := 1; i < 50000; i++ {
		item := BatItem{
			Op:    "del",
			Key:   []byte("test" + strconv.Itoa(i)),
			Value: []byte("bat values"),
		}
		items = append(items, item)
	}
	kv.BatPutOrDel(&items)
	_, err = kv.Get([]byte("test9999"))
	assert.Equal(t, err.Error(), "leveldb: not found")

	kv.Drop()
}

func TestSetTtlKvdb(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	kv, err := OpenKvdb(dir, true, true, 10)
	if err != nil {
		panic(err)
	}
	defer kv.Close()

	kv.Put([]byte("hello1"), []byte("hello value"), 3)
	kv.Put([]byte("hello2"), []byte("hello value2"), 10)

	kv.SetTTL([]byte("hello2"), 8)

	f, err := kv.GetTTL([]byte("hello2"))
	assert.NoError(t, err)
	assert.Equal(t, int(f)+1, 8)

	v, err := kv.Get([]byte("hello1"))
	assert.NoError(t, err)
	assert.Equal(t, v, []byte("hello value"))

	time.Sleep(5 * time.Second)

	_, err = kv.Get([]byte("hello1"))
	if err == nil {
		t.Error("ttl delete error")
	}

	kv.Drop()
}

//func TestTtlRunner_Run(t *testing.T) {
//	//测试kv数据库及超时TTL服务
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		panic(err)
//	}
//	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
//	kv, err := OpenKvdb(dir, true, true, 10)
//	if err != nil {
//		panic(err)
//	}
//	defer kv.Close()
//	kv.OnExpirse = func(key, value []byte) {
//		fmt.Println("exp:", string(key), string(value))
//	}
//
//	for i := 0; i <= 10; i++ {
//		kv.Put([]byte("hello"+strconv.Itoa(i)), []byte("hello value"+strconv.Itoa(i)), i+1)
//	}
//
//	for i := 1; i < 10; i++ {
//		fmt.Println("sleep")
//		time.Sleep(3 * time.Second)
//
//		all := kv.AllKeys()
//		for _, k := range all {
//			fmt.Println(k)
//		}
//	}
//
//	searchkeys, _ := kv.KeyStart([]byte("hello1"))
//	for _, k := range searchkeys {
//		fmt.Println(k)
//	}
//
//	randkeys, _ := kv.KeyRange([]byte("2017-06-01T01:01:01"), []byte("2017-07-01T01:01:01"))
//	for _, k := range randkeys {
//		fmt.Println(k)
//	}
//}

//func TestMaxKeyAndValue(t *testing.T) {
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		panic(err)
//	}
//	dir = dir + "/" + fmt.Sprintf("test_db_%d", time.Now().UnixNano())
//
//	kv, err := OpenKvdb(dir, false, false, 10)
//	if err != nil {
//		panic(err)
//	}
//
//	defer kv.Drop()
//
//	bts, err := ioutil.ReadFile("D:\\GoglandProjects\\yiyidb\\yiyidb.zip")
//	if err != nil {
//		panic(err)
//	}
//
//	err = kv.Put(bts, bts, 0)
//	if err != nil {
//		panic(err)
//	}
//	btss, err := kv.Get(bts)
//	if err != nil {
//		panic(err)
//	}
//	err = ioutil.WriteFile("D:\\GoglandProjects\\yiyidb\\yiyidb1.zip", btss, 777)
//	if err != nil {
//		panic(err)
//	}
//	//测试结果显示 key和value不超过512m为保证可用性，但性能下降，建议把key尽可能精练
//}
