package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	sourceFile := flag.String("source", "", "Path to the Kindle clippings file")
	outputFile := flag.String("output", "", "Path for generated Markdown file")
	templateFile := flag.String("template", "", "Path to the Markdown template file with Go template syntax")
	flag.Parse()

	if *sourceFile == "" || *outputFile == "" {
		fmt.Println("Please provide both source and output files.")
		os.Exit(1)
	}

	sourcePath, err := url.Parse(*sourceFile)
	if err != nil {
		fmt.Printf("Invalid source file URL: %s\n", err)
		os.Exit(1)
	}

	unescapedSourcePath, err := url.QueryUnescape(sourcePath.String())
	if err != nil {
		fmt.Println("Error unescaping URL:", err)
		os.Exit(1)
	}

	localSourcePath := strings.TrimPrefix(unescapedSourcePath, "file://")
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
		Heading   string
		Highlight string
		Text      string
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
	page.Title = strings.TrimSpace(doc.Find(".bookTitle").Text())
	page.Authors = strings.TrimSpace(doc.Find(".authors").Text())
	page.Citation = strings.TrimSpace(doc.Find(".citation").Text())
	doc.Find(".sectionHeading").Each(func(i int, s *goquery.Selection) {
		section := Section{}
		section.Title = strings.TrimSpace(s.Text())
		notesSelection := s.NextUntil(".sectionHeading")
		notesSelection.Each(func(i int, n *goquery.Selection) {
			// check only highlights
			if n.HasClass("noteHeading") || strings.HasPrefix(strings.TrimSpace(n.Prev().Text()), "Note -") {
				return
			}
			note := Note{}
			note.Highlight = strings.TrimSpace(n.Text())
			note.Heading = strings.TrimSpace(n.Prev().Text())
			nextHeading := strings.TrimSpace(n.Next().Text())
			if strings.HasPrefix(nextHeading, "Note -") {
				nextText := strings.TrimSpace(n.Next().Next().Text())
				note.Text = nextText
			}
			section.Notes = append(section.Notes, note)
		})
		page.Sections = append(page.Sections, section)
	})

	tmpl, err := template.New(*templateFile).ParseFiles(*templateFile)
	if err != nil {
		fmt.Printf("Error parsing template file: %s\n", err)
		os.Exit(1)
	}
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		os.Exit(1)
	}
	defer outFile.Close()
	err = tmpl.Execute(outFile, page)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Markdown notes file generated successfully.")
}
