package rkn

import "encoding/xml"

type RegisterSocResources struct {
	XMLName xml.Name `xml:"registerSocResources"`
	UpdateTime string `xml:"updateTime,attr"`
	FormatVersion string `xml:"formatVersion,attr"`
	Content []Content `xml:"content"`
}

type Content struct {
	XMLName xml.Name `xml:"content"`
	Hash string `xml:"hash,attr"`
	IncludeTime string `xml:"includeTime,attr"`
	Id int `xml:"id,attr"`
	ResourceName string `xml:"resourceName"`
	Domain string `xml:"domain"`
	Subnets []string `xml:"ipSubnet"`
}
