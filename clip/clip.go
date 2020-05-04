package clip

import (
	"github.com/atotto/clipboard"
	"github.com/benjamin-daniel/clippy/hash"
)

type ClipItem struct {
	Text string
	Hash string
}

func New() *ClipItem {
	text, err := clipboard.ReadAll()
	if err != nil {
		panic("There was an error getting your clipboard")
	}
	hash, err := hash.GetHash(text)
	if err != nil {
		panic(err)
	}
	return &ClipItem{Text: text, Hash: hash}
}

// // GetClipboard returns the clipboard value
// func GetClipboard (string , error) {
// 	return clipboard.ReadAll()
// }
