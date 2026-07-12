package config

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	conf, err := Load("testdata/conf.yaml")
	if err != nil {
		t.Error("Working config brackes")
	}
	if conf.Backends[0].URL != "http://localhost:8081" ||
	len(conf.Backends) != 1 ||
	conf.HealthCheckInterval != time.Second * 5 ||
	conf.Strategy != "round-robin" ||
	conf.Port != 8080 {
		t.Error("Wrong config import")
	}
}

func TestEmptyConfig(t *testing.T) {
	conf, err := Load("testdata/empty_conf.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestUnformatConfig(t *testing.T) {
	conf, err := Load("testdata/unformat_conf.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestBrokenConfig(t *testing.T) {
	conf, err := Load("testdata/broken_conf.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestNonConfig(t *testing.T) {
	conf, err := Load("testdata/broken_co.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of wrong config way")
	}
}
