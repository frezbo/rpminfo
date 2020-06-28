package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"os"

	xmlparser "github.com/tamerh/xml-stream-parser"
)

// Packages hold the json schema for the output data.json
type Packages struct {
	Name    string `json:"name"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
	Summary string `json:"summary"`
}

// Metadata holds the rpm packages metadata
type Metadata struct {
	XMLName  xml.Name  `xml:"metadata"`
	Packages []Package `xml:"package"`
}

// Package holds a single package metadata
type Package struct {
	Name        string `xml:"name"`
	Arch        string `xml:"arch"`
	VersionData VersionData
	Summary     string `xml:"summary"`
}

// VersionData holds the package version metadata
type VersionData struct {
	XMLName xml.Name `xml:"version"`
	Version string   `xml:"ver,attr"`
}

// PackageListV1 uses the go standard xml decoder
// Deprecated, kept for pedagogic purposes
// func PackageListV1(r io.Reader) (err error) {
// 	decoder := xml.NewDecoder(r)

// 	os.Remove("data.json")
// 	file, err := openFileStream("data.json")
// 	if err != nil {
// 		return err
// 	}
// 	defer closeFileStream(file)

// 	appendToFile(file, []byte("[\n  "))

// 	for {
// 		token, _ := decoder.Token()

// 		if token == nil || token == io.EOF {
// 			break
// 		}

// 		var pkg Package

// 		switch se := token.(type) {
// 		case xml.StartElement:
// 			if se.Name.Local == "package" {
// 				decoder.DecodeElement(&pkg, &se)
// 				jsonData := Packages{}
// 				jsonData.Name = pkg.Name
// 				jsonData.Arch = pkg.Arch
// 				jsonData.Version = pkg.VersionData.Version
// 				jsonData.Summary = pkg.Summary
// 				buffer := new(bytes.Buffer)
// 				encoder := json.NewEncoder(buffer)
// 				encoder.SetIndent("  ", "  ")
// 				encoder.SetEscapeHTML(false)
// 				err := encoder.Encode(&jsonData)
// 				if err != nil {
// 					return err
// 				}
// 				data := buffer.Bytes()
// 				data = data[:len(data)-1]
// 				appendToFile(file, data)
// 				appendToFile(file, []byte(",\n  "))
// 			}
// 		}

// 	}
// 	stat, err := file.Stat()
// 	if err != nil {
// 		return err
// 	}
// 	file.Truncate(stat.Size() - 4)
// 	appendToFile(file, []byte("\n]"))
// 	closeFileStream(file)
// 	return nil
// }

// PackageListV2 uses the xml-stream-parser library which seems more faster
func PackageListV2(r io.Reader) (err error) {
	br := bufio.NewReaderSize(r, 65536)

	parser := xmlparser.NewXMLParser(br, "package").SkipElements([]string{"checksum", "description", "packager", "url", "time", "size", "location", "format"})

	os.Remove("data.json")
	file, err := openFileStream("data.json")
	if err != nil {
		return err
	}
	defer closeFileStream(file)

	appendToFile(file, []byte("[\n  "))

	for xml := range parser.Stream() {
		if xml.Err != nil {
			return err
		}
		if xml.Name == "package" {
			jsonData := Packages{}
			jsonData.Name = xml.Childs["name"][0].InnerText
			jsonData.Arch = xml.Childs["arch"][0].InnerText
			jsonData.Version = xml.Childs["version"][0].Attrs["ver"]
			jsonData.Summary = xml.Childs["summary"][0].InnerText
			buffer := new(bytes.Buffer)
			encoder := json.NewEncoder(buffer)
			encoder.SetIndent("  ", "  ")
			encoder.SetEscapeHTML(false)
			err := encoder.Encode(&jsonData)
			if err != nil {
				return err
			}
			data := buffer.Bytes()
			data = data[:len(data)-1]
			appendToFile(file, data)
			appendToFile(file, []byte(",\n  "))
		}
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	file.Truncate(stat.Size() - 4)
	appendToFile(file, []byte("\n]"))
	return nil
}

func openFileStream(filepath string) (file *os.File, err error) {
	file, err = os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, err
}

func closeFileStream(file *os.File) (err error) {
	log.Println("Finished writing data.json")
	return file.Close()
}

func appendToFile(file *os.File, bytes []byte) (err error) {
	if _, err = file.Write(bytes); err != nil {
		return err
	}
	return nil
}
