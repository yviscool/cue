# Shortcut Reference

Complete, generated inventory of every shortcut Cue can interpret. This file is produced from the semantic registry itself, so it can never
describe shortcuts the app does not know about.

> Regenerate with `CUE_UPDATE_REFERENCE=1 go test ./internal/semantic -run TestShortcutReference`.
> Do not edit by hand; `TestShortcutReference` fails when this file drifts from the code.

## Coverage

| Layer | Mappings |
| :--- | ---: |
| Universal | 102 |
| Category dictionaries | 120 |
| Application profiles | 892 |
| Two-chord sequences | 24 |
| **Total** | **1138** |
| Distinct action titles | 768 |
| Profiles | 30 |

Priority when resolving a keystroke:
`user override > exact app profile > app category > universal > no label`.

## Universal

Meaning shared by virtually every Windows application. Profiles and categories override these.

| Shortcut | Action | Group |
| :--- | :--- | :--- |
| `Win + U` | Accessibility Settings | System |
| `Win + A` | Action Center | System |
| `Alt + PRINTSCREEN` | Active Window Screenshot | System |
| `Ctrl + Alt + TAB` | App Switcher | Window |
| `Win + K` | Cast | System |
| `Win + V` | Clipboard History | System |
| `Ctrl + W` | Close | File |
| `Ctrl + F4` | Close Document | File |
| `Win + Ctrl + F4` | Close Virtual Desktop | System |
| `Alt + F4` | Close Window | Window |
| `APPS` | Context Menu | Edit |
| `Shift + F10` | Context Menu | Edit |
| `Ctrl + C` | Copy | Edit |
| `Ctrl + INSERT` | Copy | Edit |
| `Ctrl + X` | Cut | Edit |
| `Shift + DELETE` | Cut | Edit |
| `Win + T` | Cycle Taskbar | System |
| `Win + .` | Emoji Picker | System |
| `Win + ;` | Emoji Picker | System |
| `Win + ESC` | Exit Magnifier | System |
| `Win + F` | Feedback Hub | System |
| `Win + E` | File Explorer | System |
| `Ctrl + F` | Find | Search |
| `Win + B` | Focus Notification Area | System |
| `Win + G` | Game Bar | System |
| `Win + Alt + PRINTSCREEN` | Game Bar Screenshot | System |
| `F1` | Help | System |
| `Win + L` | Lock PC | System |
| `Win + Shift + = (+)` | Magnifier Zoom In | System |
| `Win + -` | Magnifier Zoom Out | System |
| `Win + UP` | Maximize Window | Window |
| `Win + HOME` | Minimize All But Active | Window |
| `Win + Shift + RIGHT` | Move to Next Monitor | Window |
| `Win + Shift + LEFT` | Move to Previous Monitor | Window |
| `VOLMUTE` | Mute Audio | Media |
| `Win + Ctrl + ENTER` | Narrator | System |
| `Ctrl + N` | New | File |
| `Win + Ctrl + D` | New Virtual Desktop | System |
| `Ctrl + TAB` | Next Tab | Navigation |
| `MEDIANEXT` | Next Track | Media |
| `Win + Ctrl + RIGHT` | Next Virtual Desktop | System |
| `Win + N` | Notification Center | System |
| `Ctrl + O` | Open | File |
| `Ctrl + V` | Paste | Edit |
| `Shift + INSERT` | Paste | Edit |
| `Win + ,` | Peek Desktop | System |
| `MEDIAPLAYPAUSE` | Play/Pause | Media |
| `Win + X` | Power User Menu | System |
| `Ctrl + Shift + TAB` | Previous Tab | Navigation |
| `MEDIAPREV` | Previous Track | Media |
| `Win + Ctrl + LEFT` | Previous Virtual Desktop | System |
| `Ctrl + P` | Print | File |
| `Win + P` | Project Display | System |
| `Ctrl + Shift + Z` | Redo | Edit |
| `Ctrl + Y` | Redo | Edit |
| `Ctrl + H` | Replace | Search |
| `Win + Ctrl + Shift + B` | Reset Graphics Driver | System |
| `Ctrl + 0` | Reset Zoom | View |
| `Win + DOWN` | Restore / Minimize | Window |
| `Win + R` | Run Dialog | System |
| `Ctrl + S` | Save | File |
| `Ctrl + Shift + S` | Save As | File |
| `Win + PRINTSCREEN` | Save Screenshot | System |
| `PRINTSCREEN` | Screenshot | System |
| `Win + S` | Search | System |
| `Ctrl + A` | Select All | Edit |
| `Win + I` | Settings | System |
| `Win + D` | Show Desktop | System |
| `Win + Z` | Snap Layouts | System |
| `Win + LEFT` | Snap Left | Window |
| `Win + RIGHT` | Snap Right | Window |
| `Win + Shift + S` | Snipping Tool | System |
| `Ctrl + ESC` | Start Menu | System |
| `Win + Alt + R` | Start/Stop Recording | System |
| `MEDIASTOP` | Stop Playback | Media |
| `Alt + Shift + TAB` | Switch App | Window |
| `Alt + TAB` | Switch App | Window |
| `Win + SPACE` | Switch Input Language | System |
| `Win + PAUSE` | System Properties | System |
| `Ctrl + Shift + ESC` | Task Manager | System |
| `Win + TAB` | Task View | System |
| `Win + 1` | Taskbar App | System |
| `Win + 2` | Taskbar App | System |
| `Win + 3` | Taskbar App | System |
| `Win + 4` | Taskbar App | System |
| `Win + 5` | Taskbar App | System |
| `Win + 6` | Taskbar App | System |
| `Win + 7` | Taskbar App | System |
| `Win + 8` | Taskbar App | System |
| `Win + 9` | Taskbar App | System |
| `Ctrl + Z` | Undo | Edit |
| `Win + H` | Voice Typing | System |
| `VOLDOWN` | Volume Down | Media |
| `VOLUP` | Volume Up | Media |
| `Win + W` | Widgets | System |
| `Alt + SPACE` | Window Menu | Window |
| `Ctrl + =` | Zoom In | View |
| `Ctrl + SCROLL UP` | Zoom In | View |
| `Ctrl + Shift + = (+)` | Zoom In | View |
| `Ctrl + -` | Zoom Out | View |
| `Ctrl + SCROLL DOWN` | Zoom Out | View |
| `Ctrl + Shift + - (_)` | Zoom Out | View |

## Category · Editor

