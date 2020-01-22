package codegen

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/protocol"
	. "github.com/dave/jennifer/jen"
)

func (m *Method) actionFuncName() string {
	return ExportedIdentifier(m.Name + "Action")
}

func (m *Method) actionFuncParams(def *Group) {
	m.addEntityTypes(def)
	if len(m.Params) > 0 {
		def.Id("params").Op("*").Id(m.actionStructType())
	}
}

func (m *Method) actionStructType() string {
	return m.actionFuncName() + "Params"
}

func (m *Method) actionFuncReturnParams(def *Group) {
	if m.Return != nil {
		def.Add(m.Return.PointerType())
	}
	def.Error()
}

func (r *Resource) GenerateActionCode(a *Method) *CodeFile {
	actionName := a.Name + "Action"
	c := r.NewCodeFile(actionName)

	actionNameConst := ExportedIdentifier(actionName)
	c.Code.Const().Id(actionNameConst).Op("=").Lit(a.Name).Line()

	hasParams := len(a.Params) > 0
	if hasParams {
		record := &Record{
			NamedType: NamedType{
				Identifier: Identifier{
					Name:      a.actionStructType(),
					Namespace: r.Namespace,
				},
				Doc: fmt.Sprintf("This struct provides the parameters to the %s action", a.Name),
			},
			Fields: a.Params,
		}
		c.Code.Add(record.GenerateCode())
	}

	AddWordWrappedComment(c.Code, a.Doc).Line()
	r.addClientFunc(c.Code, a)

	c.Code.BlockFunc(func(def *Group) {
		var pathFunc string
		if a.OnEntity {
			pathFunc = ResourceEntityPath
		} else {
			pathFunc = ResourcePath
		}

		returns := a.Return != nil
		var errReturnParams []Code
		if returns {
			errReturnParams = []Code{Nil(), Err()}
		} else {
			errReturnParams = []Code{Err()}
		}

		def.List(Id(PathVar), Err()).Op(":=").Id(pathFunc).Call(a.entityParams()...)
		IfErrReturn(def, errReturnParams...).Line()
		def.Id(PathVar).Op("+=").Lit("?action=").Op("+").Id(actionNameConst)

		r.callFormatQueryUrl(def)
		IfErrReturn(def, errReturnParams...).Line()

		req := def.List(Id(ReqVar), Err()).Op(":=").Id(ClientReceiver)
		var params *Statement
		if hasParams {
			params = Id("params")
		} else {
			params = Struct().Block()
		}
		req.Dot("JsonPostRequest").Call(Id(UrlVar), RestLiMethod(protocol.Method_action), params)
		IfErrReturn(def, errReturnParams...).Line()

		if returns {
			def.Id(DoAndDecodeResult).Op(":=").Struct(Id("Value").Add(a.Return.GoType())).Block()
			callDoAndDecode(def)
			returnValue := Id(DoAndDecodeResult).Dot("Value")
			if !a.Return.IsMapOrArray() {
				returnValue = Op("&").Add(returnValue)
			}
			def.Return(returnValue, Nil())
		} else {
			def.List(Id("_"), Err()).Op("=").Id(ClientReceiver).Dot("DoAndIgnore").Call(Id(ReqVar))
			def.Return(Err())
		}
	})

	return c
}
