package service

import (
	"testing"

	"github.com/KaySar12/NextZen-Common/utils/logger"
	"go.uber.org/goleak"
)

func TestSearch(t *testing.T) {
	logger.LogInitConsoleOnly()
	goleak.VerifyNone(t)

	if d, e := NewOtherService().Search("test"); e != nil || d == nil {
		t.Error("then test search error", e)
	}
}
