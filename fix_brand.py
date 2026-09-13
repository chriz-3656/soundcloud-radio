import re

with open("ui/view.go", "r") as f:
    content = f.read()

# Replace sidebar branding
content = content.replace('sb.WriteString(brandStyle.Render("♫ SOUNDCLOUD") + "\\n\\n")', 'sb.WriteString(brandStyle.Render("⚡ ECHO RADIO") + "\\n\\n")')

# Replace default ASCII art
default_ascii = """		art = `
      ███████╗ ██████╗██╗  ██╗██████╗ 
      ██╔════╝██╔════╝██║  ██║██╔═══██╗
      █████╗  ██║     ███████║██║   ██║
      ██╔══╝  ██║     ██╔══██║██║   ██║
      ███████╗╚██████╗██║  ██║╚██████╔╝
      ╚══════╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ 
        
           No track is currently
           playing on Echo Radio.`"""

# Find where the old ascii art is in ui/view.go
content = re.sub(r'art = `[^`]+`', default_ascii, content, count=1)

with open("ui/view.go", "w") as f:
    f.write(content)
