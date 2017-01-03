package multibtree

import (
	"bufio"
	"fmt"
	"strings"

	"hash/crc64"
	"os"

	"github.com/google/btree"
)

type Item string

type Iterator func(item string)

func (s Item) Less(than btree.Item) bool {
	if strings.Compare(string(s), string(than.(Item))) == -1 {
		return true
	}
	return false
}

type MultiBtree struct {
	maxBuffer uint32
	groups    int
	path      string
	prefix    string
	buffers   map[int]*os.File
}

func NewMultiTree(path string, groups int) (*MultiBtree, error) {
	tree := new(MultiBtree)
	tree.groups = groups
	tree.buffers = make(map[int]*os.File, groups)
	tree.path = path + "/tree_buffers/"
	tree.prefix = "temp"

	err := os.MkdirAll(tree.path, os.ModePerm)

	if err != nil {
		return nil, err
	}

	for i := 0; i < groups; i++ {
		f, ferr := tree.openFile(i)

		if ferr != nil {
			return nil, ferr
		}
		tree.buffers[i] = f
	}
	return tree, nil
}

func (t *MultiBtree) Insert(item string) error {
	group := strHash(item, t.groups)
	_, err := t.buffers[group].Write([]byte(item + "\n"))
	return err
}

func (t *MultiBtree) Scan(c Iterator) {
	for _, f := range t.buffers {
		f.Seek(0, 0)
		tree := btree.New(4)
		scanner := bufio.NewScanner(f)
		buf := make([]byte, 1024)
		scanner.Buffer(buf, 1024)
		for scanner.Scan() {
			tree.ReplaceOrInsert(Item(scanner.Text()))
		}
		tree.Ascend(func(item btree.Item) bool {
			c(string(item.(Item)))
			return true
		})
	}
}

func (t *MultiBtree) openFile(group int) (*os.File, error) {
	file := fmt.Sprintf("%s/%s_%d.buffer", t.path, t.prefix, group)
	return os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
}

func (t *MultiBtree) Close() {
	for _, f := range t.buffers {
		f.Close()
	}
	fmt.Println("closed")
	//os.RemoveAll(t.path)

}

func strHash(str string, max int) int {
	crcTable := crc64.MakeTable(crc64.ECMA)
	checksum64 := crc64.Checksum([]byte(str), crcTable)
	return int(checksum64 % uint64(max))
}
