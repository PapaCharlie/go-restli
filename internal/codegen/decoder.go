package codegen

import (
	. "github.com/dave/jennifer/jen"
)

type decoder struct {
	*Statement
}

var Decoder = &decoder{Id("decoder")}

func (d *decoder) ReadObject(reader func(field *Statement, def *Group)) *Statement {
	field := Id("field")
	return Add(d).Dot("ReadObject").Call(Func().Params(Add(field).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		reader(field, def)
	}))
}

func (d *decoder) ReadMap(reader func(key *Statement, def *Group)) *Statement {
	key := Id("key")
	return Add(d).Dot("ReadMap").Call(Func().Params(Add(key).String()).Params(Err().Error()).BlockFunc(func(def *Group) {
		reader(key, def)
	}))
}

func (d *decoder) ReadArray(creator func(index *Statement, def *Group)) *Statement {
	index := Id("index")
	return Add(d).Dot("ReadObject").Call(Func().Params(Add(index).Int()).Params(Err().Error()).BlockFunc(func(def *Group) {
		creator(index, def)
	}))
}
