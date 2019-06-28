package utils

import (
	"fmt"
	"regexp"
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

func Test3(t *testing.T) {
	reg, _ := regexp.Compile("\\[|\\]")
	str := "[上海]"
	fmt.Println(str)
	str = reg.ReplaceAllString(str, "")
	fmt.Println(str)

	arr := strings.Split(str, "/")
	fmt.Println(len(arr))
	if len(arr) > 1 {
		fmt.Println(arr[1])
	}
}
