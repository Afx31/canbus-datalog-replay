package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	SETTINGS_ECU = "kpro"
	SETTINGS_CAN = "vcan0"
	SETTINGS_HZ  = 10
)

type FrameMisc struct {
	Hertz float64
}
type Frame660 struct {
	Rpm uint16
	Speed uint16
	Gear uint8
	Voltage uint8
}
type Frame661 struct {
	Iat uint16
	Ect uint16
	Mil uint8
	Vts uint8
	Cl  uint8
}
type Frame662 struct {
	Tps uint16
	Map uint16
}
type Frame663 struct {
	Inj uint16
	Ign uint16
}
type Frame664 struct {
	LambdaRatio uint16
}
type Frame665 struct {
	Knock uint16
}
type Frame666 struct {
	TargetCamAngle float64
	ActualCamAngle float64
}
type Frame667 struct {
	Analog0 uint16
	Analog1 uint16
	Analog2 uint16
	Analog3 uint16
}
type Frame668 struct {
	Analog4 uint16
	Analog5 uint16
	Analog6 uint16
	Analog7 uint16
}
type Frame669S300 struct {
	Frequency uint8
	Duty      uint8
	Content   float64
}
type Frame669KPRO struct {
	Frequency       uint8
	EthanolContent  uint8
	FuelTemperature uint16
}

var (
	frameMisc = []FrameMisc{{
		Hertz: 0,
	}}

	frame660 = []Frame660{{
		Rpm:     0,
		Speed:   0,
		Gear:    0,
		Voltage: 0,
	}}

	frame661 = []Frame661{{
		Iat: 0,
		Ect: 0,
		Mil: 0,
		Vts: 0,
		Cl:  0,
	}}

	frame662 = []Frame662{{
		Tps: 0,
		Map: 0,
	}}

	frame663 = []Frame663{{
		Inj: 0,
		Ign: 0,
	}}

	frame664 = []Frame664{{
		LambdaRatio: 0,
	}}

	frame665 = []Frame665{{
		Knock: 0,
	}}

	frame666 = []Frame666{{
		TargetCamAngle: 0,
		ActualCamAngle: 0,
	}}

	frame667 = []Frame667{{
		Analog0: 0,
		Analog1: 0,
		Analog2: 0,
		Analog3: 0,
	}}

	frame668 = []Frame668{{
		Analog4: 0,
		Analog5: 0,
		Analog6: 0,
		Analog7: 0,
	}}

	frame669S300 = []Frame669S300{{
		Frequency: 0,
		Duty:      0,
		Content:   0,
	}}

	frame669KPRO = []Frame669KPRO{{
		Frequency:       0,
		EthanolContent:  0,
		FuelTemperature: 0,
	}}
)

func toUint8(s string) uint8 {
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		fmt.Println("[ERROR] - toUint8")
		return 0
	}
	return uint8(v)
}
func toUint16(s string) uint16 {
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		fmt.Println("[ERROR] - toUint16")
		return 0
	}
	return uint16(v)
}
func toFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("[ERROR] - toFloat64")
		return 0
	}
	return float64(v)
}

func main() {
	fmt.Println("---------- CANBus Datalog Replay - Started ----------")

	// Open file
	file, err := os.Open("testdata.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create new CSV reader
	reader := csv.NewReader(file)

	lineCounter := 0
	for {
		record, err := reader.Read()
		if err != nil {
			break //EOF
		}
		if lineCounter >= 3 {
			frameMisc = append(frameMisc, FrameMisc{
				Hertz: toFloat64(record[0]),
			})

			frame660 = append(frame660, Frame660{
				Rpm:     toUint16(record[1]),
				Speed:   toUint16(record[2]),
				Gear:    toUint8(record[3]),
				Voltage: toUint8(record[4]),
			})

			frame661 = append(frame661, Frame661{
				Iat: toUint16(record[5]),
				Ect: toUint16(record[6]),
				Mil: toUint8(record[7]),
				Vts: toUint8(record[8]),
				Cl:  toUint8(record[9]),
			})

			frame662 = append(frame662, Frame662{
				Tps: toUint16(record[10]),
				Map: toUint16(record[11]),
			})

			frame663 = append(frame663, Frame663{
				Inj: toUint16(record[12]),
				Ign: toUint16(record[13]),
			})

			frame664 = append(frame664, Frame664{
				LambdaRatio: toUint16(record[14]),
			})

			frame665 = append(frame665, Frame665{
				Knock: toUint16(record[15]),
			})

			frame666 = append(frame666, Frame666{
				TargetCamAngle: toFloat64(record[16]),
				ActualCamAngle: toFloat64(record[17]),
			})

			frame667 = append(frame667, Frame667{
				Analog0: toUint16(record[18]),
				Analog1: toUint16(record[19]),
				Analog2: toUint16(record[20]),
				Analog3: toUint16(record[21]),
			})

			frame668 = append(frame668, Frame668{
				Analog4: toUint16(record[22]),
				Analog5: toUint16(record[23]),
				Analog6: toUint16(record[24]),
				Analog7: toUint16(record[25]),
			})

			if SETTINGS_ECU == "s300" {
				frame669S300 = append(frame669S300, Frame669S300{
					Frequency: toUint8(record[26]),
					Duty:      toUint8(record[27]),
					Content:   toFloat64(record[28]),
				})
			} else if SETTINGS_ECU == "kpro" {
				frame669KPRO = append(frame669KPRO, Frame669KPRO{
					Frequency:       toUint8(record[26]),
					EthanolContent:  toUint8(record[27]),
					FuelTemperature: toUint16(record[28]),
				})
			}
		}

		lineCounter++
	}

}
