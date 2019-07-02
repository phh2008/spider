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

func Test4(t *testing.T) {
	reg, _ := regexp.Compile("：|\\s+|\u3002")
	str := `：
                                    主要生产枪瞄，激光测距仪，望远镜，放大镜，医疗器材，镜头等各种高中档光学镜片。
                                    `
	fmt.Println(str)
	str = reg.ReplaceAllString(str, "")
	fmt.Println(str)
}

func Test5(t *testing.T) {
	var aa string
	arr := strings.Split(aa, "-")
	fmt.Println("len: ", len(arr))
	fmt.Println("[0]: ", arr[0])
	fmt.Println("[0]=='': ", arr[0] == "")
}

func Test6(t *testing.T) {
	ret := RegBlank.ReplaceAllString("传    真", "")
	//ret := regexp.MustCompile("(&nbsp;)+").ReplaceAllString("电&nbsp;&nbsp;&nbsp;话", "")
	fmt.Println(ret)
}
