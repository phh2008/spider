package rules

import (
	"com.phh/spider/utils"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"strconv"
	"strings"
)

func init() {
	YuzhuWoodCoporate.Register()
}

var YuzhuWoodCoporate = &spider.Spider{
	Name:         "木材王国",
	Description:  "http://www.yuzhuwood.com/enterprise/search.htm",
	EnableCookie: false,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.Aid(map[string]interface{}{
				"loop": [2]int{1, 13407}, //end: 13407
				"rule": "请求URL",
			}, "请求URL")
		},
		Trunk: map[string]*spider.Rule{
			"请求URL": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.yuzhuwood.com/enterprise/search.htm?currentPage=" + strconv.Itoa(loop[0]),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					dom.Find(".ec_xxc .ec_xxc1").Each(func(i int, div *goquery.Selection) {
						ctx.SetTemp("div", div)
						ctx.Parse("企业信息")
					})
				},
			},
			"企业信息": {
				ItemFields: []string{
					"企业名称",
					"主营大类",
					"主营小类",
					"省份",
					"城市",
					"联系人",
					"联系电话",
				},
				ParseFunc: func(ctx *spider.Context) {
					div := ctx.GetTemp("div", &goquery.Selection{}).(*goquery.Selection)
					corpDiv := div.Find(".ec_xxwb").Eq(1)
					corpName := corpDiv.Find("a").First().Find("span").Text()
					corpSpan := corpDiv.ChildrenFiltered("span")
					major2 := ""
					if corpSpan.Size() > 2 {
						major2 = corpSpan.First().Text()
					}
					region, _ := div.Find(".ec_xxwb").Eq(2).Attr("title")
					//广东      省 - 广州市市,广东省 - 东莞市
					region = strings.ReplaceAll(region, " ", "") //去空格
					arr := strings.Split(region, "-")
					province := utils.If(len(arr) > 0, arr[0], "")
					city := utils.If(len(arr) > 1, arr[1], "").(string)
					if city == "市" {
						city = ""
					} else if idx := utils.UnicodeIndex(city, "市市"); idx >= 0 {
						//删除一个
						city = utils.SubString2(city, 0, idx+1)
					}
					major, _ := div.Find(".ec_xxwb").Eq(3).Attr("title")
					contact, _ := div.Find(".ec_xxwb").Eq(4).Attr("title")
					mobileTel, _ := div.Find(".ec_xxwb").Eq(5).Attr("title")

					ctx.Output(map[int]interface{}{
						0: corpName,
						1: major,
						2: major2,
						3: province,
						4: city,
						5: contact,
						6: mobileTel,
					})
				},
			},
		},
	},
}
