package socres

import (
	"encoding/xml"
	"github.com/Mortimor1/socres/internal/rkn"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func ParseXML(filePath string) rkn.RegisterSocResources {
	xmlFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Info(err)
	}

	log.Info("Successfully Opened register.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var register rkn.RegisterSocResources
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &register)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	log.Info(register)
	return register
}
