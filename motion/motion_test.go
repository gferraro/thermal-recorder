package motion

import (
	"testing"

	"github.com/TheCacophonyProject/lepton3"
	"github.com/stretchr/testify/assert"
)

func defaultMotionParams() MotionConfig {
	return MotionConfig{
		TempThresh:        3000,
		DeltaThresh:       30,
		CountThresh:       8,
		NonzeroMaxPercent: 50,
		FrameCompareGap:   3,
		WarmerOnly:        false,
	}
}

func makeFrame(position, background, brightSpot int) *lepton3.Frame {
	frame := new(lepton3.Frame)

	if background != 0 {
		for y := 0; y < lepton3.FrameRows; y++ {
			for x := 0; x < lepton3.FrameCols; x++ {
				frame[y][x] = uint16(background)
			}
		}
	}

	brightness16 := uint16(background + brightSpot)

	frame[position][position] = brightness16
	frame[position+1][position] = brightness16
	frame[position+2][position] = brightness16
	frame[position][position+1] = brightness16
	frame[position+1][position+1] = brightness16
	frame[position+2][position+1] = brightness16
	frame[position][position+2] = brightness16
	frame[position+1][position+2] = brightness16
	frame[position+2][position+2] = brightness16

	return frame
}

func MovingBoxDetections(detector *motionDetector, frames, background, brightSpot int) ([]int, []bool) {
	results := make([]bool, frames)
	pixels := make([]int, frames)

	var detectionMask MotionState
	mask := &detectionMask

	for i := range results {
		results[i], mask = detector.pixelsChanged(makeFrame(10+i, background, i*brightSpot), &detectionMask)
		if results[i] {
			pixels[i] = mask.MaskSize()
		} else {
			pixels[i] = int(mask.DetectionState)
		}
	}
	return pixels, results
}

func TestRevertsToUsingSmallerFrameIntervalWhenNotEnoughFrames_OneFrame(t *testing.T) {
	config := defaultMotionParams()
	config.UseOneDiffOnly = true
	detector := NewMotionDetector(config)

	pixels, detecteds := MovingBoxDetections(detector, 5, 3300, MotionState_NoData)
	assert.Equal(t, []int{MotionState_NoData, 9, 9, 9, 18}, pixels)
	assert.Equal(t, []bool{false, true, true, true, true}, detecteds)
}

func TestNoMotionDetectedIfNothingHasChanged(t *testing.T) {
	config := defaultMotionParams()
	config.UseOneDiffOnly = true
	detector := NewMotionDetector(config)

	pixels, detecteds := MovingBoxDetections(detector, 5, 3300, 0)
	assert.Equal(t, []int{MotionState_NoData, 0, 0, 0, 0}, pixels)
	assert.Equal(t, []bool{false, false, false, false, false}, detecteds)
}

func TestIfUsingTwoFramesItOnlyCountsWhereBothFramesHaveChanged(t *testing.T) {
	config := defaultMotionParams()
	detector := NewMotionDetector(config)

	pixels, detecteds := MovingBoxDetections(detector, 6, 3300, MotionState_NoData)
	assert.Equal(t, []int{MotionState_NoData, 0, 4, 4, 5, 9}, pixels)
	assert.Equal(t, []bool{false, false, false, false, false, true}, detecteds)
}

func TestChangeCountThresh(t *testing.T) {
	config := defaultMotionParams()
	config.CountThresh = 4
	detector := NewMotionDetector(config)

	pixels, detecteds := MovingBoxDetections(detector, 6, 3300, MotionState_NoData)
	assert.Equal(t, []int{MotionState_NoData, 0, 4, 4, 5, 9}, pixels)
	assert.Equal(t, []bool{false, false, true, true, true, true}, detecteds)
}

func TestSomethingMovingWhileRecalculation_TwoPoints(t *testing.T) {
	config := defaultMotionParams()
	config.CountThresh = 4
	detector := NewMotionDetector(config)

	pixels, detecteds := MovingBoxDetections(detector, 6, 3300, MotionState_NoData)
	assert.Equal(t, []int{MotionState_NoData, 0, 4, 4, 5, 9}, pixels)

	pixels, detecteds = MovingBoxDetections(detector, 6, 3800, MotionState_NoData)
	assert.Equal(t, []int{MotionState_TooManyPoints, MotionState_NoData, 4, 4, 5, 9}, pixels)
	assert.Equal(t, []bool{false, false, true, true, true, true}, detecteds)
}

func TestIfRecalculationNothingMoving_TwoPoints(t *testing.T) {
	config := defaultMotionParams()
	detector := NewMotionDetector(config)

	// fill buffer
	pixels, detecteds := MovingBoxDetections(detector, 5, 3300, 0)
	assert.Equal(t, []int{MotionState_NoData, 0, 0, 0, 0}, pixels)

	// recalibration
	pixels, detecteds = MovingBoxDetections(detector, 5, 3800, 0)
	assert.Equal(t, []int{MotionState_TooManyPoints, MotionState_NoData, 0, 0, 0}, pixels)
	assert.Equal(t, []bool{false, false, false, false, false}, detecteds)
}
