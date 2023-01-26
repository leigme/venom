package tool

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"runtime"
	"strings"
)

type Vmd struct {
	vmdHeads []string
	pathDir  string
	out      func(data []byte) string
}

func NewVmd() *Vmd {
	vmd := Vmd{}
	if strings.EqualFold(runtime.GOOS, windowsOs) {
		vmd.vmdHeads = []string{windowsCmd, "/C"}
		vmd.pathDir = windowsCd
		vmd.out = convertGB18030
	} else {
		vmd.vmdHeads = []string{unixBash, "-c"}
		vmd.pathDir = unixPwd
		vmd.out = func(data []byte) string {
			return string(data)
		}
	}
	return &vmd
}

func (v *Vmd) Execute(command string) string {
	cmd := exec.Command(v.vmdHeads[0], v.vmdHeads[1], command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("execute cmd: %s, is error: %s", command, err.Error())
	}
	return v.out(output)
}

func (v *Vmd) Pwd() string {
	return v.Execute(v.pathDir)
}

func convertGB18030(bytes []byte) string {
	var decodeBytes, err = simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
	if err != nil {
		fmt.Println(err)
	}
	return string(decodeBytes)
}
