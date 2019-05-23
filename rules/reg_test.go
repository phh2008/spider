package rules

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestReg(t *testing.T) {
	str := "李小华 先生（销售 代表）"
	str = strings.ReplaceAll(str, "(", " ")
	str = strings.ReplaceAll(str, "（", " ")
	str = strings.ReplaceAll(str, ")", "")
	str = strings.ReplaceAll(str, "）", "")
	r := strings.Split(str, " ")
	fmt.Println(r)
	fmt.Println(len(r))
	fmt.Println(r[2:])
}

var reg, _ = regexp.Compile(`(var itm = '\S+)|(var itt = '\S+)`)

func Test002(t *testing.T) {
	str := `        var itm = '13938730113';
        var itt = '86-0756-7236115';
        $(".MB1").html(itm);
        $(".Tel").html(itt);
        if ($("#clt")) {
            $("#clt").html(itt);
        }
        if ($("#clm")) {
            $("#clm").html(itm);
        }
    `
	match := reg.FindAllStringSubmatch(str, 2)
	for _, v := range match {
		for _, m := range v {
			if strings.Index(m, "var itm =") >= 0 {
				tmp := strings.ReplaceAll(m, "var itm = '", "")
				tmp = strings.ReplaceAll(tmp, "';", "")
				fmt.Println(tmp)
			} else if strings.Index(m, "var itt = '") >= 0 {
				tmp := strings.ReplaceAll(m, "var itt = '", "")
				tmp = strings.ReplaceAll(tmp, "';", "")
				fmt.Println(tmp)
			}
		}
	}
}
