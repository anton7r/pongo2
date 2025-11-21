package pongo2

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/CloudyKit/fastprinter"
)

// bufferPool provides a pool of bytes.Buffer instances to reduce allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

// getBuffer retrieves a buffer from the pool
func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// putBuffer returns a buffer to the pool after resetting it
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

// templateWriterPool provides a pool of templateWriter instances to reduce allocations
var templateWriterPool = sync.Pool{
	New: func() interface{} {
		return &templateWriter{}
	},
}

// getTemplateWriter retrieves a templateWriter from the pool and sets its writer
func getTemplateWriter(w io.Writer) *templateWriter {
	tw := templateWriterPool.Get().(*templateWriter)
	tw.w = w
	return tw
}

// putTemplateWriter returns a templateWriter to the pool after clearing its writer
func putTemplateWriter(tw *templateWriter) {
	tw.w = nil
	templateWriterPool.Put(tw)
}

// bufferedTemplateWriterPool provides a pool of buffered templateWriter instances
var bufferedTemplateWriterPool = sync.Pool{
	New: func() interface{} {
		return &bufferedTemplateWriter{
			buf: bytes.NewBuffer(make([]byte, 0, 1024)),
			tw:  &templateWriter{},
		}
	},
}

// bufferedTemplateWriter wraps a templateWriter with its own buffer
type bufferedTemplateWriter struct {
	buf *bytes.Buffer
	tw  *templateWriter
}

// getBufferedTemplateWriter retrieves a buffered templateWriter from the pool
func getBufferedTemplateWriter() *bufferedTemplateWriter {
	btw := bufferedTemplateWriterPool.Get().(*bufferedTemplateWriter)
	btw.tw.w = btw.buf
	return btw
}

// putBufferedTemplateWriter returns a buffered templateWriter to the pool after resetting
func putBufferedTemplateWriter(btw *bufferedTemplateWriter) {
	btw.buf.Reset()
	btw.tw.w = nil
	bufferedTemplateWriterPool.Put(btw)
}

type TemplateWriter interface {
	io.Writer
	WriteString(string) (int, error)
	WriteAny(*Value) (int, error)
}

type templateWriter struct {
	w io.Writer
}

func (tw *templateWriter) WriteString(s string) (int, error) {
	return fastprinter.PrintString(tw.w, s)
}

func (tw *templateWriter) Write(b []byte) (int, error) {
	return tw.w.Write(b)
}

func (tw *templateWriter) WriteAny(v *Value) (int, error) {
	if v.IsNil() {
		return 0, nil
	}

	// Use fastprinter's optimized functions for basic types
	// Use type switch on resolved value to avoid reflection
	switch val := v.getResolvedValue().(type) {
	case string:
		return fastprinter.PrintString(tw.w, val)
	case int:
		return fastprinter.PrintInt(tw.w, int64(val))
	case int8:
		return fastprinter.PrintInt(tw.w, int64(val))
	case int16:
		return fastprinter.PrintInt(tw.w, int64(val))
	case int32:
		return fastprinter.PrintInt(tw.w, int64(val))
	case int64:
		return fastprinter.PrintInt(tw.w, val)
	case uint:
		return fastprinter.PrintUint(tw.w, uint64(val))
	case uint8:
		return fastprinter.PrintUint(tw.w, uint64(val))
	case uint16:
		return fastprinter.PrintUint(tw.w, uint64(val))
	case uint32:
		return fastprinter.PrintUint(tw.w, uint64(val))
	case uint64:
		return fastprinter.PrintUint(tw.w, val)
	case float32:
		// Use 6 decimal places to match fmt.Sprintf("%f", ...) behavior
		return fastprinter.PrintFloatPrecision(tw.w, float64(val), 6)
	case float64:
		return fastprinter.PrintFloatPrecision(tw.w, val, 6)
	case bool:
		// Match pongo2's capitalized boolean format (True/False)
		if val {
			return fastprinter.PrintString(tw.w, "True")
		}
		return fastprinter.PrintString(tw.w, "False")
	case fmt.Stringer:
		return fastprinter.PrintString(tw.w, val.String())
	default:
		// Fall back to String() method for unsupported types
		return fastprinter.PrintString(tw.w, v.String())
	}
}

