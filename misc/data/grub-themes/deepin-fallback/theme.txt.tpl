# GRUB2 gfxmenu Linux Deepin theme
# Designed for any resolution

# Global Property
title-text: ""
desktop-image: "background_in_theme.jpg"
desktop-color: "#000000"
terminal-font: "Unifont Regular 14"
terminal-box: "terminal_box_*.png"
terminal-left: "0"
terminal-top: "0"
terminal-width: "100%"
terminal-height: "100%"
terminal-border: "0"

# Show the boot menu
+ boot_menu {
  left = 34%
  top = 51%
  width = 32%
  height = 50%
  item_font = "Unifont Regular 16"
  item_color = "#dddddd"
  selected_item_color = "#ffffff"
  item_height = 18
  item_spacing = 25
  selected_item_pixmap_style = "selected_item_*.png"
}

# Show a countdown message using the label component
+ label {
    left = 0
    top = 97%
    width = 100%
    align = "center"
    id = "__timeout__"
  _text_en = "Booting in %d seconds"
  # zh_CN
  _text_zh_CN = "在 %d 秒内启动"
  # zh_TW
  _text_zh_TW = "在 %d 秒內啟動"
  # zh_HK
  _text_zh_HK = "在 %d 秒內啟動"
    color = "#7d7d7d"
    font = "Unifont Regular 16"
}
+ label {
    left = 0
    top = 94%
    width = 100%
    align = "center"
    color = "#7d7d7d"
    font = "Unifont Regular 16"
    # EN
    _text_en = "Use ↑ and ↓ keys to change selection, Enter to confirm, E to edit the commands before booting or C for a command-line"
    # zh_CN
    _text_zh_CN = "使用 ↑ 和 ↓ 键移动选择条，Enter 键确认，E 键编辑启动命令，C 键进入命令行"
    # zh_TW
    _text_zh_TW = "使用 ↑ 和 ↓ 鍵移動選擇條，Enter 鍵確認，E 鍵編輯啟動命令，C 鍵進入命令行"
    # zh_HK
    _text_zh_HK = "使用 ↑ 和 ↓ 鍵移動選擇條，Enter 鍵確認，E 鍵編輯啟動命令，C 鍵進入命令行"
}
