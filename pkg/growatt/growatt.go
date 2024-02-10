package growatt

import (
	"context"
	"errors"
	"growatt/atp/pkg/domain"
	"time"

	"github.com/simonvetter/modbus"
)

type strucT struct {
	setting Setting
}

type Setting struct {
	Port     string
	Baudrate uint
	Timeout  time.Duration
}

func NewRepository(setting Setting) RepositoryI {
	return &strucT{
		setting: setting,
	}
}

type RepositoryI interface {
	SPF5000(ctx context.Context, id uint8) (data domain.Data, err error)
}

func (r strucT) SPF5000(ctx context.Context, id uint8) (data domain.Data, err error) {
	url := "rtu://" + r.setting.Port
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      url,
		Speed:    r.setting.Baudrate,
		DataBits: 8,
		Parity:   modbus.PARITY_NONE,
		StopBits: 1,
		Timeout:  r.setting.Timeout,
	})
	if err != nil {
		err := errors.New("E0")
		return data, err
	}

	err = client.Open()
	if err != nil {
		err := errors.New("E1")
		return data, err
	}
	defer client.Close()

	client.SetUnitId(id)

	datas, err := client.ReadRegisters(17, 11, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(17, 11, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-1")
		return data, err
	}

	for i, value := range datas {
		switch i {
		case 0: //17
			data.Battery_Voltage = float32(value) / 100.00
		case 1: //18
			data.Battery_SOC = float32(value)

		case 3: //20
			data.AC_Input_Voltage = float32(value) / 10.00
		case 4: //21
			data.AC_Input_Frequency = float32(value) / 100.00

		case 5: //22
			data.AC_Output_Voltage = float32(value) / 10.00
		case 6: //23
			data.AC_Output_Frequency = float32(value) / 100.00

		case 8: //25
			data.Inverter_Temperature = float32(value) / 10.00
		case 10: //27
			data.Load_Percentage = float32(value) / 10.00

		}
	}

	datas, err = client.ReadRegisters(34, 2, modbus.INPUT_REGISTER)
	datas, err = client.ReadRegisters(34, 2, modbus.INPUT_REGISTER)
	if err != nil {
		err := errors.New("timeout-2")
		return data, err
	}

	for i, value := range datas {
		switch i {
		case 0: //34
			data.Output_Current = float32(value) / 100.00
		case 1: //35
			data.Inverter_Current = float32(value) / 100.00
		}
	}

	return data, nil
}
