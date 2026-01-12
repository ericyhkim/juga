# Output Settings
Output assets/demo.gif

Set FontFamily "Monaco"
Set Margin 20

# Setup (Hidden)
# We build the binary and set a temporary HOME to avoid affecting user's real data
Hide
Type "mkdir -p /tmp/juga-demo/bin" Enter
Type "export HOME=/tmp/juga-demo" Enter
Type "export PATH=$HOME/bin:$PATH" Enter
Type "go build -o $HOME/bin/juga ./cmd/juga" Enter
Type "clear" Enter
Show

# Scene 1: Basic
Type "juga 삼성전자 000660"
Sleep 500ms
Enter
Sleep 2s

Hide
Type "clear" Enter
Show

# Scene 2: Fuzzy Search
Type "juga find 삼전"
Sleep 500ms
Enter
Sleep 2s

# Scene 3: Aliases & Portfolios
Type "juga alias set sam 005930"
Sleep 500ms
Enter
Sleep 2s

Type "juga portfolio set tech sam 000660"
Sleep 500ms
Enter
Sleep 2s

Hide
Type "clear" Enter
Show

# Scene 4:
Type "juga tech"
Sleep 500ms
Enter
Sleep 2s

# Scene 5: Market
Type "juga market"
Sleep 500ms
Enter
Sleep 3s