package connections

import (
	"os/exec"
)

type Connections struct {
	Computer string
	ID       string
	Status   string
}

var List []Connections

func AddConnection(computer string, id string, status string) {

	f := Connections{
		Computer: string(computer),
		ID:       string(id),
		Status:   string(status),
	}
	List = append(List, f)

}

// get session details
func Qwinsta(computer string) ([]byte, error) {

	cmd := exec.Command("qwinsta.exe", "/server:"+computer)
	return cmd.CombinedOutput()

}

// remove a session from a computer
func RemoveSession(computer string, session string) ([]byte, error) {

	cmd := exec.Command("rwinsta.exe", session, "/v", computer)
	return cmd.CombinedOutput()

}
