package main

import "github.com/hajimehoshi/ebiten/v2"

type Keyboard struct {
	lastPressed map[ebiten.Key]uint64
}

func MakeKeyboard() Keyboard {
	return Keyboard{make(map[ebiten.Key]uint64, 0)}
}

func (kb *Keyboard) KeyPulse(key ebiten.Key, gameTick uint64, cooldown int) bool {
	lastTick, ok := kb.lastPressed[key]
	status := ebiten.IsKeyPressed(key)
	if !status {
		kb.lastPressed[key] = 0
		return false
	}

	if !ok || lastTick == 0 || gameTick-lastTick >= uint64(cooldown) {
		kb.lastPressed[key] = gameTick
		return true
	}
	return false
}
