package extensions

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var DiscordMentionKind = ast.NewNodeKind("DiscordMention")

type DiscordMentionNode struct {
	ast.BaseInline

	ID string
}

func (m *DiscordMentionNode) Dump(source []byte, level int) {
	ast.DumpHelper(
		m,
		source,
		level,
		map[string]string{
			"ID": m.ID,
		},
		nil)
}

func (m *DiscordMentionNode) Kind() ast.NodeKind {
	return DiscordMentionKind
}

var _ ast.Node = (*DiscordMentionNode)(nil)

type discordMentionParser struct{}

func (d discordMentionParser) Trigger() []byte {
	return []byte("<")
}

func (d discordMentionParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()

	lineHead := 0
	userId := ""
	inMention := false

	for ; lineHead < len(line); lineHead++ {
		char := line[lineHead]

		if char == '@' {
			inMention = true
		} else if inMention && char == '>' {
			break
		} else if inMention {
			userId += string(char)
		}
	}

	block.Advance(lineHead + 1)

	return &DiscordMentionNode{
		ID: userId,
	}
}

var _ parser.InlineParser = (*discordMentionParser)(nil)

type discordMentionRenderer struct{}

func (r *discordMentionRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	//node := n.(*DiscordMentionNode)
	w.WriteString("Pretend this is a Discord mention")

	return ast.WalkContinue, nil
}

func (r *discordMentionRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(DiscordMentionKind, r.render)
}

var _ renderer.NodeRenderer = (*discordMentionRenderer)(nil)

type DiscordMention struct{}

func (d *DiscordMention) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(util.Prioritized(&discordMentionRenderer{}, 500)))
	m.Parser().AddOptions(parser.WithInlineParsers(util.Prioritized(discordMentionParser{}, 500)))
}

func NewDiscordMention() goldmark.Extender {
	return &DiscordMention{}
}
