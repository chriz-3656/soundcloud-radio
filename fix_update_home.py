import re

with open("ui/update.go", "r") as f:
    content = f.read()

home_block = """		case ViewHome:
			switch msg.String() {
			case "j", "down":
				if m.homeCursor < len(m.homeFeed)-1 {
					m.homeCursor++
				}
			case "k", "up":
				if m.homeCursor > 0 {
					m.homeCursor--
				}
			case "enter":
				if len(m.homeFeed) > 0 {
					selected := m.homeFeed[m.homeCursor]
					m.queue.PlayNext(selected)
					if track, ok := m.queue.Pop(); ok {
						m.statusMsg = "Resolving " + track.Title + "..."
						cmds = append(cmds, m.resolveCmd(track))
					}
				}
			case "a":
				if len(m.homeFeed) > 0 {
					m.queue.Add(m.homeFeed[m.homeCursor])
					m.showNotif("✓ Added to queue")
				}
			}
			
"""

content = content.replace("		case ViewQueue:\n\t\t\tswitch msg.String() {\n\t\t\tcase \"j\", \"down\":", home_block + "		case ViewQueue:\n\t\t\tswitch msg.String() {\n\t\t\tcase \"j\", \"down\":")

with open("ui/update.go", "w") as f:
    f.write(content)
