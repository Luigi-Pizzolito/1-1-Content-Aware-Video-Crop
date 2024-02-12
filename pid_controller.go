package main

import (
	"time"
)

type PIDController struct {
	Kp, Ki, Kd float64
	setpoint   float64
	integral   float64
	lastError  float64
    delta      float64
}

func NewPIDController(Kp, Ki, Kd, setpoint float64) *PIDController {
	return &PIDController{
		Kp:       Kp,
		Ki:       Ki,
		Kd:       Kd,
		setpoint: setpoint,
        delta:    setpoint,
	}
}

func (pid *PIDController) SetSetpoint(setpoint float64) {
	pid.setpoint = setpoint
}
func (pid *PIDController) SetDelta(lastsetpoint float64) {
	pid.delta = lastsetpoint
}

func (pid *PIDController) Update() float64 {
	error := pid.setpoint - pid.delta
	now := time.Now().UnixNano() / int64(time.Millisecond)
	dt := float64(now/1000) - pid.lastError
	pid.integral += error * dt
	derivative := (error - pid.lastError) / dt
	output := pid.Kp*error + pid.Ki*pid.integral + pid.Kd*derivative
	pid.lastError = error
    // fmt.Printf("set: %.4f, last: %.4f, err: %.4f, adj: %.4f\n", pid.setpoint, pid.delta, error, output)
	return output
}