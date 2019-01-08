package xrp
import (
	"bytes"
	"log"
	"os/exec"
)

func DoCurlRequest (data string) string {
	var err error
	cmd := exec.Command("/bin/sh")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in
	var out bytes.Buffer
	cmd.Stdout = &out
	go func() {
		in.WriteString("curl -X POST -H \"Content-Type\":application/json --data '" + data + "' http://127.0.0.1:8551")
	}()
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	return out.String()
}
