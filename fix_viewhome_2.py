import re

with open("ui/view.go", "r") as f:
    content = f.read()

view_home_block = """	case ViewHome:
		mb.WriteString(titleStyle.Render(strings.ToUpper(m.activeProvider.GetName()) + " HOME") + "\\n\\n")
		
		if len(m.homeFeed) == 0 {
			if m.errorMsg != "" {
				mb.WriteString(errorStyle.Render(m.errorMsg) + "\\n")
			} else {
				mb.WriteString(mutedStyle.Render("Loading trending feed...\\n"))
			}
		} else {
			itemsPerPage := (mainHeight - 10) / 2
			if itemsPerPage < 1 { itemsPerPage = 1 }
			startIndex := 0
			if m.homeCursor >= itemsPerPage {
				startIndex = m.homeCursor - itemsPerPage + 1
			}
			endIndex := startIndex + itemsPerPage
			if endIndex > len(m.homeFeed) { endIndex = len(m.homeFeed) }

			for i := startIndex; i < endIndex; i++ {
				t := m.homeFeed[i]
				cursor := "  "
				style := itemStyle
				if i == m.homeCursor {
					cursor = "▶ "
					style = selectedItemStyle
				}
				
				lyricsIcon := ""
				if m.lyricsAvailable[t.ID] {
					lyricsIcon = " 📜"
				}
				
				line := fmt.Sprintf("%s%02d  %s%s", cursor, i+1, t.Title, lyricsIcon)
				runes := []rune(line)
				if len(runes) > mainWidth - 4 {
					line = string(runes[:mainWidth-7]) + "..."
				}
				mb.WriteString(style.Width(mainWidth-2).Render(line) + "\\n")
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("      %s", formatArtistAlbum(t))) + "\\n")
			}
			if len(m.homeFeed) > itemsPerPage {
				mb.WriteString(mutedStyle.Render(fmt.Sprintf("\\n  ... %d trending tracks", len(m.homeFeed))))
			}
		}
		
"""

content = content.replace("	case ViewQueue:\n\t\tmb.WriteString(titleStyle.Render(\"UP NEXT\") + \"\\n\\n\")", view_home_block + "	case ViewQueue:\n\t\tmb.WriteString(titleStyle.Render(\"UP NEXT\") + \"\\n\\n\")")

with open("ui/view.go", "w") as f:
    f.write(content)
