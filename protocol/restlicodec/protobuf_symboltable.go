package restlicodec

type SymbolTable interface {
	GetSymbolId(symbolName string) (int, bool)
	GetSymbolName(symbolId int) (string, bool)
}

type genericSymbolTable struct {
	byName map[string]int
	byID   map[int]string
	size   int
}

// TODO this probably isn't right
func NewSymbolTable(symbols []string) SymbolTable {
	table := &genericSymbolTable{
		byName: make(map[string]int),
		byID:   make(map[int]string),
	}
	for _, v := range symbols {
		table.byID[table.size] = v
		table.byName[v] = table.size
		table.size++
	}
	return table
}

func (t *genericSymbolTable) GetSymbolId(symbolName string) (int, bool) {
	i, f := t.byName[symbolName]
	return i, f
}
func (t *genericSymbolTable) GetSymbolName(symbolId int) (string, bool) {
	s, f := t.byID[symbolId]
	return s, f
}
