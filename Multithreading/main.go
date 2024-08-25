package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fastjson"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	cep := "67033000"

	go Cep_ViaCEP(cep, c1)
	go Cep_BrasilAPI(cep, c2)

	select {
	case msg1 := <-c1:
		fmt.Printf("Resultados da consulta utilizando a API VIACEP:\n%v\n", msg1)
	case msg2 := <-c2:
		fmt.Printf("Resultados da consulta utilizando a API BRASILAPI:\n%v\n", msg2)
	case <-time.After(time.Second):
		fmt.Printf("Timeout")
	}
}

func Cep_ViaCEP(endereco string, canal chan<- string) {
	var p fastjson.Parser

	url := "https://viacep.com.br/ws/" + endereco + "/json/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erro ao processar a requisição utilizando a API VIACEP: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		corpo, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Erro ao consultar o site VIACEP")
		}

		dados, err := p.Parse(string(corpo))
		if err != nil {
			log.Fatalf("Erro ao processar a requisição utilizando a API VIACEP: %v", err)
		}
		canal <- dados.String()
	} else {
		log.Fatalf("Erro ao processar a requisição utilizando a API VIACEP: %v", err)
	}
}

func Cep_BrasilAPI(endereco string, canal chan<- string) {
	var p fastjson.Parser

	url := "https://brasilapi.com.br/api/cep/v1/" + endereco
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erro ao processar a requisição utilizando a API BRASILAPI: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		corpo, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error ao acessar o site BRASILAPI")
		}

		dados, err := p.Parse(string(corpo))
		if err != nil {
			log.Fatalf("Erro ao processar a requisição utilizando a API BRASILAPI: %v", err)
		}
		canal <- dados.String()
	} else {
		log.Fatalf("Erro ao processar a requisição utilizando a API BRASILAPI: %v", err)
	}
}
