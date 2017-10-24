package main

import (
	"bufio"
	"encoding/json"
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

//Report Struct do Report
type Report struct {
	Node string

	Testes []Teste
}

//Teste Struct da Teste
type Teste struct {
	Item                  string
	Desc                  string
	OcorrenciasAdicionais int //OcorrenciasAdicionais Adicionais
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

	extractReportFromFile("all")

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

	var qtde = 0
	//var linha string
	re := regexp.MustCompile(`\d{1,2}[\.]\d{1,2}`)
	//flag := false

	var report Report

	report = Report{Node: k8snode}
	var item string
	var descricao string
	var linhas []string

	//report.Node = k8snode

	for fscanner.Scan() {

		linetmp := fscanner.Text()

		if strings.Contains(linetmp, "[WARN]") {
			ind := strings.Index(linetmp, "0m")
			linha := linetmp[ind+2:]
			linhas = append(linhas, linha)

		}
	}
	var indexItem int

	for index, linha := range linhas {

		item = re.FindString(linha)

		//fmt.Println("index: ", index, " item: ", item, "linha: ", linha)

		if len(item) > 0 && strings.Index(linha, "*") < 0 { //WARN com item e descricao
			descricao = linha[strings.Index(linha, "-"):len(linha)]
			//fmt.Println("Desc: ", descricao)
			qtde = 0
			report.Testes = append(report.Testes, Teste{Item: item, Desc: descricao})
			indexItem = len(report.Testes) - 1
			//fmt.Println("indexItem: ", indexItem)

		} else if strings.Index(linha, "*") > 0 { //ocorrencia adicional do item
			qtde++

			if index == len(linhas)-1 || strings.Index(linhas[index+1], "*") < 0 { //se a proxima linha nao contem *
				report.Testes[indexItem].OcorrenciasAdicionais = qtde
			}

		} else { //Linha mal formada
			panic("fudeu")

		}

	}

	b, err := json.Marshal(report)
	check(err)
	os.Stdout.Write(b)

}

func main() {
	//http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/report/", reportHandler)
	http.ListenAndServe(":8080", nil)

}
