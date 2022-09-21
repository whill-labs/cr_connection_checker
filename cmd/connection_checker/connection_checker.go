package main

import (
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

const (
	DATA_SET_0 Command = iota
	DATA_SET_1
)

const (
	DATA_SET_BODY_OFFSET = 3
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
	_, err = port.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 200)

	data = []byte{0xAF, 0x03, 0x02, 0x01, 0xAF}
	_, err = port.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func read(b []byte, port io.ReadWriteCloser) error {

	_, err := port.Read(b)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//fmt.Printf("%v\n", hex.EncodeToString(b[:n]))
	return err
}

func startSendingDataSet1(device string) error {
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
	_, err = port.Write(data_conf)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func receive(device string) error {
	err := startSendingDataSet1(device)
	if err != nil {
		return err
	}

	mode := &serial.Mode{
		BaudRate: 38400,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal("serial open error: ", err)
		os.Exit(1)
	}
	defer port.Close()

	for {
		buff := make([]byte, 256)
		for {
			time.Sleep(time.Millisecond * 1000)
			err = read(buff, port)
			if err != nil {
				return err
			}
			analyze(buff)
		}

	}
}

func calcChecksum(b []byte, len int) byte {
	var checksum byte
	for i := 0; i < len+1; i++ {
		checksum ^= b[i]
	}
	return checksum
}

func parseBool(b byte) bool {
	if b == 1 {
		return true
	} else {
		return false
	}
}

func parseDataSet1(b []byte, body DataSet1Body) error {

	buff := b[DATA_SET_BODY_OFFSET:]
	body.AccelX = int16(buff[0])<<8 | int16(buff[1])
	body.AccelY = int16(buff[2])<<8 | int16(buff[3])
	body.AccelZ = int16(buff[4])<<8 | int16(buff[5])
	body.GyroX = int16(buff[6])<<8 | int16(buff[7])
	body.GyroY = int16(buff[8])<<8 | int16(buff[9])
	body.GyroZ = int16(buff[10])<<8 | int16(buff[11])
	body.JoyFront = int8(buff[12])
	body.JoySide = int8(buff[13])
	body.BatteryPower = uint8(buff[14])
	body.BatteryCurrent = int16(buff[15])<<8 | int16(buff[16])
	body.RightMotorAngle = int16(buff[17])<<8 | int16(buff[18])
	body.LeftMotorAngle = int16(buff[19])<<8 | int16(buff[20])
	body.RightMotorSpeed = int16(buff[21])<<8 | int16(buff[22])
	body.LeftMotorSpeed = int16(buff[23])<<8 | int16(buff[24])
	body.IsPoweredOn = parseBool(buff[25])
	body.SpeedModeIndicator = uint8(buff[26])
	body.Error = uint8(buff[27])
	body.AngleDetectCounter = uint8(buff[28])
	fmt.Printf("JoyFront %d\n", body.JoyFront)
	fmt.Printf("JoySide %d\n", body.JoySide)
	fmt.Printf("BatteryPower %d\n", body.BatteryPower)
	fmt.Printf("BatteryCurrent %d\n", body.BatteryCurrent)
	fmt.Printf("isPoweredOn %t\n", body.IsPoweredOn)
	fmt.Printf("SpeedMode %x\n", body.SpeedModeIndicator)
	fmt.Printf("Error %x\n", body.Error)
	return nil
}

func analyze(b []byte) (body DataSet1Body, err error) {
	fmt.Println("analyze")

	if b[0] != 0xAF {
		return body, nil
	}

	//Data Set 1
	length := b[1]
	command := b[2]
	if length == 0x1F && Command(command) == DATA_SET_1 {
		fmt.Println("recv data set 1: ")

		if calcChecksum(b, (int)(length)) != b[length+1] {
			err := fmt.Errorf("Checksum unmatch")
			return body, err
		}
		//var body DataSet1Body
		parseDataSet1(b, body)
	}
	return body, err
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
