# Global properties
title-text: ""
desktop-image: "background.jpg"
desktop-color: "#000000"
terminal-font: "Unifont:style=Medium;1"
terminal-box: "terminal_box_*.png"
terminal-left: "0"
terminal-top: "0"
terminal-width: "100%"
terminal-height: "100%"
terminal-border: "0"

# Boot menu
+ boot_menu {
  left = "(screen_width - width) / 2"
  top = "(screen_height - height) / 2"
  width = "1.7 * height"
  height = "6*item_spacing + 8*item_height + 2*item_r + 3"
  item_font = "Noto Sans CJK SC Regular;1"
  item_color = "#dddddd"
  selected_item_color = "#ffffff"
  item_height = "font_height * 1.574"
  item_spacing = "font_height * 0.328"
  item_padding = "font_height * 0.328"
  icon_width = "font_height * 1.115"
  icon_height = "font_height * 0.787"
  item_icon_space = "font_height * 0.656"
  item_pixmap_style = "item_*.png"
  selected_item_pixmap_style = "selected_item_*.png"
  menu_pixmap_style = "menu_*.png"
  scrollbar_thumb = "scrollbar_thumb_*.png"
}

# Countdown message
+ label {
  _id = "label1"
  left = 0
  top = "screen_height - 1 * font_height"
  width = 100%
  align = "center"
  id = "__timeout__"
  # DE
  # text = "Start in %d Sekunden."
  # EN
  _text_en = "Booting in %d seconds"
  # ES
  # text = "Arranque en% d segundos"
  # FR
  # text = "Démarrage automatique dans %d secondes"
  # NO
  # text = "Starter om %d sekunder"
  # PT
  # text = "Arranque automático dentro de %d segundos"
  # RU
  # text = "Загрузка выбранного пункта через %d сек."
  # UA
  # text = "Автоматичне завантаження розпочнеться через %d сек."
  # zh_CN
  _text_zh_CN = "在 %d 秒内启动"
  # zh_TW
  _text_zh_TW = "在 %d 秒內啟動"
  # zh_HK
  _text_zh_HK = "在 %d 秒內啟動"
  color = "#99E53E"
  font = "Noto Sans CJK SC Regular;0.85"
}

# Navigation keys hint 
+ label {
  _id = "label2"
  left = 0
  top = "screen_height - 2 * font_height"
  width = 100%
  align = "center"
  # DE
  # text = "System mit ↑ und ↓ auswählen und mit Enter bestätigen."
  # EN
  _text_en = "Use ↑ and ↓ keys to change selection, Enter to confirm, E to edit the commands before booting or C for a command-line"
  # ES
  # text = "Use las teclas ↑ y ↓ para cambiar la selección, Enter para confirmar"
  # FR
  # text = "Choisissez le système avec les flèches du clavier (↑ et ↓), puis validez avec la touche Enter (↲)"
  # NO
  # text = "Bruk ↑ og ↓ for å endre menyvalg, velg med Enter"
  # PT
  # text = "Use as teclas ↑ e ↓ para mudar a seleção, e ENTER para confirmar"
  # RU
  # text = "Используйте клавиши ↑ и ↓ для изменения выбора, Enter для подтверждения"
  # UA
  # text = "Використовуйте ↑ та ↓ для вибору, Enter для підтвердження"
  # zh_CN
  _text_zh_CN = "使用 ↑ 和 ↓ 键移动选择条，Enter 键确认，E 键编辑启动命令，C 键进入命令行"
  # zh_TW
  _text_zh_TW = "使用 ↑ 和 ↓ 鍵移動選擇條，Enter 鍵確認，E 鍵編輯啟動命令，C 鍵進入命令行"
  # zh_HK
  _text_zh_HK = "使用 ↑ 和 ↓ 鍵移動選擇條，Enter 鍵確認，E 鍵編輯啟動命令，C 鍵進入命令行"
  color = "#99E53E"
  font = "Noto Sans CJK SC Regular;0.85"
}

