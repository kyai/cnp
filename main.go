package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		usage()
	}

	switch strings.ToLower(args[1]) {
	case "list":
		list()
	case "env":
	case "set":
	case "unset":
	default:
		usage()
	}
}

func usage() {
	fmt.Println("----- usage -----")
	os.Exit(0)
}

func list() {
	nodes, err := getList()
	if err != nil {
		panic(err)
	}

	output := ""
	for _, node := range nodes {
		color := 0
		speed := node.Speed / 10
		switch {
		case speed >= 9:
			color = 42
		case speed >= 6:
			color = 43
		case speed >= 0:
			color = 41
		}

		output += fmt.Sprintf("%s\t%s\t%v\n\n",
			space(node.Ip, 15),
			space(node.Port, 5),
			fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, space("", speed)),
		)
	}
	fmt.Println(output)
}

const WEB_URL = "https://cn-proxy.com"

// Node of single proxy
type Node struct {
	Ip        string
	Port      string
	Location  string
	Speed     int
	LastCheck time.Time
}

// Get list of proxy's node
func getList() (nodes []Node, err error) {
	res, err := httpGet(WEB_URL)
	if err != nil {
		return
	}

	r := regexp.MustCompile(`
<tr>
<td>(?P<ip>.+?)</td>
<td>(?P<port>.+?)</td>
<td>(?P<location>.+?)</td>
<td>
.+width: (?P<speed>.+?)%.+
</td>
<td>(?P<lastcheck>.+?)</td>
</tr>`)
	m := r.FindAllSubmatch(res, -1)

	a := make([]map[string]string, len(m))

	for k, v := range m {
		a[k] = make(map[string]string)
		for i, name := range r.SubexpNames() {
			if i != 0 && len(name) > 0 {
				a[k][name] = string(v[i])
			}
		}
	}

	nodes = make([]Node, len(a))

	for k, v := range a {
		nodes[k].Ip = v["ip"]
		nodes[k].Port = v["port"]
		nodes[k].Location = v["location"]
		nodes[k].Speed, _ = strconv.Atoi(v["speed"])
		nodes[k].LastCheck, _ = time.Parse("2006-01-02 15:04:05", v["lastcheck"])
	}

	return
}

func httpGet(url string) (res []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// possible methods

func space(s string, n int) string {
	n -= len(s)
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}
