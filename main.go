package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Token    string
	Scrapers map[string]tomlSectionScraper
}

type tomlSectionScraper struct {
	Name             string
	URL              string
	Deeplinkselector string
	Template         string
	Fields           map[string]tomlFieldDefinition
}

type tomlFieldDefinition struct {
	Selector  string
	Attribute string
}

const configFileName = "config.toml"

func readConfig() (*tomlConfig, error) {
	configcontent, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatalf("Error opening config file %s: %s", configFileName, err)
		return nil, err
	}
	var conf tomlConfig
	if _, err := toml.Decode(string(configcontent), &conf); err != nil {
		log.Fatalf("Error decoding config file %s: %s", configFileName, err)
	}
	return &conf, nil
}

func initLogging() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "trace"
	}
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
}

func main() {
	initLogging()
	conf, err := readConfig()
	if err != nil {
		logrus.Fatalf("Error reading config: %v", err)
	}

	resultsChannel := make(chan *SearchResult, 100)
	for _, scraperConf := range conf.Scrapers {
		logrus.Debugf("Creating scraper '%s'", scraperConf.Name)
		scraper := NewScraper(scraperConf)
		go scraper.Scrape(resultsChannel)
	}

	bot := NewBot(conf.Token)

	for {
		res := <-resultsChannel
		tpl := NewMessageTemplate(res.scraper.config.Template)

		bot.Send(tpl.render(res.results))
		time.Sleep(2 * time.Second)
	}

}
