package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Microsoft Office family (plus the WPS equivalents, which ship the same
// defaults for the shortcuts covered here).

func (r *Registry) initOffice() {
	r.initWord()
	r.initExcel()
	r.initPowerPoint()
	r.initOutlook()
}

func (r *Registry) initWord() {
	word := newProfile("word", "Microsoft Word", context.CategoryOffice,
		"winword.exe", "wps.exe",
	).source("Word for Windows keyboard shortcuts", "Microsoft 365 Word (Windows)")

	// Character formatting
	word.k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl, "U", "Underline").
		k(hook.ModCtrl, "D", "Font Dialog").
		k(hook.ModCtrl|hook.ModShift, "P", "Font Size").
		k(hook.ModCtrl|hook.ModShift, ">", "Increase Font Size").
		k(hook.ModCtrl|hook.ModShift, "<", "Decrease Font Size").
		k(hook.ModCtrl, "]", "Increase Font Size by 1pt").
		k(hook.ModCtrl, "[", "Decrease Font Size by 1pt").
		k(hook.ModCtrl|hook.ModShift, "A", "All Caps").
		k(hook.ModCtrl|hook.ModShift, "K", "Small Caps").
		k(hook.ModCtrl, "=", "Subscript").
		k(hook.ModCtrl|hook.ModShift, "=", "Superscript").
		k(hook.ModCtrl, "SPACE", "Clear Character Formatting").
		k(hook.ModCtrl|hook.ModShift, "C", "Copy Formatting").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste Formatting")

	// Paragraph formatting
	word.k(hook.ModCtrl, "E", "Center").
		k(hook.ModCtrl, "L", "Align Left").
		k(hook.ModCtrl, "R", "Align Right").
		k(hook.ModCtrl, "J", "Justify").
		k(hook.ModCtrl, "M", "Increase Indent").
		k(hook.ModCtrl|hook.ModShift, "M", "Decrease Indent").
		k(hook.ModCtrl, "1", "Single Line Spacing").
		k(hook.ModCtrl, "2", "Double Line Spacing").
		k(hook.ModCtrl, "5", "Line Spacing 1.5").
		k(hook.ModCtrl, "Q", "Clear Paragraph Formatting").
		k(hook.ModCtrl|hook.ModShift, "N", "Normal Style").
		k(hook.ModCtrl|hook.ModShift, "L", "Bulleted List")

	// Styles
	word.k(hook.ModCtrl|hook.ModAlt, "1", "Heading 1").
		k(hook.ModCtrl|hook.ModAlt, "2", "Heading 2").
		k(hook.ModCtrl|hook.ModAlt, "3", "Heading 3")

	// Insert & review
	word.k(hook.ModCtrl, "K", "Hyperlink").
		k(hook.ModCtrl|hook.ModAlt, "M", "Insert Comment").
		k(hook.ModCtrl|hook.ModAlt, "F", "Insert Footnote").
		k(hook.ModCtrl, "ENTER", "Page Break").
		k(hook.ModCtrl|hook.ModShift, "ENTER", "Column Break").
		k(hook.ModCtrl|hook.ModShift, "8", "Show Formatting Marks").
		k(hook.ModCtrl|hook.ModShift, "E", "Track Changes").
		k(hook.ModNone, "F7", "Spelling").
		k(hook.ModShift, "F7", "Thesaurus")

	// Navigation
	word.k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "H", "Replace").
		k(hook.ModCtrl, "G", "Go to Page").
		k(hook.ModCtrl|hook.ModShift, "G", "Word Count").
		k(hook.ModNone, "F12", "Save As")

	r.addProfile(word)
}

