package domain

type Payload struct {
	Device_ID string
	Timestamp string
	Inverter  Inverter
}

type Inverter struct {
	ModbusID uint8
	Data     Data
}

type Data struct {
	Battery_Voltage      float32
	Battery_SOC          float32
	AC_Input_Voltage     float32
	AC_Input_Frequency   float32
	AC_Output_Voltage    float32
	AC_Output_Frequency  float32
	Inverter_Temperature float32
	Load_Percentage      float32
	Output_Current       float32
	Inverter_Current     float32
}
