package mydb

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

type DBInfo struct {
	Id string
	Pwd string
	Host string
	Port int
	Name string
}

func Connect(c DBInfo) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Id, c.Pwd, c.Host, c.Port, c.Name)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func to_Name(name string) (ret string) {
	re := regexp.MustCompile("[A-Z]*[a-z0-9]*")
	s := re.FindAllString(name, -1)
	for i, v := range s {
		ret += strings.ToLower(v)
		if i < len(s)-1 {
			ret += "_"
		}
	}
	return
}

func primary_key(typ reflect.Type) (flag bool, pkConstraint string) {
	cnt := 0
	pkConstraint = ", constraint %s_pk primary key("
	for i:=0; i<typ.NumField(); i++ {
		if typ.Field(i).Tag.Get("key")  == "primarykey" {
			cnt += 1
			if cnt == 1 {
				pkConstraint += to_Name(typ.Field(i).Name)
			} else {
				pkConstraint = pkConstraint+ ", " + to_Name(typ.Field(i).Name)
			}
		}
	}
	pkConstraint += "))"
	flag = (cnt > 0)
	return
}

func get_field_name(field reflect.StructField) (name string) {
	name = to_Name(field.Name)
	if field.Tag.Get("db") != "" {
		name = field.Tag.Get("db")
	}
	return
}

func get_field(field reflect.StructField) (ret string) {
	name := get_field_name(field)
	typ := get_type(field)
	ret = fmt.Sprintf("%s %s", name, typ)
	return
}

func get_type(field reflect.StructField) (ret string) {
	if field.Tag.Get("typ") != "" {
		ret = field.Tag.Get("typ")
		return
	}

	switch field.Type.Name() {
	case "int":
		ret = "bigint"
	case "int8", "int16", "int32":
		ret = "int"
	case "int64":
		ret = "bigint"
	case "string":
		length := field.Tag.Get("length")
		if length == "" {
			length = "255"
		}
		ret = fmt.Sprintf("varchar(%s)", length)
	case "float32":
		ret = "float"
	case "float64":
		ret = "double"
	case "time", "time.Time", "Time":
		ret = "datetime"
	case "bool":
		ret = "boolean"
	}
	return
}

func create(tbl interface{}) (query string) {
	st := reflect.ValueOf(tbl).Elem().Type()
	tname := to_Name(st.Name())
	pkFlag, pkConstraint := primary_key(st)
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tname)
	num := st.NumField()
	for i:=0; i<num; i++ {
		txt := get_field(st.Field(i))
		query += txt
		if i < num-1 {
			query += ", "
		}
	}
	if pkFlag {
		pkConstraint = fmt.Sprintf(pkConstraint, tname)
		query += pkConstraint
	} else {
		query += ")"
	}
	return query
}

func Create(tbl interface{}) {
	query := insert(tbl)
	// res, err := DB.Exec(query)
	DB.Exec(query)
}

func get_keys(st reflect.Type) (keys map[string]string) {
	return 
}

func get_field_value(field reflect.StructField, value reflect.Value) (ret string) {
	typ := get_type(field)
	switch typ {
	case "int", "bigint":
		ret = fmt.Sprintf("%d", value.Int())
	case "float", "double":
		ret = fmt.Sprintf("%f", value.Float())
	case "datetime":
		ret = "'"+ value.Interface().(time.Time).Format("2006-01-02 15:04:05") +"'"
	case "boolean":
		if value.Bool() {
			ret = "true"
		} else {
			ret = "false"
		}
	default:
		ret = "'"+ value.String()+"'"
	}
	return ret
}

func insert(data interface{}) (query string) {
	st := reflect.ValueOf(data).Elem()
	tname := to_Name(st.Type().Name())
	fieldnames := ""
	values := ""
	num := st.Type().NumField()
	for i:=0; i<num; i++ {
		fieldnames += get_field_name(st.Type().Field(i))
		values += get_field_value(st.Type().Field(i), st.Field(i))
		if i < num-1 {
			fieldnames += ","
			values += ","
		}
	}
	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tname, fieldnames, values)
	return
}

func Insert(data interface{}) {
	query := insert(data)
	// res, err := DB.Exec(query)
	DB.Exec(query)
}

func get_value(v interface{}) (ret string) {
	
	switch reflect.TypeOf(v).Name() {
	case "int32":
		ret = fmt.Sprintf("%d", v.(int32))
	case "int64":
		ret = fmt.Sprintf("%d", v.(int64))
	case "int":
		ret = fmt.Sprintf("%d", v.(int))
	case "float32":
		ret = fmt.Sprintf("%f", v.(float32))
	case "float64":
		ret = fmt.Sprintf("%f", v.(float64))
	case "bool":
		if v.(bool) {
			ret = "true"
		} else {
			ret = "false"
		}
	case "time", "Time":
		ret = "'"+ v.(time.Time).Format("2006-01-02 15:04:05")+"'"
	case "string":
		ret = "'"+v.(string)+"'"
	}
	return
}

func update(tbl string, cond map[string]interface{}, value map[string]interface{}) (query string) {
	set := ""
	n1 := len(value)
	i := 0
	for k, v := range value {
		i++
		set += k + "=" + get_value(v)
		if i < n1 {
			set += ","
		} 
	}
	
	n2 := len(cond)
	if n2 == 0 {
		query = fmt.Sprintf("UPDATE %s SET %s", tbl, set)
	} else {
		condition := ""
		
		i = 0
		for k, v := range cond {
			i++
			condition += k + "=" + get_value(v)
			if i < n1 {
				condition += " and "
			} 
		}
		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", tbl, set, condition)
	}
	return
}

func Update(tbl string, cond map[string]interface{}, value map[string]interface{}) {
	query := update(tbl, cond, value)
	// res, err := DB.Exec(query)
	DB.Exec(query)
}

func upsert_compare(data interface{}) (query string) {
	return
}

func Upsert_compare(data interface{}) {
	query := insert(data)
	// res, err := DB.Exec(query)
	DB.Exec(query)
}

func upsert(data interface{}) (query string) {
	query = insert(data) + " ON DUPLICATE KEY UPDATE "
	st := reflect.ValueOf(data).Elem()
	num := st.Type().NumField()
	ustr := ""
	for i:=0; i<num; i++ {
		if st.Type().Field(i).Tag.Get("key") != "primarykey" {
			fieldname := get_field_name(st.Type().Field(i))
		value := get_field_value(st.Type().Field(i), st.Field(i))
		ustr += fieldname +"=" + value + ","
		}
	}
	n := len(ustr)
	if ustr[n-1] == ',' {
		ustr = ustr[0:n-1]
	}
	query += ustr
	return
}

func Upsert(data interface{}) {
	query := upsert(data)
	// res, err := DB.Exec(query)
	DB.Exec(query)
}