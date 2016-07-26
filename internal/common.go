package internal

import (
	"bufio"
	"github.com/sky-uk/gonsx"
	"io"
	"log"
	"os"
	"os/exec"
)

var nsxClient *gonsx.NSXClient

// GetNSXClient : method to retrieve the nsxclient singleton
func GetNSXClient() *gonsx.NSXClient {
	if nsxClient == nil {
		nsxURL := os.Getenv("NSX_URL")
		nsxUser := os.Getenv("NSX_USER")
		nsxPassword := os.Getenv("NSX_PASSWORD")
		if nsxURL == "" || nsxUser == "" || nsxPassword == "" {
			panic("either NSX_URL, NSX_USER or NSX_PASSWORD environment variables are empty!")
		}
		nsxClient = gonsx.NewNSXClient(nsxURL, nsxUser, nsxPassword, true, false)
	}
	return nsxClient
}

// CheckError : Print a FATAL log if error is != null
func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

// ExecuteCommand : Executes a shell command and returns stdOut and stdErr as byte[] and exitError
func ExecuteCommand(command string, arg ...string) ([]byte, []byte, *exec.ExitError) {
	cmd := exec.Command(command, arg...)
	stdout, err := cmd.StdoutPipe()
	CheckError(err)
	stderr, err := cmd.StderrPipe()
	CheckError(err)
	err = cmd.Start()
	CheckError(err)

	stdOutBytes := ReadAll(bufio.NewReader(stdout))
	stdErrBytes := ReadAll(bufio.NewReader(stderr))

	err = cmd.Wait()
	var exitError *exec.ExitError
	if err != nil {
		exitError = err.(*exec.ExitError)
	}

	return stdOutBytes, stdErrBytes, exitError
}

// ReadAll : Reads all data from a Reader and returns a byte[]
func ReadAll(r *bufio.Reader) []byte {
	content := []byte{}
	nBytes, nChunks := int64(0), int64(0)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))
		// process buf
		content = append(content, buf...)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	return content
}
