package schema

type Schema struct {
	Name       string
	Namespace  string
	Path       string
	Schema     string
	Doc        string
	Collection Collection
	ActionSet  struct {
		Actions []Action
	} `json:"actionSet"`
}

type Collection struct {
	Identifier struct {
		Name string
		Type string
	}
	Supports []string
	Methods  []Method
	Actions  []Action
	Finders  []Finder
	Entity   Entity
}

type Entity struct {
	Path    string
	Actions []Action
}

type Method struct {
	Method string
	Doc    string
}

type Endpoint struct {
	Name       string
	Doc        string
	Parameters []Parameter
	Returns    string
}

type Finder struct {
	Endpoint
	PagingSupported bool `json:"pagingSupported"`
}

type Parameter struct {
	Name     string
	Doc      string
	Type     string
	Optional bool
	Default  string
}

type Action struct {
	Endpoint
}
