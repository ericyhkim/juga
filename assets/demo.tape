# Output Settings
Output assets/demo.gif

# Visuals
Set Theme "Catppuccin Mocha"
Set WindowBar "Colorful"
Set WindowBarSize 40
Set BorderRadius 10
Set Padding 20
Set Margin 40
Set MarginFill "#1e1e2e"

# Interaction
Set TypingSpeed 75ms

# Setup (Hidden)
Hide
Type "export DEMO_ROOT=$(mktemp -d)" Enter
Type "mkdir -p $DEMO_ROOT/bin $DEMO_ROOT/config $DEMO_ROOT/data $DEMO_ROOT/cache" Enter
Type "export JUGA_CONFIG_HOME=$DEMO_ROOT/config" Enter
Type "export JUGA_DATA_HOME=$DEMO_ROOT/data" Enter
Type "export JUGA_CACHE_HOME=$DEMO_ROOT/cache" Enter
Type "cp pkg/storage/master_tickers.csv $DEMO_ROOT/data/" Enter
Type "go build -o $DEMO_ROOT/bin/juga ./cmd/juga" Enter
Type "export PATH=$DEMO_ROOT/bin:$PATH" Enter
Type "export PS1='> '" Enter
Ctrl+L
Show

# Scene 1: Direct Lookup
Type "# Direct lookup by stock code" Enter
Type "juga '#005930'"
Sleep 500ms
Enter
Sleep 3s

Hide
Ctrl+L
Show

# Scene 2: Interactive Picker
Type "# Interactive selection for ambiguous results" Enter
Type "juga /삼성전"
Sleep 500ms
Enter
Sleep 1s
Right
Sleep 750ms
Right
Sleep 750ms
Left
Sleep 500ms
Left
Sleep 500ms
Enter
Sleep 3s

Hide
Ctrl+L
Show

# Scene 3: Portfolios & Aliases
Type "# Setup a custom portfolio" Enter
Type "juga alias set sam 005930"
Enter
Sleep 1s
Type "juga portfolio set tech sam 000660"
Enter
Sleep 2s

Hide
Ctrl+L
Show

# Scene 4: Portfolio Resolution
Type "# Accessing portfolios using the @ prefix" Enter
Type "juga @tech"
Sleep 500ms
Enter
Sleep 3s

Hide
Ctrl+L
Show

# Scene 5: Market Overview
Type "# Real-time market status" Enter
Type "juga market"
Sleep 500ms
Enter
Sleep 5s

# Cleanup (Hidden)
Hide
Type "rm -rf $DEMO_ROOT" Enter
Show