package godrone

import "math"

import "time"

type Controller struct {
	RotationBand  float64
	ThrottleHover float64
	ThrottleMin   float64
	Pitch         PID
	Roll          PID
	Yaw           PID
	Altitude      PID
}

func (c *Controller) Control(actual, desired Placement, dt time.Duration) [4]float64 {
	var speeds [4]float64
	var (
		pitchOut = c.Roll.Update(actual.Pitch, desired.Pitch, dt)
		rollOut  = c.Roll.Update(actual.Roll, desired.Roll, dt)
		yawOut   = c.Yaw.Update(actual.Yaw, desired.Yaw, dt)
		altOut   = c.Altitude.Update(actual.Altitude, desired.Altitude, dt)
	)
	throttle := math.Max(c.ThrottleMin, math.Min(1-c.RotationBand, c.ThrottleMin+altOut))
	speeds = [4]float64{
		throttle + clipBand(+rollOut+pitchOut+yawOut, rotationBand),
		throttle + clipBand(-rollOut+pitchOut-yawOut, rotationBand),
		throttle + clipBand(-rollOut-pitchOut+yawOut, rotationBand),
		throttle + clipBand(+rollOut-pitchOut-yawOut, rotationBand),
	}
	return speeds
}

func clipBand(val, band float64) float64 {
	return band/2 + clip(val, band/2)
}

func clip(val, max float64) float64 {
	return math.Max(math.Min(val, max), -max)
}
