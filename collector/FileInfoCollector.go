/**
 * @Author: zhangyw
 * @Description:
 * @File:  FileInfoCollector
 * @Date: 2021/5/20 16:06
 */

package collector

import (
	"fileSrv/config"
	"fileSrv/model"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileInfoCollector struct {
	RootDir     string
	TrimRootDir string
	DirLevel    int
	LoadHistory bool
	Ignore      *config.IgnoreConfig
	collector   Collector

	watcher *fsnotify.Watcher

	watchDirMapLock sync.Mutex
	watchDirMap     map[string]*model.DirEntry
}

func NewFileInfoCollector(rootDir string, dirLevel int, loadHistory bool, ignore *config.IgnoreConfig) *FileInfoCollector {
	collector := &FileInfoCollector{
		RootDir:     rootDir,
		DirLevel:    dirLevel,
		LoadHistory: loadHistory,
		Ignore:      ignore,
	}
	return collector
}

func (this *FileInfoCollector) StartCollect() error {
	absDir, err := filepath.Abs(this.RootDir)
	if err != nil {
		return err
	}
	this.RootDir = filepath.Clean(absDir)
	if strings.HasSuffix(this.RootDir, "/") {
		this.TrimRootDir = this.RootDir
	} else {
		this.TrimRootDir = this.RootDir + "/"
	}
	//fmt.Println("rootDir:", this.RootDir)

	dirInfo, err := os.Stat(this.RootDir)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("rootDir %s is not Directory", this.RootDir)
	}

	if this.watcher == nil {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		this.watcher = watcher
		go this.watchLoop()
	}
	if this.watchDirMap == nil {
		this.watchDirMap = make(map[string]*model.DirEntry)
	}
	err = this.watcher.Add(this.RootDir)
	this.watchDirMapLock.Lock()
	this.watchDirMap[this.RootDir] = model.NewDirEntry(this.DirLevel)
	this.watchDirMapLock.Unlock()

	rootFlag := true
	filepath.Walk(this.RootDir, func(path string, info os.FileInfo, err error) error {
		if rootFlag {
			rootFlag = false
			return nil
		}
		//fmt.Println(path, this.LoadHistory, this.collector)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		if info.IsDir() {
			if this.DirLevel > 0 {
				err := this.WalkCollect(path, this.DirLevel-1, this.LoadHistory)
				//subCollector := NewFileInfoCollector(path, this.DirLevel-1, this.LoadHistory, nil)
				//subCollector.InitWatcher(this.watcher)
				//subCollector.InitCollector(this.collector)
				//err := subCollector.StartCollect()
				if err != nil {
					log.Fatal(err)
				}
			}
			return filepath.SkipDir
		}

		if this.LoadHistory {
			//todo load histroy
			if this.collector != nil {
				this.Collect(&model.FileInfo{
					FilePath: path,
					Op:       model.FILE_CREATE,
				})
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (this *FileInfoCollector) WalkCollect(dirPath string, dirLevel int, loadHistory bool) error {
	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}
	rootDir := filepath.Clean(absDir)
	//fmt.Println("rootDir:", rootDir)

	dirInfo, err := os.Stat(rootDir)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("rootDir %s is not Directory", rootDir)
	}

	err = this.watcher.Add(rootDir)
	this.watchDirMapLock.Lock()
	this.watchDirMap[rootDir] = model.NewDirEntry(dirLevel)
	this.watchDirMapLock.Unlock()

	rootFlag := true
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if rootFlag {
			rootFlag = false
			return nil
		}
		//fmt.Println(path, this.LoadHistory, this.collector)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		if info.IsDir() {
			if dirLevel > 0 {
				this.WalkCollect(path, dirLevel-1, true)
				if err != nil {
					log.Fatal(err)
				}
			}
			return filepath.SkipDir
		}

		if loadHistory {
			//todo load histroy
			if this.collector != nil {
				this.watchDirMapLock.Lock()
				delete(this.watchDirMap, path)
				this.watchDirMapLock.Unlock()
				this.Collect(&model.FileInfo{
					FilePath: path,
					Op:       model.FILE_CREATE,
				})
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (this *FileInfoCollector) Collect(fileInfo *model.FileInfo) {
	if this.collector != nil && !this.IsIgnored(fileInfo) {
		fileInfo.FilePath = strings.TrimPrefix(fileInfo.FilePath, this.TrimRootDir)
		this.collector.Collect(fileInfo)
	}
}

func (this *FileInfoCollector) IsIgnored(fileInfo *model.FileInfo) bool {
	if fileInfo == nil {
		return true
	}
	if this.Ignore == nil {
		return false
	}

	baseName := filepath.Base(fileInfo.FilePath)
	//fmt.Println(baseName, fmt.Sprintf("%+v", this.Ignore))
	for _, filePrefix := range this.Ignore.Files.Prefix {
		if strings.HasPrefix(baseName, filePrefix) {
			return true
		}
	}
	for _, fileSuffix := range this.Ignore.Files.Prefix {
		if strings.HasSuffix(baseName, fileSuffix) {
			return true
		}
	}

	return false
}

func (this *FileInfoCollector) watchLoop() {
	for {
		select {
		case event, ok := <-this.watcher.Events:
			if !ok {
				return
			}
			//log.Println("event:", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				stat, err := os.Stat(event.Name)
				if err != nil {
					continue
				}
				if stat.IsDir() {
					//fmt.Println("path:", event.Name)
					//fmt.Println("base:", filepath.Base(event.Name))
					//fmt.Println("dir:", filepath.Dir(event.Name))
					pathDir := filepath.Dir(event.Name)
					//pathBase := filepath.Base(event.Name)

					this.watchDirMapLock.Lock()
					dirEntry, exists := this.watchDirMap[pathDir]
					if exists {
						dirEntry.Deleted = false
					}
					this.watchDirMapLock.Unlock()

					if !exists {
						continue
					}
					//this.watcher.Add(event.Name)

					err := this.WalkCollect(event.Name, dirEntry.DirLevel-1, true)
					//subCollector := NewFileInfoCollector(event.Name, dirLevel-1, true, nil)
					//subCollector.InitWatcher(this.watcher)
					//subCollector.InitCollector(this.collector)
					//subCollector.StartCollect()
					if err != nil {
						log.Fatal(err)
					}

					//rootFlag := true
					//filepath.Walk(event.Name, func(path string, info os.FileInfo, err error) error {
					//	if rootFlag {
					//		rootFlag = false
					//		return nil
					//	}
					//	if err != nil {
					//		log.Fatal(err)
					//		return nil
					//	}
					//
					//	if info.IsDir() {
					//		if dirLevel > 0 {
					//			subCollector := NewFileInfoCollector(path, dirLevel-1, true, nil)
					//			subCollector.InitWatcher(this.watcher)
					//			subCollector.InitCollector(this.collector)
					//			subCollector.StartCollect()
					//		}
					//		return filepath.SkipDir
					//	}
					//
					//	// always loadHistory
					//	if this.LoadHistory {
					//		//todo load histroy
					//		if this.collector != nil {
					//			this.collector.Collect(&model.FileInfo{
					//				FilePath: event.Name,
					//				Op:       model.FILE_CREATE,
					//			})
					//		}
					//	}
					//
					//	return nil
					//})

					//log.Println("create dir:", event.Name)
				} else {
					//log.Println("create file:", event.Name)
					this.watchDirMapLock.Lock()
					delete(this.watchDirMap, event.Name)
					this.watchDirMapLock.Unlock()
					if this.collector != nil {
						this.Collect(&model.FileInfo{
							FilePath: event.Name,
							Op:       model.FILE_CREATE,
						})
					}
				}

			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				deleted := false
				this.watchDirMapLock.Lock()
				dirEntry, exists := this.watchDirMap[event.Name]
				if exists {
					if dirEntry.Deleted {
						deleted = true
					}
					dirEntry.Deleted = true
				}
				this.watchDirMapLock.Unlock()

				if exists {
					if !deleted {
						this.watcher.Remove(event.Name)
					}
					//log.Println("delete dir:", event.Name)
				} else {
					//log.Println("delete file:", event.Name)
					if this.collector != nil {
						this.Collect(&model.FileInfo{
							FilePath: event.Name,
							Op:       model.FILE_DELETE,
						})
					}
				}
			}
		case err, ok := <-this.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (this *FileInfoCollector) InitWatcher(watcher *fsnotify.Watcher) {
	this.watcher = watcher
}

func (this *FileInfoCollector) InitCollector(collector Collector) {
	this.collector = collector
}
