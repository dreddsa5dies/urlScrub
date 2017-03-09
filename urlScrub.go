// urlScrab
package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/opesun/goquery"
)

var opts struct {
	FileNameCompany string `short:"o" long:"open" default:"./names.txt" description:"With the names of the companies file"`
	FileFinal       string `short:"f" long:"final" default:"final.csv" description:"The file with the saved information about the companies"`
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

	// создание файла отчета в формате csv
	file, err := os.OpenFile(pwdDir+"/reports/"+opts.FileFinal, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)

	}
	defer file.Close()

	// справка по компаниям
	fileTXT, err := os.OpenFile(pwdDir+"/reports/reports.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)

	}
	defer file.Close()

	// TODO: заголовок привести к нормальному виду после корректировки вывода
	/*
		getFile, err := file.Stat()
		if err != nil {
			log.Fatalln(err)
		}

		if getFile.Size() <= 1 {
			// заголовок
			file.WriteString("Наименование;ФИО директора;Положение директора;Виды деятельности;Дата регистрации;Кол-во сотрудников;ИНН;КПП;ОГРН;ОКПО;Адрес;Сайт;Место в категории;Уставной капитал;Основной заказчик\n")
		}
	*/
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

		// запрос по url
		resp, err := http.Get("https://www.google.ru/search?q=" + search + "+inurl%3Asbis.ru")
		if err != nil {
			log.Fatalln(err)
		}

		// отложенное закрытие коннекта
		defer resp.Body.Close()

		// парсинг ответа
		x, err := goquery.Parse(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// храниение итоговых ссылок
		var urlsSearchs []string

		// формирование нормальной ссылки
		for _, l := range x.Find("h3").HtmlAll() {
			// обрезка html
			j := strings.TrimLeft(l, `<a href="/url?q=`)
			// надо убрать "левый" код в ссылке
			k := strings.Split(j, `&amp;sa=U&amp;ved=`)
			// итоговая ссылка готова
			urlsSearchs = append(urlsSearchs, "h"+k[0])
		}
		log.Printf("Найдены такие ссылки:\n")
		for _, logURL := range urlsSearchs {
			log.Printf("\t%v\n", logURL)
		}

		lenURL := 3
		if len(urlsSearchs) < 3 {
			lenURL = len(urlsSearchs)
		}
		for o := 0; o < lenURL; o++ {
			searchURL(urlsSearchs[o], file, fileTXT)
		}
	}
	log.Println("Готово")
}

func searchURL(url string, file, fileTXT *os.File) {
	x, err := goquery.ParseUrl(url)
	if err != nil {
		log.Fatalf("Ошибка парсинга страницы:\t%v\n", err)
	}
	// Ссылка на сайте
	writeString(url, file)
	// Наименование
	nameCo := strings.Split(x.Find("div.cCard__MainReq-Name").Text(), ", ")
	// только наименование
	writeString(nameCo[1], file)
	// форма организации
	writeString(nameCo[0], file)
	// ФИО директора
	writeString(x.Find("div.cCard__Director-Name").Text(), file)
	// положение директора
	// можно раскидать по количеству компаний еще
	compNum := strings.Split(strings.ToLower(x.Find("div.cCard__Director-Position").Text()), "  ")
	if len(compNum) > 1 {
		writeString(compNum[0], file)
		writeString(compNum[1], file)
	} else {
		writeString(compNum[0], file)
		writeString("Нет данных", file)
	}
	// Основная деятельность
	writeString(x.Find("div.cCard__OKVED cCard__Wide").Text(), file)
	// Адрес
	writeString(x.Find("div.cCard__Contacts-Address").Text(), file)
	// Широта и долгота
	// через API Yandex
	resp, err := http.Get("https://geocode-maps.yandex.ru/1.x/?geocode=" + x.Find("div.cCard__Contacts-Address").Text())
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// координаты точки в ответе
	point, err := regexp.Compile(`<lowerCorner>\d\d\.\d{4,6} \d\d\.\d{4,6}</lowerCorner>`)
	pointsStr := string(point.Find(body))
	// убираем теги
	pointsStr = strings.TrimLeft(pointsStr, "<lowerCorner>")
	pointsStr = strings.TrimRight(pointsStr, "</lowerCorner>")
	// разделяем на широту и долготу
	geoDATA := strings.Split(pointsStr, " ")
	if len(geoDATA) > 1 {
		writeString(geoDATA[0], file)
		writeString(geoDATA[1], file)
	} else {
		writeString(geoDATA[0], file)
		writeString("", file)
	}
	// Контакты
	writeString(x.Find("div.cCard__Contacts-Value").Text(), file)
	// Размер уставного капитала
	writeString(x.Find("div.cCard__Owners-OwnerList-Sum").Text(), file)
	// Сроки действия
	writeString(x.Find("div.cCard__Status-Value").Text(), file)
	// ИНН КПП ОГРН ОКПО
	dataINN := strings.Split(x.Find("div.cCard__MainReq-Right-Req-Line").Text(), "    ")
	// только ИНН
	writeString(strings.TrimLeft(dataINN[0], "ИНН "), file)
	// только КПП
	writeString(strings.TrimLeft(dataINN[1], "КПП "), file)
	// только ОГРН
	writeString(strings.TrimLeft(dataINN[2], "ОГРН "), file)
	// только ОКПО
	writeString(strings.TrimLeft(dataINN[3], "ОКПО "), file)

	// -----------------------------------------------------------
	// вся инфа в текстовую справку
	writeString(x.Find("div.cCard__CompanyDescription").Text(), fileTXT)
	// новая строка
	_, err = file.WriteString("\n")
	if err != nil {
		log.Fatalln(err)
	}
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
