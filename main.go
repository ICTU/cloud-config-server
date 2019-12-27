package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var port = flag.Int("port", 3000, "tcp/ip port to listen on")
var nfsIPPrefix = flag.String("nfs-ip-prefix", "192.168", "First two octets of nfs client ip")
var workDir = flag.String("work-dir", ".", "Workspace directory containing templates and configs")

func main() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
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
	octet := strings.Split(ip, ".")
	nfsIP := fmt.Sprintf("%s.%s.%s", *nfsIPPrefix, octet[2], octet[3])
	cnt, err := parseTemplate(template, config)
	if err != nil {
		logrus.Errorf("Failed to serve %v: %v", vars["environment"], err)
		return
	}
	cnt = strings.Replace(cnt, "REMOTE_IP", ip, -1)
	cnt = strings.Replace(cnt, "NFS_CLIENT_IP", nfsIP, -1)
  cnt = strings.Replace(cnt, "LAST_CLIENT_IP_OCTET", octet[3], -1) 
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
	_ = yaml.Unmarshal([]byte(string(config)), &configMap)

	fm := template.FuncMap{"substract": func(a, b int) int {
		return a - b
	}}

	t := template.Must(template.New(pathTemplate).Funcs(fm).ParseFiles(pathTemplate))

	var doc bytes.Buffer
	err = t.Execute(&doc, &configMap)
	if err != nil {
		logrus.Errorf("Failed to execute template: %s", err)
	}
	s := doc.String()

	return s, nil
}
