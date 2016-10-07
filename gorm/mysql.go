package dbannotate

import (
"fmt"
"reflect"
)

func test() {
	type S struct {
		F string `dbcolumn:"gopher" ispk:"true"`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	fmt.Println(field.Tag.Get("dbcolumn"), field.Tag.Get("ispk"))

}
