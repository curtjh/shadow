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
	var credPrompt bool
	var disconnect string

	flag.StringVar(&v, "v", "", "Computer Name - Required first for all commands")
	flag.BoolVar(&control, "control", false, "Take Control of system (see consent)")
	flag.BoolVar(&consent, "consent", false, "Get consent from user before shadowing/controlling")
	flag.BoolVar(&listUsers, "listUsers", false, "Get list of users logged in (active and disconnected)")
	flag.BoolVar(&credPrompt, "credPrompt", false, "Prompt for credentials that have permissions to shadow")
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

		// result :=
		connections.ParseOutput(string(out))

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

				prompt := ""

				if credPrompt {

					prompt = " /prompt"

				}

				if !consent {

					if control {

						// both control and no consent (default condition of no other args are provided beyond -v (computername) )
						cmd = "{" + shadowCmd + " /noConsentPrompt /control" + prompt + "}"

					} else {

						// shadow no consent
						cmd = "{" + shadowCmd + " /noConsentPrompt" + prompt + "}"

					}

				} else {

					if control {

						// both control and consent
						cmd = "{" + shadowCmd + " /control" + prompt + "}"

					} else {

						// no control with consent consent
						cmd = "{" + shadowCmd + prompt + "}"

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
