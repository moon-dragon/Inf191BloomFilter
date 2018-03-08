package main
package beego
beego.PprofOn = true


import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/vlam321/Inf191BloomFilter/payload"


	"github.com/fatih/structs"
	"github.com/spf13/viper"
)


// BloomServerIPs struct holding ips of each bloom filter server; retrieved through getBloomFilterIPs()
type BloomServerIPs struct {
	BloomFilterServer1  string
	BloomFilterServer2  string
	BloomFilterServer3  string
	BloomFilterServer4  string
	BloomFilterServer5  string
	BloomFilterServer6  string
	BloomFilterServer7  string
	BloomFilterServer8  string
	BloomFilterServer9  string
	BloomFilterServer10 string
}

type BloomContainerNames struct {
	BloomFilterContainer1  string
	BloomFilterContainer2  string
	BloomFilterContainer3  string
	BloomFilterContainer4  string
	BloomFilterContainer5  string
	BloomFilterContainer6  string
	BloomFilterContainer7  string
	BloomFilterContainer8  string
	BloomFilterContainer9  string
	BloomFilterContainer10 string
}

//updateTracker keeps track of which servers are updating currently or not 
var updateTracker[len(BloomServerIP)]bool

var bloomServerIPs BloomServerIPs
var bloomContainerNames BloomContainerNames
var routes map[int]string

func retrieveEndpoint(userid int) string {
	var endpoint string
	if viper.GetString("host") == "ecs" {
		endpoint = "http://" + os.Getenv(routes[userid]) + ":9090/filterUnsubscribed"
	} else {
		endpoint = "http://" + viper.GetString("dockerIP") + ":" + routes[userid] + "/filterUnsubscribed"
	}
	return endpoint
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	// read request data
	bbytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: Unable to read request data. %v\n", err)
		return
	}

	// unmarshal payload
	var pl payload.Payload
	err = json.Unmarshal(bbytes, &pl)
	if err != nil {
		log.Printf("Error: Unable to unmarshal Payload. %v\n", err)
		return
	}

	// determine endpoint based on host
	endpoint := retrieveEndpoint(pl.UserId)
	log.Printf("Request sent to: %s\n", endpoint)

	// make request to endpoint
	res, _ := http.Post(endpoint, "application/json; charset=utf-8", bytes.NewBuffer(bbytes))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Router: error reading response from bloom filter. %v\n", err)
	}
	w.Write(body)
}

//handleServerDown() routes requests to query the database if the bloomfilter is currently 
//being updated 
func handleServerDown(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Recived request: %v %v %v\n", r.Method, r.URL, r.Proto)
        bytes, err := ioutil.ReadAll(r.Body)
        if err != nil {
          log.Printf("Error: Unable to read request data. %v\n", err)
          return
        }

        //var to unload the recieved struct of bloomserver id and bool
        var ub downUpdate.DownUpdate
        err = json.Unmarshal(bytes, &ub)
        if err != nil{
          log.Printf("Error: Unable to unmarshal UpdateDown. %v", err)
          return
        }

        //update the array that tracks which server is updated
        updateTracker[ub.serverId] = ub.down 
        
}

// getMyIP() retrieve IP on local host
func getMyIP() (myIP string, err error) {
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "x.x.x.x", errors.New("Unable to find IP.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "x.x.x.x", errors.New("Unable to find IP.")
	}
	return string(body[:]), nil
}

// getBloomFilterIPs() retrieve IPs of each bloom filter server and store in bloomServerIPs
func getBloomFilterIPs() error {
	viper.SetConfigName("bfIPConf")
	viper.AddConfigPath("settings")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if viper.GetString("host") == "ecs" {
		log.Printf("host: ecs")
		err = viper.Unmarshal(&bloomContainerNames)
		if err != nil {
			return err
		}
	} else {
		log.Printf("host:docker")
		err = viper.Unmarshal(&bloomServerIPs)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapRouter(bloomFilterIPs BloomServerIPs) {
	routes = make(map[int]string)
	if viper.GetString("host") == "ecs" {
		containerNames := structs.Values(bloomContainerNames)
		for i := range containerNames {
			routes[i] = containerNames[i].(string)
		}
	} else {
		bloomIPs := structs.Values(bloomFilterIPs)
		for i := range bloomIPs {
			routes[i] = bloomIPs[i].(string)
		}
	}
	for k, v := range routes {
		log.Printf("key: %v	 value: %v\n", k, v)
	}
}

func main() {
	log.Printf("RETRIEVING BLOOM FILTER IP'S")
	err := getBloomFilterIPs()
	if err != nil {
		log.Println(err)
	}
	log.Printf("SUCCESSFULLY PARSED BLOOM SERVER IPS.")

	mapRouter(bloomServerIPs)
	log.Printf("SUCCESSFULLY MAPPED BLOOM SERVER IPS.")

	http.HandleFunc("/filterUnsubscribed", handleRoute)
        http.HandleFunc("/queryUnsubscribed", handleRoute)
	http.HandleFunc("/serverDown", handleServerDown)
        http.ListenAndServe(":9090", nil)
}
