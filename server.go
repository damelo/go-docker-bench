package cistest

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

//page Estrutura da Resposta
type page struct {
	Title string
	Body  []byte
}

//report Struct do report
type report struct {
	Node string

	testes []teste
}

//teste Struct da teste
type teste struct {
	Item                  string
	Desc                  string
	OcorrenciasAdicionais int //OcorrenciasAdicionais Adicionais
}

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
	//title := r.URL.Path[len("/report/"):]
	//p, _ := extractreportFromFile(title)
	//json.NewEncoder(w).Encode(p) //write json to

	extractreportFromFile("all")

	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

}

func extractreportFromFile(k8snode string) {

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

	var report report

	report = report{Node: k8snode}
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
			report.testes = append(report.testes, teste{Item: item, Desc: descricao})
			indexItem = len(report.testes) - 1
			//fmt.Println("indexItem: ", indexItem)

		} else if strings.Index(linha, "*") > 0 { //ocorrencia adicional do item
			qtde++

			if index == len(linhas)-1 || strings.Index(linhas[index+1], "*") < 0 { //se a proxima linha nao contem *
				report.testes[indexItem].OcorrenciasAdicionais = qtde
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
	http.ListenAndServe(":6669", nil)

}
