#! /bin/bash

# 設定
if [[ -n "$GHIDRA_HOME" ]]; then
    echo "Using Ghidra path from environment variable: $GHIDRA_HOME"
else
    GHIDRA_HOME="/opt/homebrew/opt/ghidra/libexec"
fi

GHIDRA_HEADLESS="${GHIDRA_HOME}/support/analyzeHeadless"
SCRIPT_DIR="$(pwd)/lifter"
TARGET_BIN="$1" # The target binary for binary lifting, passed as an argument to the script
PROJ_NAME="$(basename $TARGET_BIN)"
PROJ_DIR="$(echo "$TARGET_BIN" | sed 's|build|ghidra|g')/"
FILE_NAME="$(basename $TARGET_BIN)"

if [ -z "$TARGET_BIN" ]; then
    echo "Usage: $0 <target_binary>"
    exit 1
fi

# 一時プロジェクトディレクトリ作成
mkdir -p "$PROJ_DIR"

# 実行
# -import: バイナリを取り込む
# -postScript: 解析後にJavaスクリプトを実行
# -deleteProject: 終了後にプロジェクトを破棄
time $GHIDRA_HEADLESS "$PROJ_DIR" "$PROJ_NAME" \
    -import "$TARGET_BIN" \
    -scriptPath "$SCRIPT_DIR" \
    -postScript "$SCRIPT_DIR/HighPCodeLifter.java"
    # -deleteProject \
mv $FILE_NAME.json $PROJ_DIR

# ディレクトリ削除
# rm -rf "$PROJ_DIR"
