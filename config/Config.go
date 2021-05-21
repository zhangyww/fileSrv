/**
 * @Author: zhangyw
 * @Description:
 * @File:  Config
 * @Date: 2021/5/20 11:25
 */

package config

import (
	"encoding/xml"
	"os"
)

var configInstance Config

type Config struct {
	IP          string       `xml:"IP"`
	Port        int          `xml:"Port"`
	RootDir     string       `xml:"RootDir"`
	DirLevel    int          `xml:"DirLevel"`
	LoadHistory bool         `xml:"LoadHistory"`
	Ignore      IgnoreConfig `xml:"Ignore"`
}

type IgnoreConfig struct {
	Files IgnorePatternConfig `xml:"Files"`
	Dirs  IgnorePatternConfig `xml:"Dirs"`
}

type IgnorePatternConfig struct {
	Prefix []string `xml:"Prefix"`
	Suffix []string `xml:"Suffix"`
}

func (this *Config) Load(fpath string) error {
	fileBytes, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(fileBytes, this)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &configInstance
}
