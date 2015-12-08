package filedog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Hosts struct {
	hostName  []string
	timeStap  int64
	mux       sync.Mutex
	watchDone chan struct{}
}

var globleHosts = &Hosts{hostName: make([]string, 0)}
var initDone = make(chan struct{})

func init() {
	getHostName()
	watchDog()
}

func watchDog() {

	go func() {
		for {
			select {
			case <-globleHosts.watchDone:

				return
			case <-time.After(5 * time.Second):
				info, err := os.Stat("/etc/hosts")
				if err != nil {
					fmt.Println(err)
					return
				}

				if info.ModTime().Unix() > globleHosts.timeStap {

					// fmt.Println("file modified, reset the host again")
					getHostName()

					// fmt.Printf("new hosts : \n %#v\n", globleHosts)
				}

			}

		}
	}()

}
func getHostName() {
	f, err := os.OpenFile("/etc/hosts", os.O_RDONLY, 0444)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	info, _ := f.Stat()
	var hosts = make([]string, 0)
	bufReader := bufio.NewReader(f)
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
				return
			}
			break
		}
		if !strings.HasPrefix(line, "#") {
			hosts = append(hosts, line)
		}
	}

	globleHosts.mux.Lock()
	globleHosts.hostName = make([]string, 0)
	for _, host := range hosts {
		split := strings.Fields(host)
		if len(split) == 0 {
			fmt.Println("empty : ", split)
			continue
		}

		globleHosts.hostName = append(globleHosts.hostName, split[1])

	}
	globleHosts.timeStap = info.ModTime().Unix()

	globleHosts.mux.Unlock()
	// fmt.Printf("hosts : %#v\n ", globleHosts)
}

func HostNames() []string {
	names := make([]string, len(globleHosts.hostName))
	copy(names, globleHosts.hostName)
	return names
}
func JsonHostNames() ([]byte, error) {
	jsonData := struct {
		HostNames []string
		TimeStap  int64
	}{
		HostNames: globleHosts.hostName,
		TimeStap:  globleHosts.timeStap,
	}
	return json.Marshal(jsonData)
}

func Stop() {
	globleHosts.watchDone <- struct{}{}
}
