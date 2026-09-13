import re

with open("ui/artwork.go", "r") as f:
    content = f.read()

fallback = """			if strings.Contains(url, "-large.jpg") {
				url2 := strings.Replace(url, "-large.jpg", "-t500x500.jpg", 1)
				resp2, err2 := client.Get(url2)
				if err2 == nil && resp2.StatusCode == 200 {
					resp = resp2
				} else {
					if err2 == nil { resp2.Body.Close() }
					return artworkMsg{trackID: trackID, art: ""}
				}
			} else if strings.Contains(url, "500x500.jpg") {
				url2 := strings.Replace(url, "500x500.jpg", "150x150.jpg", 1)
				resp2, err2 := client.Get(url2)
				if err2 == nil && resp2.StatusCode == 200 {
					resp = resp2
				} else {
					if err2 == nil { resp2.Body.Close() }
					return artworkMsg{trackID: trackID, art: ""}
				}
			} else {
				return artworkMsg{trackID: trackID, art: ""}
			}"""

content = re.sub(r'			if strings\.Contains\(url, "-large\.jpg"\) \{.*?} else \{\n\t\t\t\treturn artworkMsg\{trackID: trackID, art: ""\}\n\t\t\t\}', fallback, content, flags=re.DOTALL)

with open("ui/artwork.go", "w") as f:
    f.write(content)
