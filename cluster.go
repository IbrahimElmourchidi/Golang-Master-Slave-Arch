package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type NodeInfo struct {
	Ip    string `json:"IP"`
	Port  string `json:"Port"`
	Event int    `json:"Event"`
}

type tcpMessage struct {
	Code string `json:"code"`
	Body string `json:"body"`
}

var allDestinations []NodeInfo
var master NodeInfo
var me NodeInfo
var isMaster bool = false

func initCluster() {
	startAsMaster := flag.Bool("master", false, "make this node a master if un able to connect to the given ip.")
	clusterIp := flag.String("clusterip", "127.0.0.1:8001", "ip address of the master")
	myPort := flag.String("myport", "8001", "the port on which this node is listening for tcp messages")
	flag.Parse()
	DBChannel := make(chan string)
	go connectDatabase(DBChannel)
	<-DBChannel
	me, master = generateConnectionString(*clusterIp, *myPort)
	startCluster(*startAsMaster)
}

func startCluster(startAsMaster bool) {
	ableToConnect := sendMessage(master, tcpMessage{
		Code: "000",
		Body: serialize(me),
	})
	if ableToConnect {
		fmt.Println("Connected to master")
		listenOnPort(me)
	} else if startAsMaster {
		fmt.Println("will start this node as master")
		isMaster = true
		go listenOnPort(me)
	} else {
		fmt.Println("Cannot connect to master")
	}
}

func generateConnectionString(clusterip string, myport string) (NodeInfo, NodeInfo) {
	lastEvent := Event{}
	err := dbManager.Last(&lastEvent)
	eventId := 0
	if err != nil {
		eventId = lastEvent.ID
	}
	myIp := getLocalIp().String()
	me := NodeInfo{
		Ip:    myIp,
		Port:  myport,
		Event: eventId,
	}
	master := NodeInfo{
		Ip:   strings.Split(clusterip, ":")[0],
		Port: strings.Split(clusterip, ":")[1],
	}
	return me, master
}

func getLocalIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func sendMessage(dest NodeInfo, message tcpMessage) bool {
	conn, err := net.DialTimeout("tcp", dest.Ip+":"+dest.Port, time.Duration(10)*time.Second)
	if err != nil {
		if isMaster {
			for i, node := range allDestinations {
				if dest == node {

					allDestinations = append(allDestinations[:i], allDestinations[i+1:]...)
				}
			}
		} else {

			fmt.Println("Couldn't connect to master")
		}
		return false
	} else {
		json.NewEncoder(conn).Encode(&message)
		return true
	}
}

func listenOnPort(me NodeInfo) {
	fmt.Println("Start Listening for messages")
	ln, _ := net.Listen("tcp", fmt.Sprint(":"+me.Port))
	if isMaster == false {
		sendMessage(master, tcpMessage{
			Code: "009",
			Body: serialize(me),
		})
	}
	for {
		connIn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error listening to message")
		} else {
			var incommingMessage tcpMessage
			json.NewDecoder(connIn).Decode(&incommingMessage)
			fmt.Println("Message Recieved ==>", incommingMessage.Code)
			if incommingMessage.Code == "000" {
				handleNewConnection(connIn, incommingMessage.Body)
			} else if incommingMessage.Code == "001" {
				handleCreateAuthor(connIn, incommingMessage.Body)
			} else if incommingMessage.Code == "002" {
				handleCreateBook(connIn, incommingMessage.Body)
			} else if incommingMessage.Code == "003" {
				handleUpdateBook(connIn, incommingMessage.Body)
			} else if incommingMessage.Code == "004" {
				handleDeleteBook(connIn, incommingMessage.Body)
			} else if incommingMessage.Code == "009" {
				sendBackEvents(connIn, incommingMessage.Body)
			}
		}
	}
}

func messageToAll(msg tcpMessage) {
	for i := range allDestinations {
		sendMessage(allDestinations[i], msg)
	}
}
