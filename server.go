package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//Page Estrutura da Resposta
type Page struct {
	Title string
	Body  []byte
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	check(err)
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/report/"):]
	//p, _ := extractReportFromFile(title)
	//json.NewEncoder(w).Encode(p) //write json to
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

}

func extractReportFromFile(k8snode string) {

	var filename string

	if k8snode == "all" {
		filename = "cis-docker-all.txt"
	} else {
		filename = "cis-docker-all.txt"
	}

	fileHandle, err := os.Open(filename)

	check(err)

	fscanner := bufio.NewScanner(fileHandle)

	for fscanner.Scan() {

		line := fscanner.Text()

		fmt.Println(line)
	}

}

func main() {
	//http.HandleFunc("/view/", viewHandler)
	//http.HandleFunc("/report/", reportHandler)
	//http.ListenAndServe(":8080", nil)

	extractReportFromFile("all")

}
