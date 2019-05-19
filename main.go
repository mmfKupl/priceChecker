package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"./checking"
)

var (
	hFlag        = flag.Bool("h", false, "a bool")
	helpFlag     = flag.Bool("help", false, "a bool")
	pathFlag     = flag.String("path", "", "a string")
	nameFlag     = flag.String("name", "", "a string")
	articlesFlag = flag.String("articles", "", "a string")
)

var rivals []checking.Rival

var articles []string

func init() {
	flag.Parse()

	if *hFlag || *helpFlag {
		showHelp()
		os.Exit(0)
	}

	rivalsPath := "./rivals.json"
	fileData, err := ioutil.ReadFile(rivalsPath)
	if err != nil {
		fmt.Println("Ошибка при открытии json файла с конкурентами!")
		panic(err)
	}
	err = json.Unmarshal(fileData, &rivals)
	if err != nil {
		panic(err)
	}

	if *articlesFlag != "" {
		articles = strings.Split(*articlesFlag, ",")
	}
	fmt.Println(articles)

}

func main() {
	fmt.Println("Добпро пожаловать в PriceChecking!")
	if len(articles) != 0 {
		argsStart()
	} else {
		comonStart()
	}
}

func argsStart() {
	fmt.Println("Будут скачаны следующин артикулы: ")
	for _, article := range articles {
		fmt.Printf("%s ", article)
	}
	fmt.Println("")
	var productsWA = new(checking.ProductsWithSeveralArticle)
	productsWA.Data = make([]*checking.Products, 0, len(rivals)+len(articles))
	for _, article := range articles {
		var products = new(checking.Products)
		for _, rival := range rivals {
			product, err := rival.GetProductInfo(article)
			if err != nil {
				fmt.Printf("%s:%s\t%s\n", rival.Name, rival.T, err.Error())
				continue
			}
			products.Data = append(products.Data, product)
			fmt.Printf("%s:%s\t%s\n", rival.Name, rival.T, product.ToString())
		}
		products.Amount = len(products.Data)
		products.ReqArticle = article
		productsWA.Data = append(productsWA.Data, products)
		productsWA.Amount++
	}
	err := writeToFile(productsWA, articles[0], ".")
	if err != nil {
		fmt.Println(err)
	}
}

func comonStart() {
	fmt.Printf("Конкуренты, которые просматриваются (%v):\n", len(rivals))
	for _, rival := range rivals {
		fmt.Print(rival.Name + " ")
	}
	fmt.Println("")
	var inputData string
	var products = new(checking.Products)
	for {
		products.Data = make([]*checking.Product, 0, len(rivals))
		fmt.Print("Введите артикул для поиска (-1 для выхода):\n-> ")
		fmt.Scan(&inputData)
		if inputData == "-1" {
			return
		}
		for _, rival := range rivals {
			product, err := rival.GetProductInfo(inputData)
			if err != nil {
				fmt.Printf("%s:%s\t%s\n", rival.Name, rival.T, err.Error())
				continue
			}
			products.Data = append(products.Data, product)

			fmt.Printf("%s:%s\t%s\n", rival.Name, rival.T, product.ToString())
			// fmt.Printf("%+v\n", product)
		}
		products.Amount = len(products.Data)
		products.ReqArticle = inputData
		err := writeToFile(products, inputData, ".")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func writeToFile(products interface{}, filename string, path string) error {
	if *nameFlag != "" {
		filename = *nameFlag
	}
	if *pathFlag != "" {
		path = *pathFlag
	}

	filename = fmt.Sprintf("%s/%s.xml", path, filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	xmlProd, err := xml.MarshalIndent(products, " ", "	")
	if err != nil {
		return err
	}
	file.Write(xmlProd)
	fmt.Printf("Данные успешно записаны в файл %s.\n", filename)
	return nil
}

func showHelp() {
	fmt.Println("Для вызова справки используются флаги -h или -help")
	fmt.Println(
		`
Флаги для запуска: 
	-path - определяет директорию, в которую будет записан результат работы (если флаг отцуствует, результат будет записан в директорию, где было вызвано приложение)
	-name - определяет название файла, в который будет записан результат (по умолчанию файл имеет название первого артикула, расширение файла xml)
	-articles - определяет перечисление артикулов, которые будут обработаны 
	(перечислять в виде -articles="article1,article2,article3") (при отцуствии флага будет запрошен ввод артикула по ходу работы)
	`)
}
