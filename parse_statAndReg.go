package goconf

import (
	"bytes"	
	"regexp"
	"fmt"
)

const version = "1.0.0"

const (
	stat_new_none = iota
	stat_new_section
	stat_new_key_value
	stat_new_comment
)

const isDebug = false

var sectionReg,_ = regexp.Compile(`^([a-zA-Z0-9\.\_\@]+)\]\s*$`)
var keyValueReg,_ = regexp.Compile(`^([a-zA-Z0-9\.\_]+)\s*=\s*(.+)\s*$`)


type StatAndRegParser struct {
	stat int
	buf bytes.Buffer
	line int
	config *Conf
	section string
	err error
	body string
}

func ParseByStatAndReg(config *Conf, content []byte) (bool, error) {
	var parser = new(StatAndRegParser)
	parser.stat = stat_new_none
	parser.config = config
	
	var lastIndex = len(content) - 1
	
	for i, char := range content {
		
		switch parser.stat {
			case stat_new_none:
				switch char {
					case ';', '#':
						parser.stat = stat_new_comment
					case ' ', '\r', '\v','\t','\f':
						
					case '\n':
						parser.line ++
					case '[':
						parser.stat = stat_new_section
					default:
						parser.buf.WriteByte(char)
						parser.stat = stat_new_key_value
				}
			case stat_new_comment:
				if char == '\n' {
					parser.stat = stat_new_none
					parser.line ++
				}
			case stat_new_section:
				if char == '\n'  {
					parser.handleSection()
					if parser.err != nil {
						return false, parser.err
					}
					parser.line ++
				} else if i == lastIndex {
					parser.buf.WriteByte(char)
					parser.handleSection()
					if parser.err != nil {
						return false, parser.err
					}
				} else {
					parser.buf.WriteByte(char)
				}
			case stat_new_key_value:
				if char == '\n'  {
					parser.handleKeyValue()
					if parser.err != nil {
						return false, parser.err
					}
					parser.line ++
				} else if i == lastIndex {
					parser.buf.WriteByte(char)
					parser.handleKeyValue()
					if parser.err != nil {
						return false, parser.err
					}
				} else {
					parser.buf.WriteByte(char)
				}
		}
	}
	
	return true, nil
}

func (parser *StatAndRegParser) handleKeyValue() {
	parser.body = parser.buf.String()
	parser.buf.Reset()
	
	var keyValueRegRet = keyValueReg.FindAllStringSubmatch(parser.body, -1)
	
	if len(keyValueRegRet) < 1 || len(keyValueRegRet[0]) < 3 {
		parser.err = ErrKeyValueFormatErr
		
		return
	}
	
	var key = keyValueRegRet[0][1]
	var value = keyValueRegRet[0][2]
	
	if isDebug {
		fmt.Println("key[", key, "] value[", value,"] section[", parser.section, "]")
	}
	
	if len(parser.section) == 0 {
		parser.config.defaultSectionData[key] = value
	} else {
		parser.config.groupData[parser.section][key] = value
	}
	
	if isDebug {
		fmt.Println("config", parser.config)
	}
	
	parser.stat = stat_new_none
}

func (parser *StatAndRegParser) handleSection() {
	parser.body = parser.buf.String()
	parser.buf.Reset()
	
	var sectionRegRet = sectionReg.FindAllStringSubmatch(parser.body, -1)
					
	if len(sectionRegRet) < 1 || len(sectionRegRet[0]) < 2 {
		parser.err = ErrSectionFormatErr
		
		return
	}
					
	parser.section = sectionRegRet[0][1]
	
	if isDebug {
		fmt.Println("section", parser.section)
	}
	
	parser.config.groupData[parser.section] = make(KvMap)
	
	parser.stat = stat_new_none
}