Fallback for applications in the editor family that have no exact profile.

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + / (?)` | Block Comment |
| `Ctrl + Shift + W` | Close Window |
| `Ctrl + SPACE` | Code Completion |
| `Ctrl + Shift + P` | Command Palette |
| `Ctrl + BACKSPACE` | Delete Word Left |
| `Ctrl + DELETE` | Delete Word Right |
| `Ctrl + END` | End of File |
| `F3` | Find Next |
| `Shift + F3` | Find Previous |
| `Ctrl + Shift + F` | Find in Files |
| `F12` | Go to Definition |
| `Ctrl + G` | Go to Line |
| `Ctrl + Shift + DOWN` | Move Line Down |
| `Ctrl + Shift + UP` | Move Line Up |
| `Ctrl + Shift + N` | New Window |
| `Ctrl + PAGEDOWN` | Next Editor Tab |
| `Ctrl + PAGEUP` | Previous Editor Tab |
| `Ctrl + Shift + T` | Reopen Closed Editor |
| `Ctrl + HOME` | Start of File |
| `F10` | Step Over |
| `F9` | Toggle Breakpoint |
| `Ctrl + /` | Toggle Comment |

## Category · Browser

Fallback for applications in the browser family that have no exact profile.

| Shortcut | Action |
| :--- | :--- |
| `Alt + LEFT` | Back |
| `Ctrl + Shift + D` | Bookmark All Tabs |
| `Ctrl + Shift + O` | Bookmark Manager |
| `Ctrl + D` | Bookmark Tab |
| `Shift + ESC` | Browser Task Manager |
| `F7` | Caret Browsing |
| `Ctrl + Shift + DELETE` | Clear Browsing Data |
| `Ctrl + F4` | Close Tab |
| `Ctrl + W` | Close Tab |
| `Ctrl + Shift + W` | Close Window |
| `Ctrl + Shift + J` | Console |
| `Ctrl + Shift + I` | Developer Tools |
| `F12` | Developer Tools |
| `Ctrl + J` | Downloads |
| `Ctrl + G` | Find Next |
| `Ctrl + Shift + G` | Find Previous |
| `Alt + D` | Focus Address Bar |
| `Ctrl + L` | Focus Address Bar |
| `F6` | Focus Address Bar |
| `Alt + Shift + B` | Focus Bookmarks Bar |
| `Alt + RIGHT` | Forward |
| `F11` | Full Screen |
| `Ctrl + F5` | Hard Refresh |
| `Ctrl + Shift + R` | Hard Refresh |
| `Ctrl + H` | History |
| `Alt + HOME` | Home Page |
| `Ctrl + Shift + C` | Inspect Element |
| `Ctrl + 9` | Last Tab |
| `Ctrl + Shift + PAGEUP` | Move Tab Left |
| `Ctrl + Shift + PAGEDOWN` | Move Tab Right |
| `Ctrl + Shift + N` | New Incognito Window |
| `Ctrl + T` | New Tab |
| `Ctrl + N` | New Window |
| `Ctrl + PAGEDOWN` | Next Tab |
| `Ctrl + Shift + V` | Paste and Go |
| `Ctrl + PAGEUP` | Previous Tab |
| `Ctrl + P` | Print |
| `Ctrl + Shift + M` | Profile Menu |
| `Ctrl + R` | Refresh |
| `F5` | Refresh |
| `Ctrl + Shift + T` | Reopen Closed Tab |
| `Ctrl + Shift + A` | Search Tabs |
| `Ctrl + E` | Search from Address Bar |
| `Ctrl + K` | Search from Address Bar |
| `Ctrl + 1` | Switch to Tab |
| `Ctrl + 2` | Switch to Tab |
| `Ctrl + 3` | Switch to Tab |
| `Ctrl + 4` | Switch to Tab |
| `Ctrl + 5` | Switch to Tab |
| `Ctrl + 6` | Switch to Tab |
| `Ctrl + 7` | Switch to Tab |
| `Ctrl + 8` | Switch to Tab |
| `Ctrl + Shift + B` | Toggle Bookmarks Bar |
| `Ctrl + U` | View Source |

## Category · Terminal

Fallback for applications in the terminal family that have no exact profile.

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + L` | Clear Terminal |
| `Ctrl + Shift + W` | Close Tab/Pane |
| `Ctrl + Shift + P` | Command Palette |
| `Ctrl + Shift + C` | Copy Selection |
| `Ctrl + BACKSPACE` | Delete Word Left |
| `Ctrl + DELETE` | Delete Word Right |
| `Alt + Shift + D` | Duplicate Pane |
| `Ctrl + D` | Exit / EOF |
| `Ctrl + Shift + F` | Find |
| `Alt + DOWN` | Focus Pane Down |
| `Alt + LEFT` | Focus Pane Left |
| `Alt + RIGHT` | Focus Pane Right |
| `Alt + UP` | Focus Pane Up |
| `F11` | Full Screen |
| `Ctrl + C` | Interrupt Process |
| `Ctrl + Shift + T` | New Tab |
| `Ctrl + Shift + 1 (!)` | New Tab by Profile |
| `Ctrl + Shift + 2 (@)` | New Tab by Profile |
| `Ctrl + Shift + 3 (#)` | New Tab by Profile |
| `Ctrl + Shift + 4 ($)` | New Tab by Profile |
| `Ctrl + Shift + 5 (%)` | New Tab by Profile |
| `Ctrl + Shift + 6 (^)` | New Tab by Profile |
| `Ctrl + Shift + 7 (&)` | New Tab by Profile |
| `Ctrl + Shift + 8 (*)` | New Tab by Profile |
| `Ctrl + Shift + 9 (()` | New Tab by Profile |
| `Ctrl + Shift + N` | New Window |
| `Ctrl + TAB` | Next Tab |
| `Ctrl + Shift + V` | Paste |
| `Ctrl + Shift + TAB` | Previous Tab |
| `Alt + Shift + DOWN` | Resize Pane Down |
| `Alt + Shift + LEFT` | Resize Pane Left |
| `Alt + Shift + RIGHT` | Resize Pane Right |
| `Alt + Shift + UP` | Resize Pane Up |
| `Ctrl + Shift + DOWN` | Scroll Down |
| `Ctrl + Shift + PAGEDOWN` | Scroll Page Down |
| `Shift + PAGEDOWN` | Scroll Page Down |
| `Ctrl + Shift + PAGEUP` | Scroll Page Up |
| `Ctrl + Shift + UP` | Scroll Up |
| `Shift + PAGEUP` | Scroll Up |
| `Ctrl + END` | Scroll to Bottom |
| `Ctrl + HOME` | Scroll to Top |
| `Alt + Shift + - (_)` | Split Pane Down |
| `Alt + Shift + = (+)` | Split Pane Right |
| `Alt + ENTER` | Toggle Fullscreen |

## Profile · Dev-C++

29 mappings · verified against Dev-C++ 6.3 / 5.11 (Windows) · source: Dev-C++ default shortcut list

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + F4` | Close File |
| `Ctrl + SPACE` | Code Completion |
| `Ctrl + F9` | Compile |
| `F9` | Compile |
| `Ctrl + F11` | Compile & Run |
| `F11` | Compile & Run |
| `F5` | Debug |
| `Ctrl + BACKSPACE` | Delete Word Left |
| `Ctrl + DELETE` | Delete Word Right |
| `Ctrl + F` | Find |
| `F3` | Find Next |
| `Shift + F3` | Find Previous |
| `Ctrl + Shift + F` | Find in Files |
| `Ctrl + Shift + A` | Format Code |
| `Ctrl + G` | Go to Line |
| `Ctrl + I` | Indent |
| `Ctrl + Shift + DOWN` | Move Line Down |
| `Ctrl + Shift + UP` | Move Line Up |
| `Ctrl + N` | New File |
| `Ctrl + O` | Open File |
| `Ctrl + R` | Replace |
| `Ctrl + F10` | Run |
| `F10` | Run |
| `Ctrl + S` | Save |
| `Ctrl + Shift + S` | Save All |
| `F6` | Stop Execution |
| `Ctrl + Alt + F2` | Terminate Program |
| `Ctrl + /` | Toggle Comment |
| `Ctrl + Shift + I` | Unindent |

## Profile · Eclipse

37 mappings · verified against Eclipse 2024-09 (Windows) · source: Eclipse Default Key Bindings

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + / (?)` | Block Comment |
| `Ctrl + SPACE` | Code Completion |
| `F11` | Debug |
| `Ctrl + D` | Delete Line |
| `Ctrl + Alt + DOWN` | Duplicate Line |
| `Alt + Shift + L` | Extract Local Variable |
| `Alt + Shift + M` | Extract Method |
| `Ctrl + F` | Find |
| `Ctrl + K` | Find Next |
| `Ctrl + Shift + K` | Find Previous |
| `Ctrl + Shift + G` | Find References |
| `Ctrl + Shift + F` | Format Document |
| `Ctrl + L` | Go to Line |
| `Ctrl + J` | Incremental Find |
| `Alt + Shift + I` | Inline |
| `Ctrl + M` | Maximize Editor |
| `Alt + Shift + V` | Move |
| `Alt + DOWN` | Move Line Down |
| `Alt + UP` | Move Line Up |
| `F3` | Open Declaration |
| `Ctrl + Shift + R` | Open Resource |
| `Ctrl + Shift + T` | Open Type |
| `F4` | Open Type Hierarchy |
| `Ctrl + Shift + O` | Organize Imports |
| `Ctrl + 3` | Quick Access |
| `Ctrl + 1` | Quick Fix |
| `Ctrl + O` | Quick Outline |
| `Ctrl + Shift + \ (|)` | Remove Block Comment |
| `Alt + Shift + R` | Rename Symbol |
| `F8` | Resume |
| `Ctrl + F11` | Run |
| `Ctrl + H` | Search |
| `F5` | Step Into |
| `F7` | Step Out |
| `F6` | Step Over |
| `Ctrl + Shift + B` | Toggle Breakpoint |
| `Ctrl + /` | Toggle Comment |

