/**
 * @Author: zhangyw
 * @Description:
 * @File:  FileServer
 * @Date: 2021/5/20 10:51
 */

package fileServer

import (
	"fmt"
	"log"
	"net/http"
)

type FileServer struct {
	IP      string
	Port    int
	RootDir string
}

func NewFileServer(ip string, port int, rootDir string) *FileServer {
	server := &FileServer{
		IP:      ip,
		Port:    port,
		RootDir: rootDir,
	}
	return server
}

func (this *FileServer) ListenAndServe() error {
	http.Handle("/", http.FileServer(http.Dir(this.RootDir)))

	addr := fmt.Sprintf("%s:%d", this.IP, this.Port)
	log.Println("addr", addr)
	return http.ListenAndServe(addr, nil)
}
