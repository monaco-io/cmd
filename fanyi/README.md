# Fanyi - AI-Powered CLI Translation Tool

🌐 **Fanyi** (翻译) is a dead simple command-line translation tool powered by large language models. Just type and translate!

## Features

- 🚀 **Ultra Simple**: `fanyi hello world` - that's all!
- 🌍 **Multi-Language**: Auto-translates to Chinese & English by default
- ⚡ **Flexible Config**: Customize API, model, and languages
- 💾 **Smart Cache**: No redundant API calls
- 📝 **Multiple Modes**: CLI, files, pipes, or interactive
- 🔄 **Batch Ready**: Translate files and streams with ease
- 🎯 **Fast & Accurate**: Powered by GPT-4, Claude, etc.

## Installation

```bash
# Build from source
go install

# Or manually
go build -o fanyi
sudo mv fanyi /usr/local/bin/
```

## Quick Start - 3 Steps

### 1. Configure API

```bash
mkdir -p ~/.config/fanyi
cp config.example.yaml ~/.config/fanyi/config.yaml
# Edit config.yaml with your API key
```

Or set environment variables:

```bash
export FANYI_API_KEY="sk-your-api-key-here"
export FANYI_API_MODEL="gpt-4"
```

### 2. Translate Now!

```bash
# Translate to default languages (Chinese & English)
fanyi 你好呀，今天心情咋样

# Translate specific language
fanyi -t zh hello world, how today?

# From file
fanyi -f article.txt -o output.txt

# Pipe
echo "hello world" | fanyi -t ja

# Interactive
fanyi -i
```

### 3. Done! 🎉

That's it. Start translating!

---

## Usage

### Direct Translation

```bash
# Translate text
fanyi "Hello, World"
fanyi 你好世界

# With target language
fanyi -t zh "English text"
fanyi -t en "中文文本"
fanyi -l ja "Text to Japanese"

# Shorthand - just type!
fanyi hello world
fanyi 你好
```

### File Translation

```bash
# Translate file
fanyi -f input.txt -o output.txt

# Default languages (Chinese, English)
fanyi -f README.md

# Specific language
fanyi -f article.txt -t ja -o article_ja.txt
```

### Pipe Input

```bash
# From stdin
echo "hello" | fanyi -t zh

# Chain commands
cat file.txt | fanyi -t es > output.txt

# From URL
curl https://example.com/text | fanyi -t ko
```

### Interactive Mode

```bash
# Interactive translator
fanyi -i

# Interactive with specific language
fanyi -i -t zh
```

---

## Configuration

### Config File (`~/.config/fanyi/config.yaml`)

```yaml
api:
  endpoint: "https://api.openai.com/v1/chat/completions"
  key: "sk-your-api-key"
  model: "gpt-4"
  timeout: 30
  max_tokens: 1000
  temperature: 0.3

languages:
  common:
    - zh    # Chinese
    - en    # English
    - ja    # Japanese
    - ko    # Korean
    - es    # Spanish
  priority:
    - zh    # 1st priority
    - en    # 2nd priority

cache:
  enabled: true
  directory: ".cache/fanyi"
  ttl: 720    # 30 days

advanced:
  debug: false
  log_dir: ".log/fanyi"
```

### Command-line Options

```bash
fanyi [OPTIONS] [TEXT]

OPTIONS:
  -t, --target-lang <LANG>    Target language (zh, en, ja, ko, es, etc.)
  -l, --language <LANG>       Same as --target-lang
  -f, --file <FILE>           Input file path
  -o, --output <FILE>         Output file path
  -i, --interactive           Interactive mode
  --no-cache                  Skip cache
  -h, --help                  Show help
```

### Environment Variables

```bash
FANYI_API_ENDPOINT       # API URL
FANYI_API_KEY           # API key
FANYI_API_MODEL         # Model name
FANYI_API_TIMEOUT       # Timeout (seconds)
FANYI_API_MAX_TOKENS    # Max tokens

FANYI_LANGUAGES         # Languages (zh,en,ja)
FANYI_LANGUAGE_PRIORITY # Priority (zh>en>ja)

FANYI_CACHE_ENABLED     # Enable cache (true/false)
FANYI_CACHE_DIR         # Cache dir

FANYI_DEBUG             # Debug mode
FANYI_LOG_DIR           # Log directory
```

### Configuration Priority

1. **Command-line arguments** (highest)
2. **Environment variables**
3. **Config file** (~/.config/fanyi/config.yaml)
4. **Default values** (lowest)

---

## Examples

### Example 1: Simple Translation

```bash
$ fanyi hello world
Original: hello world
---
Chinese: 你好世界
English: hello world
```

### Example 2: Chinese to English

```bash
$ fanyi -t en 你好朋友，最近怎么样
How are you, my friend? How have you been lately?
```

