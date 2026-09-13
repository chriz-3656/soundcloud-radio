import re

with open("ui/update.go", "r") as f:
    content = f.read()

pattern = re.compile(r'case 3, 4: m\.viewState = ViewHome\n\s+case 5, 6: m\.viewState = ViewSearch\n\s+case 7, 8: m\.viewState = ViewFavorites\n\s+case 9, 10: m\.viewState = ViewHistory\n\s+case 12, 13: m\.viewState = ViewHelp\n\s+case 14, 15: m\.viewState = ViewSettings')

replacement = """case 3, 4: m.viewState = ViewHome
				case 5, 6: m.viewState = ViewQueue
				case 7, 8: m.viewState = ViewSearch
				case 9, 10: m.viewState = ViewFavorites
				case 11, 12: m.viewState = ViewHistory
				case 14, 15: m.viewState = ViewHelp
				case 16, 17: m.viewState = ViewSettings"""

new_content = pattern.sub(replacement, content)

with open("ui/update.go", "w") as f:
    f.write(new_content)
