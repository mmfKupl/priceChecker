package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	hFlag        = flag.Bool("h", false, "a bool")
	helpFlag     = flag.Bool("help", false, "a bool")
	pathFlag     = flag.String("path", "", "a string")
	nameFlag     = flag.String("name", "", "a string")
	articlesFlag = flag.String("articles", "", "a string")
	rivalsFlag   = flag.String("rivals", "", "a string")
)

var rivals []Rival

var articles []string

func init() {
	flag.Parse()

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

	if *hFlag || *helpFlag {
		showHelp()
		os.Exit(0)
	}

	if *articlesFlag != "" {
		articles = strings.Split(*articlesFlag, ",")
	}

	if *rivalsFlag != "" {
		tmpRivals := strings.Split(*rivalsFlag, ",")
		newRivals := make([]Rival, 0, len(tmpRivals))
		for _, val := range tmpRivals {
			i, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("init: ", err)
				continue
			}
			newRivals = append(newRivals, rivals[i])
		}
		rivals = newRivals
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
	var productsWA = new(ProductsWithSeveralArticle)
	productsWA.Data = make([]*Products, 0, len(rivals)+len(articles))
	for _, article := range articles {
		var products = new(Products)
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
	fmt.Println("0 - выбрать всех")
	for i, rival := range rivals {
		fmt.Printf("%v - %s\n", i+1, rival.Name)
	}
	fmt.Printf("%v - %s\n", len(rivals)+1, "подтвердить выбор")
	fmt.Println("")
	rivals = chooseRivals()
	fmt.Printf("Конкуренты, которые просматриваются (%v):\n", len(rivals))
	for _, rival := range rivals {
		fmt.Print(rival.Name + " ")
	}
	fmt.Println("")
	var inputData string
	var products = new(Products)
	for {
		products.Data = make([]*Product, 0, len(rivals))
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

func chooseRivals() []Rival {
	var inputData string
	newRivals := make([]Rival, 0)
	for {
		fmt.Print("-> ")
		fmt.Scan(&inputData)
		if inputData == "0" {
			return rivals
		}

		i, err := strconv.Atoi(inputData)
		if err != nil {
			fmt.Println("choose rivals: ", err)
			continue
		}
		if i == len(rivals)+1 {
			if len(newRivals) == 0 {
				return rivals
			}
			return newRivals
		}
		if i <= 0 || i > len(rivals) {
			fmt.Println("Неверное значение!")
			continue
		}
		if contains(newRivals, rivals[i-1]) {
			continue
		}
		newRivals = append(newRivals, rivals[i-1])
		for _, val := range newRivals {
			fmt.Print(val.Name + " ")
		}
		fmt.Println("")
	}
}

func contains(arr []Rival, val Rival) bool {
	for _, _val := range arr {
		if _val.Name == val.Name {
			return true
		}
	}
	return false
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
	(перечислять в виде -articles="article1,article2,article3") (при отцуствии флага будет запрошен ввод артикула по ходу работы)`)

	fmt.Println(`
	-rivals - определяет какие магазины будут просмотрены
	(перечислять в виде -rivals="1,2,3", если ничего не введено, то будут выбраны все)`)
	for i, val := range rivals {
		fmt.Printf("\t%v - %s\n", i, val.Name)
	}
}
