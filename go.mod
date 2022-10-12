module github.com/PapaCharlie/go-restli

require (
	github.com/dave/jennifer v1.3.0
	github.com/go-zookeeper/zk v1.0.3
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/mailru/easyjson v0.7.2
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

go 1.18

replace github.com/go-zookeeper/zk v1.0.3 => github.com/PapaCharlie/zk v1.0.4-0.20221012222955-44b3748be649
