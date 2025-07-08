package server

import (
	"encoding/xml"
	"net/http"
)

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

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

func writeXML(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/xml; chatset=utf-8")
	w.WriteHeader(status)
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(v); err != nil {
		http.Error(w, "XML encoding error", http.StatusInternalServerError)
	}
}
