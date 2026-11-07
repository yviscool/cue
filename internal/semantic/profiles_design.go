package semantic

import (
	"github.com/yviscool/cue/internal/context"
	"github.com/yviscool/cue/internal/hook"
)

// Design tools. Layer and selection semantics differ wildly between Figma,
// Photoshop and Illustrator, so nothing here leaks into a shared category.

func (r *Registry) initDesign() {
	r.initFigma()
	r.initPhotoshop()
	r.initIllustrator()
}

func (r *Registry) initFigma() {
	figma := newProfile("figma", "Figma", context.CategoryDesign,
		"figma.exe",
	).source("Figma keyboard shortcuts (Windows)", "Figma desktop 2024 (Windows)")

	// Search & navigation
	figma.k(hook.ModCtrl, "/", "Search Menu").
		k(hook.ModCtrl, "K", "Quick Actions").
		k(hook.ModCtrl, "R", "Rename Layer").
		k(hook.ModShift, "1", "Zoom to Fit").
		k(hook.ModShift, "2", "Zoom to Selection").
		k(hook.ModShift, "0", "Zoom to 100%")

	// Selection & arrangement
	figma.k(hook.ModCtrl, "G", "Group Selection").
		k(hook.ModCtrl|hook.ModShift, "G", "Ungroup Selection").
		k(hook.ModCtrl, "D", "Duplicate").
		k(hook.ModCtrl|hook.ModAlt, "G", "Frame Selection").
		k(hook.ModCtrl|hook.ModAlt, "C", "Copy Properties").
		k(hook.ModCtrl|hook.ModAlt, "V", "Paste Properties").
		k(hook.ModCtrl, "[", "Send Backward").
		k(hook.ModCtrl, "]", "Bring Forward").
		k(hook.ModCtrl|hook.ModShift, "[", "Send to Back").
		k(hook.ModCtrl|hook.ModShift, "]", "Bring to Front")

	// Boolean operations
	figma.k(hook.ModCtrl|hook.ModAlt, "U", "Union Selection").
		k(hook.ModCtrl|hook.ModAlt, "S", "Subtract Selection").
		k(hook.ModCtrl|hook.ModAlt, "I", "Intersect Selection").
		k(hook.ModCtrl|hook.ModAlt, "X", "Exclude Selection").
		k(hook.ModCtrl|hook.ModAlt, "M", "Use as Mask").
		k(hook.ModShift, "A", "Auto Layout")

	// Components & assets
	figma.k(hook.ModCtrl|hook.ModAlt, "K", "Create Component").
		k(hook.ModCtrl|hook.ModShift, "K", "Place Image")

	// Export & flatten
	figma.k(hook.ModCtrl|hook.ModShift, "E", "Export").
		k(hook.ModCtrl|hook.ModShift, "O", "Outline Stroke").
		k(hook.ModCtrl, "E", "Flatten Selection").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste Over Selection")

	// View
	figma.k(hook.ModNone, "TAB", "Hide UI").
		k(hook.ModShift, "R", "Toggle Rulers").
		k(hook.ModShift, "G", "Toggle Layout Grid").
		k(hook.ModShift, "O", "Toggle Outlines").
		k(hook.ModCtrl|hook.ModShift, "'", "Toggle Pixel Grid")

	r.addProfile(figma)
}

