// Code generated by go-bindata.
// sources:
// data/CPMono_v07 Bold.ttf
// data/CPMono_v07 Plain.ttf
// data/clear.wav
// data/meshes.obj
// data/move.wav
// data/select.wav
// data/shader.frag
// data/shader.vert
// data/swap.wav
// data/texture.png
// data/thud.wav
// DO NOT EDIT!

package asset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// bindataRead reads the given file from disk. It returns an error on failure.
func bindataRead(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

// cpmono_v07BoldTtf reads file data from disk. It returns an error on failure.
func cpmono_v07BoldTtf() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/CPMono_v07 Bold.ttf"
	name := "CPMono_v07 Bold.ttf"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// cpmono_v07PlainTtf reads file data from disk. It returns an error on failure.
func cpmono_v07PlainTtf() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/CPMono_v07 Plain.ttf"
	name := "CPMono_v07 Plain.ttf"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// clearWav reads file data from disk. It returns an error on failure.
func clearWav() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/clear.wav"
	name := "clear.wav"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// meshesObj reads file data from disk. It returns an error on failure.
func meshesObj() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/meshes.obj"
	name := "meshes.obj"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// moveWav reads file data from disk. It returns an error on failure.
func moveWav() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/move.wav"
	name := "move.wav"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// selectWav reads file data from disk. It returns an error on failure.
func selectWav() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/select.wav"
	name := "select.wav"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// shaderFrag reads file data from disk. It returns an error on failure.
func shaderFrag() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/shader.frag"
	name := "shader.frag"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// shaderVert reads file data from disk. It returns an error on failure.
func shaderVert() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/shader.vert"
	name := "shader.vert"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// swapWav reads file data from disk. It returns an error on failure.
func swapWav() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/swap.wav"
	name := "swap.wav"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// texturePng reads file data from disk. It returns an error on failure.
func texturePng() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/texture.png"
	name := "texture.png"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// thudWav reads file data from disk. It returns an error on failure.
func thudWav() (*asset, error) {
	path := "/home/btmura/work/go/src/github.com/btmura/blockcillin/internal/asset/data/thud.wav"
	name := "thud.wav"
	bytes, err := bindataRead(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"CPMono_v07 Bold.ttf": cpmono_v07BoldTtf,
	"CPMono_v07 Plain.ttf": cpmono_v07PlainTtf,
	"clear.wav": clearWav,
	"meshes.obj": meshesObj,
	"move.wav": moveWav,
	"select.wav": selectWav,
	"shader.frag": shaderFrag,
	"shader.vert": shaderVert,
	"swap.wav": swapWav,
	"texture.png": texturePng,
	"thud.wav": thudWav,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"CPMono_v07 Bold.ttf": &bintree{cpmono_v07BoldTtf, map[string]*bintree{}},
	"CPMono_v07 Plain.ttf": &bintree{cpmono_v07PlainTtf, map[string]*bintree{}},
	"clear.wav": &bintree{clearWav, map[string]*bintree{}},
	"meshes.obj": &bintree{meshesObj, map[string]*bintree{}},
	"move.wav": &bintree{moveWav, map[string]*bintree{}},
	"select.wav": &bintree{selectWav, map[string]*bintree{}},
	"shader.frag": &bintree{shaderFrag, map[string]*bintree{}},
	"shader.vert": &bintree{shaderVert, map[string]*bintree{}},
	"swap.wav": &bintree{swapWav, map[string]*bintree{}},
	"texture.png": &bintree{texturePng, map[string]*bintree{}},
	"thud.wav": &bintree{thudWav, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

