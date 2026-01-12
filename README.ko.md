# juga (주가) 📈
[🇺🇸 English](README.md) | **🇰🇷 한국어**

> 한국 실시간 주가 확인을 위한 간단한 CLI 도구

복잡한 설정 없이 별칭(Alias)과 퍼지 검색(Fuzzy Search)을 통해 KOSPI/KOSDAQ 시장 데이터를 터미널에서 즉시 확인할 수 있습니다.

## ⚡️ 목적
- **Simple:** 복잡한 설정이나 API 키 없이, 실행 파일 하나로 즉시 사용 가능.
- **Smart Search:** `005930`을 외울 필요 없이 `juga 삼전`으로 검색.
- **Clean Output:** 가격과 변동폭 등 꼭 필요한 정보만 표기.
- **Mix & Match:** 종목명, 6자리 코드, 별칭을 원하는 대로 섞어서 여러 종목을 한 번에 조회.

## 📥 설치 방법

### 방법 1: Go Install (권장)
[Go 1.25](https://go.dev/dl/) 버전으로 개발되었으나, 1.21 이상 버전에서도 작동합니다.

**macOS / Linux:**
```bash
go install github.com/ericyhkim/juga/cmd/juga@latest

# Go bin 경로를 PATH에 추가 (영구 적용을 위해 ~/.zshrc 또는 ~/.bashrc에 추가)
export PATH=$PATH:$(go env GOPATH)/bin
```

**Windows (PowerShell):**
```powershell
go install github.com/ericyhkim/juga/cmd/juga@latest

# 명령어를 찾을 수 없는 경우 PATH 추가:
$env:Path += ";$(go env GOPATH)\bin"
```

### 방법 2: 직접 빌드
소스 코드를 수정하거나 특정 브랜치에서 빌드하고 싶은 경우:

**macOS / Linux:**
```bash
git clone https://github.com/ericyhkim/juga.git
cd juga

# 로컬 소스를 빌드하고 Go bin 폴더에 자동으로 설치합니다.
go install ./cmd/juga
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/ericyhkim/juga.git
cd juga

# 로컬 소스를 빌드하고 Go bin 폴더에 자동으로 설치합니다.
go install ./cmd/juga
```

### 🗑️ 삭제 방법
**macOS / Linux:**
```bash
# 1. 바이너리 삭제
rm $(go env GOPATH)/bin/juga || sudo rm /usr/local/bin/juga

# 2. 설정 파일 및 종목 데이터베이스 삭제
rm -rf ~/.juga
```

**Windows:**
GOPATH의 `bin` 폴더에서 `juga.exe`를 삭제 후, 사용자 디렉토리(`%USERPROFILE%`)에서 `.juga` 폴더 삭제.

## 🏎️ 빠른 시작
```bash
# 1. 여러 종목 확인 (종목명, 코드, 별칭을 섞어서 사용할 수 있습니다)
juga sam 005380 SK하이닉스

# 2. 종목 코드가 기억나지 않을 때 검색
juga find 삼전

# 3. 별칭 설정
juga alias set sam 005380

# 4. 설정한 별칭으로 검색
juga sam

# 5. 관심 종목들을 포트폴리오로 묶기
juga portfolio set 내주식 sam NAVER 005380

# 6. 포트폴리오 한 번에 확인
juga 내주식
```

## 💻 명령어
| 명령어 | 약칭 | 설명 |
| :--- | :--- | :--- |
| `juga [names...]` | - | **빠른 확인.** 종목명, 코드, 별칭, 포트폴리오의 실시간 가격과 변동 확인. |
| `juga alias set <nick> <tgt>` | `a set` | 별칭 등록 (별칭 → 종목명/코드 매칭) |
| `juga alias edit` | `a edit`, `a e` | 별칭 목록 편집 |
| `juga alias list` | `a list`, `a ls` | 별칭 목록 표기 |
| `juga alias remove <nick>` | `a remove`, `a rm` | 특정 별칭 삭제 |
| `juga portfolio set <name> [s...]` | `p set` | 종목 모음(포트폴리오) 생성 / 덮어쓰기 |
| `juga portfolio edit <name>` | `p edit`, `p e` | 포트폴리오 종목 수정 |
| `juga portfolio list` | `p list`, `p ls` | 포트폴리오 목록 표기 |
| `juga portfolio remove <name>` | `p remove`, `p rm` | 포트폴리오 삭제 |
| `juga find <query>` | `f`, `search` | 종목명으로 코드 검색 |
| `juga update` | `up` | 최신 종목 리스트 조회 및 갱신 |
| `juga market` | `m` | 상세 시장 지수(KOSPI/KOSDAQ) 정보 확인 |

## 🛠 기술 스택
- **Language:** Go (Golang)
- **CLI Framework:** `spf13/cobra`
- **UI/Styling:** `charmbracelet/lipgloss`
- **Fuzzy Matching:** `sahilm/fuzzy`
- **Data Source:** 네이버 금융 실시간 폴링 API (JSON).

## 📂 .dotfiles
- **`~/.juga/aliases.json`**: 별칭과 6자리 코드를 매핑한 개인 설정 파일.
- **`~/.juga/portfolios.json`**: 사용자 정의 포트폴리오 파일.
- **`~/.juga/master_tickers.csv`**: 약 3,600개의 종목 정보 라이브러리. 최초 실행 시 자동 생성.
> **참고:** Windows의 경우 `~`는 `%USERPROFILE%`을 의미합니다.

- **종목 해석 로직**:
  1. 입력값이 **포트폴리오** 명칭인지 확인 (일치 시 종목 목록으로 확장)
  2. `aliases.json`에서 정확히 일치하는 **별칭**이 있는지 확인.
  3. 유효한 **6자리 종목 코드**인지 확인.
  4. `master_tickers.csv`에서 **퍼지 검색**으로 종목명 매칭.
  5. 데이터 검색 및 수신.

## 🎨 Demo

![Demo](assets/demo.gif)

## ❓ 문제 해결
- **"Could not find stock..."**
  종목 리스트가 오래되었을 수 있습니다. `juga update` 명령어로 갱신해보세요.
- **원하지 않는 종목이 계속 검색되나요?**
  오타나 유사한 종목명으로 인해 잘못된 결과가 캐싱되었을 수 있습니다.
  - **해결 1 (추천):** `juga alias set <이름> <코드>` 명령어로 별칭을 직접 등록하세요.
  - **해결 2 (초기화):** `juga clean` 명령어로 검색 기록과 종목 데이터베이스를 초기화하세요.
