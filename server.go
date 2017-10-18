package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
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

	qtde := 0

	for fscanner.Scan() {

		line := fscanner.Text()

		//fmt.Println(line)
		var linhax string

		re := regexp.MustCompile(`\d{1,2}[\.]\d{1,2}`)

		if strings.Contains(line, "[WARN]") {

			ind := strings.Index(line, "0m")
			//fmt.Println("Index: ", ind+2)
			fmt.Println(line[ind+2:])

			linhax = line[ind+2:]

			if strings.Index(linhax, "*") > 0 {
				//fmt.Println("Tem *: ", linhax)
				qtde++
				//fmt.Println("Qtde de *: ", qtde)

			} else {
				//fmt.Println("Não Tem *: ", linhax)
				qtde = 0
				fmt.Printf("%q\n", re.FindString(line[ind+2:]))
			}
			if qtde > 0 {
				fmt.Println("Ocorrencias: ", qtde)
			}

		}

	}

}

func main() {
	//http.HandleFunc("/view/", viewHandler)
	//http.HandleFunc("/report/", reportHandler)
	//http.ListenAndServe(":8080", nil)

	extractReportFromFile("all")

}
