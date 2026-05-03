package widgets

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/factory"
	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
)

type editableNode interface {
	toTree() tree.FormulaNode
}

type eDiceNode struct {
	count int
	sides int
}

func (n *eDiceNode) toTree() tree.FormulaNode {
	return &tree.DiceNode{Count: n.count, Sides: n.sides}
}

type eFlatNode struct {
	value int
}

func (n *eFlatNode) toTree() tree.FormulaNode {
	return &tree.FlatNode{Value: n.value}
}

type eBinaryNode struct {
	op    tree.BinaryOp
	left  editableNode
	right editableNode
}

func (n *eBinaryNode) toTree() tree.FormulaNode {
	return &tree.BinaryNode{Op: n.op, Left: n.left.toTree(), Right: n.right.toTree()}
}

type eUnaryNode struct {
	op    tree.UnaryOp
	child editableNode
}

func (n *eUnaryNode) toTree() tree.FormulaNode {
	return &tree.UnaryNode{Op: n.op, Child: n.child.toTree()}
}

type eModifierNode struct {
	name string
	args factory.ModifierArgs
	child editableNode
}

func (n *eModifierNode) toTree() tree.FormulaNode {
	mod, err := factory.GlobalRegistry.Build(n.name, n.args)
	if err != nil {
		return n.child.toTree()
	}
	return &tree.ModifierNode{
		Modifier: mod,
		Child:    n.child.toTree(),
		Name:     n.name,
		Args:     n.args,
	}
}

func toEditable(n tree.FormulaNode) editableNode {
	switch node := n.(type) {
	case *tree.DiceNode:
		return &eDiceNode{count: node.Count, sides: node.Sides}
	case *tree.FlatNode:
		return &eFlatNode{value: node.Value}
	case *tree.BinaryNode:
		return &eBinaryNode{
			op:    node.Op,
			left:  toEditable(node.Left),
			right: toEditable(node.Right),
		}
	case *tree.UnaryNode:
		return &eUnaryNode{op: node.Op, child: toEditable(node.Child)}
	case *tree.ModifierNode:
		return &eModifierNode{
			name:  node.Name,
			args:  node.Args,
			child: toEditable(node.Child),
		}
	default:
		return &eFlatNode{value: 0}
	}
}

func currentType(n editableNode) string {
	switch n.(type) {
	case *eDiceNode:
		return "Dice"
	case *eFlatNode:
		return "Flat"
	case *eBinaryNode:
		return "Binary"
	case *eUnaryNode:
		return "Unary"
	case *eModifierNode:
		return "Modifier"
	default:
		return "Flat"
	}
}

func createNodeOfType(t string, old editableNode) editableNode {
	switch t {
	case "Dice":
		return &eDiceNode{count: 1, sides: 6}
	case "Flat":
		return &eFlatNode{value: 0}
	case "Binary":
		var left editableNode = &eFlatNode{value: 0}
		var right editableNode = &eFlatNode{value: 0}
		if u, ok := old.(*eUnaryNode); ok {
			left = u.child
		} else if b, ok := old.(*eBinaryNode); ok {
			left = b.left
			right = b.right
		} else if m, ok := old.(*eModifierNode); ok {
			left = m.child
		}
		return &eBinaryNode{op: tree.BinaryPlus, left: left, right: right}
	case "Unary":
		var child editableNode = &eFlatNode{value: 0}
		if u, ok := old.(*eUnaryNode); ok {
			child = u.child
		} else if b, ok := old.(*eBinaryNode); ok {
			child = b.left
		} else if m, ok := old.(*eModifierNode); ok {
			child = m.child
		}
		return &eUnaryNode{op: tree.UnaryPlus, child: child}
	case "Modifier":
		var child editableNode = &eDiceNode{count: 1, sides: 6}
		if u, ok := old.(*eUnaryNode); ok {
			child = toEditable(u.child.toTree())
		} else if b, ok := old.(*eBinaryNode); ok {
			child = toEditable(b.left.toTree())
		} else if m, ok := old.(*eModifierNode); ok {
			child = m.child
		} else if d, ok := old.(*eDiceNode); ok {
			child = d
		}
		return &eModifierNode{name: defaultModifierName(), args: factory.ModifierArgs{}, child: child}
	default:
		return &eFlatNode{value: 0}
	}
}

