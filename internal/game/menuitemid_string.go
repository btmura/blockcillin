// Code generated by "stringer -type=MenuItemID"; DO NOT EDIT

package game

import "fmt"

const _MenuItemID_name = "MenuNewGameItemMenuStatsMenuOptionsMenuCreditsMenuExitMenuSpeedMenuDifficultyMenuOKMenuContinueGameMenuQuit"

var _MenuItemID_index = [...]uint8{0, 15, 24, 35, 46, 54, 63, 77, 83, 99, 107}

func (i MenuItemID) String() string {
	if i >= MenuItemID(len(_MenuItemID_index)-1) {
		return fmt.Sprintf("MenuItemID(%d)", i)
	}
	return _MenuItemID_name[_MenuItemID_index[i]:_MenuItemID_index[i+1]]
}
