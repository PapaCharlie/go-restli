# go-restli: Golang bindings for Rest.li

## How to
Check out this repo and install the binary it contains:
```bash
% git clone git@github.com/PapaCharlie/go-restli
% cd go-restli && go install
```
You can now use this too to generate Rest.li bindings from any snapshot file! Snapshot files are generated at build time
buy the Rest.li annotations processor. You can find the snapshot files in `src/mainGeneratedRest/snapshot`. The syntax
to invoke the tool and generate the code is as follows:

```bash
go-restli PACKAGE_PREFIX OUTPUT_DIR SNAPSHOT_FILES...
```
+ **PACKAGE_PREFIX**: All files will be generated inside of this namespace (e.g. `generated/`), and the generated code
  will need to be imported accordingly.
+ **OUTPUT_DIR**: The directory in which to output the files. Any necessary subdirectories will be created.
+ **SNAPSHOT_FILES**: All remaining parameters are the paths to the snapshot files.

## TODO
There are still many missing parts to this, including documentation and polish. I first focused on the biggest pain
point in working with Rest.li in golang, which is to generate the structs that are used to send and receive requests to
Rest.li endpoints. Most of the useful constants like resource paths and action names get extracted from the spec as
well, just to make it easier to write the code against net/http.Client, or whatever your favorite HTTP client framework
might be.

+ Finder support
+ Base library for generated code (i.e. all Collection, Simple and Association methods), to be imported by generated 
  code
+ Client code generation
+ Documentation and polish
+ Cli flags for more customization
