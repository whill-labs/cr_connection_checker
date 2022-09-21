package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

func main() {
	var cmdRecv = &cobra.Command{
		Use:   "rcv [Receive data for Model CR]",
		Short: "Receive Data",
		Long:  `Send command power on and receive data to Model CR`,
		Run: func(cmd *cobra.Command, args []string) {
			receive()
		},
	}

	var rootCmd = &cobra.Command{Use: "cr_connection_checker"}
	rootCmd.AddCommand(cmdRecv)
	rootCmd.Execute()
}

func receive() {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("/dev/tty.imu", mode)
	if err != nil {
		log.Fatal("serial open error: ", err)
		os.Exit(1)
	}

	//set ODR 200Hz
	fmt.Println("set ODR 200Hz")
	//55 55 75 50 0C 04 00 00 00 C8 00 00 00 00 00 00 00 F0 D5
	data := []byte{0x55, 0x55, 0x75, 0x50, 0x0C, 0x04, 0x0, 0x0, 0x0, 0xC8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xF0, 0xD5}
	n, err := port.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	fmt.Println("save config")
	//save config
	data_conf := []byte{0x55, 0x55, 0x73, 0x43, 0x0, 0xC8, 0xCB}
	n, err = port.Write(data_conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}
}
