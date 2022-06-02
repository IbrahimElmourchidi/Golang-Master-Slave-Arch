package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func handleNewConnection(conn net.Conn, body string) {
	node := nodeParser(body)
	allDestinations = append(allDestinations, node)

	// go sendBackEvents(conn, backEvents)
	fmt.Println(node.Ip+":"+node.Port, "is connected")
	conn.Close()
}

func handleCreateAuthor(conn net.Conn, body string) {
	author := authorParser(body)
	dbOk, _ := createAuthor(author)
	if dbOk {
		saveEvent("001", body)
		conn.Close()
	} else {
		sendErrorMessage(conn)
	}
}

func handleCreateBook(conn net.Conn, body string) {
	book := bookParser(body)
	dbOk, _ := createBook(book)
	if dbOk {
		saveEvent("002", body)
		conn.Close()
	} else {
		sendErrorMessage(conn)
	}
}

func handleUpdateBook(conn net.Conn, body string) {
	book := bookParser(body)
	dbOk, _ := updateBook(book)
	if dbOk {
		saveEvent("003", body)
		conn.Close()
	} else {
		sendErrorMessage(conn)
	}

}

func handleDeleteBook(conn net.Conn, body string) {
	book := bookParser(body)
	dbOk := deleteBook(book.ID)
	if dbOk {
		saveEvent("004", body)
		conn.Close()
	} else {
		sendErrorMessage(conn)
	}
}

func nodeParser(body string) NodeInfo {
	body = strings.Replace(body, "{", "", -1)
	body = strings.Replace(body, "}", "", -1)
	Ip := strings.Split(body, " ")[0]
	Port := strings.Split(body, " ")[1]
	Event := strings.Split(body, " ")[2]
	eventInt, _ := strconv.Atoi(Event)
	return NodeInfo{
		Ip:    Ip,
		Port:  Port,
		Event: eventInt,
	}
}

func authorParser(body string) Author {
	body = trimBody(body)
	ID := strings.Split(body, " ")[0]
	idInt, _ := strconv.Atoi(ID)
	First := strings.Split(body, " ")[1]
	Last := strings.Split(body, " ")[2]
	return Author{
		ID:    idInt,
		First: First,
		Last:  Last,
	}
}

func bookParser(body string) Book {
	fmt.Println(body)
	ID := strings.Split(body, " ")[0]
	idInt, _ := strconv.Atoi(ID)
	Isbn := strings.Split(body, " ")[1]
	Title := strings.Split(body, " ")[2]
	price := strings.Split(body, " ")[3]
	priceInt, _ := strconv.Atoi(price)
	authorID := strings.Split(body, " ")[4]
	authInt, _ := strconv.Atoi(authorID)
	return Book{
		ID:       idInt,
		Isbn:     Isbn,
		Title:    Title,
		Price: priceInt,
		AuthorID: authInt,
	}
}

func trimBody(body string) string {
	body = strings.Replace(body, "{", "", -1)
	body = strings.Replace(body, "}", "", -1)
	return body
}

func sendErrorMessage(conn net.Conn) {
	message := tcpMessage{
		Code: "111",
		Body: me.Ip + me.Port + " failed",
	}
	json.NewEncoder(conn).Encode(&message)
	conn.Close()
}

func sendBackEvents(conn net.Conn, body string) {
	node := nodeParser(body)
	fmt.Println(node.Event)
	backEvents := []Event{}
	dbManager.Where("id > ?", node.Event).Find(&backEvents)
	for _, val := range backEvents {
		fmt.Println("sending ",val.Code)
		message := tcpMessage{
			Code: val.Code,
			Body: val.Body,
		}
		sendMessage(node, message)
	}
	conn.Close()
}

