// Code generated by "stringer -type=BlockState"; DO NOT EDIT

package game

import "fmt"

const _BlockState_name = "BlockStaticBlockSwappingFromLeftBlockSwappingFromRightBlockDroppingFromAboveBlockFlashingBlockCrackingBlockCrackedBlockExplodingBlockExplodedBlockClearPausingBlockCleared"

var _BlockState_index = [...]uint8{0, 11, 32, 54, 76, 89, 102, 114, 128, 141, 158, 170}

func (i BlockState) String() string {
	if i < 0 || i >= BlockState(len(_BlockState_index)-1) {
		return fmt.Sprintf("BlockState(%d)", i)
	}
	return _BlockState_name[_BlockState_index[i]:_BlockState_index[i+1]]
}