func (r *Registry) initExcel() {
	excel := newProfile("excel", "Microsoft Excel", context.CategoryOffice,
		"excel.exe", "et.exe",
	).source("Excel for Windows keyboard shortcuts", "Microsoft 365 Excel (Windows)")

	// Cell editing
	excel.k(hook.ModNone, "F2", "Edit Cell").
		k(hook.ModCtrl, "1", "Format Cells").
		k(hook.ModCtrl, "ENTER", "Fill Selection").
		k(hook.ModAlt, "ENTER", "New Line in Cell").
		k(hook.ModNone, "F4", "Repeat Last Action").
		k(hook.ModShift, "F2", "Cell Comment").
		k(hook.ModCtrl, "K", "Insert Hyperlink").
		k(hook.ModCtrl, "D", "Fill Down").
		k(hook.ModCtrl, "R", "Fill Right")
	excel.k(hook.ModCtrl|hook.ModShift, "ENTER", "Array Formula")

	// Selection & movement
	excel.k(hook.ModCtrl, "SPACE", "Select Column").
		k(hook.ModShift, "SPACE", "Select Row").
		k(hook.ModCtrl, "HOME", "Go to A1").
		k(hook.ModCtrl, "END", "Last Cell").
		k(hook.ModNone, "F5", "Go To").
		k(hook.ModCtrl, "UP", "Jump to Edge Up").
		k(hook.ModCtrl, "DOWN", "Jump to Edge Down").
		k(hook.ModCtrl, "LEFT", "Jump to Edge Left").
		k(hook.ModCtrl, "RIGHT", "Jump to Edge Right").
		k(hook.ModCtrl|hook.ModShift, "UP", "Extend Selection Up").
		k(hook.ModCtrl|hook.ModShift, "DOWN", "Extend Selection Down").
		k(hook.ModCtrl|hook.ModShift, "LEFT", "Extend Selection Left").
		k(hook.ModCtrl|hook.ModShift, "RIGHT", "Extend Selection Right")

	// Number formats
	excel.k(hook.ModCtrl|hook.ModShift, "$", "Currency Format").
		k(hook.ModCtrl|hook.ModShift, "%", "Percent Format").
		k(hook.ModCtrl|hook.ModShift, "#", "Date Format").
		k(hook.ModCtrl|hook.ModShift, "~", "General Format").
		k(hook.ModCtrl|hook.ModShift, "!", "Number Format").
		k(hook.ModCtrl|hook.ModShift, "@", "Time Format").
		k(hook.ModCtrl|hook.ModShift, "^", "Scientific Format").
		k(hook.ModCtrl|hook.ModShift, "&", "Add Outline Border").
		k(hook.ModCtrl|hook.ModShift, "_", "Remove Borders")

	// Rows, columns, sheets
	excel.k(hook.ModCtrl|hook.ModShift, "+", "Insert Cells").
		k(hook.ModCtrl, "-", "Delete Cells").
		k(hook.ModCtrl, "9", "Hide Row").
		k(hook.ModCtrl|hook.ModShift, "9", "Unhide Row").
		k(hook.ModCtrl, "0", "Hide Column").
		k(hook.ModCtrl|hook.ModShift, "0", "Unhide Column").
		k(hook.ModCtrl, "PAGEUP", "Previous Sheet").
		k(hook.ModCtrl, "PAGEDOWN", "Next Sheet").
		k(hook.ModShift, "F11", "New Worksheet")

	// Data, formulas, charts
	excel.k(hook.ModAlt, "=", "AutoSum").
		k(hook.ModCtrl|hook.ModShift, "L", "Toggle Filter").
		k(hook.ModCtrl, ";", "Insert Date").
		k(hook.ModCtrl|hook.ModShift, ";", "Insert Time").
		k(hook.ModCtrl, "`", "Show Formulas").
		k(hook.ModCtrl, "'", "Copy Formula from Above").
		k(hook.ModAlt, "F1", "Create Chart").
		k(hook.ModNone, "F11", "Create Chart Sheet")

	// Formatting & search
	excel.k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl, "U", "Underline").
		k(hook.ModCtrl, "F", "Find").
		k(hook.ModCtrl, "H", "Replace")

	r.addProfile(excel)
}

