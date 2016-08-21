package goconf

import (
	"errors"
	"time"
	"os"
	"path/filepath"
	"io/ioutil"
)

var (
	ErrKeyValueFormatErr = errors.New("goconf : key-value format err")
	ErrSectionFormatErr = errors.New("goconf : section format err")
)

const (
	DefaultSection = "-1"
)

type KvMap map[string]string

type Conf struct {
	modTime time.Time
	defaultSectionData KvMap
	groupData map[string]KvMap
}

func getNewConfig() *Conf {
	var config = new(Conf)
	config.defaultSectionData = make(KvMap)
	config.groupData = make(map[string]KvMap)
	
	return config
}

func (config *Conf) Read(confPath string) (bool, error) {
	if !filepath.IsAbs(confPath) {
		confAbsPath, err := filepath.Abs(confPath)
		
		if err != nil {
			return false, err
		}
		
		confPath = confAbsPath
	}
	
	fileInfo, err := os.Stat(confPath)
	if err != nil {
		return false, err
	}
	
	config.modTime = fileInfo.ModTime()
	
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		return false, err
	}
	
	return ParseByStatAndReg(config, content)
}

func (config *Conf) Get(section string, key string) string {
	if section == DefaultSection {
		return config.defaultSectionData[key]
	} else {
		return config.groupData[section][key]
	}
}