package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func TestBids(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@postgres/bids?sslmode=disable")
	if err != nil {
		t.Error(err)
	}

	url := httptest.NewServer(bidsHandler(db)).URL

	//balance check
	setBalance(t, db, 0)
	if getBalance(t, db) != 0 {
		t.Error("could not set balance to 0")
	}

	amount := "11.11"

	// win
	tx := BidJSON{
		BidID: fmt.Sprintf("%d", rand.Int()),
		State:         win,
		Amount:        amount,
	}
	code := sendTx(t, url, tx, "game")
	if code != http.StatusOK {
		t.Error("win. wrong status code")
	}
	if balanceString(getBalance(t, db)) != amount {
		t.Error("win. wrong balance")
	}

	// win idempotence
	code = sendTx(t, url, tx, "game")
	if code != http.StatusOK {
		t.Error("win idempotence. wrong status code")
	}
	if balanceString(getBalance(t, db)) != amount {
		t.Error("win idempotence. wrong balance")
	}

	// lose
	tx = BidJSON{
		BidID: fmt.Sprintf("%d", rand.Int()),
		State:         lose,
		Amount:        amount,
	}
	code = sendTx(t, url, tx, "server")
	if code != http.StatusOK {
		t.Error("lose. wrong status code")
	}
	if balanceString(getBalance(t, db)) != "0.00" {
		t.Error("lose. wrong balance")
	}

	// lose idempotence
	code = sendTx(t, url, tx, "server")
	if code != http.StatusOK {
		t.Error("lose idempotence. wrong status code")
	}
	if balanceString(getBalance(t, db)) != "0.00" {
		t.Error("lose idempotence. wrong balance")
	}

	// lose insufficient funds
	code = sendTx(t, url, BidJSON{
		BidID: fmt.Sprintf("%d", rand.Int()),
		State:         lose,
		Amount:        amount,
	}, "game")
	if code != http.StatusBadRequest {
		t.Error("lose insufficient funds. wrong status code")
	}
	if balanceString(getBalance(t, db)) != "0.00" {
		t.Error("check lose insufficient funds. wrong balance")
	}

	// source type payment
	code = sendTx(t, url, BidJSON{
		BidID: fmt.Sprintf("%d", rand.Int()),
		State:         win,
		Amount:        amount,
	}, "payment")
	if code != http.StatusOK {
		t.Error("source type payment. wrong status code")
	}
	if balanceString(getBalance(t, db)) != amount {
		t.Error("source type payment. wrong balance")
	}

	// wrong source type
	code = sendTx(t, url, BidJSON{
		BidID: fmt.Sprintf("%d", rand.Int()),
		State:         lose,
		Amount:        amount,
	}, "qwerty")
	if code != http.StatusInternalServerError {
		t.Error("wrong status code. check wrong source type")
	}
	if balanceString(getBalance(t, db)) != amount {
		t.Error("wrong status code. wrong balance")
	}
}

func sendTx(t *testing.T, url string, tx BidJSON, sourceType string) int {
	body, err := json.Marshal(tx)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Source-Type", sourceType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func setBalance(t *testing.T, db *sqlx.DB, balance int) {
	_, err := db.Exec("update player set balance = $1 where phone_number = '1'", balance)
	if err != nil {
		t.Error(err)
	}
}

func getBalance(t *testing.T, db *sqlx.DB) int {
	balance := 0
	if err := db.Get(&balance, "select balance from player where phone_number = '1'"); err != nil {
		t.Error(err)
	}
	return balance
}

func balanceString(balance int) string {
	return fmt.Sprintf("%.2f", float64(balance)/1000)
}
