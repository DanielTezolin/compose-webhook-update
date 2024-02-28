package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Log(message string) {
	fmt.Println(message)
	// Obter a data atual
	dataAtual := time.Now().Format("2006-01-02")
	horaAtual := time.Now().Format("15:04:05")

	absPath, _ := filepath.Abs("./logs/" + dataAtual + ".txt")

	// Verificar se o arquivo já existe
	arquivo, err := os.OpenFile(absPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// Se nao existir, criar o arquivo
		arquivo, err = os.Create(absPath)
		if err != nil {
			fmt.Println("Erro ao criar o arquivo de log:", err)
			return
		}
		defer arquivo.Close() // Fechar o arquivo no final da função
	}

	log := fmt.Sprint(horaAtual, " - ", message, "\n")

	_, err = arquivo.WriteString(log)
	if err != nil {
		fmt.Println("Erro ao escrever no log:", err)
		return
	}
}
