package modparser_test

import (
	"testing"

	"github.com/tidwall/gjson"
	"github.com/wxb/got/modparser"
	"gopkg.in/yaml.v3"
)

const (
	// 测试go.mod文件
	TestModFile = "./gomod.demo"
)

func TestParser_Parse(t *testing.T) {
	// 创建一个Parser对象
	parser, err := modparser.New(TestModFile)
	if err != nil {
		t.Fatalf("无法创建Parser对象：%s", err.Error())
	}

	// 调用Parse方法
	modFile := parser.Parse()

	// 检查解析结果是否为nil
	if modFile == nil {
		t.Error("解析结果为nil")
	} else {
		// 检查解析结果的某些属性是否符合预期
		expectedName := "github.com/wxb/example"
		if modFile.Module.Mod.Path != expectedName {
			t.Errorf("解析结果的Module.Mod.Path不符合预期。期望：%s，实际：%s", expectedName, modFile.Module.Mod.Path)
		}
	}
}

func TestParser_JsonMarshal(t *testing.T) {
	// 创建一个Parser对象
	parser, err := modparser.New(TestModFile)
	if err != nil {
		t.Fatalf("无法创建Parser对象：%s", err.Error())
	}

	// 调用JsonMarshal方法
	jsonData, err := parser.JsonMarshal()
	if err != nil {
		t.Fatalf("JsonMarshal方法出错：%s", err.Error())
	}

	result := gjson.ParseBytes(jsonData)
	// 检查结果是否符合预期
	if expected := "github.com/PuerkitoBio/goquery"; expected != result.Get("Require.0.Mod.Path").String() {
		t.Errorf("JsonMarshal结果不符合预期。期望：%s，实际：%s", expected, result.Get("Require.0.Mod.Path").String())
	}
}

func TestParser_JsonMarshalIndent(t *testing.T) {
	// 创建一个Parser对象
	parser, err := modparser.New(TestModFile)
	if err != nil {
		t.Fatalf("无法创建Parser对象：%s", err.Error())
	}

	// 调用JsonMarshalIndent方法
	jsonData, err := parser.JsonMarshalIndent("", "  ")
	if err != nil {
		t.Fatalf("JsonMarshalIndent方法出错：%s", err.Error())
	}

	result := gjson.ParseBytes(jsonData)
	// 检查结果是否符合预期
	if expected := true; expected != result.Get("Require.1.Indirect").Bool() {
		t.Errorf("JsonMarshalIndent结果不符合预期。期望：%v，实际：%v", expected, result.Get("Require.1.Indirect").Bool())
	}
}

func TestParser_YamlMarshal(t *testing.T) {
	// 创建一个Parser对象
	parser, err := modparser.New(TestModFile)
	if err != nil {
		t.Fatalf("无法创建Parser对象：%s", err.Error())
	}

	// 调用YamlMarshal方法
	yamlData, err := parser.YamlMarshal()
	if err != nil {
		t.Fatalf("YamlMarshal方法出错：%s", err.Error())
	}

	// 检查结果是否符合预期
	result := struct {
		Go struct {
			Version string `yaml:"version"`
		} `yaml:"go"`
	}{}
	err = yaml.Unmarshal(yamlData, &result)
	if err != nil {
		t.Error("YamlMarshal结果解析失败")
	}

	expected := "1.16"
	if expected != result.Go.Version {
		t.Errorf("YamlMarshal结果不符合预期。期望：%s，实际：%s", expected, yamlData)
	}
}

func TestParser_GetRequireModList(t *testing.T) {
	// 创建一个Parser对象
	parser, err := modparser.New(TestModFile)
	if err != nil {
		t.Fatalf("无法创建Parser对象：%s", err.Error())
	}

	// 调用GetRequireModList方法
	modList := parser.GetRequireModList()

	// 检查结果是否符合预期
	expected := []string{"github.com/PuerkitoBio/goquery"}
	if !equalStringSlices(modList, expected) {
		t.Errorf("GetRequireModList结果不符合预期。期望：%v，实际：%v", expected, modList)
	}
}

func TestParser_GetRequireIndirectModList(t *testing.T) {
	// 创建一个Parser对象
	parser, err := modparser.New(TestModFile)
	if err != nil {
		t.Fatalf("无法创建Parser对象：%s", err.Error())
	}

	// 调用GetRequireIndirectModList方法
	modList := parser.GetRequireIndirectModList()

	// 检查结果是否符合预期
	expected := []string{"github.com/andybalholm/cascadia"}
	if !equalStringSlices(modList, expected) {
		t.Errorf("GetRequireIndirectModList结果不符合预期。期望：%v，实际：%v", expected, modList)
	}
}

// helper function to check if two string slices are equal
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
