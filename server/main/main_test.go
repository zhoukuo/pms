package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

var baseurl = "http://127.0.0.1:8088/"

func TestIndex(t *testing.T) {

	url := baseurl
	expected := "Welcome to PMS!"

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if string(body) != expected {
		t.Errorf("Index failed. got %s, expected %s", string(body), expected)
	}

}
