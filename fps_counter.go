package main

import "time"

type RollingAverageFPS struct {
	window []float32
	pos    int
}

func MakeFPSAverage(size int) RollingAverageFPS {
	return RollingAverageFPS{
		make([]float32, size),
		0,
	}
}

func Fps(start, end time.Time) float32 {
	return 1_000_000 / float32(end.Sub(start).Microseconds())
}

func (roll *RollingAverageFPS) InsertFPS(fps float32) {
	roll.window[roll.pos] = fps
	roll.pos = (roll.pos + 1) % len(roll.window)
}

func (roll *RollingAverageFPS) AverageFPS() float32 {
	sum := float32(0)
	count := 0
	for _, fps := range roll.window {
		if fps < 0.0001 {
			break
		}
		sum += fps
		count++
	}
	if count == 0 {
		return float32(0)
	}
	return sum / float32(count)
}
