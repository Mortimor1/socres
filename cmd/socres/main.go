package main

import (
	"fmt"
	"github.com/Mortimor1/socres/internal/rkn"
	"github.com/Mortimor1/socres/internal/socres"
	"github.com/Mortimor1/socres/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
	"os"
	"strings"
	"time"
)

const requestPath = "zapros.xml"
const signaturePath = "zapros.xml.sig"
const outputPath = "registry.zip"

var noLoad = false
var host = ""
var username = ""
var password = ""

func main() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	parseArgs()

	if !noLoad {
		rkn.RknLoad(requestPath, signaturePath, outputPath)
		utils.Unzip("registry.zip", "output")
	}

	register := socres.ParseXML("output/register.xml")

	mikrot := socres.Client{
		Address:  host,
		Username: username,
		Password: password,
		Async:    true,
		UseTLS:   false,
	}

	re := mikrot.SendMikrotik("/ip/firewall/address-list/print ?list=social_resources_list")
	resources := checkResources(&register, re)

	commands := parseCommand(resources)
	for _, c := range commands {
		mikrot.SendMikrotik(c)
		time.Sleep(1 * time.Second)
	}
}

func parseArgs() {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--no-load" {
			noLoad = true
		}
		if os.Args[i] == "--host" {
			host = os.Args[i+1]
		}
		if os.Args[i] == "--user" {
			username = os.Args[i+1]
		}
		if os.Args[i] == "--password" {
			password = os.Args[i+1]
		}
		if os.Args[i] == "--help" {
			fmt.Println()
		}
	}
}

func parseCommand(resources *[]socres.Resource) []string {
	var commands []string
	for _, res := range *resources {
		if res.Flag == 1 {
			command := "/ip/firewall/address-list/add =list=social_resources_list =address=" + res.Address + " =comment=" + res.Domain
			commands = append(commands, command)
		}
		if res.Flag == -1 {
			command := "/ip/firewall/address-list/remove =.id=" + res.Id
			commands = append(commands, command)
		}
	}
	return commands
}

func checkResources(register *rkn.RegisterSocResources, re []*proto.Sentence) *[]socres.Resource {
	var resources []socres.Resource

	for _, content := range register.Content {
		for _, address := range content.Subnets {
			res := socres.Resource{
				Domain:  content.Domain,
				Address: address,
				Flag:    0,
			}
			resources = append(resources, res)
		}
	}

	//Проверяем есль ли новые
	for i, res := range resources {
		exist := false
		for _, item := range re {
			if strings.Split(res.Address, "/32")[0] == item.Map["address"] {
				exist = true
			}
		}
		if !exist {
			resources[i].Flag = 1
		}
	}

	// Проверяем что нужно удалить
	for _, item := range re {
		exist := false
		for _, res := range resources {
			if strings.Split(res.Address, "/32")[0] == item.Map["address"] {
				exist = true
			}
		}
		if !exist {
			newRes := socres.Resource{
				Id:      item.Map[".id"],
				Domain:  "",
				Address: item.Map["address"],
				Flag:    -1,
			}
			resources = append(resources, newRes)
		}
	}

	return &resources
}
