package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", updateFaveFoods)
	log.Printf("Server starting on port 3000")
	http.ListenAndServe(":3000", nil)
}

func updateFaveFoods(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:5500")
	foodCh := make(chan string)

	go spitOutFoods(r.Context(), foodCh)

	for food := range foodCh {
		event, err := formatSSE(food)
		if err != nil {
			fmt.Println(err)
			break
		}

		_, err = fmt.Fprint(w, event)
		if err != nil {
			fmt.Println(err)
			break
		}

		flusher.Flush()
	}
}

func spitOutFoods(ctx context.Context, foodChan chan<- string) {
	ticker := time.NewTicker(time.Second)
	foods :=  []string{
		"Jollof rice",
		"Fried rice",
		"Fufu and Eguisi soup",
		"Amala and Ewedu",
	}

	outerloop:
		for {
			select {
			case <-ctx.Done():
				break outerloop
			case <-ticker.C:
				food := foods[rand.Intn(len(foods))]
				foodChan <- food
			}			
		}

		ticker.Stop()
		close(foodChan)
}

func formatSSE(data string) (string, error) {
	m := map[string]string{
		"data": data,
	}

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(m)

	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}