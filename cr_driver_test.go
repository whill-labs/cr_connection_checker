package main

import (
	"encoding/binary"
	"testing"
)

func Uint2bytes(i uint64, size int) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, i)
	return bytes[8-size : 8]
}
func Int2bytes(i int, size int) []byte {
	var ui uint64
	if 0 < i {
		ui = uint64(i)
	} else {
		ui = (^uint64(-i) + 1)
	}
	return Uint2bytes(ui, size)
}

func TestChecksum(t *testing.T) {
	dummyData := []byte{0xaf, 0x1f, 0x01, 0x16, 0x8f, 0xfe, 0xfd, 0xe9, 0xd9, 0x01, 0xdb, 0xfb, 0xd6, 0xfd, 0x63, 0x00, 0x00, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10}
	length := dummyData[1]
	if calcChecksum(dummyData, int(length)) != dummyData[len(dummyData)-1] {
		t.Error("checksum")
	}
}

func TestAnalyzeSample(t *testing.T) {
	dummyData := []byte{0xaf, 0x1f, 0x01, 0x16, 0x8f, 0xfe, 0xfd, 0xe9, 0xd9, 0x01, 0xdb, 0xfb, 0xd6, 0xfd, 0x63, 0x00, 0x00, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10}
	cr := &CRDriver{}
	body, err := cr.analyze(dummyData)
	if err != nil {
		t.Error("CRDriver analyze err")
	}

	//fmt.Println("dummyData: ")
	//fmt.Printf("%v\n", hex.EncodeToString(dummyData[:]))
	//fmt.Println("body: ", body)

	if body.DataSetNumber != 1 {
		t.Errorf("CRDriver analyze body DataSetNumber %d", body.DataSetNumber)
	}
	if body.AccelX != 5775 {
		t.Errorf("CRDriver analyze body accelX %d", body.AccelX)
	}
	if body.AccelY != -259 {
		t.Errorf("CRDriver analyze body accelY %d", body.AccelY)
	}
	if body.AccelZ != -5671 {
		t.Errorf("CRDriver analyze body accelZ %d", body.AccelZ)
	}
	if body.GyroX != 475 {
		t.Errorf("CRDriver analyze body GyroX %d", body.GyroX)
	}
	if body.GyroY != -1066 {
		t.Errorf("CRDriver analyze body GyroY %d", body.GyroY)
	}
	if body.GyroZ != -669 {
		t.Errorf("CRDriver analyze body GyroZ %d", body.GyroZ)
	}
	if body.BatteryPower != 99 {
		t.Errorf("CRDriver analyze body BatteryPower %d", body.BatteryPower)
	}
}

