package config

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCoachKOSEnvironmentConfig(t *testing.T) {
	tests := []struct {
		file    string
		baseURL string
		token   string
	}{
		{file: "../../app.uat.yaml", baseURL: "https://kosuat.coach.com.cn/api", token: "REDACTED_FOR_SOURCE_SHARING"},
		{file: "../../app.prod.yaml", baseURL: "https://kos.coach.com.cn/api", token: ""},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			content, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatal(err)
			}

			var data struct {
				CoachKOS CoachKOSConfig `yaml:"coach_kos"`
			}
			if err := yaml.Unmarshal(content, &data); err != nil {
				t.Fatal(err)
			}
			if data.CoachKOS.BaseURL != tt.baseURL {
				t.Errorf("base_url = %q, want %q", data.CoachKOS.BaseURL, tt.baseURL)
			}
			if data.CoachKOS.Token != tt.token {
				t.Errorf("token = %q, want %q", data.CoachKOS.Token, tt.token)
			}
		})
	}
}

func TestCoachKOSConfigContract(t *testing.T) {
	configType := reflect.TypeOf(Config{})
	field, ok := configType.FieldByName("CoachKOS")
	if !ok {
		t.Fatal("Config 缺少 CoachKOS 字段")
	}
	if field.Type != reflect.TypeOf(CoachKOSConfig{}) {
		t.Fatalf("CoachKOS 字段类型错误：%v", field.Type)
	}
	if got := field.Tag.Get("yaml"); got != "coach_kos" {
		t.Fatalf("CoachKOS yaml 标签错误：%q", got)
	}

	wantTags := map[string]string{
		"BaseURL":  "base_url",
		"Token":    "token",
		"PageSize": "page_size",
	}
	kosType := reflect.TypeOf(CoachKOSConfig{})
	for name, wantTag := range wantTags {
		field, ok := kosType.FieldByName(name)
		if !ok {
			t.Fatalf("CoachKOSConfig 缺少 %s 字段", name)
		}
		if got := field.Tag.Get("yaml"); got != wantTag {
			t.Fatalf("%s yaml 标签错误：%q", name, got)
		}
	}
}
