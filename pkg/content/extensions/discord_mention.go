package extensions

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
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
		nil,
	)
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
			inMention = false
			break
		} else if inMention {
			userId += string(char)
		}
	}

	if inMention || userId == "" {
		return nil
	}

	block.Advance(lineHead + 1)

	return &DiscordMentionNode{
		ID: userId,
	}
}

var _ parser.InlineParser = (*discordMentionParser)(nil)

type discordMentionRenderer struct {
	resolver *DiscordUserResolver
}

func (r *discordMentionRenderer) render(
	w util.BufWriter, source []byte, n ast.Node, entering bool,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	username := fmt.Sprintf("<@%s>", n.(*DiscordMentionNode).ID)

	if r.resolver != nil {
		username = r.resolver.Resolve(n.(*DiscordMentionNode).ID)
	}

	_, err := w.WriteString(username)
	if err != nil {
		return ast.WalkStop, fmt.Errorf("failed to write string: %w", err)
	}

	return ast.WalkContinue, nil
}

func (r *discordMentionRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(DiscordMentionKind, r.render)
}

var _ renderer.NodeRenderer = (*discordMentionRenderer)(nil)

type DiscordMention struct {
	resolver *DiscordUserResolver
}

func (d *DiscordMention) Extend(m goldmark.Markdown) {
	m.Renderer().
		AddOptions(renderer.WithNodeRenderers(util.Prioritized(&discordMentionRenderer{resolver: d.resolver}, 500)))
	m.Parser().
		AddOptions(parser.WithInlineParsers(util.Prioritized(&discordMentionParser{}, 500)))
}

type CachedMention struct {
	Username  string
	CacheTime time.Time
}

type DiscordUserResolver struct {
	config        DiscordUserResolverConfig
	cache         map[string]CachedMention
	discordClient *discordgo.Session
}

func (r *DiscordUserResolver) Resolve(userId string) string {
	if cachedVal, ok := r.cache[userId]; ok && time.Since(cachedVal.CacheTime) < time.Hour*24*7 {
		return cachedVal.Username
	}

	slog.Info("Resolving user...", slog.String("user_id", userId))

	var username string
	user, err := r.discordClient.User(userId)
	if err != nil {
		username = fmt.Sprintf("<@%s>", userId)
		slog.Error("Failed to resolve user", "error", err)
	} else {
		username = user.Username
	}

	r.cache[userId] = CachedMention{
		Username:  username,
		CacheTime: time.Now(),
	}

	if r.config.CachePath != "" {
		f, err := os.Create(r.config.CachePath)
		if err != nil {
			slog.Error("Failed to open cache file", "error", err)
		} else {
			defer f.Close()

			err = json.NewEncoder(f).Encode(r.cache)
			if err != nil {
				slog.Error("Failed to encode cache file", "error", err)
			}
		}
	}

	return username
}

type DiscordUserResolverConfig struct {
	CachePath    string
	DiscordToken string
}

func NewDiscordUserResolver(config DiscordUserResolverConfig) (*DiscordUserResolver, error) {
	var cache map[string]CachedMention
	if config.CachePath == "" {
		slog.Warn("No cache path provided, Discord mention details will not be persisted across Almanac executions")
		cache = make(map[string]CachedMention)
	} else if _, err := os.Stat(config.CachePath); os.IsNotExist(err) {
		cache = make(map[string]CachedMention)
	} else {
		f, err := os.Open(config.CachePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open cache file: %w", err)
		}
		defer f.Close()

		err = json.NewDecoder(f).Decode(&cache)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cache file: %w", err)
		}
	}

	if config.DiscordToken == "" {
		return nil, fmt.Errorf("bot token cannot be empty")
	}

	discordClient, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord client: %w", err)
	}

	return &DiscordUserResolver{
		cache:         cache,
		discordClient: discordClient,
		config:        config,
	}, nil
}

func NewDiscordMention(resolver *DiscordUserResolver) goldmark.Extender {
	return &DiscordMention{resolver: resolver}
}
