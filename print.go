// urlScrab
package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/opesun/goquery"
)

func main() {
	pwdDir, _ := os.Getwd()
	// создание файла log
	fLog, err := os.OpenFile(pwdDir+`/log.txt`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatalln(err)
	}
	// запись в err в log и консоль
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	defer fLog.Close()

	// создание файла отчета
	file, err := os.OpenFile(pwdDir+`/new.csv`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	defer file.Close()

	getFile, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	if getFile.Size() <= 1 {
		// заголовок
		file.WriteString("Наименование;ФИО директора;Положение директора;Виды деятельности;Дата регистрации;Кол-во сотрудников;ИНН;КПП;ОГРН;ОКПО;Адрес;Сайт;Место в категории;Уставной капитал;Основной заказчик\n")
	}

	// разобрать названия компаний для перебора
	var massName []string
	fileOpen, err := os.Open(pwdDir + `/1.txt`)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	// построчное считывание
	scanner := bufio.NewScanner(fileOpen)
	for scanner.Scan() {
		massName = append(massName, scanner.Text())
		log.Println(scanner.Text())
		log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	defer fileOpen.Close()

	for i := 0; i < len(massName); i++ {
		search := massName[i]

		// запрос по url
		resp, err := http.Get("https://www.google.ru/search?q=" + search + "+inurl%3Asbis.ru")
		if err != nil {
			log.Fatalln(err)
		}
		log.SetOutput(io.MultiWriter(fLog, os.Stdout))
		// отложенное закрытие коннекта
		defer resp.Body.Close()

		// парсинг ответа
		x, err := goquery.Parse(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetOutput(io.MultiWriter(fLog, os.Stdout))

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
			log.Println("h" + k[0])
			log.SetOutput(io.MultiWriter(fLog, os.Stdout))
		}

		lenUrl := 3
		if len(urlsSearchs) < 3 {
			lenUrl = len(urlsSearchs)
		}
		for o := 0; o < lenUrl; o++ {
			searchUrl(urlsSearchs[o], file)
			log.Println(urlsSearchs[o])
			log.SetOutput(io.MultiWriter(fLog, os.Stdout))
		}
	}
	log.Println("Готово")
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
}

func searchUrl(url string, file *os.File) {
	x, err := goquery.ParseUrl(url)
	if err == nil {
		// обработать для записи
		massData := strings.Split(x.Find(".content").Text(), "   ")
		for j := 0; j < len(massData)-1; j++ {
			massData[j] = strings.Trim(massData[j], " ")
		}
		log.Println(massData)

		// запись строки в файл (добавление)
		// не совсем корректно, требуется фильтрация контента
		if len(massData) > 1 {
			_, err := file.WriteString(massData[7] + ";" + massData[9] + ";" + massData[10] + ";" + massData[13] + ";" + massData[15] + ";" + massData[17] + ";" + massData[19] + ";" + massData[20] + ";" + massData[21] + ";" + massData[22] + ";" + massData[32] + ";" + massData[35] + ";" + massData[54] + ";" + massData[59] + ";" + massData[101] + "\n")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	log.Println(err)
}
