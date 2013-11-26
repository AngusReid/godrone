package navboard

// From https://github.com/RoboticaTUDelft/paparazzi/blob/minor1/sw/airborne/boards/ardrone/navdata.h
type Data struct {
	Seq uint16

	// Accelerometers
	Ax int16
	Ay int16
	Az int16

	// Gyroscopes
	Gx int16
	Gy int16
	Gz int16

	// Everything below has not been confirmed to be correct yet
	TemperatureAcc  uint16
	TemperatureGyro uint16

	Ultrasound uint16

	UsDebutEcho       uint16
	UsFinEcho         uint16
	UsAssociationEcho uint16
	UsDistanceEcho    uint16

	UsCurveTime  uint16
	UsCurveValue uint16
	UsCurveRef   uint16

	NbEcho uint16

	SumEcho  uint32
	Gradient int16

	FlagEchoIni uint16

	Pressure            int32
	TemperaturePressure int16

	Mx int16
	My int16
	Mz int16

	Checksum uint16
}
