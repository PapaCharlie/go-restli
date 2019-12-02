# go-restli: Golang bindings for Rest.li

## How to
Check out this repo and install the binary it contains:
```bash
% git clone --recurse-submodules git@github.com/PapaCharlie/go-restli
% cd go-restli && go install
```
You can now use this tool to generate Rest.li bindings for any given resource. You will need to acquire all the PDSC
models that your resources depend on, as well as their restspecs. Once you do, you can run the tool as follows:

```bash
go-restli \
  --package-prefix github.com/PapaCharlie/go-restli/tests/generated \
  --output-dir ./tests/generated \
  --pdsc-dir ./pegasus \
  idl/*.restspec.json
```
+ **--package-prefix**: All files will be generated inside of this namespace (e.g. `generated/`), and the generated code
  will need to be imported accordingly.
+ **--output-dir**: The directory in which to output the files. Any necessary subdirectories will be created.
+ **--pdsc-dir**: The directory that contains all the .pdsc files that may be used by the resources we want to generate
  code for.
+ All remaining parameters are the paths to the restspec files for the resources we want to call.

## Getting the PDSCs and Restpecs
You may wish to use gradle to extract the PDSC and restspec.json files from the incoming jars. To do so, you can use a
task like this:
```gradle
task extractPdscsAndRestpecs >> {
  copy {
    project.configurations.restliSpecs.each {
      from zipTree(it)
      include "idl/"
      include "pegasus/"
    }

    into temporaryDir
  }
}
```

## Conflict resolution in cyclic packages
Java allows cyclic package imports since multiple modules can define classes for the same packages. Similarly, it's
entirely possible for PDSC models to introduce package cycles. To mitigate this, the code generator will attempt to
resolve dependency chains that introduce package cycles and move the offending models to a fixed package called
`conflictResolution`.

## TODO
There are still many missing parts to this, including documentation and polish. I first focused on the biggest pain
point in working with Rest.li in golang, which is to generate the structs that are used to send and receive requests to
Rest.li endpoints. Most of the useful constants like resource paths and action names get extracted from the spec as
well, just to make it easier to write the code against net/http.Client, or whatever your favorite HTTP client framework
might be.

## Contributing to this project
First, you have to clone this repo and all its submodules:
```bash
% git clone --recurse-submodules git@github.com/PapaCharlie/go-restli
```
There exists a testing framework for Rest.li client implementations that provide expected requests and responses. The
[gorestli_test.go](tests/gorestli_test.go) and [manifest.go](tests/manifest.go) files read the testing
[manifest](tests/rest.li-test-suite/client-testsuite/manifest.json) and load all the corresponding requests and
responses. Once a new resource type or method is implemented, please be sure to integrate all the tests for that new
feature. This is done by adding a new function on the `TestServer` struct for the corresponding name. For example, to
test the `collection-get` test (see the corresponding object in the manifest), all that's needed is to add a
correspondingly named method called `CollectionGet`, like this:
```go
func (s *TestServer) CollectionGet(t *testing.T, c *Client) {
	id := int64(1)
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.Message{Id: &id, Message: "test message"}, res, "Invalid response from server")
}
```
