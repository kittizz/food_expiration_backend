package test

import (
	"testing"
	"time"
)

func TestUserUsecase(t *testing.T) {
	t.Log((time.Hour*7 + (time.Minute * 32)).Minutes())

}
