package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//page Estrutura da Resposta
type page struct {
	Title string
	Body  []byte
}

//Report Struct do Report
type Report struct {
	Node   string
	Testes []Teste
}

//Teste Struct do Teste
type Teste struct {
	Item                  string
	Desc                  string
	OcorrenciasAdicionais int //OcorrenciasAdicionais Adicionais
}

//Config struct da Configuracao
type Config struct {
	Dir    string `yaml:"dir"`
	Port   string `yaml:"port"`
	Ipaddr string `yaml:"ipaddr"`
}

var config Config

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadpage(title string) (*page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	check(err)
	return &page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadpage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reportHandler - ", config.Dir)
	node := r.URL.Path[len("/cisreport/"):]
	p := extractReportFromFile(node, config.Dir)
	//json.NewEncoder(w).Encode//(p) //write json to
	//extractReportFromFile(node,)

	w.Header().Set("Content-Type", "application/json")
	w.Write(p)
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

}

func loadConfig(filename string) {
	fmt.Println("Carregando Configuração...")
	file, _ := filepath.Abs(filename)

	fmt.Println("File: ", file)

	yamlFile, err := ioutil.ReadFile(file)
	//fmt.Println("loadconfig y: ", config.Dir)
	check(err)

	err = yaml.Unmarshal(yamlFile, &config)
	fmt.Println("loadconfig config dir: ", config.Dir)
	check(err)

}

func extractReportFromFile(k8snode string, dir string) []byte {

	var filename string

	if k8snode == "all" {
		filename = "cis-docker-all.txt"
	} else if strings.Contains(k8snode, "spessrvvpkn") {
		filename = dir + "/cis-docker-" + k8snode + ".estaleiro.serpro.txt"
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

	//Report.Node = k8snode

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

	jsonreport, err := json.Marshal(report)
	check(err)
	//os.Stdout.Write(b)
	return jsonreport

}

func main() {
	fmt.Println("Iniciando...")
	loadConfig("config.yml")
	//fmt.Println("config dir: ", config.Dir)

	//http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/cisreport/", reportHandler)
	http.ListenAndServe("0.0.0.0:6669", nil)
	
