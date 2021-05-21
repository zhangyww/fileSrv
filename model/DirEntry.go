/**
 * @Author: zhangyw
 * @Description:
 * @File:  DirEntry
 * @Date: 2021/5/21 15:08
 */

package model

type DirEntry struct {
	DirLevel int
	Deleted  bool
}

func NewDirEntry(dirLevel int) *DirEntry {
	return &DirEntry{
		DirLevel: dirLevel,
		Deleted:  false,
	}
}
