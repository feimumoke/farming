package util

import (
	"fmt"
	"reflect"
)

var (
	StrType   = reflect.TypeOf("")
	BoolType  = reflect.TypeOf(true)
	Int64Type = reflect.TypeOf(int64(0))
	FloatType = reflect.TypeOf(float64(0))
)

type Builder struct {
	fieldList []reflect.StructField
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) AddField(field string, typ reflect.Type) *Builder {
	b.fieldList = append(b.fieldList, reflect.StructField{Name: field, Type: typ})
	return b
}

func (b *Builder) AddFieldWithTag(field string, typ reflect.Type, tagPair ...string) *Builder {
	if len(tagPair)%2 != 0 {
		panic("tagPair not have tag or value")
	}
	structField := newStructField(field, typ, tagPair...)
	b.fieldList = append(b.fieldList, *structField)
	return b
}

func (b *Builder) AddFieldBefore(base, field string, typ reflect.Type, tagPair ...string) *Builder {
	if len(tagPair)%2 != 0 {
		panic("tagPair not have tag or value")
	}
	for i, sf := range b.fieldList {
		if sf.Name == base {
			structField := newStructField(field, typ, tagPair...)
			b.fieldList = append(b.fieldList[0:i], append([]reflect.StructField{*structField}, b.fieldList[i:]...)...)
			return b
		}
	}
	return b
}

func (b *Builder) AddFieldAfter(base, field string, typ reflect.Type, tagPair ...string) *Builder {
	if len(tagPair)%2 != 0 {
		panic("tagPair not have tag or value")
	}
	for i, sf := range b.fieldList {
		if sf.Name == base {
			structField := newStructField(field, typ, tagPair...)
			b.fieldList = append(b.fieldList[0:i+1], append([]reflect.StructField{*structField}, b.fieldList[i+1:]...)...)
			return b
		}
	}
	return b
}

func newStructField(field string, typ reflect.Type, tagPair ...string) *reflect.StructField {
	structField := reflect.StructField{Name: field, Type: typ}
	if len(tagPair) > 0 {
		tags := ""
		for i := 0; i < len(tagPair)/2; i++ {
			tag := tagPair[2*i]
			val := tagPair[2*i+1]
			tags += fmt.Sprintf(`%v:"%v" `, tag, val)
		}
		structField.Tag = reflect.StructTag(tags)
	}
	return &structField
}

// 根据预先添加的字段构建出结构体
func (b *Builder) Build() *Struct {
	stu := reflect.StructOf(b.fieldList)
	index := make(map[string]int)
	for i := 0; i < stu.NumField(); i++ {
		index[stu.Field(i).Name] = i
	}
	return &Struct{stu, index}
}

// 结构体的类型
type Struct struct {
	typ   reflect.Type
	index map[string]int
}

func (s Struct) New() *Instance {
	return &Instance{reflect.New(s.typ).Elem(), s.index}
}

// 结构体的值
type Instance struct {
	instance reflect.Value
	index    map[string]int
}

func (in Instance) Field(name string) (bool, reflect.Value) {
	if i, ok := in.index[name]; ok {
		return true, in.instance.Field(i)
	} else {
		return false, reflect.Value{}
	}
}

func (in *Instance) SetValue(name string, value reflect.Value) {
	if i, ok := in.index[name]; ok {
		in.instance.Field(i).Set(value)
	}
}

func (i *Instance) Interface() interface{} {
	return i.instance.Interface()
}

func (i *Instance) Addr() interface{} {
	return i.instance.Addr().Interface()
}

func (i *Instance) Pointer() *reflect.Value {
	val := i.instance.Addr()
	return &val
}
