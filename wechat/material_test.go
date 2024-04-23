package wechat

import "testing"

func TestAdd(t *testing.T) {
	t.Run("TestAdd", func(t *testing.T) {
		Add("test.jpg")
	})
}
