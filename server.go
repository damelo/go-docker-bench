package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Page Estrutura da Resposta
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/report/"):]
	p, _ := extractReportFromFile(title)
	json.NewEncoder(w).Encode(p) //write json to
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

}

func execAnsible() (*Page, error) {

	//spessrvvpkn00001.estaleiro.serpro

	return nil, nil

}

func extractReportFromFile(k8snode string) (*Page, error) {

	var filename string

	if k8snode == "all" {
		filename = "arquivo.txt"
	} else {
		filename = "arquivo.txt"
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: filename, Body: body}, nil

}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/report/", reportHandler)
	http.ListenAndServe(":8080", nil)

}
