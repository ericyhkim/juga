# juga (주가) 📈
**🇺🇸 English** | [🇰🇷 한국어](README.ko.md)

> A minimalist CLI for real-time Korean stock prices.

It’s a simple terminal tool letting you check KOSPI/KOSDAQ market data instantly using aliases and fuzzy search.

## ⚡️ Why use it?
- **Simple:** No API keys, no heavy setup. Just a single binary.
- **Smart Search:** Type `juga 삼전` instead of memorizing `005930`.
- **Clean Output:** Shows only what matters—price and change—without the clutter.
- **Mix & Match:** Fetch multiple stocks at once using any combination of names, codes, or aliases.

## 📥 Installation

### Option 1: Go Install (Recommended)
Developed with [Go 1.25](https://go.dev/dl/), but should work with version 1.21+.

**macOS / Linux:**
```bash
go install github.com/ericyhkim/juga/cmd/juga@latest

# Add Go bin to PATH (add to ~/.zshrc or ~/.bashrc to make it permanent)
export PATH=$PATH:$(go env GOPATH)/bin
```

**Windows (PowerShell):**
```powershell
go install github.com/ericyhkim/juga/cmd/juga@latest

# Add Go bin to PATH if command not found:
$env:Path += ";$(go env GOPATH)\bin"
```

### Option 2: Manual Build
Use this if you want to modify the source or build from a specific branch:

**macOS / Linux:**
```bash
git clone https://github.com/ericyhkim/juga.git
cd juga

# Builds the local source and installs it to your Go bin automatically
go install ./cmd/juga
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/ericyhkim/juga.git
cd juga

# Builds the local source and installs it to your Go bin automatically
go install ./cmd/juga
```

### 🗑️ Uninstallation
**macOS / Linux:**
```bash
# 1. Remove the binary (check both common locations)
rm $(go env GOPATH)/bin/juga || sudo rm /usr/local/bin/juga

# 2. Remove configuration and ticker database
rm -rf ~/.juga
```

**Windows:**
Delete the `juga.exe` file from your GOPATH and the `.juga` folder from your user directory (`%USERPROFILE%`).

## 🏎️ Quick Start
```bash
# 1. Mix names, codes, and aliases freely
juga sam 005380 SK하이닉스

# 2. Find a stock code if you're unsure
juga find 삼전

# 3. Set an alias
juga alias set sam 005380

# 4. Use your alias anytime
juga sam

# 5. Create a portfolio
juga portfolio set my-tech sam NAVER 005380

# 6. Check your portfolio with one command
juga my-tech
```

## 💻 Commands
| Command | Shorthand | Description |
| :--- | :--- | :--- |
| `juga [names...]` | - | **The Quick Peek.** Fetches real-time price & change for stocks, aliases, or portfolios. |
| `juga alias set <nick> <tgt>` | `a set` | Links a nickname to a 6-digit code or name. |
| `juga alias edit` | `a edit`, `a e` | Opens all aliases in your text editor. |
| `juga alias list` | `a list`, `a ls` | Displays all your currently saved shortcuts. |
| `juga alias remove <nick>` | `a remove`, `a rm` | Removes a nickname from your private map. |
| `juga portfolio set <name> [s...]` | `p set` | Creates or overwrites a collection of stocks. |
| `juga portfolio edit <name>` | `p edit`, `p e` | Opens the portfolio in your text editor for bulk changes. |
| `juga portfolio list` | `p list`, `p ls` | Lists all your saved portfolios. |
| `juga portfolio remove <name>` | `p remove`, `p rm` | Removes a portfolio. |
| `juga find <query>` | `f`, `search` | Fuzzy searches the master ticker list to discover new stocks. |
| `juga update` | `up` | Scrapes the data source to keep the master list current. |
| `juga market` | `m` | Show detailed market index information (KOSPI/KOSDAQ). |

## 🛠 Tech Spec
- **Language:** Go (Golang)
- **CLI Framework:** `spf13/cobra`
- **UI/Styling:** `charmbracelet/lipgloss`
- **Fuzzy Matching:** `sahilm/fuzzy`
- **Data Source:** Naver Finance Real-time Polling API (JSON).

## 📂 .dotfiles
- **`~/.juga/aliases.json`**: Your private mapping of nicknames to 6-digit codes.
- **`~/.juga/portfolios.json`**: Your private collections of stock groups.
- **`~/.juga/master_tickers.csv`**: A library of ~3,600 stocks. Automatically initialized from embedded data.
> **Note:** On Windows, `~` refers to `%USERPROFILE%`.

- **Resolver Logic**:
  1. Check if the input is a **Portfolio** (expands to list of items).
  2. Check `aliases.json` for an exact match.
  3. Check if the input is a valid 6-digit stock code.
  4. Fuzzy Search `master_tickers.csv` for a name match.
  5. Fetch data.

## 🎨 Demo

![Demo](assets/demo.gif)

## ❓ Troubleshooting
- **"⚠️ Could not find stock for..."**
  Your local ticker list might be outdated. Run `juga update` to refresh the master list from Naver.
- **Unexpected Search Results?**
  If `juga <name>` keeps showing the wrong stock (e.g. due to a past typo or ambiguity), the app may have cached the result.
  - **Fix 1 (Recommended):** Set an explicit alias: `juga alias set <name> <code>`.
  - **Fix 2 (Reset):** Run `juga clean` to wipe the search history and ticker database.