func (r *Registry) initPhotoshop() {
	ps := newProfile("photoshop", "Adobe Photoshop", context.CategoryDesign,
		"photoshop.exe",
	).source("Photoshop default keyboard shortcuts", "Photoshop 2024 (Windows)")

	// Layers
	ps.k(hook.ModCtrl, "J", "Duplicate Layer").
		k(hook.ModCtrl|hook.ModShift, "N", "New Layer").
		k(hook.ModCtrl, "E", "Merge Down").
		k(hook.ModCtrl|hook.ModShift, "E", "Merge Visible").
		k(hook.ModCtrl|hook.ModAlt|hook.ModShift, "E", "Stamp Visible").
		k(hook.ModCtrl, "G", "Group Layers").
		k(hook.ModCtrl|hook.ModShift, "G", "Ungroup Layers").
		k(hook.ModCtrl|hook.ModAlt, "G", "Create Clipping Mask").
		k(hook.ModAlt, "[", "Select Layer Below").
		k(hook.ModAlt, "]", "Select Layer Above").
		k(hook.ModCtrl|hook.ModShift, "[", "Send to Back").
		k(hook.ModCtrl|hook.ModShift, "]", "Bring to Front")

	// Transform & view
	ps.k(hook.ModCtrl, "T", "Free Transform").
		k(hook.ModCtrl, "0", "Fit on Screen").
		k(hook.ModCtrl, "1", "Actual Pixels").
		k(hook.ModCtrl, "[", "Decrease Brush Size").
		k(hook.ModCtrl, "]", "Increase Brush Size")

	// Selection
	ps.k(hook.ModCtrl, "D", "Deselect").
		k(hook.ModCtrl|hook.ModShift, "I", "Inverse Selection").
		k(hook.ModCtrl|hook.ModShift, "C", "Copy Merged").
		k(hook.ModCtrl|hook.ModShift, "V", "Paste Into")

	// History
	ps.k(hook.ModCtrl|hook.ModAlt, "Z", "Step Backward").
		k(hook.ModCtrl|hook.ModShift, "Z", "Step Forward")

	// Adjustments
	ps.k(hook.ModCtrl, "L", "Levels").
		k(hook.ModCtrl, "M", "Curves").
		k(hook.ModCtrl, "U", "Hue/Saturation").
		k(hook.ModCtrl, "B", "Color Balance").
		k(hook.ModCtrl, "I", "Invert").
		k(hook.ModCtrl|hook.ModShift, "U", "Desaturate").
		k(hook.ModCtrl|hook.ModShift, "L", "Auto Tone").
		k(hook.ModCtrl|hook.ModShift, "B", "Auto Color").
		k(hook.ModCtrl|hook.ModShift|hook.ModAlt, "L", "Auto Contrast")

	// Files
	ps.k(hook.ModCtrl|hook.ModAlt, "S", "Save a Copy").
		k(hook.ModCtrl|hook.ModAlt|hook.ModShift, "S", "Save for Web").
		k(hook.ModCtrl|hook.ModAlt, "I", "Image Size").
		k(hook.ModCtrl|hook.ModAlt, "C", "Canvas Size")

	r.addProfile(ps)
}

func (r *Registry) initIllustrator() {
	ai := newProfile("illustrator", "Adobe Illustrator", context.CategoryDesign,
		"illustrator.exe",
	).source("Illustrator default keyboard shortcuts", "Illustrator 2024 (Windows)")

	// Objects
	ai.k(hook.ModCtrl, "G", "Group Selection").
		k(hook.ModCtrl|hook.ModShift, "G", "Ungroup Selection").
		k(hook.ModCtrl, "J", "Join").
		k(hook.ModCtrl, "D", "Transform Again").
		k(hook.ModCtrl, "8", "Compound Path").
		k(hook.ModCtrl|hook.ModAlt, "8", "Release Compound Path").
		k(hook.ModCtrl, "7", "Clipping Mask").
		k(hook.ModCtrl|hook.ModAlt, "7", "Release Clipping Mask").
		k(hook.ModCtrl|hook.ModShift, "O", "Create Outlines")

	// Arrangement
	ai.k(hook.ModCtrl, "[", "Send Backward").
		k(hook.ModCtrl, "]", "Bring Forward").
		k(hook.ModCtrl|hook.ModShift, "[", "Send to Back").
		k(hook.ModCtrl|hook.ModShift, "]", "Bring to Front")

	// Paste
	ai.k(hook.ModCtrl, "F", "Paste in Front").
		k(hook.ModCtrl, "B", "Paste in Back")

	// Locking & visibility
	ai.k(hook.ModCtrl, "2", "Lock Selection").
		k(hook.ModCtrl|hook.ModAlt, "2", "Unlock All").
		k(hook.ModCtrl, "3", "Hide Selection").
		k(hook.ModCtrl|hook.ModAlt, "3", "Show All")

	// View
	ai.k(hook.ModCtrl, "0", "Fit Artboard").
		k(hook.ModCtrl, "1", "Actual Size").
		k(hook.ModCtrl, "Y", "Outline Mode").
		k(hook.ModCtrl, "R", "Show Rulers").
		k(hook.ModCtrl, ";", "Show Guides").
		k(hook.ModCtrl, "U", "Smart Guides").
		k(hook.ModCtrl, "\"", "Show Grid").
		k(hook.ModCtrl|hook.ModShift, "H", "Show Edges")

	r.addProfile(ai)
}
