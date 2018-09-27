package throttle

import (
	"github.com/TheCacophonyProject/lepton3"
	"github.com/TheCacophonyProject/thermal-recorder/motion"
)

type NoveltyDetector struct {
	step   uint
	memory lepton3.Frame
}

func NewNoveltyDetector() *NoveltyDetector {
	return &NoveltyDetector{
		step: 4,
	}
}

func (nd *NoveltyDetector) AddToCounts(state motion.MotionState) {
	for y := 0; y < lepton3.FrameRows; y++ {
		for x := 0; x < lepton3.FrameCols; x++ {
			if state.Mask[y][x] {
				nd.memory[y>>nd.step][x>>nd.step]++
			}
		}
	}
}
