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

	// Define flag "dir" to specify the input directory
	dirFlag := flag.String("dir", "", "Input directory")
	flag.Parse()

	// List all .txt files in directory
	files, err := ioutil.ReadDir(*dirFlag)
	if err != nil {
		panic(err)
	}
	

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {

			// Print the .txt files found
			fmt.Println(filepath.Join(*dirFlag, file.Name())) 

			// Read the contents of the text file
			fileContents, err := ioutil.ReadFile(filepath.Join(*dirFlag, file.Name()))
			if err != nil {
				panic(err)
			}


			// Prepare the Page struct for each .txt file found.
			page := Page{
				TextFilePath: filepath.Join(*dirFlag, file.Name()),
				TextFileName: file.Name(),
				HTMLPagePath: file.Name() + ".html",
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
	}
}