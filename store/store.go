package store

import (
	// "fmt"

	// "github.com/benjamin-daniel/clippy/clip"
	"fmt"

	"github.com/atotto/clipboard"
	hash "github.com/benjamin-daniel/clippy/hash"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// ClipBoardItem holds the string of the clipboard and hash
type ClipBoardItem struct {
	gorm.Model
	Text string `gorm:"type:MEDIUMTEXT"`
	Hash string
}

func (clip *ClipBoardItem) String() string {
	return fmt.Sprintf("Text: %s \nHash: %s\nCreated: %s\n", clip.Text, clip.Hash, clip.CreatedAt)
}

// TruncateText is used to truncate text
func (clip *ClipBoardItem) TruncateText(num int) (bnoden string) {
	str := clip.Text
	// bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

// GetLast returns last clipboard item
func GetLast(db *gorm.DB) *ClipBoardItem {
	clip := &ClipBoardItem{}
	db.Last(clip)
	return clip
}

// AddIfNotPresent added the text to the db if the text isn't the last in the db
func AddIfNotPresent() bool {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&ClipBoardItem{})

	// clipItem := New()
	// Create
	currentClip := New()

	// this handles when we copy images
	if currentClip.Text == "" {
		return false
	}
	lastClip := GetLast(db)
	if currentClip.Hash != lastClip.Hash {
		db.Create(currentClip)
		// comment this out to stop go-daemon from coping every clipboard action to the log file
		// fmt.Println(currentClip)
		return true
	}
	return false
}

// New Create and returns a new ClipBoardItem
func New() *ClipBoardItem {
	text, err := clipboard.ReadAll()
	if err != nil {
		panic("There was an error getting your clipboard")
	}
	hash, err := hash.GetHash(text)
	if err != nil {
		panic(err)
	}
	return &ClipBoardItem{Text: text, Hash: hash}
}
