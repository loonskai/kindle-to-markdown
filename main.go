package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// read arguments:
	// --source <path to kindle clippings file>
	// --output <path to markdown file>
	// --template <path to markdown template file with Go template syntax>

	sourceFile := flag.String("source", "", "Path to the Kindle clippings file")
	outputFile := flag.String("output", "", "Path for generated Markdown file")
	// templateFile := flag.String("template", "", "Path to the Markdown template file with Go template syntax")
	flag.Parse()

	// check if source and output files are provided
	if *sourceFile == "" || *outputFile == "" {
		fmt.Println("Please provide both source and output files.")
		os.Exit(1)
	}

	sourceURL, err := url.Parse(*sourceFile)
	if err != nil {
		fmt.Printf("Invalid source file URL: %s\n", err)
		os.Exit(1)
	}

	unescapedSourceURL, err := url.QueryUnescape(sourceURL.String())
	if err != nil {
		fmt.Println("Error unescaping URL:", err)
	} else {
		fmt.Println("Unescaped URL:", unescapedSourceURL)
	}

	if sourceURL.Scheme == "file" {
		localSourcePath := strings.TrimPrefix(unescapedSourceURL, "file://")
		file, err := os.Open(localSourcePath)
		if err != nil {
			fmt.Printf("Error reading kindle clippings file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			fmt.Printf("Error parsing kindle clippings file: %s\n", err)
			os.Exit(1)
		}

		type Note struct {
			Heading string
			Text    string
		}

		type Section struct {
			Title string
			Notes []Note
		}

		type SourcePage struct {
			Title    string
			Authors  string
			Citation string
			Sections []Section
		}

		page := SourcePage{}
		page.Title = doc.Find(".bookTitle").Text()
		page.Authors = doc.Find(".authors").Text()
		page.Citation = doc.Find(".citation").Text()
		doc.Find(".sectionHeading").Each(func(i int, s *goquery.Selection) {
			section := Section{}
			section.Title = s.Text()
			notesSelection := s.NextUntil(".sectionHeading")
			notesSelection.Each(func(i int, n *goquery.Selection) {
				if n.HasClass("noteText") {
					return
				}
				note := Note{}
				note.Heading = n.Text()
				if n.Next().HasClass("noteText") {
					note.Text = n.Next().Text()
				}
				section.Notes = append(section.Notes, note)
			})
			page.Sections = append(page.Sections, section)
		})
		fmt.Printf("%+v\n", page)

	} else {
		fmt.Println("Only local files are supported for now.")
		os.Exit(1)

	}

	// parse source HTML file
	// convert to markdown
	// write to markdown file
	// read kindle clippings file
	// convert to markdown
	// write to markdown file
}

func parseArgs() {

}
