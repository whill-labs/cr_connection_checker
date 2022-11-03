package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type config struct {
	Device string
}

type receiveUI struct {
	app                *tview.Application
	conf               *config
	cr                 *CRDriver
	power              *tview.TableCell
	rightMotorAngle    *tview.TableCell
	leftMotorAngle     *tview.TableCell
	rightMotorSpeed    *tview.TableCell
	leftMotorSpeed     *tview.TableCell
	angleDetectCounter *tview.TableCell
	batteryPower       *tview.TableCell
	batteryCurrent     *tview.TableCell
	joyFront           *tview.TableCell
	joySide            *tview.TableCell
	errorCode          *tview.TableCell
	speedSetting       *tview.TableCell
	device             *tview.TableCell
	deviceError        *tview.TableCell
}

const refreshInterval uint16 = 100

func queueUpdateAndDraw(app *tview.Application, f func()) {
	app.QueueUpdateDraw(f)
}

func parseOnOff(onOff bool) string {
	if onOff {
		return "ON"
	} else {
		return "OFF"
	}
}

func (rUI *receiveUI) updateReceivedData(data DataSet1Body) {
	queueUpdateAndDraw(rUI.app, func() {
		rUI.power.SetText(parseOnOff(data.IsPoweredOn))
		rUI.rightMotorAngle.SetText(fmt.Sprintf("%d", data.RightMotorAngle))
		rUI.leftMotorAngle.SetText(fmt.Sprintf("%d", data.LeftMotorAngle))
		rUI.rightMotorSpeed.SetText(fmt.Sprintf("%d", data.RightMotorSpeed))
		rUI.leftMotorSpeed.SetText(fmt.Sprintf("%d", data.LeftMotorSpeed))
		rUI.angleDetectCounter.SetText(fmt.Sprintf("%d", data.AngleDetectCounter))
		rUI.batteryPower.SetText(fmt.Sprintf("%d", data.BatteryPower))
		rUI.batteryCurrent.SetText(fmt.Sprintf("%d", data.BatteryCurrent))
		rUI.joyFront.SetText(fmt.Sprintf("%d", data.JoyFront))
		rUI.joySide.SetText(fmt.Sprintf("%d", data.JoySide))
		rUI.errorCode.SetText(fmt.Sprintf("%d", data.Error))
		rUI.speedSetting.SetText(fmt.Sprintf("%d", data.SpeedModeIndicator))
	})
}

func receive(rUI *receiveUI) error {
	err := rUI.cr.startSendingDataSet1(refreshInterval)
	if err != nil {
		return err
	}

	for {
		time.Sleep((time.Duration)(refreshInterval) * time.Millisecond)
		var data DataSet1Body
		data, err = rUI.cr.receive()
		if err != nil {
		} else {
			rUI.updateReceivedData(data)
		}
	}
}

func createCommandList() (commandList *tview.List) {
	commandList = tview.NewList()
	commandList.SetBorder(true).SetTitle("Command")
	return commandList
}

func createTextViewPanel(app *tview.Application, name string) (panel *tview.TextView) {
	panel = tview.NewTextView()
	panel.SetBorder(true).SetTitle(name)
	panel.SetChangedFunc(func() {
		app.Draw()
	})
	return panel
}

