package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	"net"
	"os"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"bytes"
)

var port = flag.Int("port", 3000, "tcp/ip port to listen on")
var workDir = flag.String("work-dir", ".", "Workspace directory containing templates and configs")

func main() {
	flag.Parse()
	logrus.Infof("ICTU CloudConfig Service started")
	logrus.Infof("Listening on port: %v", *port)
	os.Chdir(*workDir)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/cloud-init/{environment}", renderYaml)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

func renderYaml(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	template := vars["environment"] + ".yml"
	config := vars["environment"] + "-vars.yml"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	cnt, err := parseTemplate(template, config)
	if err != nil {
		logrus.Errorf("Failed to serve %v: %v", vars["environment"], err)
		return
	}
	cnt = strings.Replace(cnt, "REMOTE_IP", ip, -1)
	fmt.Fprintf(w, cnt)
	logrus.Infof("Served %s cloud-init to %s", vars["environment"], ip)
}

func parseTemplate(pathTemplate string, pathConfig string) (content string, err error) {
	logrus.Debugf("Parsing template %v", pathTemplate)
	config, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		return "", err
	}
	configMap := make(map[interface{}]interface{})
  err = yaml.Unmarshal([]byte(string(config)), &configMap)

	t, err := template.ParseFiles(pathTemplate)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	t.Execute(&doc, &configMap)
	s := doc.String()

	return s, nil
}
