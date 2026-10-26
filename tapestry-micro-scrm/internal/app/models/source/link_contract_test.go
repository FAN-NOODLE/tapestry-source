package source

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLinkModelAndMigrationSupportLongKOSURLs(t *testing.T) {
	field, ok := reflect.TypeOf(LinkModel{}).FieldByName("CoverURL")
	if !ok || !strings.Contains(field.Tag.Get("gorm"), "varchar(1000)") {
		t.Fatalf("CoverURL gorm tag = %q，期望 varchar(1000)", field.Tag.Get("gorm"))
	}
	for _, migration := range []struct {
		file string
		size string
	}{
		{file: "../../../../db/migrations/000036_3.6.5_coach_kos_source_link_url.up.sql", size: "varchar(1000)"},
		{file: "../../../../db/migrations/000036_3.6.5_coach_kos_source_link_url.down.sql", size: "varchar(255)"},
	} {
		content, err := os.ReadFile(migration.file)
		if err != nil {
			t.Fatalf("读取迁移 %s 失败：%v", migration.file, err)
		}
		if strings.Count(string(content), migration.size) != 2 || !strings.Contains(string(content), "`link_url`") || !strings.Contains(string(content), "`cover_url`") {
			t.Fatalf("迁移 %s 未把 link_url 和 cover_url 同时设置为 %s", migration.file, migration.size)
		}
	}

	downContent, err := os.ReadFile("../../../../db/migrations/000036_3.6.5_coach_kos_source_link_url.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	downSQL := string(downContent)
	if updateIndex, alterIndex := strings.Index(downSQL, "UPDATE `source_link`"), strings.Index(downSQL, "ALTER TABLE `source_link`"); updateIndex < 0 || alterIndex < 0 || updateIndex >= alterIndex {
		t.Fatal("000036 down 必须先截断 URL，再缩短字段")
	}
	if strings.Count(downSQL, "LEFT(") != 2 || !strings.Contains(downSQL, "CHAR_LENGTH(`link_url`)") || !strings.Contains(downSQL, "CHAR_LENGTH(`cover_url`)") {
		t.Fatal("000036 down 必须按字符将 link_url、cover_url 截到 255")
	}
}

func TestLinkModelAndMigrationSupportKOSUtf8mb4(t *testing.T) {
	modelType := reflect.TypeOf(LinkModel{})
	for _, fieldExpectation := range []struct {
		name string
		size string
	}{
		{name: "Title", size: "varchar(255)"},
		{name: "Description", size: "varchar(255)"},
		{name: "Content", size: "type:text"},
		{name: "Author", size: "varchar(255)"},
	} {
		field, ok := modelType.FieldByName(fieldExpectation.name)
		if !ok || !strings.Contains(field.Tag.Get("gorm"), fieldExpectation.size) {
			t.Fatalf("%s gorm tag = %q，期望包含 %s", fieldExpectation.name, field.Tag.Get("gorm"), fieldExpectation.size)
		}
	}

	for _, file := range []string{
		"../../../../db/migrations/000037_3.6.5_coach_kos_source_link_utf8mb4.up.sql",
		"../../../../db/migrations/000037_3.6.5_coach_kos_source_link_utf8mb4.down.sql",
	} {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("读取迁移 %s 失败：%v", file, err)
		}
		sql := string(content)
		if strings.Count(sql, "varchar(255)") != 3 || !strings.Contains(sql, "`content` text") {
			t.Fatalf("迁移 %s 必须保持 title、description、author 为 varchar(255)，content 为 text", file)
		}
		for _, column := range []string{"`title`", "`description`", "`content`", "`author`"} {
			if !strings.Contains(sql, column) {
				t.Fatalf("迁移 %s 缺少字段 %s", file, column)
			}
		}
		if strings.Count(sql, "CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci") != 4 {
			t.Fatalf("迁移 %s 未将 4 个 KOS 文本字段全部设置为 utf8mb4", file)
		}
	}
}
