# go-restli: Golang bindings for Rest.li

## How to use a `restli.Client`
```go
restLiClient := &Client{
	// This uses a standard http.Client. Most clients will need to be configured with the right timeouts, TLS
	// contexts, CA certs etc. The following is a recommended configuration that aggressively times out connections
	// and TLS handshakes but not the actual request time, allowing servers to block as long as they want (this is
	// common for actions).
	Client: &http.Client{
		Transport: &http.Transport{
			// This times out connecting to the server
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
			// This times out the TLS handshake
			TLSHandshakeTimeout: 10 * time.Second,
		},
		// Deliberately not set to ensure the client does not have an overall timeout
		Timeout: 0,
	},

	// The HostnameResolver implements locating the host a request should go to. This example always returns the
	// same hostname, but can be implemented in any way shape or form.
	HostnameResolver: &SimpleHostnameResolver{Hostname: hostname},
	// For example, a D2 client can be used to look up hostnames dynamically
	HostnameResolver: d2.Client{Conn: zkConn},

	// This flag toggles returning errors on incomplete/illegal responses from servers. This happens if for whatever
	// reason an endpoint implementation forgot to set a field, or if a field is removed as the API evolves. Because
	// of how often this happens, strict deserialization is disabled by default.
	StrictResponseDeserialization: false,
}

// Now that we have a restli.Client, we can use it to call some resources. Every resource defines a NewClient method
// that takes in a restli.Client
collectionClient := collection.NewClient(restLiClient)
// The returned client wraps the restli client and exposes that resource's methods. All network errors will be
// url.Errors, otherwise they will be corresponding error types declared in the restli package.
msg, err := collectionClient.Get(123)
...

// restli clients can be reused across multiple resources
actionsetClient := actionset.NewClient(restLiClient)
```

## How to use a `restli.Server`
Each resource will generate a `Resource` interface that needs to be implemented. Suppose the following generated
interface:
```go
type Resource interface {
	Create(ctx *restli.RequestContext, entity *Message) (createdEntity *CreatedEntity, err error)
	Get(ctx *restli.RequestContext, collectionId int64) (entity *Message, err error)
	Update(ctx *restli.RequestContext, collectionId int64, entity *Message) (err error)
	Delete(ctx *restli.RequestContext, collectionId int64) (err error)
}
```

Let's implement it in the simplest way possible:
```go
type impl map[int64]*Message

func (i impl) Create(ctx *restli.RequestContext, entity *Message) (createdEntity *CreatedEntity, err error) {
	id := rand.Int63()
	i[id] = entity
	// If the Status field of CreatedEntity is not set, it will default to http.StatusCreated
	createdEntity = &CreatedEntity{Id: id}
	// The SetLocation method is used to set (optionally) the Location field, and returns an error if the id cannot be
	// serialized
	_ = restli.SetLocation(ctx, createdEntity)
	return createdEntity, nil
}

func (i impl) Get(ctx *restli.RequestContext, collectionId int64) (entity *Message, err error) {
	if m, ok := i[collectionId]; ok {
		return m, nil
	} else {
		// Returning a *restlidata.ErrorResponse will automatically set the response code and format a restli service
		// error.
		return nil, &restlidata.ErrorResponse{
			Status:  restli.Int32Pointer(http.StatusNotFound),
			Message: restli.StringPointerf("No such message %q", collectionId),
		}
	}
}

func (i impl) Update(ctx *restli.RequestContext, collectionId int64, entity *Message) (err error) {
	// The original request is accessible via the request context. For example, to read any extra headers that may have
	// been added or to get the client certificates to check if the client has access to specific methods
	if err = checkUserCanUpdate(ctx.Request.TLS.PeerCertificates); err != nil {
		return err
	}

	if _, ok := i[collectionId]; ok {
		i[collectionId] = entity
		return nil
	} else {
		return &restlidata.ErrorResponse{
			Status:  restli.Int32Pointer(http.StatusNotFound),
			Message: restli.StringPointerf("No such message %q", collectionId),
		}
	}
}

func (i impl) Delete(ctx *restli.RequestContext, collectionId int64) (err error) {
	// If ctx.ResponseStatus isn't set, Delete will default to http.StatusNoContent. In this case we want to return
	// another status if the key didn't actually exist, but without failing
	if _, ok := i[collectionId]; !ok {
		ctx.ResponseStatus = http.StatusOK
	}
	delete(i, collectionId)
	return nil
}
```
And finally we can go ahead and register this resource against a `restli.Server`. Alongside the generated `Resource`
interface, a `RegisterResource` is also generated. The usage pattern is as follows:
```go
// When a prefix is specified, the Server will ignore any requests whose path do not start with said prefix. Otherwise
// NewServer defaults to "/"
server := restli.NewPrefixedServer("/api/v1", /* Filters can be added here */)
RegisterResource(server, make(impl))

...

// Once all the resources have been registered, the server can be added to a normal *http.ServeMux:
mux := http.NewServeMux()
server.AddToMux(mux)

// Or if no other endpoints need to be added, the Server can be used to directly serve requests
http.ListenAndServe("localhost:8080", server.Handler())
```