func createReceivePanel(app *tview.Application, rUI *receiveUI) (receivePanel *tview.Flex) {

	receiveInfo := tview.NewTable()
	receiveInfo.SetBorder(true).SetTitle("Receive")

	cnt := 0
	receiveInfo.SetCellSimple(cnt, 0, "Power:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.power = tview.NewTableCell("OFF")
	receiveInfo.SetCell(cnt, 1, rUI.power)
	cnt++

	//Battery
	receiveInfo.SetCellSimple(cnt, 0, "Battery Power:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.batteryPower = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.batteryPower)
	cnt++

	receiveInfo.SetCellSimple(cnt, 0, "Battery Current:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.batteryCurrent = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.batteryCurrent)
	cnt++

	//Joystick
	receiveInfo.SetCellSimple(cnt, 0, "Joystick Front:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.joyFront = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.joyFront)
	cnt++

	receiveInfo.SetCellSimple(cnt, 0, "Joystick Side:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.joySide = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.joySide)
	cnt++

	//Error
	receiveInfo.SetCellSimple(cnt, 0, "Error Code:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.errorCode = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.errorCode)
	cnt++

	//Speed
	receiveInfo.SetCellSimple(cnt, 0, "Speed Setting:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.speedSetting = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.speedSetting)
	cnt++

	//Motor
	receiveInfo.SetCellSimple(cnt, 0, "Right Motor Angle:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.rightMotorAngle = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.rightMotorAngle)
	cnt++

	receiveInfo.SetCellSimple(cnt, 0, "Left Motor Angle:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.leftMotorAngle = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.leftMotorAngle)
	cnt++

	receiveInfo.SetCellSimple(cnt, 0, "Right Motor Speed:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.rightMotorSpeed = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.rightMotorSpeed)
	cnt++

	receiveInfo.SetCellSimple(cnt, 0, "Left Motor Speed:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.leftMotorSpeed = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.leftMotorSpeed)
	cnt++

	receiveInfo.SetCellSimple(cnt, 0, "Angle Detect Counter:")
	receiveInfo.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	rUI.angleDetectCounter = tview.NewTableCell("0")
	receiveInfo.SetCell(cnt, 1, rUI.angleDetectCounter)
	cnt++

	configInfo := tview.NewTable()
	configInfo.SetBorder(true).SetTitle("Config")

	configInfo.SetCellSimple(0, 0, "Device:")
	configInfo.GetCell(0, 0).SetAlign(tview.AlignRight)
	rUI.device = tview.NewTableCell(rUI.conf.Device)
	configInfo.SetCell(0, 1, rUI.device)
	if rUI.cr.openPortError != nil {
		rUI.deviceError = tview.NewTableCell(rUI.cr.openPortError.Error())
		configInfo.SetCell(0, 2, rUI.deviceError)
	}

	receivePanel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(receiveInfo, 0, 1, false).
		AddItem(configInfo, 3, 1, false)

	return receivePanel
}

func createModalForm(pages *tview.Pages, form tview.Primitive, height int, width int) tview.Primitive {
	modal := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
	return modal
}

func powerCommand(pages *tview.Pages, rUI *receiveUI) func() {
	return func() {
		cancelFunc := func() {
			pages.SwitchToPage("main")
			pages.RemovePage("modal")
		}

		onFunc := func() {
			pages.SwitchToPage("main")
			pages.RemovePage("modal")
			rUI.cr.turnOn()
		}

		offFunc := func() {
			pages.SwitchToPage("main")
			pages.RemovePage("modal")
			rUI.cr.turnOff()
		}

		form := tview.NewForm()
		form.AddButton("ON", onFunc)
		form.AddButton("OFF", offFunc)
		form.AddButton("Cancel", cancelFunc)
		form.SetCancelFunc(cancelFunc)
		form.SetButtonsAlign(tview.AlignCenter)
		form.SetBorder(true).SetTitle("Power")
		modal := createModalForm(pages, form, 13, 55)
		pages.AddPage("modal", modal, true, true)
	}
}

func createLayout(cList tview.Primitive, recvPanel tview.Primitive, logPanel tview.Primitive) (layout *tview.Flex) {
	bodyLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(cList, 20, 1, true).
		AddItem(recvPanel, 0, 1, false)

	header := tview.NewTextView()
	header.SetBorder(false)
	header.SetText("=== Model CR connection checker === ")
	header.SetTextAlign(tview.AlignCenter)

	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 1, false).
		AddItem(bodyLayout, 0, 1, true).
		AddItem(logPanel, 3, 1, false)

	return layout
}

func createApplication(conf *config) (app *tview.Application) {
	app = tview.NewApplication()
	pages := tview.NewPages()

	cr := &CRDriver{}

	rUI := &receiveUI{}
	rUI.app = app
	rUI.conf = conf
	rUI.cr = cr

	logPanel := createTextViewPanel(app, "Log")
	log.SetOutput(logPanel)

	err := rUI.cr.open(conf.Device)
	if err == nil {
		rUI.cr.setReadTimeout((time.Duration)(refreshInterval) * time.Microsecond * 2)
	}

	receivePanel := createReceivePanel(app, rUI)

	commandList := createCommandList()
	commandList.AddItem("Power", "", 'p', powerCommand(pages, rUI))
	commandList.AddItem("Quit", "", 'q', func() {
		shutdown(rUI)
	})

	layout := createLayout(commandList, receivePanel, logPanel)
	pages.AddPage("main", layout, true, true)

	go receive(rUI)
	app.SetRoot(pages, true)
	return app
}

func loadConfig() (conf *config) {
	conf = &config{}
	conf.Device = "/dev/ttyUSB0"

	p, _ := os.Getwd()
	filename := filepath.Join(p, "device.json")
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Warning: config file not founde, use default: ", conf.Device)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		log.Println("Error: couldn't decode config file")
	}
	return conf
}

func shutdown(rUI *receiveUI) {
	rUI.cr.close()
	rUI.app.Stop()
}

func main() {
	runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}

	conf := loadConfig()
	app := createApplication(conf)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
