package configure

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type setting struct {
	MaxLine        int // 最大支持列表行数
	HttpWaitPeriod int // 请求最大等待时间 秒
}

var Conf = setting{
	MaxLine:        100,
	HttpWaitPeriod: 5,
}

func Init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("未发现配置文件")
		return
	}
	maxLine, _ := cfg.Section("").Key("max_line").Int()
	if maxLine > 0 {
		Conf.MaxLine = maxLine
	}
	waitPeriod, _ := cfg.Section("").Key("wait_period").Int()
	if waitPeriod > 0 {
		Conf.HttpWaitPeriod = waitPeriod
	}
}
