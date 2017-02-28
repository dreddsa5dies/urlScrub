// urlScrab
package main

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/opesun/goquery"
)

func main() {
	pwdDir, _ := os.Getwd()
	// создание файла log
	fLog, err := os.OpenFile(pwdDir+`\log.txt`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatalln(err)
	}
	// запись в err в log и консоль
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	defer fLog.Close()

	// создание файла отчета
	file, err := os.OpenFile(pwdDir+`\new.csv`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
		file.WriteString("Наименование;ФИО директора;Положение директора;Виды деятельности;Дата регистрации;Кол-во сотрудников;ИНН;КПП;ОГРН;ОКПО;Адрес;Сайт;Место в категории;Уставной капитал;Основной заказчик\r\n")
	}

	// разобрать url'ы для перебора
	massUrls := readCsv(pwdDir+`\url.csv`, fLog)
	for i := 0; i < len(massUrls); i++ {
		// [i][1] - пропуск номера по порядку
		x, err := goquery.ParseUrl(massUrls[i][1])
		if err == nil {
			// обработать для записи
			massData := strings.Split(x.Find(".content").Text(), "   ")
			for j := 0; j < len(massData)-1; j++ {
				massData[j] = strings.Trim(massData[j], " ")
			}

			// запись строки в файл (добавление)
			if len(massData) > 1 {
				_, err := file.WriteString(massData[7] + ";" + massData[9] + ";" + massData[10] + ";" + massData[13] + ";" + massData[15] + ";" + massData[17] + ";" + massData[19] + ";" + massData[20] + ";" + massData[21] + ";" + massData[22] + ";" + massData[32] + ";" + massData[35] + ";" + massData[54] + ";" + massData[59] + ";" + massData[101] + "\r\n")
				log.Fatalln(err)
				log.SetOutput(io.MultiWriter(fLog, os.Stdout))
			}
		}
		log.Println(err)
		log.SetOutput(io.MultiWriter(fLog, os.Stdout))
	}
	log.Println("Готово")
	log.SetOutput(io.MultiWriter(fLog, os.Stdout))
}

// считывание csv для разбора url
func readCsv(addr string, fileLog *os.File) [][]string {
	// read file
	dat, err := ioutil.ReadFile(addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(fileLog, os.Stdout))

	in := string(dat)

	// encoding csv
	r := csv.NewReader(strings.NewReader(in))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(fileLog, os.Stdout))

	return records
}
