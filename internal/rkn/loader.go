package rkn

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/Mortimor1/socres/internal/wsdl"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

func RknLoad(request string, signed string, output string) {
	basicAuth := wsdl.BasicAuth{
		Login:    "",
		Password: "",
	}

	lastUpdate := read()

	operator := wsdl.NewOperatorRequestPortType("http://vigruzki.rkn.gov.ru/services/OperatorRequest/?wsdl", true, &basicAuth)
	response := getDumpDateEx(operator)

	lastDate := getLastDate(response)
	dumpFormat := getDumpFormat(response)

	tm := time.Unix(lastDate/1000, 0)
	log.Infof("Dump last time: %s", tm)

	if lastUpdate == tm.String() {
		log.Info("Данные в реестре не изменились.")
		return
	}
	write(tm.String())

	log.Infof("Dump format: %s", dumpFormat)

	code := sendRequest(operator, request, signed, dumpFormat)

	result := false

	for i := 1; i < 5; i++ {
		log.Info("Выгрузка еще не готова, ждем еще минуту.")
		time.Sleep(1 * time.Minute)
		if getResult(operator, code) {
			result = true
			break
		}
	}

	if result != true {
		log.Fatal("ОШИБКА: Не удалось получть выгрузку")
	}

	registry := getResultSocResources(operator, code)
	decodeBase64(registry, output)
	log.Info("УСПЕХ: Выгрузка завершена")
}

func read() string {
	filerc, err := os.Open("lastDate.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer filerc.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(filerc)
	contents := buf.String()
	return contents
}

func write(lastDate string) {
	f, err := os.Create("lastDate.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(lastDate)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func getDumpDateEx(operator *wsdl.OperatorRequestPortType) *wsdl.GetLastDumpDateExResponse {
	dump := wsdl.GetLastDumpDateEx{}
	response, err := operator.GetLastDumpDateEx(&dump)

	if err != nil {
		log.Info(err)
	}
	return response
}

func getLastDate(response *wsdl.GetLastDumpDateExResponse) int64 {
	return response.LastDumpDateSocResources
}

func getDumpFormat(response *wsdl.GetLastDumpDateExResponse) string {
	return response.DumpFormatVersion
}

func sendRequest(operator *wsdl.OperatorRequestPortType, requestPath string, signaturePath string, dumpFormat string) string {
	requestFile := encodeBase64(requestPath)
	signatureFile := encodeBase64(signaturePath)

	request := wsdl.SendRequest{
		RequestFile:       *requestFile,
		SignatureFile:     *signatureFile,
		DumpFormatVersion: dumpFormat,
	}
	response, err := operator.SendRequest(&request)

	if err != nil {
		log.Error(err)
	}

	if response.Result != true {
		log.Fatal(response.ResultComment)
	}

	return response.Code
}

func getResult(operator *wsdl.OperatorRequestPortType, code string) bool {
	request := wsdl.GetResult{
		Code: code,
	}

	response, err := operator.GetResult(&request)
	if err != nil {
		log.Warn(err)
	}
	log.Info(response.ResultComment)
	return response.Result
}

func getResultSocResources(operator *wsdl.OperatorRequestPortType, code string) []byte {
	request := wsdl.GetResultSocResources{
		Code: code,
	}

	response, err := operator.GetResultSocResources(&request)
	if err != nil {
		log.Info(err)
	}
	if response.Result != true {
		log.Fatal(response.ResultComment)
	}

	log.Info(response.ResultComment)
	return response.RegisterZipArchive
}

func encodeBase64(path string) *[]byte {
	file, _ := os.Open(path)
	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	array := []byte(encoded)
	return &array
}

func decodeBase64(array []byte, path string) {
	dec, err := base64.StdEncoding.DecodeString(string(array))
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}
