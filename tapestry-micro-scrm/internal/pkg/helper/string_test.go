package helper

import "testing"

func TestCleanUserInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes control characters from source friend content",
			input: "Coach新款Jade菱格小羊皮\v 搭配娃娃裙，清甜显贵\v 夏日百搭首选👗💐🌸",
			want:  "Coach新款Jade菱格小羊皮 搭配娃娃裙，清甜显贵 夏日百搭首选👗💐🌸",
		},
		{
			name:  "removes invalid UTF-8 bytes",
			input: string([]byte{'a', 0xe0, 0xb2, 'b'}),
			want:  "ab",
		},
		{
			name:  "preserves allowed whitespace emoji and replacement rune",
			input: "a\t\n\rb�😀",
			want:  "a\t\n\rb�😀",
		},
		{
			name:  "removes ASCII and Unicode control characters",
			input: "a\x00\x08\x0b\x0c\x1f\x7f\u0085b",
			want:  "ab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanUserInput(tt.input); got != tt.want {
				t.Fatalf("CleanUserInput() = %q, want %q", got, tt.want)
			}
		})
	}
}
