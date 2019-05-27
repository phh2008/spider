package rules

import (
	"encoding/json"
	"fmt"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"github.com/henrylee2cn/pholcus/logs"
	"io/ioutil"
	"net/http"
)

func init() {
	PingPaiHongMuWang.Register()
}

var PingPaiHongMuWang = &spider.Spider{

	Name:        "品牌红木网",
	Description: "http://www.328f.cn/company/",
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			//获取企业信息
			res, err := http.Get("http://www.328f.cn/company/api/company.ashx?method=list&pNum=1&pSize=60")
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			defer res.Body.Close()
			content, _ := ioutil.ReadAll(res.Body)
			text := string(content)
			logs.Log.Informational(text)
			var result map[string]interface{}
			json.Unmarshal(content, &result)
			list := result["list"].([]interface{})
			var companyIdList []string
			for _, v := range list {
				m := v.(map[string]interface{})
				companyIdList = append(companyIdList, fmt.Sprintf("%d", int64(m["company_id"].(float64))))
			}
			ctx.Aid(map[string]interface{}{
				"companyIdList": companyIdList,
				"rule":          "企业大全",
			}, "企业大全")
		},
		Trunk: map[string]*spider.Rule{
			"企业大全": {
				AidFunc: func(ctx *spider.Context, aid map[string]interface{}) interface{} {
					for _, companyId := range aid["companyIdList"].([]string) {
						ctx.AddQueue(&request.Request{
							Url:  fmt.Sprintf("http://www.328f.cn/b2b/%s.html", companyId),
							Rule: aid["rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					dom := ctx.GetDom()
					div := dom.Find(".brand_left")
					ctx.SetTemp("div", div)
					ctx.Parse("企业信息")
				},
			},
			"企业信息": {
				ItemFields: []string{
					"企业",
					"简介",
					"联系人",
					"联系电话",
					"地区",
				},
				ParseFunc: func(ctx *spider.Context) {
					div := ctx.GetTemp("div", &goquery.Selection{}).(*goquery.Selection)
					corpName := div.Find(".item_pic strong").Text()
					info := div.Find(".item_text span").Text()
					p := div.Find(".item_info p")
					address := p.Eq(0).Text()
					tel := p.Eq(1).Text()
					contact := p.Eq(2).Text()

					ctx.Output(map[int]interface{}{
						0: corpName,
						1: info,
						2: contact,
						3: tel,
						4: address,
					})
				},
			},
		},
	},
}
