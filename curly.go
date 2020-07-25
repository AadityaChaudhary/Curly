package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args)!= 2 {
		log.Fatal("You need 2 Arguments!!")
	}

	var queryArgs map[string]string

	jsonFile := args[1]

	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &queryArgs)
	if err != nil {
		log.Fatal(err)
	}

	address := args[0]

	//construct query string from map of args
	var query string
	for key, val := range queryArgs {
		if len(query) != 0 {
			query += "&"
		}
		query = query + key + "=" + val
	}
	url := fmt.Sprintf("ws://%s/ws", address) + "?" + query
	log.Println("connecting to " + url)
	log.Println(" . . . ")

	ws, err := websocket.Dial(url, "", fmt.Sprintf("http://%s/", address))
	if err != nil {
		log.Fatal(err)
	}

	go readClientMessages(ws)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected!")
	fmt.Println("---------------------")

	for {
		fmt.Println("Enter filepath for json-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("quit", text) == 0 {
			return
		}



		data, err := ioutil.ReadFile(text)
		if err != nil {
			log.Fatal(err)
		}

		_, err = ws.Write(data)
		if err != nil {
			log.Fatal(err)
		}

	}


}

func readClientMessages(ws *websocket.Conn ) {
	for {
		var message string
		// err := websocket.JSON.Receive(ws, &message)
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Messages Recieved --->", message)
		fmt.Println("Enter filepath for json-> ")
	}
}
