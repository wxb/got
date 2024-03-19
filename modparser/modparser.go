package modparser

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"

	"golang.org/x/mod/modfile"
)

type Parser struct {
	// go.mod 文件路径
	modPath string
	// go.mod 文件内容
	modContent []byte
	// go.mod 文件解析结果
	modFile *modfile.File
}

// New 实例化一个mod文件解析器
func New(modPath string) (*Parser, error) {
	p := Parser{
		modPath: modPath,
	}

	// 读取go.mod文件内容
	content, err := os.ReadFile(modPath)
	if err != nil {
		return &p, fmt.Errorf("无法读取go.mod文件: %s", err.Error())
	}
	p.modContent = content

	// 解析go.mod文件内容
	modFile, err := modfile.Parse(p.modPath, p.modContent, nil)
	if err != nil {
		return &p, fmt.Errorf("无法解析go.mod文件内容: %s", err.Error())
	}
	p.modFile = modFile

	return &p, nil
}

// Parse go.mod文件解析成结构话结果modfile
func (p *Parser) Parse() *modfile.File {
	return p.modFile
}

// JsonMarshal adapts to json/encoding Marshal API
func (p *Parser) JsonMarshal() ([]byte, error) {
	return jsoniter.Marshal(p.modFile)
}

// JsonMarshalIndent same as json.MarshalIndent. Prefix is not supported.
func (p *Parser) JsonMarshalIndent(prefix, indent string) ([]byte, error) {
	return jsoniter.MarshalIndent(p.modFile, prefix, indent)
}

// JsonMarshalToString convenient method to write as string instead of []byte
func (p *Parser) JsonMarshalToString() (string, error) {
	return jsoniter.MarshalToString(p.modFile)
}

// YamlMarshal 转换为YAML格式
func (p *Parser) YamlMarshal() ([]byte, error) {
	return yaml.Marshal(p.modFile)
}

// GetModPath 包名
func (p *Parser) GetModPath() string {
	return p.modFile.Module.Mod.Path
}

// GetGoVersion go版本
func (p *Parser) GetGoVersion() string {
	return p.modFile.Go.Version
}

// GetRequireModList 直接引用的包列表
func (p *Parser) GetRequireModList() []string {
	res := []string{}
	for _, v := range p.modFile.Require {
		if v.Indirect {
			continue
		}

		res = append(res, v.Mod.Path)
	}

	return res
}

// GetRequireIndirectModList 间接引入的包列表
func (p *Parser) GetRequireIndirectModList() []string {
	res := []string{}
	for _, v := range p.modFile.Require {
		if !v.Indirect {
			continue
		}

		res = append(res, v.Mod.Path)
	}

	return res
}
