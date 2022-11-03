package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"go.bug.st/serial"
)

type CRHeader struct {
	ProtocolHead byte
	Length       byte
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

type CRDriver struct {
	port          serial.Port
	openPortError error
}

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

func (cr *CRDriver) turnOn() error {
	if cr.openPortError != nil {
		return cr.openPortError
	}

	//power on (it needs to be set twirce)
	data := []byte{0xAF, 0x03, 0x02, 0x01, 0xAF}
	_, err := cr.port.Write(data)
	if err != nil {
		log.Println(err)
		return err
	}

	time.Sleep(time.Millisecond * 200)

	data = []byte{0xAF, 0x03, 0x02, 0x01, 0xAF}
	_, err = cr.port.Write(data)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func (cr *CRDriver) turnOff() error {
	if cr.openPortError != nil {
		return cr.openPortError
	}

	data := []byte{0xAF, 0x03, 0x02, 0x0, 0xAE}
	_, err := cr.port.Write(data)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func (cr *CRDriver) read(b []byte, port io.ReadWriteCloser) error {
	if cr.openPortError != nil {
		return cr.openPortError
	}
	_, err := cr.port.Read(b)
	if err != nil {
		log.Println(err)
		return err
	}

	//fmt.Printf("%v\n", hex.EncodeToString(b[:n]))
	return err
}

func (cr *CRDriver) startSendingDataSet1(interval uint16) error {
	if cr.openPortError != nil {
		return cr.openPortError
	}
	//send data
	//fmt.Println("sending data set 1")
	data := []byte{0xaf, 0x6}
	data = append(data, 0x0) //command
	data = append(data, 0x1) //data set 1
	data = append(data, Uint2bytes(uint64(interval), 2)...)
	data = append(data, 0x0) //speed mode 0
	length := data[1]
	checksum := calcChecksum(data, int(length))
	data = append(data, Uint2bytes(uint64(checksum), 1)...)

	_, err := cr.port.Write(data)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (cr *CRDriver) open(device string) error {
	mode := &serial.Mode{
		BaudRate: 38400,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Println("serial open error: ", err)
		cr.openPortError = err
		return err
	}
	cr.port = port
	return nil
}
func (cr *CRDriver) close() {
	if cr.openPortError != nil {
		return
	}
	cr.port.Close()
}

func (cr *CRDriver) setReadTimeout(t time.Duration) error {
	err := cr.port.SetReadTimeout(t)
	if err != nil {
		return err
	}
	return nil
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

func (cr *CRDriver) parseDataSet1(b []byte) (body DataSet1Body) {

	body.DataSetNumber = b[2]

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

	return body
}

func (cr *CRDriver) analyze(b []byte) (body DataSet1Body, err error) {

	if b[0] != 0xAF {
		return body, nil
	}

	length := b[1]
	command := b[2]

	if len(b) < int(length+2) {
		log.Printf("buffer size less than received size b[1]:%d, len{buff) %d", length, len(b))
		return body, nil
	}

	if length == 0x1F && Command(command) == DATA_SET_1 {
		if calcChecksum(b, (int)(length)) != b[length+1] {
			err := fmt.Errorf("Checksum unmatch")
			log.Println(err)
			return body, err
		}
		body = cr.parseDataSet1(b)
	}
	return body, err
}

func (cr *CRDriver) receive() (body DataSet1Body, err error) {
	buff := make([]byte, 256)
	err = cr.read(buff, cr.port)
	if err != nil {
		log.Println("read error")
		return body, err
	}
	body, err = cr.analyze(buff)
	return body, err
}

func Uint2bytes(i uint64, size int) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, i)
	return bytes[8-size : 8]
}
