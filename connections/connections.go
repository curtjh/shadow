package connections

import (
	"os/exec"
	"strings"
)

type Connections struct {
	User   string
	ID     string
	Status string
}

var List []Connections

func addConnection(user string, id string, status string) {

	f := Connections{
		User:   string(user),
		ID:     string(id),
		Status: string(status),
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

	cmd := exec.Command("rwinsta.exe", session, "/server:"+computer)
	return cmd.CombinedOutput()

}

func ParseOutput(output string) {

	sli := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")

	for _, line := range sli {

		vals := strings.Fields(line)

		if len(vals) > 0 {

			if vals[0] != "services" && vals[0] != "SESSIONNAME" && vals[0] != " " {

				// check for active console session
				if strings.Contains(line, "console") && strings.Contains(line, "Active") {

					addConnection(vals[1], vals[2], vals[3])

				} else if strings.Contains(line, "rdp-tpc#") && strings.Contains(line, "Active") {

					// RDP session is the active session.  values 1, 2, and 3 are user, id, and status
					addConnection(vals[0], vals[1], vals[2])

				} else {

					// if it's not the other 2, it should be a slice that contains values 1, 2, and 3 are user, id, and status
					// they could be either active or disc at this point
					if vals[0] != "rdp-tcp" && !strings.Contains(line, "console") { // the rdp-tcp line is of no interest, and at this point we know the console is not active

						if strings.Contains(vals[0], "rdp-tcp#") {

							addConnection(vals[1], vals[2], vals[3])

						} else {

							addConnection(vals[0], vals[1], vals[2])
						}

					}

				}

			}

		}

	}

}
