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
	Node   string
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
	var linha string
	re := regexp.MustCompile(`\d{1,2}[\.]\d{1,2}`)
	//flag := false

	var report Report

	report.Node = k8snode

	for fscanner.Scan() {

		linetmp := fscanner.Text()
		//fmt.Println(line)

		if strings.Contains(linetmp, "[WARN]") {

			//removes special characters at beggining of line
			ind := strings.Index(linetmp, "0m")
			//fmt.Println("Index: ", ind+2)
			//fmt.Println(line[ind+2:])
			//Line whithout special formatting characters
			linha = linetmp[ind+2:]

			if id := re.FindString(linha); len(id) > 0 {
				nome := linha[strings.Index(linha, "-"):len(linha)]
				//fmt.Println(nome)

				m := Teste{
					Item: id,
					Desc: nome,
					OcorrenciasAdicionais: qtde,
				}

				b, err := json.Marshal(m)
				check(err)
				os.Stdout.Write(b)
				qtde = 0

			} else if strings.Index(linha, "*") > 0 {

			} else {
				panic(0)
			}

			if strings.Index(linha, "*") > 0 {
				// Ocorrencia quando o teste obtem múltiplos resultados
				// Ex: Item "5.1"
				// - Ensure AppArmor Profile is Enabled
				// OcorrenciasAdicionais:  78
				// 78 containers sem o profile AppArmor
				// fmt.Println("Tem *: ", linhax)

				qtde++
				//fmt.Println("Qtde de *: ", qtde)

			} else {
				//Item e Descricao
				id := re.FindString(linha)

				/* if qtde > 0 {
					fmt.Println("OcorrenciasAdicionais: ", qtde)
				}
				*/

				//fmt.Printf("%q\n", id)

				nome := linha[strings.Index(linha, "-"):len(linha)]
				//fmt.Println(nome)

				m := Teste{
					Item: id,
					Desc: nome,
					OcorrenciasAdicionais: qtde,
				}

				b, err := json.Marshal(m)

				check(err)

				//fmt.Println("JSON: ", b)

				os.Stdout.Write(b)
				//Item        string
				//Desc        string
				//OcorrenciasAdicionais int
				qtde = 0

			}

		}

	} //fim Scan()

}

func main() {
	//http.HandleFunc("/view/", viewHandler)
	//http.HandleFunc("/report/", reportHandler)
	//http.ListenAndServe(":8080", nil)

	extractReportFromFile("all")

}
