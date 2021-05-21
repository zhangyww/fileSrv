/**
 * @Author: zhangyw
 * @Description:
 * @File:  FileInfo
 * @Date: 2021/5/20 16:20
 */

package model

type FileOp int

const (
	FILE_CREATE FileOp = 1
	FILE_DELETE FileOp = 2
)

type FileInfo struct {
	FilePath string
	Op       FileOp
}
