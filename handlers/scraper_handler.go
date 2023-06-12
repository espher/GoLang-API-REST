package handlers

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScraperHandler struct {
	db *gorm.DB
}

func ScraperHandlerRouter(db *gorm.DB) *ScraperHandler {
	return &ScraperHandler{db: db}
}

type ProductDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func (sh *ScraperHandler) DoScrapeProduct(c *gin.Context) {
	type URLBody struct {
		URL string
	}

	var requestBodyURL URLBody
	c.BindJSON(&requestBodyURL)

	productDetails, err := getProductDetail(requestBodyURL.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, productDetails)
}

func getProductDetail(url string) (*ProductDetails, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find("h1").Text()
	description := doc.Find(".x-item-condition-desc").Text()
	price := doc.Find(".x-price-primary").Text()

	productDetails := &ProductDetails{
		Title:       title,
		Description: description,
		Price:       price,
	}

	return productDetails, nil
}