## Profile · IDLE

26 mappings · verified against Python 3.12 IDLE (Windows) · source: IDLE keys.py keydefs

| Shortcut | Action |
| :--- | :--- |
| `Alt + X` | Check Module |
| `Ctrl + W` | Close |
| `Alt + 3` | Comment Region |
| `Shift + TAB` | Dedent |
| `Ctrl + [` | Dedent Region |
| `Ctrl + F` | Find |
| `Ctrl + G` | Find Next |
| `Ctrl + Shift + G` | Find Previous |
| `Alt + G` | Go to Line |
| `F1` | Help |
| `TAB` | Indent |
| `Ctrl + ]` | Indent Region |
| `Ctrl + N` | New File |
| `Alt + N` | Next Command |
| `Ctrl + O` | Open File |
| `Alt + P` | Previous Command |
| `Ctrl + Q` | Quit |
| `Ctrl + Shift + Z` | Redo |
| `Ctrl + H` | Replace |
| `Shift + F5` | Run Customized |
| `F5` | Run Module |
| `Ctrl + S` | Save |
| `Ctrl + Shift + S` | Save As |
| `Ctrl + A` | Select All |
| `Alt + 4` | Uncomment Region |
| `Ctrl + Z` | Undo |

## Profile · JetBrains IDE

64 mappings · verified against IntelliJ Platform 2024.x (Windows) · source: JetBrains Default Keymap (Windows)

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Alt + I` | Auto-Indent Lines |
| `Ctrl + Shift + / (?)` | Block Comment |
| `Ctrl + Alt + H` | Call Hierarchy |
| `Ctrl + F4` | Close Editor |
| `Ctrl + SPACE` | Code Completion |
| `Ctrl + K` | Commit Changes |
| `Ctrl + Shift + ENTER` | Complete Statement |
| `Shift + F9` | Debug |
| `Ctrl + Y` | Delete Line |
| `Ctrl + D` | Duplicate Line |
| `Ctrl + W` | Extend Selection |
| `Ctrl + Alt + M` | Extract Method |
| `Ctrl + Alt + P` | Extract Parameter |
| `Ctrl + Alt + V` | Extract Variable |
| `Ctrl + F12` | File Structure |
| `Ctrl + F` | Find |
| `Ctrl + Shift + A` | Find Action |
| `Alt + F7` | Find Usages |
| `Ctrl + F7` | Find Usages in File |
| `Ctrl + Shift + F` | Find in Files |
| `Alt + INSERT` | Generate Code |
| `Ctrl + N` | Go to Class |
| `Ctrl + B` | Go to Declaration |
| `Ctrl + Shift + N` | Go to File |
| `Ctrl + Alt + B` | Go to Implementation |
| `Ctrl + U` | Go to Super Method |
| `Ctrl + Alt + Shift + N` | Go to Symbol |
| `Ctrl + Shift + B` | Go to Type Declaration |
| `Shift + ESC` | Hide Active Tool Window |
| `Ctrl + Shift + F12` | Hide All Windows |
| `Ctrl + I` | Implement Methods |
| `Ctrl + Alt + N` | Inline |
| `Ctrl + J` | Insert Live Template |
| `Ctrl + Shift + DOWN` | Move Line Down |
| `Ctrl + Shift + UP` | Move Line Up |
| `F2` | Next Highlighted Error |
| `Ctrl + Alt + O` | Optimize Imports |
| `Ctrl + O` | Override Methods |
| `Ctrl + P` | Parameter Info |
| `Shift + F2` | Previous Highlighted Error |
| `Ctrl + Shift + K` | Push Commits |
| `Ctrl + Q` | Quick Documentation |
| `Ctrl + E` | Recent Files |
| `Ctrl + Shift + E` | Recent Locations |
| `Ctrl + Alt + Shift + T` | Refactor This |
| `Ctrl + Alt + L` | Reformat Code |
| `Shift + F6` | Rename Symbol |
| `Ctrl + R` | Replace |
| `Ctrl + Shift + R` | Replace in Files |
| `Shift + F10` | Run |
| `Shift + F11` | Show Bookmarks |
| `Alt + ENTER` | Show Context Actions |
| `Ctrl + Shift + W` | Shrink Selection |
| `Ctrl + Shift + SPACE` | Smart Completion |
| `F7` | Step Into |
| `Shift + F8` | Step Out |
| `F8` | Step Over |
| `Ctrl + F2` | Stop |
| `Ctrl + Alt + T` | Surround With |
| `F11` | Toggle Bookmark |
| `Ctrl + /` | Toggle Comment |
| `Ctrl + F8` | Toggle Line Breakpoint |
| `Ctrl + H` | Type Hierarchy |
| `Ctrl + T` | Update Project |

## Profile · Notepad++

30 mappings · verified against Notepad++ 8.6 (Windows) · source: Notepad++ Shortcut Mapper defaults

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + W` | Close |
| `Ctrl + Shift + T` | Copy Current Line |
| `Ctrl + L` | Delete Line |
| `Ctrl + BACKSPACE` | Delete Word Left |
| `Ctrl + DELETE` | Delete Word Right |
| `Ctrl + D` | Duplicate Line |
| `Ctrl + F` | Find |
| `F3` | Find Next |
| `Shift + F3` | Find Previous |
| `Ctrl + Shift + F` | Find in Files |
| `Alt + 0` | Fold All |
| `Ctrl + Alt + F` | Fold Current Level |
| `F11` | Full Screen |
| `Ctrl + G` | Go to Line |
| `Ctrl + U` | Lower Case |
| `Ctrl + Shift + DOWN` | Move Line Down |
| `Ctrl + Shift + UP` | Move Line Up |
| `Ctrl + N` | New File |
| `Ctrl + O` | Open File |
| `Ctrl + P` | Print |
| `Ctrl + H` | Replace |
| `Ctrl + S` | Save |
| `Ctrl + Alt + S` | Save All |
| `Ctrl + Shift + S` | Save As |
| `Ctrl + Q` | Toggle Block Comment |
| `Ctrl + Shift + Q` | Toggle Stream Comment |
| `Ctrl + T` | Transpose Line |
| `Alt + Shift + 0 ())` | Unfold All |
| `Ctrl + Alt + Shift + F` | Unfold Current Level |
| `Ctrl + Shift + U` | Upper Case |

## Profile · Sublime Text

