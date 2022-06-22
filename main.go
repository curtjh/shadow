package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"shadow/connections"
	"strings"
)

func main() {

	var v string
	var listUsers bool
	var control bool
	var consent bool
	var disconnect string

	flag.StringVar(&v, "v", "", "Computer Name - Required first for all commands")
	flag.BoolVar(&control, "control", false, "Take Control of system (see consent)")
	flag.BoolVar(&consent, "consent", false, "Get consent from user before shadowing/controlling")
	flag.BoolVar(&listUsers, "listUsers", false, "Get list of users logged in (active and disconnected)")
	flag.StringVar(&disconnect, "disconnect", "0", "Remove user session providing session ID (use listUsers to get session")
	flag.Parse()

	if strings.Trim(v, " ") == "" {

		fmt.Println()
		fmt.Println("Computer Name argument must be passed: -v [computername] ...")
		fmt.Println()

		os.Exit(0)
	}

	if disconnect != "0" {

		connections.RemoveSession(v, disconnect)

		os.Exit(0)
	}

	out, err := connections.Qwinsta(v)
	if err != nil {

		fmt.Println()
		fmt.Println(err)
		fmt.Printf("%s is mgiht be offline...\n", v)
		fmt.Println()

	} else {

		result := string(out)

		sli := strings.Split(strings.ReplaceAll(result, "\r\n", "\n"), "\n")

		for _, line := range sli {

			vals := strings.Fields(line)

			if len(vals) > 0 {

				if vals[0] != "services" && vals[0] != "SESSIONNAME" && vals[0] != " " {

					// check for active console session
					if strings.Contains(line, "console") && strings.Contains(line, "Active") {

						connections.AddConnection(vals[1], vals[2], vals[3])

					} else if strings.Contains(line, "rdp-tpc#") && strings.Contains(line, "Active") {

						// RDP session is the active session.  values 1, 2, and 3 are user, id, and status
						connections.AddConnection(vals[0], vals[1], vals[2])

					} else {

						// if it's not the other 2, it should be a slice that contains values 1, 2, and 3 are user, id, and status
						// they could be either active or disc at this point
						if vals[0] != "rdp-tcp" && !strings.Contains(line, "console") { // the rdp-tcp line is of no interest, and at this point we know the console is not active

							if strings.Contains(vals[0], "rdp-tcp#") {

								connections.AddConnection(vals[1], vals[2], vals[3])

							} else {

								connections.AddConnection(vals[0], vals[1], vals[2])
							}

						}

					}

				}

			}

		}

		x := 0
		activeSession := false
		for _, thisConn := range connections.List {

			if listUsers {

				if x == 0 {

					fmt.Println()
					fmt.Println("Logged in Users:")
					fmt.Println()

				}

				fmt.Println(thisConn.Computer + " " + thisConn.ID + " " + thisConn.Status)

			}

			if thisConn.Status == "Active" {

				activeSession = true

				var cmd string
				shadowCmd := "c:\\windows\\system32\\mstsc.exe /v " + v + " /shadow:" + thisConn.ID

				if !consent {

					if control {

						// both control and no consent (default condition of no other args are provided beyond -v (computername) )
						cmd = "{" + shadowCmd + " /noConsentPrompt /control}"

					} else {

						// shadow no consent
						cmd = "{" + shadowCmd + " /noConsentPrompt}"

					}

				} else {

					if control {

						// both control and consent
						cmd = "{" + shadowCmd + " /control}"

					} else {

						// no control with consent consent
						cmd = "{" + shadowCmd + "}"

					}

				}

				fmt.Println()

				if !listUsers { // if listing users, don't initiate a shadow session, just list

					_, err := exec.Command("powershell.exe", "-command", "Invoke-Command", "-scriptBlock ", cmd).CombinedOutput()
					if err != nil {

						fmt.Println(err)

					}

				}

			}

			x++

		}

		if !activeSession { // there were no active sessions

			fmt.Println()
			fmt.Println("No actives sessions found to shadow...")

		}

		fmt.Println()

	}

}
