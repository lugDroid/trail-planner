package model

import "encoding/xml"

type Gpx struct {
	XMLName        xml.Name `xml:"gpx"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Creator        string   `xml:"creator,attr"`
	Version        string   `xml:"version,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Rte            struct {
		Text       string `xml:",chardata"`
		Name       string `xml:"name"`
		Desc       string `xml:"desc"`
		Extensions struct {
			Text string `xml:",chardata"`
			Line struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
				Color string `xml:"color"`
			} `xml:"line"`
		} `xml:"extensions"`
		Rtept []struct {
			Text string `xml:",chardata"`
			Lat  string `xml:"lat,attr"`
			Lon  string `xml:"lon,attr"`
			Ele  string `xml:"ele"`
			Time string `xml:"time"`
		} `xml:"rtept"`
	} `xml:"rte"`
}