42 mappings · verified against Sublime Text 4 (Build 4180, Windows) · source: Default (Windows).sublime-keymap

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + D` | Add Selection to Next Find Match |
| `Ctrl + W` | Close File |
| `Ctrl + Shift + W` | Close Window |
| `Ctrl + Shift + P` | Command Palette |
| `Ctrl + Shift + K` | Delete Line |
| `Ctrl + K, Ctrl + K` | Delete to End of Line |
| `Shift + F11` | Distraction Free Mode |
| `Ctrl + Shift + D` | Duplicate Line |
| `Ctrl + Shift + M` | Expand Selection to Brackets |
| `Ctrl + Shift + J` | Expand Selection to Indentation |
| `Ctrl + L` | Expand Selection to Line |
| `Ctrl + Shift + SPACE` | Expand Selection to Scope |
| `Ctrl + F` | Find |
| `Ctrl + Shift + F` | Find in Files |
| `F11` | Full Screen |
| `Ctrl + P` | Goto Anything |
| `Ctrl + G` | Goto Line |
| `Ctrl + R` | Goto Symbol |
| `Ctrl + Shift + R` | Goto Symbol in Project |
| `Ctrl + ENTER` | Insert Line After |
| `Ctrl + Shift + ENTER` | Insert Line Before |
| `Ctrl + J` | Join Lines |
| `Ctrl + M` | Jump to Matching Bracket |
| `Ctrl + K, Ctrl + L` | Lower Case |
| `Ctrl + Shift + DOWN` | Move Line Down |
| `Ctrl + Shift + UP` | Move Line Up |
| `Ctrl + N` | New File |
| `Ctrl + Shift + N` | New Window |
| `F4` | Next Result |
| `Ctrl + Shift + V` | Paste and Indent |
| `Shift + F4` | Previous Result |
| `Ctrl + Shift + T` | Reopen Closed File |
| `Ctrl + H` | Replace |
| `Ctrl + Shift + S` | Save As |
| `Alt + F3` | Select All Occurrences |
| `Ctrl + `` | Show Console |
| `F9` | Sort Lines |
| `Ctrl + Shift + L` | Split Selection into Lines |
| `Ctrl + Shift + / (?)` | Toggle Block Comment |
| `Ctrl + /` | Toggle Comment |
| `Ctrl + K, Ctrl + B` | Toggle Sidebar |
| `Ctrl + K, Ctrl + U` | Upper Case |

## Profile · VS Code