## How to generate bindings
Grab a binary from the latest [release](https://github.com/PapaCharlie/go-restli/releases) for your platform and put it
on your path. You can now use this tool to generate Rest.li bindings for any given resource. You will need to acquire
all the PDSC/PDL models that your resources depend on, as well as their restspecs. Once you do, you can run the tool as
follows:

```bash
go-restli \
	--output-dir internal/tests/generated \
	--resolver-path internal/tests/rest.li-test-suite/client-testsuite/schemas \
	--package-prefix github.com/PapaCharlie/go-restli/internal/tests/generated \
	--named-schemas-to-generate testsuite.Primitives \
	--named-schemas-to-generate testsuite.ComplexTypes \
	internal/tests/rest.li-test-suite/client-testsuite/restspecs/*
```

+ **-p/--package-prefix**: All files will be generated inside of this namespace (e.g. `generated/`), and the generated
  code will need to be imported accordingly.
+ **-o/--output-dir**: The directory in which to output the files. Any necessary subdirectories will be created.
+ **-r/--resolver-path**: The directory that contains all the `.pdsc` and `.pdl` files that may be used by the resources
  you want to generate code for.
+ **-n/--named-schemas-to-generate**: Generate bindings for these named schemas alongside the schemas required to call
  the given resources. Note that it's not required to specify `.restspec.json` files if at least one named schema is
  specified. Useful when rest.li schemas need to be used without calling rest.li resources. In other words, bindings can
  be generated without the need for `.restspec.json` files!
+ **--raw-records**: Any record listed in this flag will be replaced with a [`restlidata.RawRecord`](
  restlidata/RawRecord.go). Some rest.li records are deliberately untyped and are used to send data without a schema.
  A `RawRecord` is used to capture that by deserializing the entire object as an interface, which can then be read back
  into a normal record using `UnmarshalTo`.
+ All remaining parameters are the paths to the restspec files for the resources you want to call. At least one restspec
  or one `--named-schemas-to-generate` must be provided.

#### Getting the PDSCs and Restpecs
You may wish to use gradle to extract the schema and restspec from the incoming jars. To do so, you can use a task like
this:

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
entirely possible for schemas to introduce package cycles. To mitigate this, the code generator will attempt to resolve
dependency chains that introduce package cycles and move the offending models to a fixed package called
`conflictResolution`.

## Typeref support
The Java code generator simply drops typerefs and directly references the underlying type, unless a coercer is
registered. The plan here would be to support a notion of type coercion that can be plugged into the generator. This way
types with native bindings can be deserialized from their raw values (e.g. a UUID type serialized as a string that can
be deserialized into an actual `github.com/google/uuid`). For the time being, typerefs will altogether not be supported,
and the raw underlying type will be used instead.

## Note on Java dependency
The owners of Rest.li recommended against implementing a custom PDSC/PDL/RESTSPEC parser and instead recommend using the
existing Java code to parse everything. This is not only because the .pdsc format is going to be replaced by a new DSL
called PDL, which will be much harder to parse than JSON (incidentally, the .pdsc format allows comments and other
nonsense, which makes it not standard JSON either). Therefore this code actually uses Java to parse everything, then
outputs a simpler intermediary JSON file where every schema and spec is fully resolved, making the code generation step
significantly less complicated.

In order to parse the schemas and restspecs, the binaries have an embedded jar. They will unpack the jar and attempt to
execute it with `java -jar`. This jar has no dependencies, but you _must_ have a JRE installed. Please make sure that
`java` is on your PATH (though setting a correct `JAVA_HOME` isn't 100% necessary). This has been tested with Java 1.8.

## Contributing to this project
First, you have to clone this repo and all its submodules:

```bash
% git clone --recurse-submodules git@github.com:PapaCharlie/go-restli
```

There exists a testing framework for Rest.li client implementations that provide expected requests and responses. The
[gorestli_test.go](tests/gorestli_test.go) and [manifest.go](tests/manifest.go) files read the testing
[manifest](tests/rest.li-test-suite/client-testsuite/manifest.json) and load all the corresponding requests and
responses. Once a new resource type or method is implemented, please be sure to integrate all the tests for that new
feature. This is done by adding a new function on the `Operation` struct for the corresponding name. For example, to
test the `collection-get` test (see the corresponding object in the manifest), all that's needed is to add a
correspondingly named method called `CollectionGet`, like this:

```go
func (o *Operation) CollectionGet(t *testing.T, c Client) func(*testing.T) *MockResource {
	id := int64(1)
	expected := newMessage(id, "test message")
	res, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, expected, res)

	return func(t *testing.T) *MockResource {
		return &MockResource{
			MockGet: func(ctx *restli.RequestContext, collectionId int64) (entity *conflictresolution.Message, err error) {
				require.Equal(t, id, collectionId)
				return expected, nil
			},
		}
	}
}
```
Once you have written your tests, just run `make` in the root directory and all the tests will be run.

## Building from source
This project uses ``make`` as the build tool and requires following dependencies:

* Golang
* Java 1.8+
* ``goimports``
* ``stringer``
