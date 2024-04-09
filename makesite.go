package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/color"
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

	// start time
	start := time.Now()

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
			fmt.Println("🔍 Text file found:", filepath.Join(*dirFlag, fileName))

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

			successMessage := fmt.Sprintf("HTML file '%s' generated successfully!\n", page.HTMLPagePath)
			color.Println(color.Blue(successMessage, color.In))
		}
	}

	kbSize := float64(totalFileSize) / 1024

	// Record the end time
	end := time.Now()

	// Calculate how long it ran
	duration := end.Sub(start)

	// Print formatted text with color and styling using github.com/labstack/gommon/color package
	color.Println(color.Green("Success!", color.B), "Generated", color.Bold(fileCount), "pages", color.Inverse(fmt.Sprintf("(%.1fKB total)", kbSize)))
	color.Println(color.Bold("Total execution time:"), color.Green((duration), color.In))
}
