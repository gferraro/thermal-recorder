package motion

import "github.com/TheCacophonyProject/lepton3"

const MotionState_Unknown = 1000
const MotionState_Detected = 2000
const MotionState_NoData = 3000
const MotionState_TooManyPoints = 4000

type MotionState struct {
	DetectionState int
	Mask           [lepton3.FrameRows][lepton3.FrameCols]bool
}

func (ms *MotionState) Zero() {
	ms.DetectionState = MotionState_Unknown
	for y := 0; y < lepton3.FrameRows; y++ {
		for x := 0; x < lepton3.FrameCols; x++ {
			ms.Mask[y][x] = false
		}
	}
}

func (ms *MotionState) MaskSize() int {
	var count int
	for y := 0; y < lepton3.FrameRows; y++ {
		for x := 0; x < lepton3.FrameCols; x++ {
			if ms.Mask[y][x] {
				count++
			}
		}
	}
	return count
}
