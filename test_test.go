package chuanglan_test

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var testConfig = struct {
	Origin   string   `yaml:"origin"`
	Account  string   `yaml:"account"`
	Password string   `yaml:"password"`
	Tels     []string `yaml:"tels"`
	MsgV1SendJsonMessage string `yaml:"msgV1SendJsonMessage"`
}{}

func init () {
	data, err := ioutil.ReadFile("./env.yaml") ; if err != nil {
	    panic(err)
	}
	err = yaml.Unmarshal(data, &testConfig) ; if err != nil {
		panic(err)
	}
}