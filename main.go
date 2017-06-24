package main

import "encoding/csv"
import "errors"
import "flag"
import "fmt"
import "github.com/BurntSushi/toml"
import "io/ioutil"
import "os"
import "sort"
import "strconv"

type config struct {
	IsDescendingOrder bool
	Data              []*Datapoint
}

type configFile struct {
	Sort     string
	Datafile string
}

func parseConfigFile(configPath string) (*configFile, error) {
	if configPath == "" {
		return nil, nil
	}

	tomlData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var conf configFile
	_, err = toml.Decode(string(tomlData), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func parseSort(sort string) (bool, error) {
	var isDescendingOrder bool
	if sort == "ascending" { // lower is better
		isDescendingOrder = false
	} else if sort == "descending" { // higher is better
		isDescendingOrder = true
	} else {
		return true, errors.New(fmt.Sprintf("\"%s\" is not a recognized sort. Valid values are \"ascending\" and \"descending\".", sort))
	}
	return isDescendingOrder, nil
}

func parseData(dataFilePath string) ([]*Datapoint, error) {
	fp, err := os.Open(dataFilePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	datapoints := make([]*Datapoint, len(records))
	for i, record := range records {
		datapoints[i] = new(Datapoint)
		datapoints[i].name = record[0]
		datapoints[i].score, err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
	}
	return datapoints, nil
}

func makeConfig(configPathPtr *string) (*config, error) {
	conf, err := parseConfigFile(*configPathPtr)
	if err != nil {
		return nil, err
	}

	isDescendingOrder, err := parseSort(conf.Sort)
	if err != nil {
		return nil, err
	}

	datapoints, err := parseData(conf.Datafile)
	if err != nil {
		return nil, err
	}
	sort.Stable(DatapointSort(datapoints))
	if isDescendingOrder {
		for i, j := 0, len(datapoints)-1; i < j; i, j = i+1, j-1 {
			datapoints[i], datapoints[j] = datapoints[j], datapoints[i]
		}
	}

	options := config{isDescendingOrder, datapoints}
	return &options, nil
}

func main() {
	configPathPtr := flag.String("path", "", "The path to the config file.")
	flag.Parse()

	config, err := makeConfig(configPathPtr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, datapoint := range config.Data {
		fmt.Printf("%s %f\n", datapoint.name, datapoint.score)
	}
}
