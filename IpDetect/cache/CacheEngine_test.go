package cache

import (
	"fmt"
	"strings"
	"testing"
)

func TestNormal(t *testing.T) {
	var msgToTest []string
	engine := NewDefault(3, func(msg string) {
		msgToTest = append(msgToTest, fmt.Sprintf("[Handle] ==> %s", msg))
	})

	engine.CacheInternal("abc", func(msg string) {
		msgToTest = append(msgToTest, fmt.Sprintf("[Match] ==> %s", msg))
	})

	var lines = strings.Split("[Handle] ==> init engine\n[Match] ==> abc\n[Handle] ==> cached for 1", "\n")

	loopTest(msgToTest, lines, t)
}

/*
Test case for reaching max detection cache period
*/
func TestCacheCase1(t *testing.T) {
	var msgToTest []string
	engine := NewDefault(3, func(msg string) {
		msgToTest = append(msgToTest, fmt.Sprintf("[Handle] ==> %s", msg))
	})

	var logFunc = func(msg string) {
		msgToTest = append(msgToTest, fmt.Sprintf("[Match] ==> %s", msg))
	}
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("abc", logFunc)

	var lines = strings.Split("[Handle] ==> init engine\n[Match] ==> abc\n[Handle] ==> cached for 1\n[Handle] ==> ip address cached\n[Handle] ==> cached for 2\n[Handle] ==> ip address cached\n[Handle] ==> cached for 3\n[Handle] ==> ip address send\n[Match] ==> abc\n[Handle] ==> cached for 0\n[Handle] ==> ip address cached\n[Handle] ==> cached for 1", "\n")

	loopTest(msgToTest, lines, t)
}

/*
Test case for the ip really changed
*/
func TestCacheCase2(t *testing.T) {
	var msgToTest []string
	engine := NewDefault(3, func(msg string) {
		msgToTest = append(msgToTest, fmt.Sprintf("[Handle] ==> %s", msg))
	})

	var logFunc = func(msg string) {
		msgToTest = append(msgToTest, fmt.Sprintf("[Match] ==> %s", msg))
	}
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("xxxx", logFunc)
	engine.CacheInternal("abc", logFunc)
	engine.CacheInternal("abc", logFunc)

	var lines = strings.Split("[Handle] ==> init engine\n[Match] ==> abc\n[Handle] ==> cached for 1\n[Handle] ==> ip address cached\n[Handle] ==> cached for 2\n[Handle] ==> ip address changed\n[Match] ==> xxxx\n[Handle] ==> cached for 1\n[Handle] ==> ip address changed\n[Match] ==> abc\n[Handle] ==> cached for 1\n[Handle] ==> ip address cached\n[Handle] ==> cached for 2", "\n")

	loopTest(lines, msgToTest, t)
}

func loopTest(expecting []string, target []string, t *testing.T) {
	if len(expecting) != len(target) {
		t.Fatalf("Expecting size [%d] but found [%d]", len(expecting), len(target))
	}
	for i := 0; i < len(expecting); i++ {
		line := target[i]
		log := expecting[i]
		if log != line {
			t.Fatalf("Expecting [%s] but found [%s]", log, line)
		}
	}
}
