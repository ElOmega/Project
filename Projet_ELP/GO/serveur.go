package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// La matrice initiale
func matrice_init(N int) [][]int {
	matrice := make([][]int, N)
	for i := 0; i < N; i++ {
		matrice[i] = make([]int, N)
	}
	return matrice
}

// voir les compte les voisins
func compte_voisins(matrice [][]int, i int, j int, rows int, cols int) int {
	sum := 0

	// Haut-gauche
	if i > 0 && j > 0 && matrice[i-1][j-1] == 1 {
		sum++
	}
	// Haut
	if i > 0 && matrice[i-1][j] == 1 {
		sum++
	}
	// Haut-droite
	if i > 0 && j < cols-1 && matrice[i-1][j+1] == 1 {
		sum++
	}
	// Gauche
	if j > 0 && matrice[i][j-1] == 1 {
		sum++
	}
	// Droite
	if j < cols-1 && matrice[i][j+1] == 1 {
		sum++
	}
	// Bas-gauche
	if i < rows-1 && j > 0 && matrice[i+1][j-1] == 1 {
		sum++
	}
	// Bas
	if i < rows-1 && matrice[i+1][j] == 1 {
		sum++
	}
	// Bas-droite
	if i < rows-1 && j < cols-1 && matrice[i+1][j+1] == 1 {
		sum++
	}

	return sum
}

// Regarde si il y a trop de voisins ou non
func vivant(matrice [][]int, i int, j int, sum int) int {
	if sum == 3 || (matrice[i][j] == 1 && sum == 2) {
		return 1
	}
	return 0
}

// Affiche la matrice
func SendMatrice(conn net.Conn, matrice [][]int, rows int) {
	for i := 0; i < rows; i++ {
		line := "" // On construit la ligne avant de l'envoyer
		for j := 0; j < rows; j++ {
			if matrice[i][j] == 1 {
				line += "O "
			} else {
				line += ". "
			}
		}
		fmt.Fprintln(conn, line) // Envoi de la ligne avec un retour à la ligne propre
	}
	fmt.Fprintln(conn, "---") // Séparateur clair après la matrice
}

// repartire les fonctions
func ProcessSection(startRow int, endRow int, matrice [][]int, tab_copie [][]int, rows int, cols int, wg *sync.WaitGroup, ch chan<- [][]int) {
	defer wg.Done()
	localResults := matrice_init(rows)

	for i := startRow; i < endRow; i++ {
		for j := 0; j < cols; j++ {
			sum := compte_voisins(tab_copie, i, j, rows, cols)
			localResults[i][j] = vivant(tab_copie, i, j, sum)
		}
	}
	ch <- localResults
}

func copie(matrice [][]int, rows int, cols int) [][]int {

	matrice_copie := matrice_init(rows)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrice_copie[i][j] = matrice[i][j]
		}
	}
	return matrice_copie
}

func clientdecision(conn net.Conn, matrice [][]int, rows int) {
	reader := bufio.NewReader(conn)

	for {
		// Envoyer la matrice au client pour qu'il la voit dès le début
		SendMatrice(conn, matrice, rows)

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		message = strings.TrimSpace(message)
		if message == "end" {
			fmt.Println("Client finished setup.")
			break // Ne ferme pas la connexion, juste quitte la boucle
		}

		coords := strings.Split(message, ",")
		if len(coords) != 2 {
			fmt.Fprintln(conn, "Error: Invalid format. Please send row,col")
			continue
		}

		row, err1 := strconv.Atoi(coords[0])
		col, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil || row >= rows || col >= rows || row < 0 || col < 0 {
			fmt.Fprintln(conn, "Error: Invalid coordinates. Please try again.")
			continue
		}

		if matrice[row][col] == 1 {
			matrice[row][col] = 0
			fmt.Fprintln(conn, "Cell removed")
		} else {
			matrice[row][col] = 1
			fmt.Fprintln(conn, "Cell added")
		}
	}
}

func main() {
	N := 20
	nombre_tour := 100
	rows, cols := N, N
	matrice := matrice_init(N)

	// Démarrer le serveur TCP
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error when starting the server:", err)
		return
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error of connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connecté :", conn.RemoteAddr())

	// Phase d'initialisation
	clientdecision(conn, matrice, rows)

	// Simulation
	for round := 0; round < nombre_tour; round++ {
		start_time := time.Now()
		fmt.Println("Round :", round+1)
		tab_copie := copie(matrice, rows, cols)

		// Concurrence
		var wg sync.WaitGroup
		NumWorkers := 4
		ch := make(chan [][]int, NumWorkers)
		sectionSize := rows / NumWorkers

		for i := 0; i < NumWorkers; i++ {
			startRow := i * sectionSize
			endRow := startRow + sectionSize
			if i == NumWorkers-1 {
				endRow = rows
			}
			wg.Add(1)
			go ProcessSection(startRow, endRow, matrice, tab_copie, rows, cols, &wg, ch)
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		NewMatrice := matrice_init(rows)
		for results := range ch {
			for i := 0; i < rows; i++ {
				for j := 0; j < cols; j++ {
					if results[i][j] == 1 {
						NewMatrice[i][j] = 1
					}
				}
			}
		}
		matrice = NewMatrice

		// Envoyer la matrice au client après chaque round
		SendMatrice(conn, matrice, rows)

		end_time := time.Now()
		time_took := end_time.Sub(start_time)
		fmt.Println("The Time took to do the code is :", time_took)
		time.Sleep(1 * time.Second)
	}

	// Fermer la connexion uniquement à la fin de la simulation
	fmt.Println("Simulation finished. Closing connection.")

}
