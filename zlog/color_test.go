/*
 * @Author: seekwe
 * @Date:   2019-06-06 15:10:39
 * @Last Modified by:   seekwe
 * @Last Modified time: 2019-06-06 15:20:55
 */

package zlog

import (
	"fmt"
	"os"
	"testing"

	zls "github.com/sohaha/zlsgo"
)

func TestColor(t *testing.T) {
	T := zls.NewTest(t)
	testText := "ok"
	_ = os.Setenv("ConEmuANSI", "ON")
	bl := IsSupportColor()
	T.Equal(true, bl)
	OutAllColor()
	T.Equal(fmt.Sprintf("\x1b[%dm%s\x1b[0m", ColorGreen, testText), ColorTextWrap(ColorGreen, testText))

	DisableColor = true
	bl = IsSupportColor()
	T.Equal(false, bl)
	OutAllColor()
	T.Equal(testText, ColorTextWrap(ColorGreen, testText))
}