func (r *Registry) initPowerPoint() {
	ppt := newProfile("powerpoint", "Microsoft PowerPoint", context.CategoryOffice,
		"powerpnt.exe", "wpp.exe",
	).source("PowerPoint for Windows keyboard shortcuts", "Microsoft 365 PowerPoint (Windows)")

	// Slideshow
	ppt.k(hook.ModNone, "F5", "Start Slideshow").
		k(hook.ModShift, "F5", "From Current Slide").
		k(hook.ModAlt, "F5", "Presenter View").
		k(hook.ModNone, "B", "Black Screen").
		k(hook.ModNone, "W", "White Screen").
		k(hook.ModNone, "ESC", "End Slideshow")

	// Slides
	ppt.k(hook.ModCtrl, "M", "New Slide").
		k(hook.ModCtrl, "D", "Duplicate Slide").
		k(hook.ModCtrl|hook.ModShift, "D", "Duplicate Slide").
		k(hook.ModCtrl, "PAGEUP", "Previous Slide").
		k(hook.ModCtrl, "PAGEDOWN", "Next Slide").
		k(hook.ModCtrl, "HOME", "First Slide").
		k(hook.ModCtrl, "END", "Last Slide").
		k(hook.ModAlt|hook.ModShift, "UP", "Move Slide Up").
		k(hook.ModAlt|hook.ModShift, "DOWN", "Move Slide Down")

	// Text formatting
	ppt.k(hook.ModCtrl, "B", "Bold").
		k(hook.ModCtrl, "I", "Italic").
		k(hook.ModCtrl, "U", "Underline").
		k(hook.ModCtrl, "T", "Font Dialog").
		k(hook.ModCtrl|hook.ModShift, ">", "Increase Font Size").
		k(hook.ModCtrl|hook.ModShift, "<", "Decrease Font Size").
		k(hook.ModCtrl, "E", "Center").
		k(hook.ModCtrl, "L", "Align Left").
		k(hook.ModCtrl, "R", "Align Right").
		k(hook.ModCtrl, "J", "Justify").
		k(hook.ModCtrl|hook.ModShift, "UP", "Move Paragraph Up").
		k(hook.ModCtrl|hook.ModShift, "DOWN", "Move Paragraph Down")

	// Objects
	ppt.k(hook.ModCtrl, "G", "Group Objects").
		k(hook.ModCtrl|hook.ModShift, "G", "Ungroup Objects").
		k(hook.ModCtrl, "K", "Hyperlink").
		k(hook.ModCtrl|hook.ModShift, "C", "Copy Formatting").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste Formatting").
		k(hook.ModCtrl, "ENTER", "Next Placeholder")

	// Files
	ppt.k(hook.ModCtrl, "S", "Save").
		k(hook.ModNone, "F12", "Save As")

	r.addProfile(ppt)
}

func (r *Registry) initOutlook() {
	outlook := newProfile("outlook", "Microsoft Outlook", context.CategoryOffice,
		"outlook.exe",
	).source("Outlook for Windows keyboard shortcuts", "Microsoft 365 Outlook (Windows)")

	// Compose
	outlook.k(hook.ModCtrl, "N", "New Email").
		k(hook.ModCtrl, "R", "Reply").
		k(hook.ModCtrl|hook.ModShift, "R", "Reply All").
		k(hook.ModCtrl, "F", "Forward").
		k(hook.ModCtrl, "ENTER", "Send Message").
		k(hook.ModAlt, "S", "Send Message").
		k(hook.ModNone, "F9", "Send/Receive")

	// Create items
	outlook.k(hook.ModCtrl|hook.ModShift, "M", "New Message").
		k(hook.ModCtrl|hook.ModShift, "A", "New Appointment").
		k(hook.ModCtrl|hook.ModShift, "Q", "New Meeting").
		k(hook.ModCtrl|hook.ModShift, "C", "New Contact").
		k(hook.ModCtrl|hook.ModShift, "K", "New Task").
		k(hook.ModCtrl|hook.ModShift, "E", "New Folder")

	// Module navigation
	outlook.k(hook.ModCtrl, "1", "Go to Mail").
		k(hook.ModCtrl, "2", "Go to Calendar").
		k(hook.ModCtrl, "3", "Go to Contacts").
		k(hook.ModCtrl, "4", "Go to Tasks").
		k(hook.ModCtrl, "5", "Go to Notes").
		k(hook.ModCtrl, "6", "Folder List").
		k(hook.ModCtrl|hook.ModShift, "I", "Go to Inbox").
		k(hook.ModCtrl|hook.ModShift, "B", "Address Book")

	// Message handling
	outlook.k(hook.ModCtrl, "Q", "Mark as Read").
		k(hook.ModCtrl, "U", "Mark as Unread").
		k(hook.ModCtrl|hook.ModShift, "V", "Move to Folder").
		k(hook.ModCtrl, "E", "Search").
		k(hook.ModNone, "F3", "Find")

	r.addProfile(outlook)
}
