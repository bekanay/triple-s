package server

import "encoding/xml"

type bucketXML struct {
	XMLName          xml.Name `xml:"Bucket"`
	Name             string   `xml:"Name"`
	CreationTime     string   `xml:"CreationDate"`
	LastModifiedTime string   `xml:"LastModifiedDate"`
}

type listBucketsResultXML struct {
	XMLName xml.Name    `xml:"ListAllMyBucketsResult"`
	Buckets []bucketXML `xml:"Buckets>Bucket"`
}
