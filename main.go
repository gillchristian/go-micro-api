// the go-micro-api runs a RESTful microservice to handle a todo list
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	conf, err := ReadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := Connect(conf.DB.Url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()
	EnsureIndex(session, conf.DB)

	router := Router(&routes, session, conf.DB)

	err = runServer(conf.Port, router)
	fmt.Println(err)
	return
}

// runServer runs the http.ListenAndServe server and
// prompts for a new port if the provided one is
// already in use then tries again with the new port
func runServer(port string, router http.Handler) error {
	fmt.Printf("Running on localhost:%v\n", port)

	err := http.ListenAndServe(":"+port, router)

	if err.Error() == usedPortError(port) {
		fmt.Printf("Port %v is busy\n\n", port)
		newPort, err := prompNewPort()
		if err != nil {
			return err
		}
		fmt.Println(newPort)
		return runServer(strings.TrimSpace(newPort), router)
	}
	return err
}

// usedPortError returns http error message for a already in use port
func usedPortError(port string) string {
	return fmt.Sprintf("listen tcp :%v: bind: address already in use", port)
}

// prompNewPort prompts the user to input a port
// on the terminal and returns it
func prompNewPort() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter a new port: ")
	return reader.ReadString('\n')
}
