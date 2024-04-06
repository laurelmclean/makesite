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

	// Initialize a counter for the number of files
	fileCount := 0

		// Initialize a counter for the total size of files
	var totalFileSize int64 = 0


	for _, file := range files {

		fileName := file.Name()

		if filepath.Ext(fileName) == ".txt" {

			// Print the .txt files found
			fmt.Println(filepath.Join(*dirFlag, fileName)) 

			fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

			// Read the contents of the text file
			fileContents, err := ioutil.ReadFile(filepath.Join(*dirFlag, fileName))
			if err != nil {
				panic(err)
			}


			// Prepare the Page struct for each .txt file found.
			page := Page{
				TextFilePath: filepath.Join(*dirFlag, fileName),
				TextFileName: fileName,
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

			fileCount++

			totalFileSize += file.Size()

			fmt.Printf("HTML file '%s' generated successfully!\n", page.HTMLPagePath)
		}
	}
	fmt.Fprintf(os.Stdout, "\033[0;32m \033[1m%s\033[0m %s \033[1m%d\033[0m %s (%.2fKB total).\n", "Success!", "Generated", fileCount, "pages", float64(totalFileSize)/1024)
}