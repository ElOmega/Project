package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func isValidInput(input string) bool {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return false
	}
	for i := 0; i < len(parts); i++ {
		if _, err := strconv.Atoi(strings.TrimSpace(parts[i])); err != nil {
			return false
		}
	}
	return true
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection Error:", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte("start\n"))
	fmt.Println("Successful Connection. This is the Matrix")

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Server closed the connection.")
					os.Exit(0) // Quitte proprement le programme
				}
				fmt.Println("Connexion failed :", err)
				return
			}
			fmt.Printf("%s", message)
		}
	}()
	reader := bufio.NewReader(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter the coordinates (line,column) or 'end' to finish :")

		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if strings.ToLower(input) == "end" {
			conn.Write([]byte("end\n")) // Envoi du signal de fin au serveur
			fmt.Println("Setup completed. Starting simulation...")
			for {
				message, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Server closed the connection.")
					return
				}
				fmt.Print(message) // Affiche la matrice envoyée par le serveur
			}
		}
		if isValidInput(input) {
			valid, err := conn.Write([]byte(input + "\n"))
			if err != nil {
				fmt.Println("Error with the sending of the message :", err)
				return
			} else {
				fmt.Println(valid, "Valid input")
			}
		}

	}
}
