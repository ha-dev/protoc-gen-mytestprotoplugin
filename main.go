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
    
	input, _ := ioutil.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	proto.Unmarshal(input, &req)

	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}

	for _, file := range plugin.Files {

		var buf bytes.Buffer

		pkg := fmt.Sprintf("package %s", file.GoPackageName)
		buf.Write([]byte(pkg))

		for _, msg := range file.Proto.MessageType {
			buf.Write([]byte(fmt.Sprintf(`
			 func (x %s) Myfunc() string {
				return "Ok Good"
			 }`, *msg.Name)))
		}

		filename := file.GeneratedFilenamePrefix + ".mytestprotoplugin.go"
		file := plugin.NewGeneratedFile(filename, ".")

		file.Write(buf.Bytes())
	}

	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, string(out))
}
