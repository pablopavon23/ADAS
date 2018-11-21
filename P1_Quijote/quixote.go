package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Vertex struct {
	Lines_Appear []int
	Count        int
}

var m map[string]Vertex



func insert_in_map(s []string, lines int) {
	var count int = 0

	for index := 0; index < len(s); index++ { 
		word := strings.ToLower(s[index]) 
		if m[word].Count == 0 {           
			count = 1
		} else { 
			count = m[word].Count 
			count = count + 1
		}
		m[word] = Vertex{append(m[word].Lines_Appear, lines), count} // Añado en Lines_Appear la linea
	}
}

func print_map(m map[string]Vertex) {
	for key, value := range m {
		fmt.Println(key, ":", value.Count)
		for i := range value.Lines_Appear {
			fmt.Println("\tquijote.txt:", value.Lines_Appear[i])
		}
	}
}

func main() {
	fileName := os.Args[1]         // Leo el nombre del fichero desde linea de comandos
	file, err := os.Open(fileName) // Abro el fichero
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Esto es para que luego se cierre

	m = make(map[string]Vertex)
	scanner := bufio.NewScanner(file)
	var lines int = 0 
	for scanner.Scan() {
		actual_line := scanner.Text()
		if actual_line != "" { 
			lines = lines + 1                    
			s := strings.Split(actual_line, " ") 
			insert_in_map(s, lines)
		}

	}

	print_map(m)

}
