package utils

import (
	"fmt"
	"strings"
	"testing"
)

func Test1(t *testing.T) {

	str := "广东      省 - 广州市市"
	str = strings.ReplaceAll(str, " ", "") //去空格
	idx := UnicodeIndex(str, "市市")
	tmp := SubString2(str, 0, idx+1)
	fmt.Println(tmp)

}

func Test2(t *testing.T) {
	str := fmt.Sprintf("xxxx%daaaa", 12)
	fmt.Println(str)
}
