package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Adventure map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handler struct {
	adventure Adventure
	template  *template.Template
}

var fileName string
var port int
var tpl *template.Template

func init() {
	flag.StringVar(&fileName, "fileName", "gopher.json", "a json file with story content.")
	flag.IntVar(&port, "port", 3030, "server port")
	tpl = template.Must(template.ParseFiles("views/index.html"))

}

func main() {

	chapterMap, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Unable to read file: %s", err)
	}

	adventure := parseAdventure(chapterMap)

	h := newHandler(adventure, tpl)
	fmt.Printf("Started server at Port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))

}

func parseAdventure(chapterMap []byte) (adventure Adventure) {
	if err := json.Unmarshal(chapterMap, &adventure); err != nil {
		log.Fatalf("Error parsing chapter map: %s", err)
	}
	return
}

func newHandler(adventure Adventure, template *template.Template) http.Handler {
	if template == nil {
		template = tpl
	}
	return handler{adventure, template}
}

// method: func(t ReceiverType) FuncName (Para...) ReturnTypes
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}
	path = path[1:]
	if chapter, ok := h.adventure[path]; ok {
		err := h.template.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "chapter not found.", http.StatusNotFound)

}

/*
func getChapter(adventure Adventure) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get url variable
		vars := mux.Vars(r)
		//since url set as "/cyoa/{chapter}/", this will get {chapter} part
		chapterKey := vars["chapter"]
		chapter, exist := adventure[chapterKey]
		if !exist {
			http.Error(w, "Chapter not found", http.StatusNotFound)
			return
		}
		tmpl := template.Must(template.ParseFiles("views/index.html"))
		if err := tmpl.Execute(w, chapter); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
*/
