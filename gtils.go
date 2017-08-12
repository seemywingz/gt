package gtils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"
)

// MultiDimentionalSliceEq : compare two slices
func MultiDimentionalSliceEq(a, b [][]int) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

// Loop2D : loop through 2 nested for loops
func Loop2D(height, width int, logic func(row, col int)) {
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			logic(row, col)
		}
	}
}

// SendRequest : send http request to provided url
func SendRequest(req *http.Request) []byte {
	client := http.Client{}
	res, err := client.Do(req)
	EoE("Error Getting HTTP Response", err)

	resData, err := ioutil.ReadAll(res.Body)
	EoE("Error Parsing Response", err)
	return resData
}

// EoE : exit with error code 1 and print if err is notnull
func EoE(msg string, err error) {
	if err != nil {
		fmt.Printf("\n❌  %s\n   %v\n", msg, err)
		os.Exit(1)
	}
}

// GetHomeDir : returns a full path to user's home dorectory
func GetHomeDir() string {
	usr, err := user.Current()
	EoE("Failed to get Current User", err)
	if usr.HomeDir != "" {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// Confirm : return confirmation based on user input
func Confirm(q string) bool {
	a := GetInput(q + " (Y/n) ")
	var res bool
	switch a {
	case "":
		fallthrough
	case "y":
		fallthrough
	case "Y":
		res = true
	case "n":
	case "N":
		res = false
	default:
		return Confirm(q)
	}
	return res
}

// GetInput : return string of user input
func GetInput(q string) string {
	if q != "" {
		print(q)
	}
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	return strings.TrimRight(ans, "\n")
}

// SetFromInput : set value of provided var to the value of user input
func SetFromInput(a *string, q string) {
	*a = strings.TrimRight(GetInput(q), "\n")
}

// GetIP : get local ip address
func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	EoE("Failed to Get Inet Address", err)
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// LineCounter : count number of lines `\n`
func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}