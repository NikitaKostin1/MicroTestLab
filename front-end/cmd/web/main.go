package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)



const serverPort = "1025"



func main() {
	log.Printf("Front-end service startup on web port %s", serverPort)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		renderTemplate(writer, "test.page.gohtml")
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func renderTemplate(writer http.ResponseWriter, templateFileName string) {
	templatePartials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	templatePaths := append([]string{fmt.Sprintf("./cmd/web/templates/%s", templateFileName)}, templatePartials...)

	parsedTemplates, err := template.ParseFiles(templatePaths...)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	if err := parsedTemplates.Execute(writer, nil); err != nil {
		http.Error(writer, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}
