/**
 * @Author: zhangyw
 * @Description:
 * @File:  UploadCollector
 * @Date: 2021/5/21 10:28
 */

package collector

import (
	"context"
	"fileSrv/config"
	"fileSrv/grpc"
	"fileSrv/grpcSrv"
	"fileSrv/model"
	"fileSrv/repo/cacheRepo"
	"fileSrv/repo/grpcRepo"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
	"time"
)

type UploadCollector struct {
	storeQueue    chan *model.FileInfo
	uploadQueueDB *leveldb.DB

	lastUnixNano int64
	keyIdx       int64
}

func NewUploadCollector() *UploadCollector {
	queueDB, err := leveldb.OpenFile("./leveldb/upload", nil)
	if err != nil {
		return nil
	}
	collector := &UploadCollector{
		uploadQueueDB: queueDB,
		storeQueue:    make(chan *model.FileInfo, 100000),
	}

	return collector
}

func (this *UploadCollector) Collect(fileInfo *model.FileInfo) {
	this.storeQueue <- fileInfo
}

func (this *UploadCollector) Start() error {
	if this.storeQueue == nil {
		return fmt.Errorf("storeQueue is nil")
	}
	if this.uploadQueueDB == nil {
		return fmt.Errorf("uploadQueueDB is nil")
	}

	go this.storeLoop()
	go this.uploadLoop()

	return nil
}

func (this *UploadCollector) newDBKey() []byte {
	nowNano := time.Now().UnixNano()

	if nowNano == this.lastUnixNano {
		this.keyIdx++
	} else {
		this.lastUnixNano = nowNano
		this.keyIdx = 0
	}

	keyStr := fmt.Sprintf("%013s%07s",
		strconv.FormatInt(nowNano, 36),
		strconv.FormatInt(this.keyIdx, 36))
	return []byte(keyStr)
}

func (this *UploadCollector) storeLoop() {
	cache := cacheRepo.GetRepo()

	for fileInfo := range this.storeQueue {

		cacheValue, exists := cache.Get(fileInfo.FilePath)
		if exists {
			lastestOp := cacheValue.(model.FileOp)
			if lastestOp == fileInfo.Op {
				continue
			}
		}
		cache.SetWithTTL(fileInfo.FilePath, fileInfo.Op, time.Second)

		fmt.Printf("op:%d path:%s\n", fileInfo.Op, fileInfo.FilePath)

		jsonBytes, err := fileInfo.MarshalJSON()
		if err != nil {
			log.Fatal(err)
			continue
		}
		err = this.uploadQueueDB.Put(this.newDBKey(), jsonBytes, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (this *UploadCollector) uploadLoop() {

	conf := config.GetConfig()

	iter := this.uploadQueueDB.NewIterator(nil, nil)
	for iter.Next() {
		keyBytes := iter.Key()
		valueBytes := iter.Value()

		fileInfo := model.FileInfo{}
		err := fileInfo.UnmarshalJSON(valueBytes)
		if err != nil {
			err = this.uploadQueueDB.Delete(keyBytes, nil)
			log.Fatal(err)
			continue
		}

		uploadFileInfo := grpcSrv.FileInfo{
			Ip:   conf.IP,
			Port: int32(conf.Port),
			Path: fileInfo.FilePath,
			Op:   int32(fileInfo.Op),
		}

		client := grpcRepo.GetRepo().GetUploadClient()
		response, err := client.UploadFileInfo(context.Background(), &uploadFileInfo)
		if err != nil {
			break
		}
		if response.Code != 0 {
			break
		}
	}
	iter.Release()
}
