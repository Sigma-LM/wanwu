package util

import "fmt"

const (
	MD_URL = "![%s](%s)"
)

// MarkdownUrl 构造markdown url
func MarkdownUrl(fileUrl string, fileName string) string {
	return fmt.Sprintf(MD_URL, fileName, fileUrl)
}
