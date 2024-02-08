/* TP U ARE HERE:
- re-organize the channels given new cb shell
- decide if the shell should be an option vs automatic
*/

package main

import (
	"bufio"
	"strings"
	"time"

	//"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"

	// NetPuppy modules:
	"netpuppy/utils"
)

func sum(a int, b int) int {
	s := a + b
	return s
}

func readUserInput(ioReader chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if len(text) > 0 {
			ioReader <- text
		}
	}
}

func readFromSocket(socketReader chan<- []byte, connection net.Conn) {
	// Read from connection socket:
	for {
		dataBytes, err := bufio.NewReader(connection).ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error reading from socket: %v\n", err)
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
			return
		}
		socketReader <- dataBytes
	}
}

func startHelperShell() (*exec.Cmd, error) { // @Trauma_X_Sella 'connection'
	bashPath, err := exec.LookPath(`/bin/bash`)
	if err != nil {
		fmt.Printf("Error finding bash path: %v\n", err)
		os.Stderr.WriteString(" " + err.Error() + "\n")
		return nil, err
	}
	bCmd := exec.Command(bashPath)
	var erR error = bCmd.Start()

	return bCmd, erR
}

func main() {
	s := sum(2, 3)
	st := fmt.Sprintf("%v", s)
	fmt.Printf("sum: %v\n", st)

	flagStruct := utils.GetFlags()

	fmt.Printf("Flags = %v\n", flagStruct.Host)

	// Print banner:
	fmt.Printf("%s", utils.Banner())

	// Get STDIN and save to a variable we can use if we need:
	// stdinReader := bufio.NewReader(os.Stdin)
	// stdin, _ := stdinReader.ReadString('\n')
	// fmt.Printf("STDIN = %v", stdin) // Keep for now to avoid golang complaints about unused vars.

	// Initiate peer struct:
	thisPeer := utils.CreatePeer(flagStruct.Port, flagStruct.Host, flagStruct.Listen)

	// Now that we have our peer: try to make connection
	var asyncio_rocks net.Conn // connection @0xtib3rius
	var err error

	if thisPeer.ConnectionType == "offense" {
		listener, err1 := net.Listen("tcp", fmt.Sprintf(":%v", thisPeer.RPort))
		if err1 != nil {
			fmt.Printf("Error when creating listener: %v\n", err1)
			os.Stderr.WriteString(" " + err.Error() + "\n")
			os.Exit(1)
		}

		defer listener.Close() // Ensure the listener closes when main() returns

		asyncio_rocks, err = listener.Accept()
		if err != nil {
			os.Stderr.WriteString(" " + err.Error() + "\n")
			os.Exit(1)
			//  log.Fatal(err1.Error()
		}
	} else {
		remoteHost := fmt.Sprintf("%v:%v", thisPeer.Address, thisPeer.RPort)
		asyncio_rocks, err = net.Dial("tcp", remoteHost)

		// If there is an err, try the host address as ipv6 (need to add [] around string):
		if err != nil {
			remoteHost := fmt.Sprintf("[%v]:%v", thisPeer.Address, thisPeer.RPort)
			asyncio_rocks, err = net.Dial("tcp", remoteHost)

			if err != nil {
				os.Stderr.WriteString(" " + err.Error() + "\n")
				os.Exit(1)
			}
		}
	}

	// Attach connection to peer struct:
	thisPeer.Connection = asyncio_rocks
	localPortArr := strings.Split(thisPeer.Connection.LocalAddr().String(), ":")
	localPort := localPortArr[len(localPortArr)-1]
	thisPeer.LPort = localPort

	// Update user:
	var updateUserBanner string = utils.UserSelectionBanner(thisPeer.ConnectionType, thisPeer.Address, thisPeer.RPort, thisPeer.LPort)
	fmt.Println(updateUserBanner)

	// Start channel to listen for SIGINT:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		// If SIGINT: close connection, exit w/ code 2
		for sig := range signalChan {
			if sig.String() == "interrupt" {
				fmt.Printf("signal: %v\n", sig)
				thisPeer.Connection.Close()
				os.Exit(2)
			}
		}
	}()

	// If we're the connect_back peer, start 'helper' shell:
	if thisPeer.ConnectionType == "connect_back" {
		connectBackShell, shellStartErr := startHelperShell()
		if shellStartErr != nil {
			fmt.Printf("Error starting shell process: %v\n", shellStartErr)
			thisPeer.Connection.Close()
			os.Stderr.WriteString(" " + shellStartErr.Error() + "\n")
			os.Exit(1)
		}
		thisPeer.CbShell = connectBackShell
	}

	// IO read & socket write channels (user input will be written to socket)
	ioReader := make(chan string)

	// IO write & socket read channels (messages from socket will be printed to stdout)
	socketReader := make(chan []byte)

	go readUserInput(ioReader)
	go readFromSocket(socketReader, thisPeer.Connection)

	for {
		select {
		case userInput := <-ioReader:
			_, err := thisPeer.Connection.Write([]byte(userInput))
			if err != nil {
				// Quit here?
				fmt.Printf("Error in userInput select: %v\n", err)
				os.Stderr.WriteString(" " + err.Error() + "\n")
			}
		case socketIncoming := <-socketReader:
			_, err := os.Stdout.Write(socketIncoming)
			if err != nil {
				// Quit here?
				fmt.Printf("Error in writing to stdout: %v\n", err)
				os.Stderr.WriteString(" " + err.Error() + "\n")
			}
		default:
			time.Sleep(300 * time.Millisecond)
		}
	}
	return
}
