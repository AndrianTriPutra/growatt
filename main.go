package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"growatt/atp/pkg/domain"
	"growatt/atp/pkg/growatt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	flag.Usage = func() {
		log.Printf("Usage: go run . port modbusID")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	port := flag.Args()[0]
	slaveID := flag.Args()[1]
	modbusID, err := strconv.Atoi(slaveID)
	if err != nil {
		log.Fatalf("failed parse slaveID -> %s", err.Error())
	}

	setting := growatt.Setting{
		Port:     port,
		Baudrate: 9600,
		Timeout:  300 * time.Millisecond,
	}
	gWt := growatt.NewRepository(setting)
	ctx := context.Background()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	ts := time.Now().In(loc).Format(time.RFC3339)
	payload := domain.Payload{
		Device_ID: "dev_001",
		Timestamp: ts,
	}

	data, err := gWt.SPF5000(ctx, uint8(modbusID))
	if err != nil {
		newErr := fmt.Sprintf("failed gWt.SPF5000 id [%v]-> %s", modbusID, err.Error())
		log.Fatal(newErr)
	}

	inverter := domain.Inverter{
		ModbusID: uint8(modbusID),
		Data:     data,
	}
	payload.Inverter = inverter

	js, err := json.MarshalIndent(payload, " ", " ")
	if err != nil {
		log.Fatalf("failed MarshalIndent -> %s", err.Error())
	}
	msg := string(js)
	fmt.Printf("payload:\n%s\n", msg)
}
