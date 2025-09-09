package main

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

const (
	SETTINGS_ECU = "kpro"
	SETTINGS_CAN = "vcan0"
	SETTINGS_HZ  = 10
)

type FrameGpsLapData struct {
	Latitude  float64
	Longitude float64
}
type FrameMisc struct {
	Hertz float64
}
type Frame660 struct {
	Rpm     uint16
	Speed   uint16
	Gear    uint8
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
	frameGpsLapData = []FrameGpsLapData{{
		Latitude:  0,
		Longitude: 0,
	}}

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

func toUint8(s string, frame string) uint8 {
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		fmt.Println("[ERROR] - toUint8 | ", s, " | ", frame)
		return 0
	}
	return uint8(v)
}
func toUint16(s string, frame string) uint16 {
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		fmt.Println("[ERROR] - toUint16 | ", s, " | ", frame)
		return 0
	}
	return uint16(v)
}
func toFloat64(s string, frame string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("[ERROR] - toFloat64 | ", s, " | ", frame)
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

	fmt.Println("[INFO] Reading in datalog file")
	lineCounter := 0
	for {
		record, err := reader.Read()
		if err != nil {
			break //EOF
		}
		if lineCounter >= 3 {
			frameMisc = append(frameMisc, FrameMisc{
				Hertz: toFloat64(record[0], "frameMisc - Hertz"),
			})

			frame660 = append(frame660, Frame660{
				Rpm:   toUint16(record[1], "frame660 - Rpm"),
				Speed: toUint16(record[2], "frame660 - Speed"),
				Gear:  toUint8(record[3], "frame660 - Gear"),
				// Voltage: toUint8(record[4], "frame660 - Voltage"),
				Voltage: uint8(toFloat64(record[4], "frame660 - Voltage")),
			})

			frame661 = append(frame661, Frame661{
				Iat: toUint16(record[5], "frame661 - Iat"),
				Ect: toUint16(record[6], "frame661 - Ect"),
				Mil: toUint8(record[7], "frame661 - Mil"),
				Vts: toUint8(record[8], "frame661 - Vts"),
				Cl:  toUint8(record[9], "frame661 - Cl"),
			})

			frame662 = append(frame662, Frame662{
				Tps: toUint16(record[10], "frame662 - Tps"),
				Map: toUint16(record[11], "frame662 - Map"),
			})

			frame663 = append(frame663, Frame663{
				Inj: toUint16(record[12], "frame663 - Inj"),
				Ign: toUint16(record[13], "frame663 - Ign"),
			})

			frame664 = append(frame664, Frame664{
				// LambdaRatio: toUint16(record[14], "frame664 - LambdaRatio"),
				LambdaRatio: uint16(toFloat64(record[14], "frame664 - LambdaRatio")),
			})

			frame665 = append(frame665, Frame665{
				Knock: toUint16(record[15], "frame665 - Knock"),
			})

			frame666 = append(frame666, Frame666{
				TargetCamAngle: toFloat64(record[16], "frame666 - TargetCamAngle"),
				ActualCamAngle: toFloat64(record[17], "frame666 - ActualCamAngle"),
			})

			frame667 = append(frame667, Frame667{
				Analog0: toUint16(record[18], "frame667 - Analog0"),
				Analog1: toUint16(record[19], "frame667 - Analog1"),
				Analog2: toUint16(record[20], "frame667 - Analog2"),
				Analog3: toUint16(record[21], "frame667 - Analog3"),
			})

			frame668 = append(frame668, Frame668{
				Analog4: toUint16(record[22], "frame668 - Analog4"),
				Analog5: toUint16(record[23], "frame668 - Analog5"),
				Analog6: toUint16(record[24], "frame668 - Analog6"),
				Analog7: toUint16(record[25], "frame668 - Analog7"),
			})

			if SETTINGS_ECU == "s300" {
				frame669S300 = append(frame669S300, Frame669S300{
					Frequency: toUint8(record[26], "frame669S300 - Frequency"),
					Duty:      toUint8(record[27], "frame669S300 - Duty"),
					Content:   toFloat64(record[28], "frame669S300 - Content"),
				})
			} else if SETTINGS_ECU == "kpro" {
				frame669KPRO = append(frame669KPRO, Frame669KPRO{
					Frequency:       toUint8(record[26], "frame669KPRO - Frequency"),
					EthanolContent:  toUint8(record[27], "frame669KPRO - EthanolContent"),
					FuelTemperature: toUint16(record[28], "frame669KPRO - FuelTemperature"),
				})
			}

			frameGpsLapData = append(frameGpsLapData, FrameGpsLapData{
				Latitude:  toFloat64(record[29], "frameGpsLapData - Latitude"),
				Longitude: toFloat64(record[30], "frameGpsLapData - Longitude"),
			})
		}

		lineCounter++
	}

	// Connect to CANBus
	conn, err := socketcan.DialContext(context.Background(), "can", SETTINGS_CAN)
	if err != nil {
		log.Fatal("[ERROR] Cannot connect to: ", SETTINGS_CAN)
	}
	defer conn.Close()
	fmt.Println("[INFO] Connected to vcan0")

	tx := socketcan.NewTransmitter(conn)

	ticker := time.NewTicker(time.Second / time.Duration(SETTINGS_HZ))
	defer ticker.Stop()
	i := 0

	fmt.Println("[INFO] Transmitting data to CANBus...")

	for range ticker.C {
		// Create all the can.Frame variables to transmit
		
		var b660 [8]byte
		binary.BigEndian.PutUint16(b660[0:2], frame660[i].Rpm)
		binary.BigEndian.PutUint16(b660[2:4], frame660[i].Speed)
		b660[4] = frame660[i].Gear
		b660[5] = frame660[i].Voltage
		canFrame660 := can.Frame{
			ID:     660,
			Length: 6,
			Data: b660,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame660)


		var b661 [8]byte
		binary.BigEndian.PutUint16(b661[0:2], frame661[i].Iat)
		binary.BigEndian.PutUint16(b661[2:4], frame661[i].Ect)
		b661[4] = frame661[i].Mil
		b661[5] = frame661[i].Vts
		b661[6] = frame661[i].Cl
		canFrame661 := can.Frame{
			ID:     661,
			Length: 7,
			Data:   b661,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame661)


		var b662 [8]byte
		binary.BigEndian.PutUint16(b662[0:2], frame662[i].Tps)
		binary.BigEndian.PutUint16(b662[2:4], frame662[i].Map)
		canFrame662 := can.Frame{
			ID: 	662,
			Length: 4,
			Data: 	b662,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame662)


		var b663 [8]byte
		binary.BigEndian.PutUint16(b663[0:2], frame663[i].Inj)
		binary.BigEndian.PutUint16(b663[2:4], frame663[i].Ign)
		canFrame663 := can.Frame{
			ID: 	663,
			Length: 4,
			Data: 	b663,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame663)


		var b664 [8]byte
		binary.BigEndian.PutUint16(b664[0:2], frame664[i].LambdaRatio)
		canFrame664 := can.Frame{
			ID: 	664,
			Length: 2,
			Data: 	b664,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame664)


		var b665 [8]byte
		binary.BigEndian.PutUint16(b665[0:2], frame665[i].Knock)
		canFrame665 := can.Frame{
			ID: 	665,
			Length: 2,
			Data: 	b665,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame665)


		var b666 [8]byte
		binary.BigEndian.PutUint16(b666[0:2], uint16(frame666[i].TargetCamAngle))
		binary.BigEndian.PutUint16(b666[2:4], uint16(frame666[i].ActualCamAngle))
		canFrame666 := can.Frame{
			ID: 	666,
			Length: 4,
			Data: 	b666,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame666)


		var b667 [8]byte
		binary.BigEndian.PutUint16(b667[0:2], frame667[i].Analog0)
		binary.BigEndian.PutUint16(b667[2:4], frame667[i].Analog1)
		binary.BigEndian.PutUint16(b667[4:6], frame667[i].Analog2)
		binary.BigEndian.PutUint16(b667[6:8], frame667[i].Analog3)
		canFrame667 := can.Frame{
			ID: 	667,
			Length: 8,
			Data: 	b667,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame667)


		var b668 [8]byte
		binary.BigEndian.PutUint16(b668[0:2], frame668[i].Analog4)
		binary.BigEndian.PutUint16(b668[2:4], frame668[i].Analog5)
		binary.BigEndian.PutUint16(b668[4:6], frame668[i].Analog6)
		binary.BigEndian.PutUint16(b668[6:8], frame668[i].Analog7)
		canFrame668 := can.Frame{
			ID: 	668,
			Length: 8,
			Data: 	b668,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame668)


		var b669 [8]byte

		if SETTINGS_ECU == "s300" {
			b669[0] = frame669S300[i].Frequency
			b669[1] = frame669S300[i].Duty
			binary.BigEndian.PutUint16(b669[2:4], uint16(frame669S300[i].Content))
		} else if SETTINGS_ECU == "kpro" {
			b669[0] = frame669KPRO[i].Frequency
			b669[1] = frame669KPRO[i].EthanolContent
			binary.BigEndian.PutUint16(b669[2:4], frame669KPRO[i].FuelTemperature)
		}
		canFrame669 := can.Frame{
			ID: 	669,
			Length: 8,
			Data: 	b669,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame669)

		var b700 [8]byte
		// TODO: Test data
		fmt.Println("-----------------------")
		fmt.Println(frameGpsLapData[i].Latitude)
		fmt.Println(frameGpsLapData[i].Longitude)
		binary.BigEndian.PutUint32(b700[0:4], math.Float32bits(float32(frameGpsLapData[i].Latitude)))
		binary.BigEndian.PutUint32(b700[4:8], math.Float32bits(float32(frameGpsLapData[i].Longitude)))
		canFrame700 := can.Frame{
			ID:     6969,
			Length: 8,
			Data:   b700,
		}
		_ = tx.TransmitFrame(context.Background(), canFrame700)

		// Once we've run through the data (lineCounter minus headers), restart
		if i == lineCounter - 3 {
			i = 0
		} else {
			i++
		}
	}
}