77 mappings · verified against VS Code 1.9x (Windows) · source: VS Code Default Keybindings (Windows)

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Alt + UP` | Add Cursor Above |
| `Ctrl + Alt + DOWN` | Add Cursor Below |
| `Ctrl + K, Ctrl + C` | Add Line Comment |
| `Ctrl + D` | Add Selection to Next Find Match |
| `Ctrl + K, Ctrl + W` | Close All Editors |
| `Ctrl + W` | Close Editor |
| `Ctrl + Shift + W` | Close Window |
| `Ctrl + K, Ctrl + T` | Color Theme |
| `Ctrl + Shift + P` | Command Palette |
| `F1` | Command Palette |
| `Alt + Shift + DOWN` | Copy Line Down |
| `Alt + Shift + UP` | Copy Line Up |
| `Ctrl + Shift + K` | Delete Line |
| `Ctrl + F` | Find |
| `Shift + F12` | Find References |
| `Ctrl + Shift + F` | Find in Files |
| `Ctrl + K, Ctrl + 0` | Fold All |
| `Ctrl + Shift + [ ({)` | Fold Region |
| `Alt + Shift + F` | Format Document |
| `Ctrl + K, Ctrl + F` | Format Selection |
| `F12` | Go to Definition |
| `Ctrl + G` | Go to Line |
| `Ctrl + Shift + O` | Go to Symbol |
| `Ctrl + ]` | Indent Line |
| `Alt + Shift + I` | Insert Cursor at End of Each Line |
| `Ctrl + Shift + ENTER` | Insert Line Above |
| `Ctrl + ENTER` | Insert Line Below |
| `Ctrl + Shift + \ (|)` | Jump to Matching Bracket |
| `Ctrl + K, Ctrl + S` | Keyboard Shortcuts |
| `Ctrl + K, Ctrl + V` | Markdown Preview to the Side |
| `Alt + DOWN` | Move Line Down |
| `Alt + UP` | Move Line Up |
| `Ctrl + Shift + ` (~)` | New Terminal |
| `Ctrl + Shift + N` | New Window |
| `Ctrl + TAB` | Next Editor |
| `F8` | Next Problem |
| `Ctrl + ,` | Open Settings |
| `Ctrl + [` | Outdent Line |
| `Alt + F12` | Peek Definition |
| `Ctrl + Shift + TAB` | Previous Editor |
| `Shift + F8` | Previous Problem |
| `Ctrl + .` | Quick Fix |
| `Ctrl + P` | Quick Open |
| `Ctrl + K, Ctrl + U` | Remove Line Comment |
| `F2` | Rename Symbol |
| `Ctrl + Shift + T` | Reopen Closed Editor |
| `Ctrl + H` | Replace |
| `Ctrl + Shift + H` | Replace in Files |
| `Ctrl + F5` | Run Without Debugging |
| `Ctrl + Shift + L` | Select All Occurrences |
| `Ctrl + F2` | Select All Occurrences of Current Word |
| `Ctrl + T` | Show All Symbols |
| `Ctrl + Shift + E` | Show Explorer |
| `Ctrl + Shift + X` | Show Extensions |
| `Ctrl + K, Ctrl + I` | Show Hover |
| `Ctrl + Shift + U` | Show Output |
| `Ctrl + Shift + M` | Show Problems |
| `Ctrl + Shift + D` | Show Run and Debug |
| `Ctrl + Shift + G` | Show Source Control |
| `Ctrl + \` | Split Editor |
| `F5` | Start Debugging |
| `F11` | Step Into |
| `Shift + F11` | Step Out |
| `F10` | Step Over |
| `Shift + F5` | Stop Debugging |
| `Ctrl + Shift + A` | Toggle Block Comment |
| `F9` | Toggle Breakpoint |
| `Ctrl + /` | Toggle Line Comment |
| `Ctrl + J` | Toggle Panel |
| `Ctrl + B` | Toggle Sidebar |
| `Ctrl + `` | Toggle Terminal |
| `Ctrl + Shift + SPACE` | Trigger Parameter Hints |
| `Ctrl + SPACE` | Trigger Suggestion |
| `Ctrl + U` | Undo Last Cursor Operation |
| `Ctrl + K, Ctrl + J` | Unfold All |
| `Ctrl + Shift + ] (})` | Unfold Region |
| `Ctrl + K, Ctrl + Z` | Zen Mode |

## Profile · Visual Studio

33 mappings · verified against Visual Studio 2022 (C# profile) · source: Visual Studio Default Keyboard Shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + A` | Add New Item |
| `Ctrl + Alt + P` | Attach to Process |
| `Ctrl + Shift + B` | Build Solution |
| `Ctrl + M, Ctrl + O` | Collapse to Definitions |
| `Ctrl + K, Ctrl + C` | Comment Selection |
| `Ctrl + Alt + E` | Exception Settings |
| `Ctrl + R, Ctrl + M` | Extract Method |
| `Ctrl + F` | Find |
| `Shift + F12` | Find References |
| `Ctrl + Shift + F` | Find in Files |
| `Ctrl + K, Ctrl + D` | Format Document |
| `Ctrl + K, Ctrl + F` | Format Selection |
| `Ctrl + T` | Go to All |
| `Ctrl + F12` | Go to Declaration |
| `F12` | Go to Definition |
| `Ctrl + G` | Go to Line |
| `Ctrl + Shift + N` | New Project |
| `Alt + F12` | Peek Definition |
| `Ctrl + R, Ctrl + R` | Rename Symbol |
| `Ctrl + H` | Replace |
| `Ctrl + Shift + H` | Replace in Files |
| `Ctrl + Shift + F5` | Restart Debugging |
| `Ctrl + Alt + L` | Solution Explorer |
| `F5` | Start Debugging |
| `Ctrl + F5` | Start Without Debugging |
| `F11` | Step Into |
| `Shift + F11` | Step Out |
| `F10` | Step Over |
| `Shift + F5` | Stop Debugging |
| `Ctrl + K, Ctrl + K` | Toggle Bookmark |
| `F9` | Toggle Breakpoint |
| `Ctrl + M, Ctrl + M` | Toggle Outlining |
| `Ctrl + K, Ctrl + U` | Uncomment Selection |

## Profile · Zed

24 mappings · verified against Zed 0.16x (Windows) · source: Zed default.json keymap (Windows)

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + D` | Add Selection to Next Find Match |
| `Ctrl + W` | Close Tab |
| `Ctrl + Shift + W` | Close Window |
| `Ctrl + Shift + P` | Command Palette |
| `Ctrl + Shift + K` | Delete Line |
| `Ctrl + F` | Find |
| `F11` | Full Screen |
| `F12` | Go to Definition |
| `Ctrl + G` | Go to Line |
| `Alt + DOWN` | Move Line Down |
| `Alt + UP` | Move Line Up |
| `Ctrl + Shift + N` | New Window |
| `Ctrl + ,` | Open Settings |
| `Ctrl + Shift + F` | Project Search |
| `Ctrl + P` | Quick Open |
| `Ctrl + Q` | Quit |
| `Ctrl + Shift + T` | Reopen Closed Item |
| `Ctrl + Shift + S` | Save As |
| `Ctrl + Shift + L` | Select All Occurrences |
| `Ctrl + J` | Toggle Bottom Dock |
| `Ctrl + /` | Toggle Comment |
| `Ctrl + B` | Toggle Left Dock |
| `Ctrl + Shift + E` | Toggle Project Panel |
| `Ctrl + `` | Toggle Terminal |

## Profile · Microsoft Excel

54 mappings · verified against Microsoft 365 Excel (Windows) · source: Excel for Windows keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + 7 (&)` | Add Outline Border |
| `Ctrl + Shift + ENTER` | Array Formula |
| `Alt + =` | AutoSum |
| `Ctrl + B` | Bold |
| `Shift + F2` | Cell Comment |
| `Ctrl + '` | Copy Formula from Above |
| `Alt + F1` | Create Chart |
| `F11` | Create Chart Sheet |
| `Ctrl + Shift + 4 ($)` | Currency Format |
| `Ctrl + Shift + 3 (#)` | Date Format |
| `Ctrl + -` | Delete Cells |
| `F2` | Edit Cell |
| `Ctrl + Shift + DOWN` | Extend Selection Down |
| `Ctrl + Shift + LEFT` | Extend Selection Left |
| `Ctrl + Shift + RIGHT` | Extend Selection Right |
| `Ctrl + Shift + UP` | Extend Selection Up |
| `Ctrl + D` | Fill Down |
| `Ctrl + R` | Fill Right |
| `Ctrl + ENTER` | Fill Selection |
| `Ctrl + F` | Find |
| `Ctrl + 1` | Format Cells |
| `Ctrl + Shift + ` (~)` | General Format |
| `F5` | Go To |
| `Ctrl + HOME` | Go to A1 |
| `Ctrl + 0` | Hide Column |
| `Ctrl + 9` | Hide Row |
| `Ctrl + Shift + = (+)` | Insert Cells |
| `Ctrl + ;` | Insert Date |
| `Ctrl + K` | Insert Hyperlink |
| `Ctrl + Shift + ; (:)` | Insert Time |
| `Ctrl + I` | Italic |
| `Ctrl + DOWN` | Jump to Edge Down |
| `Ctrl + LEFT` | Jump to Edge Left |
| `Ctrl + RIGHT` | Jump to Edge Right |
| `Ctrl + UP` | Jump to Edge Up |
| `Ctrl + END` | Last Cell |
| `Alt + ENTER` | New Line in Cell |
| `Shift + F11` | New Worksheet |
| `Ctrl + PAGEDOWN` | Next Sheet |
| `Ctrl + Shift + 1 (!)` | Number Format |
| `Ctrl + Shift + 5 (%)` | Percent Format |
| `Ctrl + PAGEUP` | Previous Sheet |
| `Ctrl + Shift + - (_)` | Remove Borders |
| `F4` | Repeat Last Action |
| `Ctrl + H` | Replace |
| `Ctrl + Shift + 6 (^)` | Scientific Format |
| `Ctrl + SPACE` | Select Column |
| `Shift + SPACE` | Select Row |
| `Ctrl + `` | Show Formulas |
| `Ctrl + Shift + 2 (@)` | Time Format |
| `Ctrl + Shift + L` | Toggle Filter |
| `Ctrl + U` | Underline |
| `Ctrl + Shift + 0 ())` | Unhide Column |
| `Ctrl + Shift + 9 (()` | Unhide Row |

## Profile · Microsoft Outlook

26 mappings · verified against Microsoft 365 Outlook (Windows) · source: Outlook for Windows keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + B` | Address Book |
| `F3` | Find |
| `Ctrl + 6` | Folder List |
| `Ctrl + F` | Forward |
| `Ctrl + 2` | Go to Calendar |
| `Ctrl + 3` | Go to Contacts |
| `Ctrl + Shift + I` | Go to Inbox |
| `Ctrl + 1` | Go to Mail |
| `Ctrl + 5` | Go to Notes |
| `Ctrl + 4` | Go to Tasks |
| `Ctrl + Q` | Mark as Read |
| `Ctrl + U` | Mark as Unread |
| `Ctrl + Shift + V` | Move to Folder |
| `Ctrl + Shift + A` | New Appointment |
| `Ctrl + Shift + C` | New Contact |
| `Ctrl + N` | New Email |
| `Ctrl + Shift + E` | New Folder |
| `Ctrl + Shift + Q` | New Meeting |
| `Ctrl + Shift + M` | New Message |
| `Ctrl + Shift + K` | New Task |
| `Ctrl + R` | Reply |
| `Ctrl + Shift + R` | Reply All |
| `Ctrl + E` | Search |
| `Alt + S` | Send Message |
| `Ctrl + ENTER` | Send Message |
| `F9` | Send/Receive |

## Profile · Microsoft PowerPoint

35 mappings · verified against Microsoft 365 PowerPoint (Windows) · source: PowerPoint for Windows keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + L` | Align Left |
| `Ctrl + R` | Align Right |
| `B` | Black Screen |
| `Ctrl + B` | Bold |
| `Ctrl + E` | Center |
| `Ctrl + Shift + C` | Copy Formatting |
| `Ctrl + Shift + , (<)` | Decrease Font Size |
| `Ctrl + D` | Duplicate Slide |
| `Ctrl + Shift + D` | Duplicate Slide |
| `ESC` | End Slideshow |
| `Ctrl + HOME` | First Slide |
| `Ctrl + T` | Font Dialog |
| `Shift + F5` | From Current Slide |
| `Ctrl + G` | Group Objects |
| `Ctrl + K` | Hyperlink |
| `Ctrl + Shift + . (>)` | Increase Font Size |
| `Ctrl + I` | Italic |
| `Ctrl + J` | Justify |
| `Ctrl + END` | Last Slide |
| `Ctrl + Shift + DOWN` | Move Paragraph Down |
| `Ctrl + Shift + UP` | Move Paragraph Up |
| `Alt + Shift + DOWN` | Move Slide Down |
| `Alt + Shift + UP` | Move Slide Up |
| `Ctrl + M` | New Slide |
| `Ctrl + ENTER` | Next Placeholder |
| `Ctrl + PAGEDOWN` | Next Slide |
| `Ctrl + Shift + V` | Paste Formatting |
| `Alt + F5` | Presenter View |
| `Ctrl + PAGEUP` | Previous Slide |
| `Ctrl + S` | Save |
| `F12` | Save As |
| `F5` | Start Slideshow |
| `Ctrl + U` | Underline |
| `Ctrl + Shift + G` | Ungroup Objects |
| `W` | White Screen |

## Profile · Microsoft Word

45 mappings · verified against Microsoft 365 Word (Windows) · source: Word for Windows keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + L` | Align Left |
| `Ctrl + R` | Align Right |
| `Ctrl + Shift + A` | All Caps |
| `Ctrl + B` | Bold |
| `Ctrl + Shift + L` | Bulleted List |
| `Ctrl + E` | Center |
| `Ctrl + SPACE` | Clear Character Formatting |
| `Ctrl + Q` | Clear Paragraph Formatting |
| `Ctrl + Shift + ENTER` | Column Break |
| `Ctrl + Shift + C` | Copy Formatting |
| `Ctrl + Shift + , (<)` | Decrease Font Size |
| `Ctrl + [` | Decrease Font Size by 1pt |
| `Ctrl + Shift + M` | Decrease Indent |
| `Ctrl + 2` | Double Line Spacing |
| `Ctrl + F` | Find |
| `Ctrl + D` | Font Dialog |
| `Ctrl + Shift + P` | Font Size |
| `Ctrl + G` | Go to Page |
| `Ctrl + Alt + 1` | Heading 1 |
| `Ctrl + Alt + 2` | Heading 2 |
| `Ctrl + Alt + 3` | Heading 3 |
| `Ctrl + K` | Hyperlink |
| `Ctrl + Shift + . (>)` | Increase Font Size |
| `Ctrl + ]` | Increase Font Size by 1pt |
| `Ctrl + M` | Increase Indent |
| `Ctrl + Alt + M` | Insert Comment |
| `Ctrl + Alt + F` | Insert Footnote |
| `Ctrl + I` | Italic |
| `Ctrl + J` | Justify |
| `Ctrl + 5` | Line Spacing 1.5 |
| `Ctrl + Shift + N` | Normal Style |
| `Ctrl + ENTER` | Page Break |
| `Ctrl + Shift + V` | Paste Formatting |
| `Ctrl + H` | Replace |
| `F12` | Save As |
| `Ctrl + Shift + 8 (*)` | Show Formatting Marks |
| `Ctrl + 1` | Single Line Spacing |
| `Ctrl + Shift + K` | Small Caps |
| `F7` | Spelling |
| `Ctrl + =` | Subscript |
| `Ctrl + Shift + = (+)` | Superscript |
| `Shift + F7` | Thesaurus |
| `Ctrl + Shift + E` | Track Changes |
| `Ctrl + U` | Underline |
| `Ctrl + Shift + G` | Word Count |

## Profile · Adobe Illustrator

27 mappings · verified against Illustrator 2024 (Windows) · source: Illustrator default keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + 1` | Actual Size |
| `Ctrl + ]` | Bring Forward |
| `Ctrl + Shift + ] (})` | Bring to Front |
| `Ctrl + 7` | Clipping Mask |
| `Ctrl + 8` | Compound Path |
| `Ctrl + Shift + O` | Create Outlines |
| `Ctrl + 0` | Fit Artboard |
| `Ctrl + G` | Group Selection |
| `Ctrl + 3` | Hide Selection |
| `Ctrl + J` | Join |
| `Ctrl + 2` | Lock Selection |
| `Ctrl + Y` | Outline Mode |
| `Ctrl + B` | Paste in Back |
| `Ctrl + F` | Paste in Front |
| `Ctrl + Alt + 7` | Release Clipping Mask |
| `Ctrl + Alt + 8` | Release Compound Path |
| `Ctrl + [` | Send Backward |
| `Ctrl + Shift + [ ({)` | Send to Back |
| `Ctrl + Alt + 3` | Show All |
| `Ctrl + Shift + H` | Show Edges |
| `Ctrl + Shift + ' (")` | Show Grid |
| `Ctrl + ;` | Show Guides |
| `Ctrl + R` | Show Rulers |
| `Ctrl + U` | Smart Guides |
| `Ctrl + D` | Transform Again |
| `Ctrl + Shift + G` | Ungroup Selection |
| `Ctrl + Alt + 2` | Unlock All |

## Profile · Adobe Photoshop

36 mappings · verified against Photoshop 2024 (Windows) · source: Photoshop default keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + 1` | Actual Pixels |
| `Ctrl + Shift + B` | Auto Color |
| `Ctrl + Alt + Shift + L` | Auto Contrast |
| `Ctrl + Shift + L` | Auto Tone |
| `Ctrl + Shift + ] (})` | Bring to Front |
| `Ctrl + Alt + C` | Canvas Size |
| `Ctrl + B` | Color Balance |
| `Ctrl + Shift + C` | Copy Merged |
| `Ctrl + Alt + G` | Create Clipping Mask |
| `Ctrl + M` | Curves |
| `Ctrl + [` | Decrease Brush Size |
| `Ctrl + Shift + U` | Desaturate |
| `Ctrl + D` | Deselect |
| `Ctrl + J` | Duplicate Layer |
| `Ctrl + 0` | Fit on Screen |
| `Ctrl + T` | Free Transform |
| `Ctrl + G` | Group Layers |
| `Ctrl + U` | Hue/Saturation |
| `Ctrl + Alt + I` | Image Size |
| `Ctrl + ]` | Increase Brush Size |
| `Ctrl + Shift + I` | Inverse Selection |
| `Ctrl + I` | Invert |
| `Ctrl + L` | Levels |
| `Ctrl + E` | Merge Down |
| `Ctrl + Shift + E` | Merge Visible |
| `Ctrl + Shift + N` | New Layer |
| `Ctrl + Shift + V` | Paste Into |
| `Ctrl + Alt + S` | Save a Copy |
| `Ctrl + Alt + Shift + S` | Save for Web |
| `Alt + ]` | Select Layer Above |
| `Alt + [` | Select Layer Below |
| `Ctrl + Shift + [ ({)` | Send to Back |
| `Ctrl + Alt + Shift + E` | Stamp Visible |
| `Ctrl + Alt + Z` | Step Backward |
| `Ctrl + Shift + Z` | Step Forward |
| `Ctrl + Shift + G` | Ungroup Layers |

## Profile · Figma

33 mappings · verified against Figma desktop 2024 (Windows) · source: Figma keyboard shortcuts (Windows)

| Shortcut | Action |
| :--- | :--- |
| `Shift + A` | Auto Layout |
| `Ctrl + ]` | Bring Forward |
| `Ctrl + Shift + ] (})` | Bring to Front |
| `Ctrl + Alt + C` | Copy Properties |
| `Ctrl + Alt + K` | Create Component |
| `Ctrl + D` | Duplicate |
| `Ctrl + Alt + X` | Exclude Selection |
| `Ctrl + Shift + E` | Export |
| `Ctrl + E` | Flatten Selection |
| `Ctrl + Alt + G` | Frame Selection |
| `Ctrl + G` | Group Selection |
| `TAB` | Hide UI |
| `Ctrl + Alt + I` | Intersect Selection |
| `Ctrl + Shift + O` | Outline Stroke |
| `Ctrl + Shift + V` | Paste Over Selection |
| `Ctrl + Alt + V` | Paste Properties |
| `Ctrl + Shift + K` | Place Image |
| `Ctrl + K` | Quick Actions |
| `Ctrl + R` | Rename Layer |
| `Ctrl + /` | Search Menu |
| `Ctrl + [` | Send Backward |
| `Ctrl + Shift + [ ({)` | Send to Back |
| `Ctrl + Alt + S` | Subtract Selection |
| `Shift + G` | Toggle Layout Grid |
| `Shift + O` | Toggle Outlines |
| `Ctrl + Shift + ' (")` | Toggle Pixel Grid |
| `Shift + R` | Toggle Rulers |
| `Ctrl + Shift + G` | Ungroup Selection |
| `Ctrl + Alt + U` | Union Selection |
| `Ctrl + Alt + M` | Use as Mask |
| `Shift + 0 ())` | Zoom to 100% |
| `Shift + 1 (!)` | Zoom to Fit |
| `Shift + 2 (@)` | Zoom to Selection |

## Profile · Adobe Premiere Pro

27 mappings · verified against Premiere Pro 2024 (Windows) · source: Premiere Pro default keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + K` | Add Edit |
| `Ctrl + Shift + K` | Add Edit to All Tracks |
| `M` | Add Marker |
| `Ctrl + Shift + D` | Apply Audio Transition |
| `Ctrl + D` | Apply Video Transition |
| `Ctrl + Shift + E` | Export Frame |
| `Ctrl + M` | Export Media |
| `Ctrl + I` | Import Media |
| `X` | Mark Clip |
| `I` | Mark In |
| `O` | Mark Out |
| `Ctrl + N` | New Sequence |
| `Shift + M` | Next Marker |
| `K` | Pause |
| `J` | Play Backward |
| `L` | Play Forward |
| `SPACE` | Play/Pause |
| `Ctrl + Shift + M` | Previous Marker |
| `Ctrl + S` | Save |
| `Ctrl + R` | Speed/Duration |
| `LEFT` | Step Back |
| `Shift + LEFT` | Step Back 5 Frames |
| `RIGHT` | Step Forward |
| `Shift + RIGHT` | Step Forward 5 Frames |
| `Shift + = (+)` | Zoom In Timeline |
| `-` | Zoom Out Timeline |
| `\` | Zoom to Sequence |

## Profile · OBS Studio

9 mappings · verified against OBS Studio 30.x (Windows) · source: OBS Studio default hotkeys (Settings > Hotkeys)

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + C` | Copy |
| `Ctrl + D` | Duplicate |
| `Ctrl + V` | Paste |
| `Ctrl + Shift + Z` | Redo |
| `Ctrl + Y` | Redo |
| `DELETE` | Remove |
| `F2` | Rename |
| `Ctrl + S` | Save Scene Collection |
| `Ctrl + Z` | Undo |

## Profile · VLC

26 mappings · verified against VLC 3.0 (Windows) · source: VLC hotkeys (Preferences > Hotkeys)

| Shortcut | Action |
| :--- | :--- |
| `A` | Cycle Aspect Ratio |
| `B` | Cycle Audio Track |
| `C` | Cycle Crop |
| `V` | Cycle Subtitles |
| `ESC` | Exit Full Screen |
| `F` | Full Screen |
| `Ctrl + T` | Go to Time |
| `I` | Media Information |
| `M` | Mute |
| `E` | Next Frame |
| `N` | Next Track |
| `=` | Normal Speed |
| `Ctrl + O` | Open File |
| `Ctrl + N` | Open Network |
| `SPACE` | Play/Pause |
| `Ctrl + P` | Preferences |
| `P` | Previous Track |
| `Ctrl + LEFT` | Seek Backward |
| `Ctrl + RIGHT` | Seek Forward |
| `T` | Show Time |
| `-` | Slow Down |
| `Shift + = (+)` | Speed Up |
| `S` | Stop |
| `Ctrl + L` | Toggle Playlist |
| `Ctrl + DOWN` | Volume Down |
| `Ctrl + UP` | Volume Up |

## Profile · Discord

16 mappings · verified against Discord desktop (Windows) · source: Discord keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + B` | Bold |
| `Ctrl + F` | Find |
| `Ctrl + I` | Italic |
| `Ctrl + /` | Keyboard Shortcuts |
| `ESC` | Mark Channel Read |
| `Shift + ESC` | Mark Server Read |
| `Shift + ENTER` | New Line |
| `Alt + DOWN` | Next Channel |
| `Alt + Shift + DOWN` | Next Server |
| `Alt + UP` | Previous Channel |
| `Alt + Shift + UP` | Previous Server |
| `Ctrl + K` | Quick Switcher |
| `ENTER` | Send Message |
| `Ctrl + Shift + D` | Toggle Deafen |
| `Ctrl + Shift + M` | Toggle Mute |
| `Ctrl + U` | Underline |

## Profile · Microsoft Teams

22 mappings · verified against Microsoft Teams (work or school, Windows) · source: Microsoft Teams keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + A` | Accept Video Call |
| `Ctrl + B` | Bold |
| `Ctrl + Shift + D` | Decline Video Call |
| `Ctrl + 1` | Go to Activity |
| `Ctrl + 4` | Go to Calendar |
| `Ctrl + 5` | Go to Calls |
| `Ctrl + 2` | Go to Chat |
| `Ctrl + 6` | Go to Files |
| `Ctrl + 3` | Go to Teams |
| `Ctrl + Shift + H` | Hang Up |
| `Ctrl + I` | Italic |
| `Ctrl + N` | New Chat |
| `Shift + ENTER` | New Line |
| `Ctrl + Shift + K` | Raise Hand |
| `Ctrl + Y` | Redo |
| `Ctrl + K` | Search/Command |
| `ENTER` | Send Message |
| `Ctrl + Shift + E` | Start Screen Share |
| `Ctrl + Shift + M` | Toggle Mute |
| `Ctrl + Shift + O` | Toggle Video |
| `Ctrl + U` | Underline |
| `Ctrl + Z` | Undo |

## Profile · Slack

24 mappings · verified against Slack 4.x desktop (Windows) · source: Slack keyboard shortcuts (Windows)

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + M` | Activity |
| `Ctrl + Shift + A` | All Unreads |
| `Ctrl + Shift + . (>)` | Block Quote |
| `Ctrl + B` | Bold |
| `Ctrl + Shift + L` | Browse Channels |
| `Ctrl + Shift + K` | Browse DMs |
| `Ctrl + Shift + C` | Code Block |
| `Ctrl + UP` | Edit Last Message |
| `Ctrl + I` | Italic |
| `Shift + ESC` | Mark All as Read |
| `Shift + ENTER` | New Line |
| `Ctrl + N` | New Message |
| `Alt + DOWN` | Next Channel |
| `Ctrl + Shift + E` | People |
| `Ctrl + ,` | Preferences |
| `Alt + UP` | Previous Channel |
| `Ctrl + K` | Quick Switcher |
| `Ctrl + Shift + Z` | Redo |
| `Ctrl + F` | Search |
| `Ctrl + Shift + F` | Search All Channels |
| `ENTER` | Send Message |
| `Ctrl + Shift + Y` | Set Status |
| `Ctrl + Shift + T` | Threads |
| `Ctrl + Z` | Undo |

## Profile · WeChat

6 mappings · verified against WeChat 3.9 (Windows) · source: WeChat for Windows shortcuts

| Shortcut | Action |
| :--- | :--- |
| `ESC` | Close Window |
| `Ctrl + ENTER` | New Line |
| `Ctrl + Alt + W` | Open WeChat |
| `Alt + A` | Screenshot |
| `Ctrl + F` | Search |
| `ENTER` | Send Message |

## Profile · Zoom

13 mappings · verified against Zoom Workplace 6.x (Windows) · source: Zoom keyboard shortcuts (Windows)

| Shortcut | Action |
| :--- | :--- |
| `Alt + Q` | End Meeting |
| `Alt + F` | Enter Full Screen |
| `ESC` | Exit Full Screen |
| `Alt + M` | Mute All |
| `Alt + U` | Participants Panel |
| `Alt + P` | Pause Recording |
| `Alt + H` | Raise Hand |
| `Alt + C` | Start Cloud Recording |
| `Alt + R` | Start Recording |
| `Alt + S` | Start Screen Share |
| `Alt + N` | Switch Camera |
| `Alt + A` | Toggle Mute |
| `Alt + V` | Toggle Video |

## Profile · Notion

20 mappings · verified against Notion desktop (Windows) · source: Notion keyboard shortcuts

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + Shift + M` | Add Comment |
| `Ctrl + B` | Bold |
| `Ctrl + D` | Duplicate Block |
| `TAB` | Indent Block |
| `Ctrl + E` | Inline Code |
| `Ctrl + K` | Insert Link |
| `Ctrl + I` | Italic |
| `Ctrl + [` | Navigate Back |
| `Ctrl + ]` | Navigate Forward |
| `ENTER` | New Block |
| `Shift + ENTER` | New Line |
| `Ctrl + N` | New Page |
| `Ctrl + Shift + N` | New Window |
| `Shift + TAB` | Outdent Block |
| `Ctrl + P` | Quick Find |
| `Ctrl + A` | Select All |
| `Ctrl + Shift + S` | Strikethrough |
| `Ctrl + Shift + L` | Toggle Dark Mode |
| `Ctrl + \` | Toggle Sidebar |
| `Ctrl + U` | Underline |

## Profile · Obsidian

24 mappings · verified against Obsidian 1.6 (Windows) · source: Obsidian default hotkeys

| Shortcut | Action |
| :--- | :--- |
| `Ctrl + W` | Close Tab |
| `Ctrl + Shift + W` | Close Window |
| `Ctrl + P` | Command Palette |
| `Ctrl + F` | Find in Note |
| `Alt + ENTER` | Follow Link |
| `TAB` | Indent |
| `Ctrl + K` | Insert Link |
| `Ctrl + Alt + LEFT` | Navigate Back |
| `Ctrl + Alt + RIGHT` | Navigate Forward |
| `Ctrl + N` | New Note |
| `Ctrl + G` | Open Graph View |
| `Ctrl + Shift + G` | Open Local Graph |
| `Ctrl + ,` | Open Settings |
| `Shift + TAB` | Outdent |
| `Ctrl + Shift + V` | Paste and Match Style |
| `Ctrl + O` | Quick Open |
| `F2` | Rename File |
| `Ctrl + Shift + T` | Reopen Closed Tab |
| `Ctrl + H` | Replace in Note |
| `Ctrl + Shift + F` | Search in Vault |
| `Ctrl + B` | Toggle Bold |
| `Ctrl + ENTER` | Toggle Checkbox |
| `Ctrl + I` | Toggle Italic |
| `Ctrl + E` | Toggle Reading View |

## Profile · File Explorer

34 mappings · verified against Windows 11 23H2 · source: Windows 11 File Explorer shortcuts

| Shortcut | Action |
| :--- | :--- |
| `F4` | Address Bar History |
| `Alt + LEFT` | Back |
| `BACKSPACE` | Back |
| `NUM-` | Collapse Folder |
| `Ctrl + Shift + 8 (*)` | Content View |
| `Ctrl + Shift + C` | Copy Path |
| `DELETE` | Delete |
| `Alt + Shift + P` | Details Pane |
| `Ctrl + Shift + 6 (^)` | Details View |
| `NUM*` | Expand All Folders |
| `NUM+` | Expand Folder |
| `Ctrl + Shift + 1 (!)` | Extra Large Icons |
| `Alt + D` | Focus Address Bar |
| `Ctrl + L` | Focus Address Bar |
| `Alt + RIGHT` | Forward |
| `F11` | Full Screen |
| `Ctrl + Shift + 2 (@)` | Large Icons |
| `Ctrl + Shift + 5 (%)` | List View |
| `Ctrl + Shift + 3 (#)` | Medium Icons |
| `Ctrl + Shift + N` | New Folder |
| `Ctrl + N` | New Window |
| `Ctrl + Shift + V` | Paste Shortcut |
| `Shift + DELETE` | Permanently Delete |
| `Alt + P` | Preview Pane |
| `Alt + ENTER` | Properties |
| `Ctrl + Y` | Redo |
| `F5` | Refresh |
| `F2` | Rename |
| `Ctrl + F` | Search |
| `F3` | Search |
| `Ctrl + A` | Select All |
| `Ctrl + Shift + 4 ($)` | Small Icons |
| `Ctrl + Z` | Undo |
| `Alt + UP` | Up One Level |

## Profile · Total Commander

77 mappings · verified against Total Commander 11 · source: Total Commander Keyboard.txt (installation directory)

| Shortcut | Action |
| :--- | :--- |
| `F10` | Activate Left Menu |
| `F9` | Activate Menu |
| `Alt + LEFT` | Back |
| `Ctrl + B` | Branch View |
| `Ctrl + Shift + B` | Branch View (Selected) |
| `Ctrl + F1` | Brief View |
| `Ctrl + L` | Calculate Size |
| `Alt + F1` | Change Left Drive |
| `Alt + F2` | Change Right Drive |
| `Ctrl + Shift + W` | Close All Tabs |
| `Ctrl + W` | Close Tab |
| `Ctrl + Shift + F2` | Comments View |
| `Ctrl + F` | Connect to FTP |
| `Shift + F10` | Context Menu |
| `F5` | Copy Files |
| `Ctrl + P` | Copy Path to Command Line |
| `Shift + F5` | Copy in Place |
| `Ctrl + C` | Copy to Clipboard |
| `Ctrl + X` | Cut to Clipboard |
| `DELETE` | Delete Files |
| `F8` | Delete Files |
| `NUM-` | Deselect by Mask |
| `Alt + DOWN` | Directory History |
| `Ctrl + D` | Directory Hotlist |
| `Ctrl + F8` | Directory Tree |
| `Ctrl + Shift + F` | Disconnect from FTP |
| `Ctrl + Z` | Edit Comment |
| `F4` | Edit File |
| `Alt + F7` | Find Files |
| `Alt + RIGHT` | Forward |
| `Ctrl + F2` | Full View |
| `F1` | Help |
| `NUM*` | Invert Selection |
| `Ctrl + \` | Jump to Root |
| `Ctrl + M` | Multi-Rename Tool |
| `F7` | New Folder |
| `Shift + F7` | New Folder in Other Panel |
| `Ctrl + T` | New Tab |
| `Ctrl + Shift + T` | New Tab in Background |
| `Shift + F4` | New Text File |
| `Ctrl + TAB` | Next Tab |
| `Ctrl + PAGEDOWN` | Open Directory / Archive |
| `Ctrl + LEFT` | Open in Target Panel |
| `Ctrl + RIGHT` | Open in Target Panel |
| `Alt + F5` | Pack Files |
| `Ctrl + V` | Paste to Current Directory |
| `Ctrl + Shift + TAB` | Previous Tab |
| `Alt + ENTER` | Properties |
| `Ctrl + S` | Quick Filter |
| `Ctrl + Shift + S` | Quick Filter (Last Used) |
| `Ctrl + Q` | Quick View |
| `F6` | Rename / Move Files |
| `Shift + F6` | Rename in Place |
| `F2` | Reread Source |
| `SPACE` | Select / Count Size |
| `Ctrl + A` | Select All |
| `INSERT` | Select File |
| `NUM+` | Select by Mask |
| `Ctrl + F10` | Show All Files |
| `Ctrl + F11` | Show Only Programs |
| `Ctrl + F5` | Sort by Date |
| `Ctrl + F4` | Sort by Extension |
| `Ctrl + F3` | Sort by Name |
| `Ctrl + F6` | Sort by Size |
| `Ctrl + U` | Swap Panels |
| `Ctrl + Shift + U` | Swap Panels and Tabs |
| `TAB` | Switch Panel |
| `Ctrl + I` | Switch to Target Directory |
| `Shift + TAB` | Switch to Tree |
| `Ctrl + Shift + F1` | Thumbnails View |
| `Alt + F6` | Unpack Files |
| `Alt + F9` | Unpack Files |
| `Ctrl + F7` | Unsorted |
| `BACKSPACE` | Up One Level |
| `Ctrl + PAGEUP` | Up One Level |
| `Ctrl + F12` | User-Defined Files |
| `F3` | View File |

