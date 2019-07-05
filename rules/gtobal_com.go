package rules

import (
	"bytes"
	"com.phh/spider/utils"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strconv"
	"strings"
)

func init() {
	gtobal_com.Register()
}

var gtobal_com = &spider.Spider{
	Name:        "际通宝",
	Description: "http://www.gtobal.com/company/search-ap_440000-k_%E5%AE%B6%E5%85%B7%E4%BA%94%E9%87%91-p_1.html",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 61},
				"rule": "企业分页",
			}, "企业分页")
		},
		Trunk: map[string]*spider.Rule{
			"企业分页": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.gtobal.com/company/search-ap_440000-k_%E5%AE%B6%E5%85%B7%E4%BA%94%E9%87%91-p_" + strconv.Itoa(loop[0]) + ".html",
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					ctx.GetDom().Find(".productComList3 .clearfix").Each(func(i int, selection *goquery.Selection) {
						dd := selection.Find("dd")
						dd1 := dd.First()
						a := dd1.Find("b").First().Find("a")
						url, _ := a.Attr("href")
						name := strings.TrimSpace(a.Text())
						var majorBuf bytes.Buffer
						dd1.Find("p").Last().Find("a").Each(func(i int, ap *goquery.Selection) {
							majorBuf.WriteString(strings.TrimSpace(ap.Text()))
							majorBuf.WriteString(" ")
						})
						//e.g [广东-佛山市]
						region := strings.TrimSpace(dd.Eq(1).Text())
						region = utils.RegBracket.ReplaceAllString(region, "")
						arr := strings.Split(region, "-")
						var province, city string
						if len(arr) > 0 {
							province = arr[0]
						}
						if len(arr) > 1 {
							city = arr[1]
						}
						ctx.AddQueue(&request.Request{
							Url:  url + "/contactus.html",
							Rule: "企业联系",
							Temp: map[string]interface{}{
								"name":     name,
								"major":    majorBuf.String(),
								"province": province,
								"city":     city,
							},
						})
					})
				},
			},
			"企业联系": {
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					var contact, tel, mobile, fax, address string
					classContact := dom.Find(".contactText")
					if classContact.Size() > 0 {
						classContact.Find("p").Each(func(i int, selection *goquery.Selection) {
							//e.g 真实姓名：庞先生 先生
							text := strings.TrimSpace(selection.Text())
							arr := strings.Split(text, "：")
							if len(arr) < 2 {
								return
							}
							switch arr[0] {
							case "姓名", "真实姓名":
								contact = arr[1]
							case "电话", "公司电话":
								tel = arr[1]
							case "传真", "联系传真":
								fax = arr[1]
							case "手机", "手机号码":
								mobile = arr[1]
							case "地址", "公司地址":
								address = arr[1]
							}
						})
					} else {
						dom.Find(".c_ul_class li").Each(func(i int, selection *goquery.Selection) {
							text := strings.TrimSpace(selection.Text())
							arr := strings.Split(text, "：")
							if len(arr) < 2 {
								return
							}
							switch arr[0] {
							case "联系人":
								contact = arr[1]
							case "联系电话":
								tel = arr[1]
							case "传真号码":
								fax = arr[1]
							case "手机号码":
								mobile = arr[1]
							case "公司地址":
								address = arr[1]
							}
						})
					}
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("name", ""),
						1: ctx.GetTemp("major", ""),
						2: ctx.GetTemp("province", ""),
						3: ctx.GetTemp("city", ""),
						4: contact,
						5: tel,
						6: mobile,
						7: fax,
						8: address,
					})
				},
				ItemFields: []string{
					"名称",
					"主营",
					"省",
					"市",
					"联系人",
					"电话",
					"手机",
					"传真",
					"地址",
				},
			},
		},
	},
}
