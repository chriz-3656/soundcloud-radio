import os
import re

def replace_in_file(filepath, replacements):
    with open(filepath, 'r') as f:
        content = f.read()
    
    new_content = content
    for old, new in replacements:
        new_content = new_content.replace(old, new)
        
    if new_content != content:
        # If models.Track is used, ensure we import it
        if "models.Track" in new_content and "soundcloud-radio/internal/models" not in new_content:
            # simple injection after package declaration
            new_content = re.sub(r'(package \w+)', r'\1\n\nimport "soundcloud-radio/internal/models"', new_content, count=1)
        
        with open(filepath, 'w') as f:
            f.write(new_content)
        print(f"Updated {filepath}")

def main():
    base_dir = "/home/demo/repo"
    
    replacements = [
        ("soundcloud.Track", "models.Track"),
        ("trackID int64", "trackID string"),
        ("trackID: int64", "trackID: string"),
        ("ID: int64", "ID: string"),
        ("TrackID int64", "TrackID string"),
        ("int64(track.ID)", "track.ID"), # remove int64 casts if any
    ]
    
    for root, dirs, files in os.walk(base_dir):
        # skip internal/soundcloud because we already fixed it
        if "internal/soundcloud" in root:
            continue
            
        for file in files:
            if file.endswith(".go"):
                filepath = os.path.join(root, file)
                replace_in_file(filepath, replacements)

if __name__ == "__main__":
    main()