type Template struct {
	set *TemplateSet

	// Input
	isTplString bool
	name        string
	tpl         string
	size        int

	// Calculation
	tokens []*Token
	parser *Parser

	// first come, first serve (it's important to not override existing entries in here)
	level          int
	parent         *Template
	child          *Template
	blocks         map[string]*NodeWrapper
	exportedMacros map[string]*tagMacroNode

	// Output
	root *nodeDocument

	// Options allow you to change the behavior of template-engine.
	// You can change the options before calling the Execute method.
	Options *Options
}

func newTemplateString(set *TemplateSet, tpl []byte) (*Template, error) {
	return newTemplate(set, "<string>", true, tpl)
}

func newTemplate(set *TemplateSet, name string, isTplString bool, tpl []byte) (*Template, error) {
	strTpl := string(tpl)

	// Create the template
	t := &Template{
		set:            set,
		isTplString:    isTplString,
		name:           name,
		tpl:            strTpl,
		size:           len(strTpl),
		blocks:         make(map[string]*NodeWrapper),
		exportedMacros: make(map[string]*tagMacroNode),
		Options:        newOptions(),
	}
	// Copy all settings from another Options.
	t.Options.Update(set.Options)

	// Tokenize it
	tokens, err := lex(name, strTpl)
	if err != nil {
		return nil, err
	}
	t.tokens = tokens

	// For debugging purposes, show all tokens:
	/*for i, t := range tokens {
		fmt.Printf("%3d. %s\n", i, t)
	}*/

	// Parse it
	err = t.parse()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tpl *Template) newContextForExecution(context Context) (*Template, *ExecutionContext, error) {
	if tpl.Options.TrimBlocks || tpl.Options.LStripBlocks {
		// Issue #94 https://github.com/flosch/pongo2/issues/94
		// If an application configures pongo2 template to trim_blocks,
		// the first newline after a template tag is removed automatically (like in PHP).
		prev := &Token{
			Typ: TokenHTML,
			Val: "\n",
		}

		for _, t := range tpl.tokens {
			if tpl.Options.LStripBlocks {
				if prev.Typ == TokenHTML && t.Typ != TokenHTML && t.Val == "{%" {
					prev.Val = strings.TrimRight(prev.Val, "\t ")
				}
			}

			if tpl.Options.TrimBlocks {
				if prev.Typ != TokenHTML && t.Typ == TokenHTML && prev.Val == "%}" {
					if len(t.Val) > 0 && t.Val[0] == '\n' {
						t.Val = t.Val[1:len(t.Val)]
					}
				}
			}

			prev = t
		}
	}

	// Determine the parent to be executed (for template inheritance)
	parent := tpl
	for parent.parent != nil {
		parent = parent.parent
	}

	// Create context if none is given
	newContext := make(Context)
	newContext.Update(tpl.set.Globals)

	if context != nil {
		newContext.Update(context)

		if len(newContext) > 0 {
			// Check for context name syntax
			err := newContext.checkForValidIdentifiers()
			if err != nil {
				return parent, nil, err
			}

			// Check for clashes with macro names
			for k := range newContext {
				_, has := tpl.exportedMacros[k]
				if has {
					return parent, nil, &Error{
						Filename:  tpl.name,
						Sender:    "execution",
						OrigError: fmt.Errorf("context key name '%s' clashes with macro '%s'", k, k),
					}
				}
			}
		}
	}

	// Create operational context
	ctx := newExecutionContext(parent, newContext)

	return parent, ctx, nil
}

