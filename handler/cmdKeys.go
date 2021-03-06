package handler

import (
	"github.com/dalebao/gRedis-cli/analysis"
	"github.com/dalebao/gRedis-cli/pkg"
	"strconv"
	"strings"
)

type Keys struct {
	ST     [][]string //String 类型
	HT     [][]string //Hash 类型
	SetT   [][]string //set 类型
	LT     [][]string //list 类型
	ZST    [][]string //Zset 类型
	Expect map[string]bool
	Only   map[string]bool
	Sort   int
	Limit  int
	Export string
}

func (keys *Keys) Set(e map[string]string) {
	if e["only"] != "" {
		only := make(map[string]bool)
		x := strings.Split(e["only"], ",")
		for _, v := range x {
			only[v] = true
		}
		keys.Only = only
	}

	if e["expect"] != "" {
		expect := make(map[string]bool)
		x := strings.Split(e["expect"], ",")
		for _, v := range x {
			expect[v] = true
		}
		keys.Expect = expect
	}

	if e["sort"] != "" {
		if e["sort"] == "desc" {
			keys.Sort = 1
		}
	}

	if e["limit"] != "" {
		var err error
		keys.Limit, err = strconv.Atoi(e["limit"])
		if err != nil {
			keys.Limit = -1
		}
	} else {
		keys.Limit = -1
	}

	if e["export"] != "" {
		keys.Export = e["export"]
	}
}

func (keys *Keys) DiffType(res []string) (data [][]string) {
	lOnly := len(keys.Only)
	lExpect := len(keys.Expect)
	var l int
	if keys.Limit != -1 {
		l = keys.Limit
	} else {
		l = len(res)
	}
	for _, v := range res[:l] {
		rType, _ := pkg.Client.Type(v).Result()

		if lExpect != 0 && keys.Expect[rType] == true {
			continue
		}

		if lOnly != 0 && keys.Only[rType] == false {
			continue
		}
		ttl, _ := pkg.Client.TTL(v).Result()
		mem := analysis.Analysis(v)

		info := []string{rType, v, ttl.String(), mem}

		switch rType {
		case "string":
			keys.ST = append(keys.ST, info)
		case "hash":
			keys.HT = append(keys.HT, info)
		case "list":
			keys.LT = append(keys.LT, info)
		case "set":
			keys.SetT = append(keys.SetT, info)
		case "zset":
			keys.ZST = append(keys.ZST, info)
		}
	}
	//desc
	if keys.Sort == 1 {
		data = append(data, keys.ZST...)
		data = append(data, keys.ST...)
		data = append(data, keys.SetT...)
		data = append(data, keys.LT...)
		data = append(data, keys.HT...)
	} else {
		data = append(data, keys.HT...)
		data = append(data, keys.LT...)
		data = append(data, keys.SetT...)
		data = append(data, keys.ST...)
		data = append(data, keys.ZST...)
	}
	return data
}