### Example 3: File Translation

```bash
$ fanyi -f README.md -t ja -o README_ja.md
$ cat README_ja.md
```

### Example 4: Batch Translate

```bash
# Translate multiple files
for f in *.txt; do
  fanyi -f "$f" -t zh -o "${f%.txt}_zh.txt"
done
```

### Example 5: With Pipes

```bash
# Chain translations
cat article.txt | fanyi -t es | head -5

# From web
curl -s https://example.com/api/text | fanyi -t ko
```

### Example 6: Interactive Session

```bash
$ fanyi -i -t zh
Fanyi - Interactive Mode
Enter text (type 'exit' to quit):
> good morning
早上好

> how are you?
你好吗

> exit
```

---

## Supported LLM Providers

Fanyi works with any OpenAI-compatible API:

- **OpenAI**: `https://api.openai.com/v1/chat/completions`
- **Azure OpenAI**: Azure endpoints
- **Ollama** (Local): `http://localhost:11434/api/chat`
- **Other Compatible Services**: Any OpenAI-compatible API

### Switch Providers

```bash
# Use Ollama locally
export FANYI_API_ENDPOINT="http://localhost:11434/api/chat"
export FANYI_API_MODEL="mistral"
fanyi "hello world"

# Use Azure
export FANYI_API_ENDPOINT="https://{resource}.openai.azure.com/..."
export FANYI_API_MODEL="gpt-4"
```

---

## Tips & Tricks

### Performance

```bash
# Enable caching (default)
fanyi "text"  # Fast on 2nd run

# Skip cache
fanyi "text" --no-cache

# Clear cache
fanyi cache clear
```

### Cost Reduction

```bash
# Use cheaper model
export FANYI_API_MODEL="gpt-3.5-turbo"

# Batch translate
fanyi -f large_file.txt > output.txt
```

### Better Translations

```bash
# Lower temperature for consistent results
# (Edit config.yaml: temperature: 0.3)

# Use GPT-4 for better quality
export FANYI_API_MODEL="gpt-4"
```

---

## Troubleshooting

### API Connection Failed

```bash
# Check config
fanyi config show

# Test connectivity
ping api.openai.com

# Debug mode
FANYI_DEBUG=true fanyi "test"
```

### Invalid API Key

```bash
# Update key
fanyi config set api.key "sk-new-key"

# Or via env
export FANYI_API_KEY="sk-new-key"
```

### Cache Issues

```bash
# Clear cache
fanyi cache clear

# Check cache
fanyi cache status

# Disable cache
fanyi "text" --no-cache
```

---

## Configuration Commands

```bash
# Set values
fanyi config set api.key "sk-xxxxx"
fanyi config set api.model "gpt-4"
fanyi config set languages zh,en,ja,ko
fanyi config set language.priority zh>en

# Get values
fanyi config get api.model
fanyi config show

# View cache
fanyi cache status

# Clear cache
fanyi cache clear
```

---

## FAQ

**Q: Can I change default languages?**
```bash
fanyi config set language.priority zh>en>ja
```

**Q: How to use local LLM (Ollama)?**
```bash
export FANYI_API_ENDPOINT="http://localhost:11434/api/chat"
export FANYI_API_MODEL="mistral"
fanyi "hello"
```

**Q: What languages are supported?**
A: 100+ languages depending on your LLM model.

**Q: How to reduce costs?**
- Enable caching (default)
- Use `gpt-3.5-turbo` instead of `gpt-4`
- Batch translate files

**Q: Custom translation prompts?**
Edit `advanced.prompt_template` in config.yaml

**Q: How to see debug info?**
```bash
FANYI_DEBUG=true fanyi "text"
```

---

## Development

### Project Structure

```
fanyi/
├── README.md              # This file
├── config.example.yaml    # Config template
├── main.go               # Entry point
├── src/
│   ├── config/          # Config management
│   ├── translator/      # Translation logic
│   ├── api/             # LLM API client
│   ├── cache/           # Caching system
│   ├── logger/          # Logging
│   └── utils/           # Utilities
└── cmd/                 # CLI commands
    ├── root.go
    ├── translate.go
    ├── config.go
    └── cache.go
```

### Build

```bash
# Install dependencies
go mod download

# Build
go build -o fanyi

# Install
go install

# Run tests
go test ./...
```

### Code Guidelines

- Follow Go conventions
- Use standard library where possible
- Keep dependencies minimal
- Write tests for new features

---

## Contributing

Contributions welcome! Please:
- Follow Go coding standards
- Add tests for new features
- Update README.md
- Keep commits focused

---

## License

See LICENSE file.

---

## Support

Issues? Check troubleshooting above or review the configuration.

**Last Updated**: 2026-01-13
**Status**: Production Ready ✅

