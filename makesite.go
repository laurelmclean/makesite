package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
}

func main() {
	// Define a command-line flag named "file" to specify the input filename.
	fileFlag := flag.String("file", "", "Input filename")
	flag.Parse()

	// Check if the filename flag is provided.
	if *fileFlag == "" {
		fmt.Println("Please provide a filename using --file flag")
		return
	}

	// Read the contents of the specified text file.
	fileContents, err := ioutil.ReadFile(*fileFlag)
    if err != nil {
        panic(err)
    }

	// Extract the filename without extension.
	fileName := filepath.Base(*fileFlag)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

	// Prepare the Page struct.
	page := Page{
		TextFilePath: *fileFlag,
		TextFileName: fileNameWithoutExt,
		HTMLPagePath: fileNameWithoutExt + ".html",
		Content:      string(fileContents),
	}

	// Create a new template in memory named "template.tmpl".
	// When the template is executed, it will parse template.tmpl,
	// looking for {{ }} where we can inject content.
	t := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	// Create a new HTML file with the same name as the input file.
	newFile, err := os.Create(page.HTMLPagePath)
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return
	}
	defer newFile.Close()

	// Executing the template injects the Page instance's data,
	// allowing us to render the content of our text file.
	// Furthermore, upon execution, the rendered template will be
	// saved inside the new file we created earlier.
	err = t.Execute(newFile, page)
    if err != nil {
        panic(err)
    }

	fmt.Printf("HTML file '%s' generated successfully!\n", page.HTMLPagePath)
}
