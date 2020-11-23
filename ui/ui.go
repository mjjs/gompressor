package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/mjjs/gompressor/fileio"
	"github.com/mjjs/gompressor/huffman"
	"github.com/mjjs/gompressor/lzw"
	"github.com/mjjs/gompressor/vector"
	"github.com/rivo/tview"
)

type algorithm int

const (
	algorithmLZW algorithm = iota
	algorithmHuffman
)

type action int

const (
	actionCompress action = iota
	actionDecompress
)

// UI handles running the user interface
type UI struct {
	application *tview.Application
}

func New() *UI {
	return &UI{
		application: tview.NewApplication(),
	}
}

// Run starts the user interface main loop
func (u UI) Run() {
	list := u.getStartMenu()

	err := u.application.SetRoot(list, true).SetFocus(list).Run()
	if err != nil {
		panic(err)
	}
}

func (u UI) getStartMenu() *tview.List {
	list := tview.NewList().
		AddItem("Compress", "Compress a file", 'c', u.algorithmSelect(actionCompress)).
		AddItem("Decompress", "Decompress a file", 'd', u.algorithmSelect(actionDecompress)).
		AddItem("Quit", "Exit the application", 'q', u.application.Stop)

	return list
}

func (u UI) fileSelect(action action, algorithm algorithm) func() {
	return func() {
		rootDir := "."
		root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorRed)
		tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

		if action == actionCompress {
			u.buildFileTree(root, rootDir, action, algorithm)
		} else if algorithm == algorithmHuffman {
			u.buildFileTree(root, rootDir, action, algorithm)
		} else if algorithm == algorithmLZW {
			u.buildFileTree(root, rootDir, action, algorithm)
		} else {
			panic("This should not happen!")
		}

		u.application.SetRoot(tree, true)
	}
}

func (u UI) algorithmSelect(action action) func() {
	return func() {
		list := tview.NewList().
			AddItem("Huffman", "Huffman coding", 'h', u.fileSelect(action, algorithmHuffman)).
			AddItem("LZW", "Lempel-Ziv-Welch", 'l', u.fileSelect(action, algorithmLZW))

		u.application.SetRoot(list, true)
	}
}

func (u UI) buildFileTree(target *tview.TreeNode, path string, action action, algorithm algorithm) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetSelectable(true)

		if file.IsDir() {
			node.SetColor(tcell.ColorBlue)
			node.SetSelectedFunc(u.selectDir(node, path, file, action, algorithm))
		} else {
			node.SetSelectedFunc(u.selectFile(filepath.Join(path, file.Name()), action, algorithm))
		}

		ext := algorithmToExtension(algorithm)

		if file.IsDir() || action == actionCompress || ext == filepath.Ext(file.Name()) {
			target.AddChild(node)
		}
	}
}

func (u UI) selectDir(node *tview.TreeNode, path string, file os.FileInfo, action action, algorithm algorithm) func() {
	return func() {
		children := node.GetChildren()

		if len(children) == 0 {
			u.buildFileTree(node, filepath.Join(path, file.Name()), action, algorithm)
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	}
}

func (u UI) selectFile(filepath string, action action, algorithm algorithm) func() {
	return func() {
		var outFilename string

		if action == actionCompress {
			outFilename = fmt.Sprintf("%s.lzw", filepath)
			if algorithm == algorithmHuffman {
				outFilename = fmt.Sprintf("%s.huff", filepath)
			}

			bytes, err := fileio.ReadFile(filepath)
			if err != nil {
				panic(err)
			}

			var compressed *vector.Vector

			if algorithm == algorithmLZW {
				compressed, err = lzw.Compress(bytes)
				if err != nil {
					panic(err)
				}

				err = fileio.WriteLZWFile(compressed, outFilename)
				if err != nil {
					panic(err)
				}
			} else {
				compressed = huffman.Compress(bytes)
				err = fileio.WriteFile(compressed, outFilename)
				if err != nil {
					panic(err)
				}
			}
		} else {
			outFilename = fmt.Sprintf("%s.decompressed", filepath)
			if algorithm == algorithmHuffman {
				outFilename = fmt.Sprintf("%s.decompressed", filepath)
			}

			if algorithm == algorithmLZW {
				codes, err := fileio.ReadLZWFile(filepath)
				if err != nil {
					panic(err)
				}

				decompressed, err := lzw.Decompress(codes)
				if err != nil {
					panic(err)
				}

				err = fileio.WriteFile(decompressed, outFilename)
				if err != nil {
					panic(err)
				}
			} else {
				bytes, err := fileio.ReadFile(filepath)
				if err != nil {
					panic(err)
				}

				decompressed, err := huffman.Decompress(bytes)
				if err != nil {
					panic(err)
				}

				err = fileio.WriteFile(decompressed, outFilename)
				if err != nil {
					panic(err)
				}
			}
		}

		u.fileWrittenView(outFilename)
	}
}

func (u UI) fileWrittenView(filename string) {
	textView := tview.NewTextView().SetDoneFunc(func(tcell.Key) {
		list := u.getStartMenu()
		u.application.SetRoot(list, true)
	})

	fmt.Fprintf(textView, "Wrote to file %s\n\nPress any key to return to main menu", filename)

	u.application.SetRoot(textView, true)
}

func algorithmToExtension(a algorithm) string {
	if a == algorithmHuffman {
		return ".huff"
	}
	return ".lzw"
}
