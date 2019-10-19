package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	flags "github.com/jessevdk/go-flags"
	goq "github.com/opesun/goquery"
)

var opts struct {
	FileNameCompany string `short:"o" long:"open" default:"./names.txt" description:"With the names of the companies file"`
}

func main() {
	// разбор флагов
	flags.Parse(&opts)

	// в какой папке исполняемы файл
	pwdDir, _ := os.Getwd()

	// создание файла log для записи ошибок
	fLog, err := os.OpenFile(pwdDir+`/.log`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatalln(err)
	}
	defer fLog.Close()

	// запись ошибок и инфы в файл
	log.SetOutput(fLog)

	// создание папки с отчетами
	os.Mkdir(pwdDir+"/reports", 0755)

	// справка по компаниям
	fileTXT, err := os.OpenFile(pwdDir+"/reports/reports.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)

	}
	defer fileTXT.Close()

	// разобрать названия компаний для перебора
	var massName []string
	fileOpen, err := os.Open(opts.FileNameCompany)
	if err != nil {
		log.Fatalln(err)
	}
	// построчное считывание
	scanner := bufio.NewScanner(fileOpen)
	for scanner.Scan() {
		massName = append(massName, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	defer fileOpen.Close()

	for i := 0; i < len(massName); i++ {
		search := massName[i]
		log.Printf("Ищу данные по компании:\t%v\n", search)

		// Request the HTML page.
		doc, err := goquery.NewDocument("https://www.google.ru/search?q=" + search + "+inurl%3Asbis.ru")

		if err != nil {
			log.Fatal(err)
		}

		// храниение итоговых ссылок
		var urlsSearchs []string

		doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
			href, _ := item.Attr("href")
			if strings.Contains(href, "https://sbis.ru/") {
				// обрезка html
				j := strings.TrimLeft(href, `/url?q=`)
				if strings.HasPrefix(j, "https://sbis.ru/") {
					urlsSearchs = append(urlsSearchs, j)
				}
			}
		})

		log.Printf("Найдены такие ссылки:\n")
		for _, logURL := range urlsSearchs {
			log.Printf("\t%v\n", logURL)
		}

		lenURL := 3
		if len(urlsSearchs) < 3 {
			lenURL = len(urlsSearchs)
		}
		for o := 0; o < lenURL; o++ {
			searchURL(urlsSearchs[o], fileTXT)
		}
	}
	log.Println("Готово")
}

func searchURL(url string, fileTXT *os.File) {
	x, err := goq.ParseUrl(url)
	if err != nil {
		log.Fatalf("Ошибка парсинга страницы:\t%v\n", err)
	}
	// Наименование
	nameCo := strings.Split(x.Find("div.cCard__MainReq-Name").Text(), ", ")
	// -----------------------------------------------------------
	// вся инфа в текстовую справку
	writeString(strings.TrimLeft(x.Find("div.cCard__CompanyDescription").Text(), "Краткая справка"), fileTXT)

	_, err = fileTXT.WriteString("\n-------------------------------------------\n")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Обработал по данным на:\t%v\n", nameCo[0])
}

// запись в файл строки
func writeString(x string, file *os.File) {
	// tab в виде разделителя
	_, err := file.WriteString(x + "	")
	if err != nil {
		log.Fatalln(err)
	}
}
