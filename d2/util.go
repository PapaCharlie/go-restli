package d2

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

var Logger = log.New(ioutil.Discard, "[D2] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC)

func EnableD2Logging() {
	Logger.SetOutput(os.Stderr)
}

func ClustersPath(cluster string) string {
	return path.Join("/", "d2", "clusters", cluster)
}

func ServicesPath(service string) string {
	return path.Join("/", "d2", "services", service)
}

func UrisPath(clusterName string) string {
	return path.Join("/", "d2", "uris", clusterName)
}
