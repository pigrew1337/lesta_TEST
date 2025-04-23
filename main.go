package main

import (
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type WordData struct {
	Word string
	TF   int
	IDF  float64
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("textfile")
		if err != nil {
			c.String(http.StatusBadRequest, "Ошибка получения файла: %v", err)
			return
		}
		v, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка открытия файла: %v", err)
			return
		}
		defer v.Close()
		content, err := io.ReadAll(v)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка чтения файла: %v", err)
			return
		}
		text := string(content)
		wordData := procText(text)
		sort.Slice(wordData, func(i, j int) bool {
			return wordData[i].IDF > wordData[j].IDF
		})
		if len(wordData) > 50 {
			wordData = wordData[:50]
		}
		c.HTML(http.StatusOK, "result.html", gin.H{"words": wordData})
	})
	if err := router.Run(":8086"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func procText(text string) []WordData {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я'))
	})
	freq := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word)
		if word != "" {
			freq[word]++
		}
	}
	totalWords := len(words)
	var results []WordData
	for w, count := range freq {
		idf := math.Log(float64(totalWords) / float64(count+1))
		results = append(results, WordData{
			Word: w,
			TF:   count,
			IDF:  idf,
		})
	}
	return results
}