func (tpl *Template) execute(context Context, writer TemplateWriter) error {
	parent, ctx, err := tpl.newContextForExecution(context)
	if err != nil {
		return err
	}

	// Run the selected document
	if err := parent.root.Execute(ctx, writer); err != nil {
		return err
	}

	return nil
}

func (tpl *Template) newTemplateWriterAndExecute(context Context, writer io.Writer) error {
	tw := getTemplateWriter(writer)
	defer putTemplateWriter(tw)
	return tpl.execute(context, tw)
}

func (tpl *Template) newBufferAndExecute(context Context) (*bytes.Buffer, error) {
	// Get buffered template writer from pool
	btw := getBufferedTemplateWriter()
	defer putBufferedTemplateWriter(btw)
	if err := tpl.execute(context, btw.tw); err != nil {
		return nil, err
	}
	// Return a copy of the buffer contents since we're returning it to the pool
	result := getBuffer()
	result.Write(btw.buf.Bytes())
	return result, nil
}

// Executes the template with the given context and writes to writer (io.Writer)
// on success. Context can be nil. Nothing is written on error; instead the error
// is being returned.
func (tpl *Template) ExecuteWriter(context Context, writer io.Writer) error {
	buf, err := tpl.newBufferAndExecute(context)
	if err != nil {
		return err
	}
	defer putBuffer(buf)
	_, err = buf.WriteTo(writer)
	if err != nil {
		return err
	}
	return nil
}

// Same as ExecuteWriter. The only difference between both functions is that
// this function might already have written parts of the generated template in the
// case of an execution error because there's no intermediate buffer involved for
// performance reasons. This is handy if you need high performance template
// generation or if you want to manage your own pool of buffers.
func (tpl *Template) ExecuteWriterUnbuffered(context Context, writer io.Writer) error {
	return tpl.newTemplateWriterAndExecute(context, writer)
}

// Executes the template and returns the rendered template as a []byte
func (tpl *Template) ExecuteBytes(context Context) ([]byte, error) {
	// Execute template
	buffer, err := tpl.newBufferAndExecute(context)
	if err != nil {
		return nil, err
	}
	defer putBuffer(buffer)
	// Make a copy since we're returning the buffer to the pool
	result := make([]byte, buffer.Len())
	copy(result, buffer.Bytes())
	return result, nil
}

// Executes the template and returns the rendered template as a string
func (tpl *Template) Execute(context Context) (string, error) {
	// Execute template
	buffer, err := tpl.newBufferAndExecute(context)
	if err != nil {
		return "", err
	}
	defer putBuffer(buffer)
	return buffer.String(), nil
}

func (tpl *Template) ExecuteBlocks(context Context, blocks []string) (map[string]string, error) {
	var parents []*Template
	result := make(map[string]string)

	parent := tpl
	for parent != nil {
		// We only want to execute the template if it has a block we want
		for _, block := range blocks {
			if _, ok := tpl.blocks[block]; ok {
				parents = append(parents, parent)
				break
			}
		}
		parent = parent.parent
	}

	for _, t := range parents {
		var btw *bufferedTemplateWriter
		var ctx *ExecutionContext
		var err error
		for _, blockName := range blocks {
			if _, ok := result[blockName]; ok {
				continue
			}
			if blockWrapper, ok := t.blocks[blockName]; ok {
				// assign the buffered writer if we haven't done so
				if btw == nil {
					btw = getBufferedTemplateWriter()
				}
				// assign the context if we haven't done so
				if ctx == nil {
					_, ctx, err = t.newContextForExecution(context)
					if err != nil {
						return nil, err
					}
				}
				bErr := blockWrapper.Execute(ctx, btw.tw)
				if bErr != nil {
					return nil, bErr
				}
				result[blockName] = btw.buf.String()
				btw.buf.Reset()
			}
		}
		// Return buffered writer to pool if we used one
		if btw != nil {
			putBufferedTemplateWriter(btw)
		}
		// We have found all blocks
		if len(blocks) == len(result) {
			break
		}
	}

	return result, nil
}