func TestAnalyze(t *testing.T) {
	dummyData := []byte{0xaf, 0x1f}

	dataSetNumber := byte(1)
	dummyData = append(dummyData, dataSetNumber)
	accelX := 980
	dummyData = append(dummyData, Int2bytes(accelX, 2)...)
	accelY := -980
	dummyData = append(dummyData, Int2bytes(accelY, 2)...)
	accelZ := -300
	dummyData = append(dummyData, Int2bytes(accelZ, 2)...)
	gyroX := -100
	dummyData = append(dummyData, Int2bytes(gyroX, 2)...)
	gyroY := 100
	dummyData = append(dummyData, Int2bytes(gyroY, 2)...)
	gyroZ := -1032
	dummyData = append(dummyData, Int2bytes(gyroZ, 2)...)
	joyFront := -100
	dummyData = append(dummyData, Int2bytes(joyFront, 1)...)
	joySide := 100
	dummyData = append(dummyData, Int2bytes(joySide, 1)...)
	batteryPower := 100
	dummyData = append(dummyData, Uint2bytes(uint64(batteryPower), 1)...)
	batteryCurrent := -1100
	dummyData = append(dummyData, Int2bytes(batteryCurrent, 2)...)
	rightMotorAngle := 500
	dummyData = append(dummyData, Int2bytes(rightMotorAngle, 2)...)
	leftMotorAngle := -500
	dummyData = append(dummyData, Int2bytes(leftMotorAngle, 2)...)
	rightMotorSpeed := -800
	dummyData = append(dummyData, Int2bytes(rightMotorSpeed, 2)...)
	leftMotorSpeed := 1023
	dummyData = append(dummyData, Int2bytes(leftMotorSpeed, 2)...)
	//isPoweredOn := true
	dummyData = append(dummyData, byte(1)) //true : 1
	speedModeIndicator := 4
	dummyData = append(dummyData, Uint2bytes(uint64(speedModeIndicator), 1)...)
	errorNumber := 99
	dummyData = append(dummyData, Uint2bytes(uint64(errorNumber), 1)...)
	angleDetectCounter := 255
	dummyData = append(dummyData, Uint2bytes(uint64(angleDetectCounter), 1)...)

	length := dummyData[1]
	checksum := calcChecksum(dummyData, int(length))
	dummyData = append(dummyData, Uint2bytes(uint64(checksum), 1)...)

	cr := &CRDriver{}
	body, err := cr.analyze(dummyData)
	if err != nil {
		t.Error("CRDriver analyze err")
	}

	if body.DataSetNumber != dataSetNumber {
		t.Errorf("CRDriver analyze body DataSetNumber %d", body.DataSetNumber)
	}
	if body.AccelX != int16(accelX) {
		t.Errorf("CRDriver analyze body AccelX %d", body.AccelX)
	}
	if body.AccelY != int16(accelY) {
		t.Errorf("CRDriver analyze body AccelY %d", body.AccelY)
	}
	if body.AccelZ != int16(accelZ) {
		t.Errorf("CRDriver analyze body AccelZ %d", body.AccelZ)
	}
	if body.GyroX != int16(gyroX) {
		t.Errorf("CRDriver analyze body GyroX %d", body.GyroX)
	}
	if body.GyroY != int16(gyroY) {
		t.Errorf("CRDriver analyze body GyroY %d", body.GyroY)
	}
	if body.GyroZ != int16(gyroZ) {
		t.Errorf("CRDriver analyze body GyroZ %d", body.GyroZ)
	}
	if body.JoyFront != int8(joyFront) {
		t.Errorf("CRDriver analyze body JoyFront %d", body.JoyFront)
	}
	if body.JoySide != int8(joySide) {
		t.Errorf("CRDriver analyze body JoySide %d", body.JoySide)
	}
	if body.BatteryPower != uint8(batteryPower) {
		t.Errorf("CRDriver analyze body BatteryPower %d", body.BatteryPower)
	}
	if body.BatteryCurrent != int16(batteryCurrent) {
		t.Errorf("CRDriver analyze body BatteryCurrent %d", body.BatteryCurrent)
	}
	if body.RightMotorAngle != int16(rightMotorAngle) {
		t.Errorf("CRDriver analyze body RightMotorAngle %d", body.RightMotorAngle)
	}
	if body.LeftMotorAngle != int16(leftMotorAngle) {
		t.Errorf("CRDriver analyze body LeftMotorAngle %d", body.LeftMotorAngle)
	}
	if body.RightMotorSpeed != int16(rightMotorSpeed) {
		t.Errorf("CRDriver analyze body RightMotorSpeed %d", body.RightMotorSpeed)
	}
	if body.LeftMotorSpeed != int16(leftMotorSpeed) {
		t.Errorf("CRDriver analyze body LeftMotorSpeed %d", body.LeftMotorSpeed)
	}
	if body.IsPoweredOn != true {
		t.Errorf("CRDriver analyze body IsPoweredOn %v", body.IsPoweredOn)
	}
	if body.SpeedModeIndicator != uint8(speedModeIndicator) {
		t.Errorf("CRDriver analyze body IsPoweredOn %d", body.SpeedModeIndicator)
	}
	if body.Error != uint8(errorNumber) {
		t.Errorf("CRDriver analyze body Error %d", body.Error)
	}
	if body.AngleDetectCounter != uint8(angleDetectCounter) {
		t.Errorf("CRDriver analyze body AngleDetectCounter %d", body.AngleDetectCounter)
	}
}
