package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/reneManqueros/httpclient"
	"gorm.io/gorm"
)

type ScraperHandler struct {
	db *gorm.DB
}

func ScraperHandlerRouter(db *gorm.DB) *ScraperHandler {
	return &ScraperHandler{db: db}
}

type URLBody struct {
	URL string
}
type ProductDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type CarBody struct {
	Origin string
	URL    string
}

type CarsDetails struct {
	Title string `json:"title"`
	Vin   string `json:"vin"`
	Price string `json:"price"`
}

func (sh *ScraperHandler) DoScrapeProduct(c *gin.Context) {
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

func (sh *ScraperHandler) DoScrapeCar(c *gin.Context) {
	var requestBodyURL CarBody
	c.BindJSON(&requestBodyURL)

	response, err := getCarsList(requestBodyURL.Origin, requestBodyURL.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	dataSlice := []CarsDetails{}
	doc.Find("h4 a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			href = strings.Trim(href, `"\\\"`)
			href = strings.Trim(href, `"\\\"`)
			data, err := getCarsDetails(href)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			dataSlice = append(dataSlice, *data)
		}
	})

	c.JSON(http.StatusOK, dataSlice)
}

func getCarsList(origin string, url string) (string, error) {
	req := httpclient.Request{
		URL:     url,
		Verb:    "GET",
		Timeout: (8 * time.Second),
		Headers: [][]string{
			{"Accept-Encoding", " gzip, deflate"},
			{"Accept", "*/*"},
			{"Accept-Language", " en-US,en;q=0.9,es;q=0.8"},
			{"Cache-Control", " no-cache"},
			{"Pragma", " no-cache"},
			{"Sec-Fetch-Dest", " empty"},
			{"Sec-Fetch-Mode", " cors"},
			{"Sec-Fetch-Site", " cross-site"},
			{"User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
			{"Sec-Ch-Ua", "^^Not.A/Brand^^;v=^^8^^, ^^Chromium^^;v=^^114^^, ^^Google"},
			{"Sec-Ch-Ua-Mobile", " ?0"},
			{"Connection", " close"},
		},
	}

	response, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}

	return response, nil
}

func getCarsDetails(url string) (*CarsDetails, error) {
	last_url := "?&uri=view%2FconsumerBlock%3FlinkPath%3D%2Fmain%26fields%3Dhtml%2Cscripts%2Cstyles%2Cjsimports%2CstyleClasses&handler=blockProxyHandler&format=deferred&workflowType=block-component&&cfCacheType=&respondBlockError=true&signature=356087658&siteVersion=202237c918ec4c99077530fd1593d39960bc5ba0cfd992af76da47da43134241_1"

	req := httpclient.Request{
		URL:     url + last_url,
		Verb:    "GET",
		Timeout: (8 * time.Second),
		Headers: [][]string{
			{"Accept-Encoding", " gzip, deflate"},
			{"Accept", "*/*"},
			{"Accept-Language", " en-US,en;q=0.9,es;q=0.8"},
			{"Cache-Control", " no-cache"},
			{"Pragma", " no-cache"},
			{"Sec-Fetch-Dest", " empty"},
			{"Sec-Fetch-Mode", " cors"},
			{"Sec-Fetch-Site", " cross-site"},
			{"User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
			{"Sec-Ch-Ua", "^^Not.A/Brand^^;v=^^8^^, ^^Chromium^^;v=^^114^^, ^^Google"},
			{"Sec-Ch-Ua-Mobile", " ?0"},
			{"Connection", " close"},
		},
	}

	response_data, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	type Response struct {
		HTML string `json:"html"`
	}

	var response Response
	err = json.Unmarshal([]byte(response_data), &response)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.HTML))
	if err != nil {
		return nil, err
	}

	title := doc.Find("h1").Text()
	vin := doc.Find("ul.vehicleIdentitySpecs li").Last().Text()
	price := doc.Find(".value").First().Text()

	println(title)
	println(vin)
	println(price)

	carsDetails := &CarsDetails{
		Title: strings.Trim(title, `\n`),
		Vin:   vin,
		Price: price,
	}

	return carsDetails, nil
}
