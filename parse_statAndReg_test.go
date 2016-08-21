package goconf

import (
	"testing"
	"fmt"
)


func TestDefaultSection(t *testing.T) {
	var content = []byte("key_0 = value_0\n[section_1]\n key_1 = value_1\n")
	var config = new(Conf)
	config.defaultSectionData = make(KvMap)
	config.groupData = make(map[string]KvMap)
	
	status,err := ParseByStatAndReg(config, content)
	
	fmt.Println(status, err)
}

func TestErrKey(t *testing.T) {
	var contents = [][]byte{
		[]byte("key 0 = value"),
		[]byte("*aaa = value"),
	}
	
	batchTest(t, contents, ErrKeyValueFormatErr)
}

func TestSpaceKeyValue(t *testing.T) {
	var contents = [][]byte{
		[]byte("dftKey=dftValue"),
		[]byte("dftKey= dftValue"),
		[]byte("dftKey = dftValue"),
		[]byte("dftKey =dftValue"),
	}
	
	batchTest(t, contents, nil)
}

func TestQuoteValue(t *testing.T) {
	var contents = [][]byte{
		[]byte("dftKey='dftValue'"),
		[]byte("dftKey= 'dftValue'"),
		[]byte("dftKey= ' dftValue'"),
		[]byte("dftKey=' dftValue'"),
		[]byte("dftKey='dftValue '"),
		[]byte("dftKey=' dftValue '"),
		[]byte("dftKey=' dftValue ' "),
		[]byte("dftKey=' dft V alue ' "),
	}
	
	batchTest(t, contents, nil)
}

func batchTest(t *testing.T, contents [][]byte, expectErr error) {
	var config = new(Conf)
	config.defaultSectionData = make(KvMap)
	config.groupData = make(map[string]KvMap)
	
	for _, content := range contents {
		status, err := ParseByStatAndReg(config, content)
		
		t.Log("content", string(content))
		
		if err == expectErr {
			t.Log(status, err)
		} else {
			t.Fatal(status, err)
		}
	}
}