package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.bug.st/serial"
)

var appVersion string

type CRHeader struct {
	ProtocolHead byte
	DataLength   byte
	Body         []byte
	CheckSum     byte
}

type Command int

const (
	START_SENDING_DATA Command = iota
	STOP_SENDING_DATA
	SET_POWER
)

type SetPowerBody struct {
	CommandID byte
	OnOff     byte
}

type DataSet1Body struct {
	DataSetNumber      byte
	AccelX             int16
	AccelY             int16
	AccelZ             int16
	GyroX              int16
	GyroY              int16
	GyroZ              int16
	JoyFront           int8
	JoySide            int8
	BatteryPower       uint8
	BatteryCurrent     int16
	RightMotorAngle    int16
	LeftMotorAngle     int16
	RightMotorSpeed    int16
	LeftMotorSpeed     int16
	IsPoweredOn        bool
	SpeedModeIndicator uint8
	Error              uint8
	AngleDetectCounter uint8
}

func turnOn(device string) error {
	fmt.Println("device: ", device)
	mode := &serial.Mode{
		BaudRate: 38400,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal("serial open error: ", err)
		os.Exit(1)
	}
	defer port.Close()

	//power on (it needs to be set twirce)
	fmt.Println("power on")
	data := []byte{0xAF, 0x03, 0x02, 0x01, 0xAF}
	n, err := port.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	time.Sleep(time.Millisecond * 200)

	fmt.Println("power on")
	data = []byte{0xAF, 0x03, 0x02, 0x01, 0xAF}
	n, err = port.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	return err
}

func read(b []byte, port io.ReadWriteCloser) error {

	fmt.Println("read")
	buff := make([]byte, 1024)
	n, err := port.Read(buff)
	fmt.Printf("n:%d\n", n)
	if err != nil {
		log.Fatal(err)
		return err
	}
	/*if n == 0 {
		fmt.Println("\nEOF")
		return nil
	}*/

	//only check header of receive data
	/*
		if reflect.DeepEqual(buff[0:3], expect_head) {
			fmt.Println("Success to receive data!")
			return nil
		}
	*/
	fmt.Printf("%v", hex.EncodeToString(buff[:n]))

	return err
}

func receive(device string) error {

	mode := &serial.Mode{
		BaudRate: 38400,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal("serial open error: ", err)
		os.Exit(1)
	}
	defer port.Close()

	//send data
	fmt.Println("sending data set 1")
	data_conf := []byte{0xAF, 0x06, 0x0, 0x01, 0x03, 0xE8, 0x0, 0x43}
	n, err := port.Write(data_conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	//expect_head := []byte{0xAF, 0x1E, 0x01}
	buff := make([]byte, 1024)
	//go func() {

	for {
		time.Sleep(time.Millisecond * 1000)
		err = read(buff, port)
		if err != nil {
			//parse
		}
	}
	//}()

	return err
}

func main() {
	var (
		device string
	)

	app := cli.NewApp()
	app.Name = "Model CR connection checker"
	app.Usage = ""

	if appVersion == "" {
		appVersion = "test-version"
	}
	app.Version = appVersion

	subFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "device",
			Aliases:     []string{"d"},
			Value:       "",
			Usage:       "device name",
			Destination: &device,
			Required:    true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "recv",
			Usage: "receive data from Model CR",
			Flags: append(subFlags),
			Action: func(c *cli.Context) (err error) {
				//main process
				turnOn(device)
				receive(device)

				return
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
