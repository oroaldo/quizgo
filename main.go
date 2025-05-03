package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja bem vindo(a) ao quiz")
	fmt.Println("Escreva o seu nome: ")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Erro ao ler a string nome")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo %s", g.Name)
}

func (g *GameState) ProcessCSV() {
	f, err := os.Open("quiz.csv")
	if err != nil {
		panic("Erro ao ler arquivo")
	}

	defer f.Close() //fecha ao final da função defer indica ao GO para fazer ao final, pode estar em qualquer ponto aqui do método

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Erro ao ler o csv")
	}

	for index, record := range records {
		//fmt.Println(index, record)
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	//Exibir a pergunta para o usuario
	for index, question := range g.Questions {
		fmt.Printf("\033[33m %d. %s \033[0m\n", index+1, question.Text)

		//iterar as opcoes para responder
		for j, option := range question.Options {
			fmt.Printf("[%d] - %s\n", j+1, option)
		}

		//solicita digitar alternativa
		println("Digite uma alternativa: ")

		//coletar entradado usuario e validar caracter, se for errado precisa inserir novamente
		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			answer, err = toInt(read[:len(read)-2])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			break
		}

		//validar resposta, exibir se esta corretoou não e somar a pontuação
		if answer == question.Answer {
			fmt.Println("Parabéns você acertou")
			g.Points += 10
		} else {
			fmt.Println("Ops, errou!")
			fmt.Println("--------------")
		}

	}
}

func main() {
	game := &GameState{}
	go game.ProcessCSV()
	game.Init()
	game.Run()

	fmt.Printf("Fim de Jogo, você fez %d", game.Points)
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("insira um número")
	}

	return i, nil
}
