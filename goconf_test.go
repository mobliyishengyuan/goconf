package goconf

import (
	"testing"
)

func TestSimple(t *testing.T) {	
	var config = GetNewConfig()
	var status, err = config.Read("simple.ini")
	
	if !status {
		t.Fatal(err)
	}

	valueCase("section_1", "key_a", config, t, true)
	valueCase("section_1", "key_b", config, t, true)
	valueCase("section_3", "key_c", config, t, true)
	
	valueCase("section_1", "key_c", config, t, false)
	valueCase("section_2", "key_a", config, t, false)
}

func valueCase(section string, key string, config *Conf, t *testing.T, expectStatus bool) {
	var value, status = config.Get(section, key)
	
	if expectStatus == status {
		t.Log("section", section, "key", key, "value", value, "status", status)
	} else {
		t.Fatal("section", section, "key", key, "value", value, "status", status)
	}
}