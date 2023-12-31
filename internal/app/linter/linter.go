package linter

import (
	"ck3-parser/internal/app/parser"
	"fmt"
	"os"
	"strings"
)

type Linter struct {
	ParseTree []*parser.Node `json:"tree"`
	Level     int            `json:"level"`
	towrite   []byte
	config    LintConfig
}

func New(parsetree []*parser.Node, config LintConfig) *Linter {
	return &Linter{
		ParseTree: parsetree,
		Level:     0,
		towrite:   []byte{},
		config:    config,
	}
}

func (l *Linter) Lint() []byte {
	for _, node := range l.ParseTree {
		l.lintNode(node)
		// if i != len(l.ParseTree)-1 {
		// 	l.towrite = append(l.towrite, byte('\n'))
		// }
	}

	return l.towrite
}

func (l *Linter) lintNode(node *parser.Node) {
	switch node.Type {
	case parser.Comment:
		l.lintComment(node)
	case parser.Property:
		l.lintProperty(node)
	case parser.Comparison:
		l.lintComparison(node)
	case parser.Block:
		l.lintBlock(node)
	default:
		panic(fmt.Sprintf("[Linter] Unexpected NodeType: %q, with value of: %s",
			node.Type, node.Value))

	}
}

func (l *Linter) lintComment(node *parser.Node) {
	l.intend()

	l.towrite = append(l.towrite, node.DataLiteral()...)

	l.nextLine()

	// if l.singleline {
	// 	l.towrite = append(l.towrite, byte(' '))
	// } else if l.Level != 0 {
	// 	l.towrite = append(l.towrite, byte('\n'))
	// }
}

func (l *Linter) lintProperty(node *parser.Node) {
	l.intend()

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.operator("=")
	l.towrite = append(l.towrite, node.DataLiteral()...)

	l.nextLine()
}

func (l *Linter) lintComparison(node *parser.Node) {
	l.intend()

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.operator(node.Operator)
	l.towrite = append(l.towrite, node.DataLiteral()...)

	l.nextLine()
}

func (l *Linter) lintBlock(node *parser.Node) {
	l.intend()

	l.Level++

	l.towrite = append(l.towrite, node.KeyLiteral()...)
	l.operator("=")
	l.towrite = append(l.towrite, byte('{'))

	l.nextLine()

	children := node.Value.([]*parser.Node)
	for _, c := range children {
		l.lintNode(c)
	}
	l.Level--
	l.intend()
	l.towrite = append(l.towrite, byte('}'))

	l.nextLine()
}

func (l *Linter) intend() {
	for i := 0; i < l.Level; i++ {
		if l.config.IntendStyle == SPACES {
			intend := strings.Repeat(" ", l.config.IntendSize)
			l.towrite = append(l.towrite, []byte(intend)...)
		} else {
			l.towrite = append(l.towrite, byte('\t'))
		}
	}
}

func (l *Linter) operator(operator string) {
	l.towrite = append(l.towrite, []byte(" "+operator+" ")...)
}

func (l *Linter) nextLine() {
	l.towrite = append(l.towrite, l.config.EndOfLine...)
}

// saves file with utf8bom encoding:
func (l *Linter) Save(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write UTF-8 BOM
	if l.config.CharSet == "utf-8-bom" {
		bom := []byte{0xEF, 0xBB, 0xBF}
		file.Write(bom)
	}

	_, err = file.Write(l.towrite)
	if err != nil {
		return err
	}

	return nil
}
