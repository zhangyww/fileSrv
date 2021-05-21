/**
 * @Author: zhangyw
 * @Description:
 * @File:  main
 * @Date: 2021/5/20 10:42
 */

package main

import (
	"fileSrv/collector"
	"fileSrv/config"
	"fileSrv/fileServer"
	"log"
)

func main() {
	err := config.GetConfig().Load("./conf/fileSrv.xml")
	if err != nil {
		log.Fatal(err)
		return
	}

	conf := config.GetConfig()

	fileInfoCollector := collector.NewFileInfoCollector(conf.RootDir, conf.DirLevel, conf.LoadHistory, &conf.Ignore)
	fileInfoCollector.InitCollector(&collector.UploadCollector{})
	err = fileInfoCollector.StartCollect()
	if err != nil {
		log.Fatal(err)
	}

	server := fileServer.NewFileServer(conf.IP, conf.Port, conf.RootDir)
	server.ListenAndServe()
}
