package main

import (
	"encoding/json"
	tt "github.com/huyoufu/go-timetracker"
	"io/fs"
	"io/ioutil"
	"testing"
)

func TestClean(t *testing.T) {
	//清理一下数据
	tracker := tt.NewTimeTracker("清理数据")
	tracker.StepStart("步骤一:加载数据")
	bytes, _ := ioutil.ReadFile("data3.json")
	region := &Region{}
	json.Unmarshal(bytes, region)
	tracker.StepEnd()
	tracker.StepStart("步骤二:重新写出数据")
	buffer, _ := json.MarshalIndent(region, "", "\t")
	ioutil.WriteFile("data3.json", buffer, fs.ModePerm)
	tracker.StepEnd()
	tracker.Close()
	tracker.PrintBeautiful()
}
func TestCleanMin(t *testing.T) {
	//清理一下数据
	tracker := tt.NewTimeTracker("清理数据")
	tracker.StepStart("步骤一:加载数据")
	bytes, _ := ioutil.ReadFile("data3.json")
	region := &Region{}
	json.Unmarshal(bytes, region)
	tracker.StepEnd()
	tracker.StepStart("步骤二:重新写出数据")
	buffer, _ := json.Marshal(region)
	ioutil.WriteFile("data3.min.json", buffer, fs.ModePerm)
	tracker.StepEnd()
	tracker.Close()
	tracker.PrintBeautiful()
}
