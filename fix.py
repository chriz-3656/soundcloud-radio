import re

with open("ui/view.go", "r") as f:
    content = f.read()

pattern = re.compile(r'runes := \[\]rune\(line\)\n\s+if len\(runes\) > mainWidth - 4 \{\n\s+line = string\(runes\[:mainWidth-7\]\) \+ "\.\.\."\n\s+\} line = line\[:mainWidth-7\] \+ "\.\.\." \}')

new_content = pattern.sub(r'runes := []rune(line)\n\t\t\t\tif len(runes) > mainWidth - 4 {\n\t\t\t\t\tline = string(runes[:mainWidth-7]) + "..."\n\t\t\t\t}', content)

with open("ui/view.go", "w") as f:
    f.write(new_content)
