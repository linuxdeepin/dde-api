# GRUB2 gfxmenu Linux Deepin theme
# Designed for any resolution

# Global Property
title-text: ""
desktop-image: "background.jpg"
desktop-color: "#000000"
terminal-font: "Unifont Regular 16"
terminal-box: "terminal_box_*.png"
terminal-left: "0"
terminal-top: "0"
terminal-width: "100%"
terminal-height: "100%"
terminal-border: "0"

# Show the boot menu
+ boot_menu {
  left = 24%
  top = 27%
  width = 52%
  height = 46%
  item_font = "Unifont Regular 16"
  item_color = "#cccccc"
  selected_item_color = "#0099ff"
  item_height = 68
  item_spacing = 14
  item_padding = 14
  icon_width = 48
  icon_height = 34
  item_icon_space = 28
  selected_item_pixmap_style = "select_*.png"
  scrollbar_thumb = "scrollbar_thumb_*.png"
  scrollbar_width = 12
  menu_pixmap_style = "menu_*.png"
}

# Show a countdown message using the label component
+ label {
  left = 0
  top = 97%
  width = 100%
  align = "center"
  id = "__timeout__"
  _text_en = "Booting in %d seconds"
  text = "在 %d 秒内启动"
  color = "#99E53E"
  font = "Unifont Regular 16"
}

+ label {
    left = 0
    top = 95%
    width = 100%
    align = "center"
    color = "#99E53E"
    font = "Unifont Regular 16"
    # EN
    _text_en = "Use ↑ and ↓ keys to change selection, Enter to confirm, E to edit the commands before booting or C for a command-line"
    # zh_CN
    text = "使用 ↑ 和 ↓ 键移动选择条，Enter 键确认，E 键编辑启动命令，C 键进入命令行"
}
