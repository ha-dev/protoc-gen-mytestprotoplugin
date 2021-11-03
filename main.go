package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Protoc passes pluginpb.CodeGeneratorRequest in via stdin
	// marshalled with Protobuf
	input, _ := ioutil.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	proto.Unmarshal(input, &req)

	// Initialise our plugin with default options
	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}

	// Protoc passes a slice of File structs for us to process
	for _, file := range plugin.Files {

		// Time to generate code...!

		// 1. Initialise a buffer to hold the generated code
		var buf bytes.Buffer

		// 2. Write the package name
		pkg := fmt.Sprintf("package %s", file.GoPackageName)
		buf.Write([]byte(pkg))

		// 3. For each message add our Foo() method
		for _, msg := range file.Proto.MessageType {
			buf.Write([]byte(fmt.Sprintf(`
			 func (x %s) Myfunc() string {
				return "Ok Good"
			 }`, *msg.Name)))
		}

		// 4. Specify the output filename, in this case test.foo.go
		filename := file.GeneratedFilenamePrefix + ".mytestprotoplugin.go"
		file := plugin.NewGeneratedFile(filename, ".")

		// 5. Pass the data from our buffer to the plugin file struct
		file.Write(buf.Bytes())
	}

	// Generate a response from our plugin and marshall as protobuf
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	// Write the response to stdout, to be picked up by protoc
	fmt.Fprintf(os.Stdout, string(out))
}
