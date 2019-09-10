package d2

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var Logger = log.New(ioutil.Discard, "[D2] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC)

func EnableD2Logging() {
	Logger.SetOutput(os.Stderr)
}

func ClustersPath(cluster string) string {
	return filepath.Join("/", "d2", "clusters", cluster)
}

func ServicesPath(service string) string {
	return filepath.Join("/", "d2", "services", service)
}

func UrisPath(uri string) string {
	return filepath.Join("/", "d2", "uris", uri)
}