func defaultModifierName() string {
	names := factory.GlobalRegistry.Names()
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

type BuilderPanel struct {
	root     editableNode
	onChange func(tree.FormulaNode)
	content  *fyne.Container
}

func NewBuilderPanel(node tree.FormulaNode, onChange func(tree.FormulaNode)) *BuilderPanel {
	bp := &BuilderPanel{
		root:     toEditable(node),
		onChange: onChange,
		content:  container.NewVBox(),
	}
	bp.refreshContent()
	return bp
}

func (bp *BuilderPanel) refreshContent() {
	bp.content.Objects = []fyne.CanvasObject{buildNodeEditor(bp.root, func(n editableNode) {
		bp.root = n
		bp.refreshContent()
		if bp.onChange != nil {
			bp.onChange(bp.root.toTree())
		}
	}, nil)}
	bp.content.Refresh()
}

func (bp *BuilderPanel) CanvasObject() fyne.CanvasObject {
	return bp.content
}

func buildNodeEditor(node editableNode, replace func(editableNode), remove func()) fyne.CanvasObject {
	types := []string{"Dice", "Flat", "Binary", "Unary", "Modifier"}
	typeSelect := widget.NewSelect(types, nil)
	typeSelect.SetSelected(currentType(node))
	typeSelect.OnChanged = func(selected string) {
		if selected == "" {
			return
		}
		newNode := createNodeOfType(selected, node)
		replace(newNode)
	}

	var specific fyne.CanvasObject
	switch n := node.(type) {
	case *eDiceNode:
		specific = buildDiceEditor(n, replace)
	case *eFlatNode:
		specific = buildFlatEditor(n, replace)
	case *eBinaryNode:
		specific = buildBinaryEditor(n, replace)
	case *eUnaryNode:
		specific = buildUnaryEditor(n, replace)
	case *eModifierNode:
		specific = buildModifierEditor(n, replace)
	default:
		specific = widget.NewLabel("Unknown node")
	}

	btns := container.NewHBox(widget.NewLabel("Type:"), typeSelect)
	if remove != nil {
		delBtn := widget.NewButton("Delete", func() {
			remove()
		})
		btns.Add(delBtn)
	}

	return widget.NewCard("Node Editor", "", container.NewVBox(btns, specific))
}

func buildDiceEditor(n *eDiceNode, replace func(editableNode)) fyne.CanvasObject {
	countEntry := widget.NewEntry()
	countEntry.SetText(strconv.Itoa(n.count))
	sidesEntry := widget.NewEntry()
	sidesEntry.SetText(strconv.Itoa(n.sides))

	onChange := func() {
		c, _ := strconv.Atoi(countEntry.Text)
		s, _ := strconv.Atoi(sidesEntry.Text)
		if c < 1 {
			c = 1
		}
		if s < 1 {
			s = 1
		}
		n.count = c
		n.sides = s
		replace(n)
	}

	countEntry.OnChanged = func(_ string) { onChange() }
	sidesEntry.OnChanged = func(_ string) { onChange() }

	return container.NewHBox(
		widget.NewLabel("Count"),
		countEntry,
		widget.NewLabel("Sides"),
		sidesEntry,
	)
}

func buildFlatEditor(n *eFlatNode, replace func(editableNode)) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetText(strconv.Itoa(n.value))
	entry.OnChanged = func(s string) {
		v, _ := strconv.Atoi(s)
		n.value = v
		replace(n)
	}
	return container.NewHBox(widget.NewLabel("Value"), entry)
}

func buildBinaryEditor(n *eBinaryNode, replace func(editableNode)) fyne.CanvasObject {
	ops := []string{"+", "-", "*", "/"}
	opSelect := widget.NewSelect(ops, nil)
	switch n.op {
	case tree.BinaryPlus:
		opSelect.SetSelected("+")
	case tree.BinaryMinus:
		opSelect.SetSelected("-")
	case tree.BinaryMul:
		opSelect.SetSelected("*")
	case tree.BinaryDiv:
		opSelect.SetSelected("/")
	}
	opSelect.OnChanged = func(s string) {
		switch s {
		case "+":
			n.op = tree.BinaryPlus
		case "-":
			n.op = tree.BinaryMinus
		case "*":
			n.op = tree.BinaryMul
		case "/":
			n.op = tree.BinaryDiv
		}
		replace(n)
	}

	leftEditor := buildNodeEditor(n.left, func(newLeft editableNode) {
		n.left = newLeft
		replace(n)
	}, func() {
		n.left = &eFlatNode{value: 0}
		replace(n)
	})

	rightEditor := buildNodeEditor(n.right, func(newRight editableNode) {
		n.right = newRight
		replace(n)
	}, func() {
		n.right = &eFlatNode{value: 0}
		replace(n)
	})

	return container.NewVBox(
		container.NewHBox(widget.NewLabel("Op"), opSelect),
		widget.NewLabel("Left"),
		leftEditor,
		widget.NewLabel("Right"),
		rightEditor,
	)
}

func buildUnaryEditor(n *eUnaryNode, replace func(editableNode)) fyne.CanvasObject {
	ops := []string{"+", "-"}
	opSelect := widget.NewSelect(ops, nil)
	if n.op == tree.UnaryMinus {
		opSelect.SetSelected("-")
	} else {
		opSelect.SetSelected("+")
	}
	opSelect.OnChanged = func(s string) {
		if s == "-" {
			n.op = tree.UnaryMinus
		} else {
			n.op = tree.UnaryPlus
		}
		replace(n)
	}

	childEditor := buildNodeEditor(n.child, func(newChild editableNode) {
		n.child = newChild
		replace(n)
	}, func() {
		n.child = &eFlatNode{value: 0}
		replace(n)
	})

	return container.NewVBox(
		container.NewHBox(widget.NewLabel("Op"), opSelect),
		widget.NewLabel("Child"),
		childEditor,
	)
}

