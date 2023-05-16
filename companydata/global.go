package companydata

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	lib "github.com/wopta/goworkspace/lib"
)

var config = lib.SftpConfig{
	Username:     os.Getenv("GLOBAL_SFTP_USER"),
	Password:     os.Getenv("GLOBAL_SFTP_PSW"), // required only if password authentication is to be used
	PrivateKey:   "",                           // required only if private key authentication is to be used
	Server:       "ftps.globalassistance.it:222",
	KeyExchanges: []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1"}, // optional
	Timeout:      time.Second * 30,                                                                                            // 0 for not timeout
}

func GlobalSftpDownload(filename string, bucket string, folder string) ([]byte, io.ReadCloser, error) {
	localPath := "../tmp/" + filename
	client, e := lib.NewSftpclient(config)
	client.ListFiles(".")
	println("folder +filename: ", folder+filename)
	println("GlobalSftpDownload error: ", fmt.Errorf("unable to open remote file: %v", e))
	srcFile, e := client.Download(folder + filename)
	if e != nil {
		fmt.Errorf("unable to open remote file: %v", e)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localPath)
	if err != nil {
		fmt.Errorf("unable to open local file: %v", err)
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Errorf("unable to download remote file: %v", err)
	}
	log.Printf("%d bytes copied to %v", bytes, localPath)

	sourceByte, e := io.ReadAll(srcFile)
	log.Println(e)
	//buf := new(bytes.Buffer)
	//_, e = buf.ReadFrom(reader)

	return sourceByte, srcFile, e
}
func GlobalSftpUpload(filename string, folder string) error {
	client, e := lib.NewSftpclient(config)
	println("filename: ", filename)
	defer client.Close()
	log.Println("Open local file for reading.:")
	source, e := os.Open("../tmp/" + filename)
	lib.CheckError(e)
	//defer source.Close()
	log.Println("Create remote file for writing:")
	// Create remote file for writing.
	lib.Files("../tmp")
	//destination, e := client.Create(filePath)
	destination, e := client.Create(folder + filename)

	defer destination.Close()
	log.Println("Upload local file to a remote location as in 1MB (byte) chunks.")
	info, e := source.Stat()
	log.Println(info.Size())
	// Upload local file to a remote location as in 1MB (byte) chunks.
	e = client.Upload(source, destination, int(info.Size()))
	return e
}
func GlobalSftpDelete(filename string) error {
	client, e := lib.NewSftpclient(config)
	println("filename: ", filename)
	defer client.Close()
	log.Println("Open local file for reading.:")
	e = client.Remove(filename)
	defer client.Close()

	return e
}
