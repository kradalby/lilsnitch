package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
	"io/ioutil"
	"net/http"
	"strings"
)

type LittleSnitch struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Rules       []LittleSnitchRule `json:"rules"`
}

type LittleSnitchRule struct {
	Action      string `json:"action"`
	Process     string `json:"process"`
	RemoteHosts string `json:"remote-hosts"`
	Direction   string `json:"direction"`
}

func main() {
	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"alive": true})
	})

	r.GET("/hosts.lsrules", func(c *gin.Context) {

		hostMap, err := GetHostMap("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		rules := CreateLittleSnitch("Steven Black's hosts", "Host list created by Steven Black, https://github.com/StevenBlack/hosts", hostMap["0.0.0.0"])

		c.JSON(200, rules)
	})

	r.Run()
}

func CreateLittleSnitch(name string, description string, hosts []string) LittleSnitch {
	rules := make([]LittleSnitchRule, len(hosts))
	for index, host := range hosts {
		rules[index] = CreateLittleSnitchRule(host)
	}
	return LittleSnitch{
		Name:        name,
		Description: description,
		Rules:       rules,
	}
}

func CreateLittleSnitchRule(host string) LittleSnitchRule {
	return LittleSnitchRule{
		Action:      "deny",
		Process:     "any",
		RemoteHosts: host,
		Direction:   "outgoing",
	}
}

func GetHostMap(hostsURL string) (map[string][]string, error) {
	response, err := http.Get(hostsURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	hosts, err := ParseHosts(content)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// Taken from:
// https://github.com/jaytaylor/go-hostsfile/blob/master/hosts.go
// ParseHosts takes in hosts file content and returns a map of parsed results.
func ParseHosts(hostsFileContent []byte) (map[string][]string, error) {
	hostsMap := map[string][]string{}
	for _, line := range strings.Split(strings.Trim(string(hostsFileContent), " \t\r\n"), "\n") {
		line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) > 1 && len(pieces[0]) > 0 {
			if names := strings.Fields(pieces[1]); len(names) > 0 {
				if _, ok := hostsMap[pieces[0]]; ok {
					hostsMap[pieces[0]] = append(hostsMap[pieces[0]], names...)
				} else {
					hostsMap[pieces[0]] = names
				}
			}
		}
	}
	return hostsMap, nil
}
