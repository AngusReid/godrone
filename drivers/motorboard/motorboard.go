package motorboard

import (
	"fmt"
	"github.com/felixge/godrone/log"
	"os"
	"sync"
	"time"
)

const (
	PWM_MAX    = float64(511)
	DefaultTTY = "/dev/ttyO0"
)

type Motorboard struct {
	file        *os.File
	pwms        [4]int
	leds        [4]LedColor
	ledsChanged bool
	lock        sync.RWMutex
	log         log.Interface
	hz          int
}

func NewMotorboard(ttyPath string, log log.Interface) (*Motorboard, error) {
	motorboard := &Motorboard{log: log, hz: 20}
	err := motorboard.open(ttyPath)
	if err != nil {
		return nil, err
	}
	go motorboard.loop()
	return motorboard, nil
}

func (m *Motorboard) open(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	m.file = file
	return nil
}

func (m *Motorboard) loop() {
	maxSleep := time.Second / time.Duration(m.hz)
	for {
		start := time.Now()
		m.lock.RLock()
		m.updateSpeeds()
		if m.ledsChanged {
			m.updateLeds()
			m.ledsChanged = false
		}
		m.lock.RUnlock()

		sleep := maxSleep - time.Since(start)
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}

func (m *Motorboard) Leds() []LedColor {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.leds[:]
}

func (m *Motorboard) SetLed(led int, color LedColor) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.leds[led] = color
	m.ledsChanged = true
}

func (m *Motorboard) SetLeds(colors []LedColor) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i := 0; i < len(m.leds) && i < len(colors); i++ {
		m.leds[i] = colors[i]
	}
	m.ledsChanged = true
}

func (m *Motorboard) Speed(motorId int) (float64, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if motorId >= len(m.pwms) {
		return 0, fmt.Errorf("unknown motor: %d", motorId)
	}

	return float64(m.pwms[motorId]) / PWM_MAX, nil
}

func (m *Motorboard) SetSpeed(motorId int, speed float64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if motorId >= len(m.pwms) {
		return fmt.Errorf("unknown motor: %d", motorId)
	}

	m.pwms[motorId] = int(speed * PWM_MAX)
	return nil
}

func (m *Motorboard) MotorCount() int {
	return len(m.pwms)
}

type LedColor int

const (
	LedOff    LedColor = iota
	LedRed             = 1
	LedGreen           = 2
	LedOrange          = 3
)

// cmd = 011rrrrx xxxggggx (used to be 011grgrg rgrxxxxx in AR Drone 1.0)
// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c#L243
func (m *Motorboard) ledCmd() []byte {
	cmd := make([]byte, 2)
	cmd[0] = 0x60

	for i, color := range m.leds {
		if color == LedRed || color == LedOrange {
			cmd[0] = cmd[0] | (1 << (byte(i) + 1))
		}

		if color == LedGreen || color == LedOrange {
			cmd[1] = cmd[1] | (1 << (byte(i) + 1))
		}
	}
	return cmd
}

// see: https://github.com/ardrone/ardrone/blob/master/ardrone/motorboard/motorboard.c
func (m *Motorboard) pwmCmd() []byte {
	cmd := make([]byte, 5)
	cmd[0] = byte(0x20 | ((m.pwms[0] & 0x1ff) >> 4))
	cmd[1] = byte(((m.pwms[0] & 0x1ff) << 4) | ((m.pwms[1] & 0x1ff) >> 5))
	cmd[2] = byte(((m.pwms[1] & 0x1ff) << 3) | ((m.pwms[2] & 0x1ff) >> 6))
	cmd[3] = byte(((m.pwms[2] & 0x1ff) << 2) | ((m.pwms[3] & 0x1ff) >> 7))
	cmd[4] = byte(((m.pwms[3] & 0x1ff) << 1))
	return cmd
}

func (m *Motorboard) updateSpeeds() error {
	_, err := m.file.Write(m.pwmCmd())
	return err
}

func (m *Motorboard) updateLeds() error {
	_, err := m.file.Write(m.ledCmd())
	return err
}
