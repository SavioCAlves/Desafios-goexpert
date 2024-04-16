package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DolarToReal struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type DolarToRealResult struct {
	VarBid string `json:"varBid"`
}

type CotacaoBR struct {
	ID         int `gorm:primaryKey`
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
}

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", nil)
	// db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	// Criando Base de dados
	// db.AutoMigrate(&CotacaoBR{})
}

// Servidor web e rotas.
func CotacaoHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := Cotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

}

// Processamento da API
func Cotacao() (*DolarToRealResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var cot DolarToReal
	err = json.Unmarshal(body, &cot)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	SalvaCotacao(&cot)
	resultado := DolarToRealResult{
		VarBid: cot.Usdbrl.VarBid,
	}
	return &resultado, nil

}

// Insercao no banco de dados.
func SalvaCotacao(dados *DolarToReal) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if err = db.WithContext(ctx).Create(&CotacaoBR{
		Code:       dados.Usdbrl.Code,
		Codein:     dados.Usdbrl.Codein,
		Name:       dados.Usdbrl.Name,
		High:       dados.Usdbrl.High,
		Low:        dados.Usdbrl.Low,
		VarBid:     dados.Usdbrl.VarBid,
		PctChange:  dados.Usdbrl.PctChange,
		Bid:        dados.Usdbrl.Bid,
		Ask:        dados.Usdbrl.Ask,
		Timestamp:  dados.Usdbrl.Timestamp,
		CreateDate: dados.Usdbrl.CreateDate,
	}).Error; err != nil {
		log.Println(err)
	}
	return nil
}
