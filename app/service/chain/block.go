package chain

import (
	"crypto/sha256"
	"strconv"
	"encoding/json"
	"log"
	"fmt"
)

type BoomBlock struct {
	prevHash string
	index int
	author BlockAuthor
	timestamp int64
	data []byte
	hash []byte
}


type BlockAuthor struct {
	UserCode string
	GroupCode string
	MasterCode string
}


type BlockHandler interface {
	Init() BlockHandler
	Index() int
	SetIndex(ind int) BlockHandler
	Timestamp() int64
	SetTimestamp(timestamp int64) BlockHandler
	Data() []byte
	SetData(data []byte) BlockHandler
	Author() BlockAuthor
	SetAuthor(person string, group string, master string) BlockHandler
	GenHash() []byte
	ValidateHash() bool
	PreviousHash() string
	SetPreviousHash(hash string) BlockHandler
	Hash() []byte
	SetHash(hash []byte) BlockHandler
	HashString() string
}

func (b BoomBlock) Init() BlockHandler {
	return b
}

func (b BoomBlock) Index() int {
	return b.index
}

func (b BoomBlock) SetIndex(ind int) BlockHandler {
	b.index = ind

	return b
}

func (b BoomBlock) Timestamp() int64 {
	return b.timestamp
}

func (b BoomBlock) SetTimestamp(time int64) BlockHandler {
	b.timestamp = time

	return b
}

func (b BoomBlock) Data() []byte {
	return b.data
}

func (b BoomBlock) SetData(data []byte) BlockHandler {
	b.data = data

	return b
}


func (b BoomBlock) Author() BlockAuthor {
	return b.author
}

func (b BoomBlock) SetAuthor(p string, g string, m string) BlockHandler {
	if m == "" {
		log.Panic("User and Group may be empty but master must include a string as its the primary author identifier")
	}

	b.author = BlockAuthor{}
	b.author.MasterCode = m
	b.author.GroupCode = g
	b.author.UserCode = p

	return b
}

func (block BoomBlock) GenHash() []byte {
	sha := sha256.New()
	index := []byte(strconv.Itoa(block.Index()))
	auth, _ := json.Marshal(block.Author())
	sha.Write(index)
	sha.Write([]byte(":"))
	sha.Write(auth)
	sha.Write([]byte(":"))
	sha.Write([]byte(block.PreviousHash()))
	sha.Write([]byte(":"))
	sha.Write(block.data)
	log.Print(sha.Sum(nil))

	return sha.Sum(nil)
}

func (b BoomBlock) ValidateHash() bool {
	return true
}

func (b BoomBlock) SetPreviousHash(hash string) BlockHandler {
	b.prevHash = hash

	return b
}

func (b BoomBlock) PreviousHash() string {
	return b.prevHash
}

func (b BoomBlock) Hash() []byte {
	return b.hash
}

func (b BoomBlock) SetHash(hash []byte) BlockHandler {
	b.hash = hash

	return b
}

func (b BoomBlock) HashString() string {
	return fmt.Sprintf("%x",b.GenHash())
}
