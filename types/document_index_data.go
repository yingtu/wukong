package types

import (
	"sort"
	"strings"
)

type DocumentIndexData struct {
	// 文档全文（必须是UTF-8格式），用于生成待索引的关键词
	Content string

	// 文档的关键词
	// 当Content不为空的时候，优先从Content中分词得到关键词。
	// Tokens存在的意义在于绕过悟空内置的分词器，在引擎外部
	// 进行分词和预处理。
	Tokens TokenDataList

	// 文档标签（必须是UTF-8格式），比如文档的类别属性等，这些标签并不出现在文档文本中
	Labels []string

	// 文档的评分字段，可以接纳任何类型的结构体
	Fields interface{}
}

type TokenDataList []TokenData

// 文档的一个关键词
type TokenData struct {
	// 关键词的字符串
	Text string

	// 关键词的首字节在文档中出现的位置
	Locations []int
}

// 方便通过 TokenData 来计算文档哈希值
func (tokens TokenDataList) GetTokensContent() string {
	var content string
	for _, token := range tokens {
		content += token.Text
	}
	return content
}
func (tokens TokenDataList) Len() int {
	return len(tokens)
}
func (tokens TokenDataList) Swap(i, j int) {
	tokens[i], tokens[j] = tokens[j], tokens[i]
}
func (tokens TokenDataList) Less(i, j int) bool {
	return tokens[i].Text < tokens[j].Text
}

func (documentIndexData *DocumentIndexData) GetContent() string {
	if documentIndexData.Content != "" {
		return documentIndexData.Content
	}

	if len(documentIndexData.Tokens) != 0 {
		sort.Sort(documentIndexData.Tokens)
		return documentIndexData.Tokens.GetTokensContent()
	}

	sort.Strings(documentIndexData.Labels)
	return strings.Join(documentIndexData.Labels, "")
}