func buildModifierEditor(n *eModifierNode, replace func(editableNode)) fyne.CanvasObject {
	names := factory.GlobalRegistry.Names()
	nameSelect := widget.NewSelect(names, nil)
	if n.name != "" {
		nameSelect.SetSelected(n.name)
	} else if len(names) > 0 {
		nameSelect.SetSelected(names[0])
		n.name = names[0]
	}

	argsBox := container.NewVBox()

	var refreshArgs func()
	refreshArgs = func() {
		argsBox.Objects = nil
		schema, ok := factory.GlobalRegistry.Schema(n.name)
		if !ok {
			argsBox.Add(widget.NewLabel("No schema"))
			argsBox.Refresh()
			return
		}

		currentNamed := n.args.Named
		if currentNamed == nil {
			currentNamed = make(map[string]any)
		}
		currentPos := n.args.Positional

		newArgs := factory.ModifierArgs{
			Named: make(map[string]any),
		}

		for i, arg := range schema.Args {
			labelText := arg.Name
			if arg.Required {
				labelText += " *"
			}

			var initial string
			if i < len(currentPos) {
				initial = fmt.Sprintf("%v", currentPos[i])
			} else if v, ok := currentNamed[arg.Name]; ok {
				initial = fmt.Sprintf("%v", v)
			}

			switch arg.Type {
			case factory.ArgEnum:
				sel := widget.NewSelect(arg.Options, nil)
				sel.SetSelected(initial)
				if sel.Selected == "" && len(arg.Options) > 0 {
					sel.SetSelected(arg.Options[0])
				}
				sel.OnChanged = func(name string, s *widget.Select) func(string) {
					return func(val string) {
						newArgs.Named[name] = val
						if name == schema.Args[0].Name && schema.Args[0].Type == factory.ArgEnum {
							newArgs.Positional = []any{val}
						}
						n.args = newArgs
						replace(n)
					}
				}(arg.Name, sel)
				newArgs.Named[arg.Name] = sel.Selected
				if arg.Name == schema.Args[0].Name && schema.Args[0].Type == factory.ArgEnum {
					newArgs.Positional = []any{sel.Selected}
				}
				argsBox.Add(container.NewHBox(widget.NewLabel(labelText), sel))

			case factory.ArgInt:
				entry := widget.NewEntry()
				entry.SetText(initial)
				entry.OnChanged = func(name string) func(string) {
					return func(s string) {
						v, _ := strconv.Atoi(s)
						newArgs.Named[name] = v
						n.args = newArgs
						replace(n)
					}
				}(arg.Name)
				if initial != "" {
					v, _ := strconv.Atoi(initial)
					newArgs.Named[arg.Name] = v
				}
				argsBox.Add(container.NewHBox(widget.NewLabel(labelText), entry))

			case factory.ArgIntSlice:
				entry := widget.NewEntry()
				entry.SetPlaceHolder("1,2,3")
				entry.SetText(initial)
				entry.OnChanged = func(name string) func(string) {
					return func(s string) {
						parts := strings.Split(s, ",")
						var vals []any
						for _, p := range parts {
							p = strings.TrimSpace(p)
							if p == "" {
								continue
							}
							v, _ := strconv.Atoi(p)
							vals = append(vals, v)
						}
						newArgs.Named[name] = vals
						n.args = newArgs
						replace(n)
					}
				}(arg.Name)
				if initial != "" {
					parts := strings.Split(initial, ",")
					var vals []any
					for _, p := range parts {
						p = strings.TrimSpace(p)
						if p == "" {
							continue
						}
						v, _ := strconv.Atoi(p)
						vals = append(vals, v)
					}
					newArgs.Named[arg.Name] = vals
				}
				argsBox.Add(container.NewHBox(widget.NewLabel(labelText), entry))

			case factory.ArgString:
				entry := widget.NewEntry()
				entry.SetText(initial)
				entry.OnChanged = func(name string) func(string) {
					return func(s string) {
						newArgs.Named[name] = s
						n.args = newArgs
						replace(n)
					}
				}(arg.Name)
				if initial != "" {
					newArgs.Named[arg.Name] = initial
				}
				argsBox.Add(container.NewHBox(widget.NewLabel(labelText), entry))
			}
		}

		n.args = newArgs
		argsBox.Refresh()
	}

	nameSelect.OnChanged = func(s string) {
		n.name = s
		n.args = factory.ModifierArgs{}
		refreshArgs()
		replace(n)
	}

	refreshArgs()

	childEditor := buildNodeEditor(n.child, func(newChild editableNode) {
		n.child = newChild
		replace(n)
	}, func() {
		n.child = &eFlatNode{value: 0}
		replace(n)
	})

	return container.NewVBox(
		container.NewHBox(widget.NewLabel("Modifier"), nameSelect),
		argsBox,
		widget.NewLabel("Child"),
		childEditor,
	)
}
