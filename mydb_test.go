package mydb

import (
	"testing"
	"time"
)

type TestStruct1 struct {
	VarInt8 int8 `key:"primarykey"`
	VarInt32 int32
	VarInt int
	VarInt64 int64
	VarFloat32 float32
	VarFloat64 float64
	VarBool bool
	VarTime time.Time
	VarString string
}

type TestStruct2 struct {
	VarInt8 int8 `db:"vint8"`
	VarInt32 int32 `db:"vint32"`
	VarInt int `db:"vint"`
	VarInt64 int64 `db:"vint64"`
	VarFloat32 float32 `db:"vfloat"`
	VarFloat64 float64 `db:"vdouble"`
	VarBool bool `db:"vbool"`
	VarTime time.Time `db:"vtime"`
	VarString string `db:"vstring"`
}

type TestStruct3 struct {
	VarInt8 int8 `db:"vint8" typ:"int"`
	VarInt32 int32 `db:"vint32" typ:"bigint"`
	VarInt int `db:"vint" typ:"int"`
	VarInt64 int64 `db:"vint64" typ:"int"`
	VarFloat32 float32 `db:"vfloat" typ:"string"`
	VarFloat64 float64 `db:"vdouble"` 
	VarBool bool `db:"vbool"`
	VarTime time.Time `db:"vtime"`
	VarString string `db:"vstring"`
}
func TestCreate(t *testing.T) {
	q := create(&TestStruct1{})
	t1 := "CREATE TABLE IF NOT EXISTS test_struct1 (var_int8 int, var_int32 int, var_int bigint, var_int64 bigint, "
	t1 += "var_float32 float, var_float64 double, var_bool boolean, var_time datetime, var_string varchar(255), "
	t1 += "constraint test_struct1_pk primary key(var_int8))"
	if q != t1 {
		t.Errorf("Error\n예상값\n%s\n출력값\n%s", t1, q)
	} else {
		t.Log("TestStruct1 Create OK")
	}

	q = create(&TestStruct2{})
	t1 = "CREATE TABLE IF NOT EXISTS test_struct2 (vint8 int, vint32 int, vint bigint, vint64 bigint, "
	t1 += "vfloat float, vdouble double, vbool boolean, vtime datetime, vstring varchar(255))"
	if q != t1 {
		t.Errorf("Error\n예상값\n%s\n출력값\n%s", t1, q)
	} else {
		t.Log("TestStruct2 Create OK")
	}

	q = create(&TestStruct3{})
	t1 = "CREATE TABLE IF NOT EXISTS test_struct3 (vint8 int, vint32 bigint, vint int, vint64 int, "
	t1 += "vfloat string, vdouble double, vbool boolean, vtime datetime, vstring varchar(255))"
	if q != t1 {
		t.Errorf("Error\n예상값\n%s\n출력값\n%s", t1, q)
	} else {
		t.Log("TestStruct3 Create OK")
	}
}

func TestInsert(t *testing.T) {
	b, _ := time.Parse("2006-01-02 15:04:05", "2021-12-20 11:43:00")
	var a TestStruct1 = TestStruct1{1, 2, 3, 4, 5.0, 6.0, true, b, "abcd"}
	q := insert(&a)
	t1 := "INSERT INTO test_struct1 (var_int8,var_int32,var_int,var_int64,var_float32,var_float64,"
	t1 += "var_bool,var_time,var_string) VALUES (1,2,3,4,5.000000,6.000000,true,'2021-12-20 11:43:00','abcd')"
	if q != t1 {
		t.Errorf("Error\n예상값\n%s\n출력값\n%s", t1, q)
	} else {
		t.Log("TestStruct1 Insert OK")
	}
	var c TestStruct2 = TestStruct2{1, 2, 3, 4, 5.0, 6.0, true, b, "abcd"}
	q = insert(&c)
	t1 = "INSERT INTO test_struct2 (vint8,vint32,vint,vint64,vfloat,vdouble,"
	t1 += "vbool,vtime,vstring) VALUES (1,2,3,4,5.000000,6.000000,true,'2021-12-20 11:43:00','abcd')"
	if q != t1 {
		t.Errorf("Error\n예상값\n%s\n출력값\n%s", t1, q)
	} else {
		t.Log("TestStruct2 Insert OK")
	}
}

func TestUpdate(t *testing.T) {
	var cond, value map[string]interface{}
	cond = make(map[string]interface{})
	value = make(map[string]interface{})
	cond["a"] = 1
	cond["b"] = "b"
	cond["c"] = time.Now()
	value["a"] = 2
	value["d"] = 3
	value["e"] ="c"
	str := update("abc", cond, value)
	t.Error(str)
}

func TestUpsert(t *testing.T) {
	// b, _ := time.Parse("2006-01-02 15:04:05", "2021-12-20 11:43:00")
	// var a TestStruct1 = TestStruct1{1, 2, 3, 4, 5.0, 6.0, true, b, "abcd"}
	// q := upsert(&a)
	// t.Error(q)
	type Schedule struct {
		Item string `db:"item" key:"primarykey" length:"20"`
		Date time.Time `db:"date"`
	}
	s := Schedule{"abc", time.Now()}
	q := upsert(&s)
	t.Error(q)


}