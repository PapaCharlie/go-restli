package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/dave/jennifer/jen"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("Usage: %s PROTOCOL_SOURCE_DIR OUTPUT_ZIPPED_PROTOCOL_FILE", os.Args[0])
	}

	if len(os.Args) < 2 {
		log.Fatalf("Must specify the protocol source dir")
	}
	protocolDirectory := os.Args[1]

	if len(os.Args) < 3 {
		log.Fatalf("Must specify the output file")
	}
	outputFile := os.Args[2]

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := filepath.Walk(protocolDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		writer, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
		fmt.Printf("Adding %s to %s\n", path, outputFile)
		return nil
	})

	if err != nil {
		panic(err)
	}

	err = zipWriter.Close()
	if err != nil {
		panic(err)
	}

	f := jen.NewFile("codegen")
	f.HeaderComment(`THIS FILE WAS AUTOMATICALLY GENERATED
It contains a b64 encoded zip file of the "protocol" directory in the root of this project
DO NOT EDIT BY HAND, USE go generate TO REFRESH THIS FILE`)

	f.Var().Id("ProtocolZip").Index().Byte()
	f.Func().Id("init").Params().BlockFunc(func (g *jen.Group) {
		k := jen.Empty()
		b64Zip := base64.StdEncoding.EncodeToString(buf.Bytes())
		for len(b64Zip) > 80 {
			k.Line().Lit(b64Zip[:80]).Op("+")
			b64Zip = b64Zip[80:]
		}
		k.Line().Lit(b64Zip)
		g.Var().Err().Error()
		g.List(jen.Id("ProtocolZip"), jen.Err()).Op("=").Qual("encoding/base64", "StdEncoding").Dot("DecodeString").Call(k)
		g.If(jen.Err().Op("!=").Nil()).Block(jen.Panic(jen.Err()))
	})

	out, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = f.Render(out)
	if err != nil {
		panic(err)
	}
}
