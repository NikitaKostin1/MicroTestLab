package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)



const serverPort = "1025"


func main() {
	// Set up route for the root path
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		renderTemplate(writer, "test.page.gohtml")
	})

	// Log the startup message
	log.Printf("Starting front-end service on port %s", serverPort)

	// Start the HTTP server
	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil)
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}

// renderTemplate loads and executes the specified template along with its partials
func renderTemplate(writer http.ResponseWriter, templateFileName string) {
	// Define the list of template partials
	templatePartials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	// Add the main template to the list of templates to parse
	templatePaths := append([]string{fmt.Sprintf("./cmd/web/templates/%s", templateFileName)}, templatePartials...)

	// Parse the templates
	parsedTemplates, err := template.ParseFiles(templatePaths...)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	// Execute the parsed templates
	if err := parsedTemplates.Execute(writer, nil); err != nil {
		http.Error(writer, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}
