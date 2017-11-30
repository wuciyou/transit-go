package config

import (
	"bufio"
	"flag"
	"log"
	"io"
	"os"
	"strings"
)
var defaultConfigFile = "./config/default.ini"
var configName  = flag.String("configName", defaultConfigFile, "自定义配置文件路径")
var is_print_config  = flag.Bool("is_print_config", false, "是否打印配置信息")
var isFinishParse = false
func Parse() {
	flag.Parse()
	if *configName != "" && isFinishParse == false {
		isFinishParse = true
		configFile, err := os.Open(*configName)
		if err != nil {
			if *configName != defaultConfigFile {
				panic(err)
			}
			return
		}
		fileReader := bufio.NewReader(configFile)
		var argsData []string
		var argsDataMap = make(map[string]int)
		argsData = append(argsData, os.Args[0])
		for {
			data, _, err := fileReader.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			dataStr := strings.TrimSpace(string(data))
			if dataStr == "" || strings.Index(dataStr, "#") == 0 {
				continue
			}
			argsName := strings.TrimSpace(dataStr[0:strings.Index(dataStr, "=")])
			if index, exist := argsDataMap[argsName]; exist {
				argsData[index] = "-" + dataStr
			} else {
				argsData = append(argsData, "-"+dataStr)
				argsDataMap[argsName] = len(argsData) - 1
			}
		}


		var allFlag = make(map[string]string)
		flag.VisitAll(func(f *flag.Flag) {
			allFlag[f.Name] = f.Name
		})
		for name, index := range argsDataMap {
			if _, exist := allFlag[name]; !exist {

				argsData[index] = ""

			}
		}
		os.Args = append(argsData, os.Args[1:]...)

		for index, args := range os.Args {
			if args == "" {
				if len(os.Args)-1 < index{
					continue
				}
				os.Args = append(os.Args[0:index], os.Args[index+1:]...)
			}
		}

		flag.Parse()

	}

	if *is_print_config {
		flag.VisitAll(func(f *flag.Flag){
			log.Printf("config item:%+v \n",f)
		})
	}

}
