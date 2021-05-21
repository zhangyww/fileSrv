/**
 * @Author: zhangyw
 * @Description:
 * @File:  grpcRepo
 * @Date: 2021/5/21 18:22
 */

package grpcRepo

import (
	"fileSrv/grpcSrv"
	"fmt"
	"google.golang.org/grpc"
)

var repoInstance = newGrpcRepo()

type GrpcRepo struct {
	uploadClient grpcSrv.UploadServiceClient
}

func newGrpcRepo() *GrpcRepo {
	return &GrpcRepo{}
}

func GetRepo() *GrpcRepo {
	return repoInstance
}

func (this *GrpcRepo) Init(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := grpc.Dial(addr, nil)
	if err != nil {
		return err
	}
	this.uploadClient = grpcSrv.NewUploadServiceClient(conn)
	return nil
}

func (this *GrpcRepo) GetUploadClient() grpcSrv.UploadServiceClient {
	return this.uploadClient
}
