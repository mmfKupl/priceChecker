package checking

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

// Rival contains a main URL and SerachURL
type Rival struct {
	Name             string `json:"name"`
	URL              string `json:"URL"`
	SearchURL        string `json:"searchURL"`
	ProductsSelector string `json:"productsSelector"`
	PriceSelector    string `json:"priceSelector"`
	TitleSelector    string `json:"titleSelector"`
	T                string `json:"t"`
}

// Products contains a couple of the Product
type Products struct {
	XMLName    xml.Name `xml:"products"`
	Amount     int      `xml:"amount,attr"`
	ReqArticle string   `xml:"reqArticle,attr"`
	Data       []*Product
}

// ProductsWithSeveralArticle contains a couple of the Products
type ProductsWithSeveralArticle struct {
	XMLName xml.Name `xml:"articles"`
	Data    []*Products
	Amount  int `xml:"amount,attr"`
}

// Product contains a main info of a product
type Product struct {
	XMLName    xml.Name `xml:"product"`
	Title      string   `xml:"title"`
	ReqArticle string   `xml:"-"`
	Price      string   `xml:"price,attr"`
	URL        string   `xml:"url"`
	RivalName  string   `xml:"rivalName,attr"`
}

// ToString return format string of product
func (p *Product) ToString() string {
	return fmt.Sprintf("%s (%s). %s", p.Title, p.Price, p.URL)
}

// Links to checking websites
// const (
// 	InstBy        = "https://instr.by/"
// 	InstrumentBy  = "https://instrument.by/"
// 	DewaltBDBy    = "http://dewalt-bd.by/"
// 	ToolsBy       = "http://www.tools.by/"
// 	TProBy        = "https://tpro.by/"
// 	By7745        = "https://7745.by/"
// 	Vek21By       = "https://www.21vek.by/"
// 	DelomasteraBy = "https://delomastera.by/"
// 	OmaBy         = "https://www.oma.by/"
// )

// GetProductInfo return Product corresponding to the recived article
func (r *Rival) GetProductInfo(article string) (*Product, error) {
	searchURL := fmt.Sprintf(r.SearchURL, article)
	// fmt.Println(searchURL)
	page, err := downloadPage(searchURL)
	if err != nil {
		return nil, err
	}
	nodes := page.Find(r.ProductsSelector)
	if nodes.Size() == 0 {
		return nil, fmt.Errorf("Подходящего товара не найдено")
	}
	nodesArr := make([]*goquery.Selection, 0, nodes.Size())
	nodes = nodes.Each(func(i int, s *goquery.Selection) {
		title := s.Find(r.TitleSelector).Text()
		rg := fmt.Sprintf(`%s(\s|\z|\))`, article)
		matched, err := regexp.Match(rg, []byte(title))
		if err != nil {
			fmt.Println(err)
		}
		if matched {
			nodesArr = append([]*goquery.Selection{s}, nodesArr...)
		}
		nodesArr = append(nodesArr, s)
	})
	product := new(Product)
	prodNode := nodesArr[0]

	product.Title = prodNode.Find(r.TitleSelector).Text()
	if product.Title == "" || product.Title == " " {
		product.Title = "Заголовок отцувствует."
	}
	if prURL, ok := prodNode.Find(r.TitleSelector).Attr("href"); ok {
		if strings.Contains(prURL, "http") {
			product.URL = prURL
		} else {
			product.URL = r.URL + prURL
		}
	} else {
		product.URL = "none"
	}
	product.ReqArticle = article
	product.Price = strings.ReplaceAll(prodNode.Find(r.PriceSelector).Text(), "\n", " ")
	product.Price = strings.Trim(product.Price, " \n\t")
	product.RivalName = r.Name
	return product, nil
}

func downloadPage(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке данных сайта!: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка при загрузке данных сайта!: %d %s", res.StatusCode, res.Status)
	}

	utf8, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("Ошибка шифровки: %s", err.Error())
	}

	page, err := goquery.NewDocumentFromReader(utf8)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при загрузке данных сайта!: %s", err.Error())
	}

	return page, nil
}
