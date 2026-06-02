Output scripts/demo.gif

Set FontSize 14
Set Width 900
Set Height 550
Set Theme "Dracula"

# Launch certui with a sample cert file
Type "certui"
Enter
Sleep 1s

# Show the main menu - pause so it's readable
Sleep 1s

# Navigate down to "Inspect" (first item under Validate)
# Inspect is the first menu item, already highlighted
Sleep 500ms

# Select Inspect
Enter
Sleep 800ms

# The file picker appears — type the filename
Type "test.pem"
Sleep 300ms
Enter
Sleep 1s

# Certificate details are shown
# Scroll down to reveal more fields
Down
Sleep 300ms
Down
Sleep 300ms

# Toggle full view to show extra fields (OCSP, CRL, policies...)
Type "f"
Sleep 800ms

# Scroll through full details
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms

# Copy to clipboard
Type "y"
Sleep 500ms

# Go back to the main menu
Escape
Sleep 400ms
Escape
Sleep 400ms

# Navigate down to "Batch inspect"
Down
Sleep 300ms

# Navigate down to "Convert" section — keep going
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms

# Back up to top and quit
Type "q"
Sleep 500ms
