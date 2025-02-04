package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// La matrice initiale est carré et est de dimension fini. 1 = cellule vivante, 0 = cellule morte

// initialisation de la matrice
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

// En fonction de sum, voir si la cellule survit ou meurt
func vivant(matrice [][]int, i int, j int, sum int) int {
	if sum == 3 || (matrice[i][j] == 1 && sum == 2) {
		return 1
	}
	return 0
}

// copie du tableau pour modifier plus tard
func copie(matrice [][]int, rows int, cols int) [][]int {

	matrice_copie := matrice_init(rows)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrice_copie[i][j] = matrice[i][j]
		}
	}
	return matrice_copie
}

// fonction d'affichage
func SendMatrice(conn net.Conn, matrice [][]int, rows int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			if matrice[i][j] == 1 {
				fmt.Fprintf(conn, "O ")
			} else {
				fmt.Fprintf(conn, ". ")
			}
		}
		fmt.Fprintln(conn)
	}
	fmt.Fprintln(conn, "---")
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

// serveur
func clientdecision(conn net.Conn, matrice [][]int, rows int) {
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		SendMatrice(conn, matrice, rows)

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Reading Error :", err)
			return
		}
		message = strings.TrimSpace(message)
		coords := strings.Split(message, ";")
		if len(coords) != 2 {
			fmt.Print("Error, send rows;line")
			continue
		}
		row, err1 := strconv.Atoi(coords[0])
		col, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil || row >= rows || col >= rows || row < 0 || col < 0 {
			fmt.Fprintln(conn, "Invalide Coordinates, try again")
			continue
		}

		matrice[row][col] = 1

		fmt.Fprintln(conn, "Cell added")

	}
}

// lancer le code et initialisations des variable
// depart des cellules
// boucle principale
func main() {
	//initialisation valeur
	N := 20
	nombre_tour := 100
	rows, cols := N, N
	matrice := matrice_init(N)
	var conn net.Conn

	//démarrage du serveur TCP
	go func() {
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			fmt.Println("Error when starting the server :", err)
			return
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error of connection:", err)
				continue
			}
			fmt.Println("Client connecté :", conn.RemoteAddr())

			// Gestion de la connexion client
			go clientdecision(conn, matrice, rows)
		}

	}()

	//démarrage du jeu
	start_time := time.Now()
	for round := 0; round < nombre_tour; round++ {

		fmt.Println("Round :", round+1)
		tab_copie := copie(matrice, rows, cols)

		//Concurrence
		var wg sync.WaitGroup
		NumWorkers := 4
		ch := make(chan [][]int, NumWorkers)
		sectionSize := rows / NumWorkers

		//lancer les Goroutines
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

		//reccupérer les resultats
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
		SendMatrice(conn, matrice, rows)
		end_time := time.Now()
		time_took := end_time.Sub(start_time)
		fmt.Println("The Time took to do the code is :", time_took)

	}
}
