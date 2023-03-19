package processor

import (
	"github.com/golang/mock/gomock"
	"github.com/mtvarkovsky/goodjob/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_V1Processor_StartStop(t *testing.T) {
	ctrl := gomock.NewController(t)

	queue := mocks.NewMockQueue(ctrl)

	processor := NewV1Processor(queue, nil, nil)

	queue.EXPECT().GetNextJob().Return(nil, nil).MinTimes(1)

	err := processor.Start()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	err = processor.Stop()
	assert.NoError(t, err)
}
