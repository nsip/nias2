package naprr

import (
	"archive/zip"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nuid"
	"io"
	"log"
	"os"
)

// connection to the Streaming NATS server
func CreateSTANConnection() stan.Conn {

	clusterID := "nap-rr"
	clientID := nuid.Next()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:5222"))
	if err != nil {
		log.Fatalln("Cannot connect to STAN server: ", err)
	}

	return sc

}

func isZipFile(fname string) bool {

	xmlZipFile, err := zip.OpenReader(fname)
	if err != nil {
		return false
	}
	defer xmlZipFile.Close()

	return true

}

func openDataFileZip(fname string) (io.Reader, error) {

	xmlZipFile, err := zip.OpenReader(fname)
	if err != nil {
		log.Println("Unable to open zip file: ", err)
		return nil, err
	}
	// assume only one file in the archive
	xmlFile, err := xmlZipFile.File[0].Open()
	if err != nil {
		return xmlFile, err
	}

	return xmlFile, nil

}

func openDataFile(fname string) (io.Reader, error) {

	xmlFile, err := os.Open(fname)
	if err != nil {
		return xmlFile, err
	}
	return xmlFile, nil

}

func OpenResultsFile(fname string) (io.Reader, error) {
	var xmlFile io.Reader
	var ferr error

	if isZipFile(fname) {
		xmlFile, ferr = openDataFileZip(fname)
	} else {
		xmlFile, ferr = openDataFile(fname)
	}
	if ferr != nil {
		return xmlFile, ferr
	}

	return xmlFile, ferr

}
