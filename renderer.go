package markdownxml

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
	"html/template"
)

func NewRenderer() blackfriday.Renderer {
	return &renderer{}
}

type renderer struct {
	indent int
}

func (r *renderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Document:
		if entering {
			r.writeWithIndent(w, `<document xmlns="http://commonmark.org/xml/1.0">`)
			r.indent++
		} else {
			r.indent--
			r.writeWithIndent(w, `</document>`)
		}
	case blackfriday.BlockQuote:
		r.tag(w, "block_quote", entering)
	case blackfriday.List:
		var attrs map[string]interface{}
		if node.ListFlags&blackfriday.ListTypeOrdered != 0 {
			attrs = map[string]interface{}{
				"type":  "ordered",
				"tight": node.ListData.Tight,
			}
			switch node.ListData.Delimiter {
			case ')':
				attrs["delimiter"] = "paren"
			case '.':
				fallthrough
			default:
				attrs["delimiter"] = "period"
			}
		} else {
			attrs = map[string]interface{}{
				"type":  "bullet",
				"tight": node.ListData.Tight,
			}
		}
		r.tag(w, "list", entering, attrs)
	case blackfriday.Item:
		r.tag(w, "item", entering, map[string]interface{}{
			"flags": node.ListFlags,
		})
	case blackfriday.Paragraph:
		r.tag(w, "paragraph", entering)
	case blackfriday.Heading:
		r.tag(w, "heading", entering, map[string]interface{}{
			"level": node.HeadingData.Level,
		})
	case blackfriday.HorizontalRule:
	case blackfriday.Emph:
		r.tag(w, "emph", entering)
	case blackfriday.Strong:
		r.tag(w, "strong", entering)
	case blackfriday.Del:
		r.tag(w, "strikethrough", entering)
	case blackfriday.Link:
		r.tag(w, "link", entering, map[string]interface{}{
			"destination": string(node.LinkData.Destination),
			"title": string(node.LinkData.Title),
		})
	case blackfriday.Image:
		r.tag(w, "image", entering, map[string]interface{}{
			"destination": string(node.LinkData.Destination),
			"title": string(node.LinkData.Title),
		})
	case blackfriday.Text:
		r.writeWithIndent(w, `<text>%s</text>`, string(node.Literal))
	case blackfriday.HTMLBlock:
		r.writeWithIndent(w, "<html_block>%s</html_block>", template.HTMLEscapeString(string(node.Literal)))
	case blackfriday.CodeBlock:
		if len(node.CodeBlockData.Info) > 0 {
			r.writeWithIndent(w, "<code_block info=\"%s\">", string(node.CodeBlockData.Info))
		} else {
			r.writeWithIndent(w, "<code_block>")
		}
		w.Write(node.Literal)
		io.WriteString(w, "</code_block>\n")
	case blackfriday.Softbreak:
		r.writeWithIndent(w, "<softbreak />")
	case blackfriday.Hardbreak:
		r.writeWithIndent(w, "<linebreak />")
	case blackfriday.Code:
		r.writeWithIndent(w, "<code>%s</code>", string(node.Literal))
	case blackfriday.HTMLSpan:
	case blackfriday.Table:
		r.tag(w, "table", entering)
	case blackfriday.TableCell:
		r.tag(w, "td", entering)
	case blackfriday.TableHead:
		r.tag(w, "thead", entering)
	case blackfriday.TableBody:
		r.tag(w, "tbody", entering)
	case blackfriday.TableRow:
		r.tag(w, "tr", entering)
	}
	return blackfriday.GoToNext
}

func (r *renderer) tag(w io.Writer, name string, entering bool, attrs ...map[string]interface{}) {
	if entering {
		var tags []string
		for _, m := range attrs {
			for k, v := range m {
				tags = append(tags, fmt.Sprintf("%s=%#v", k, fmt.Sprint(v)))
			}
		}

		if len(tags) > 0 {
			r.writeWithIndent(w, "<%s %s>", name, strings.Join(tags, " "))
		} else {
			r.writeWithIndent(w, "<%s>", name)
		}
		r.indent++
	} else {
		r.indent--
		r.writeWithIndent(w, "</%s>", name)
	}
}

func (r *renderer) writeWithIndent(w io.Writer, fmtStr string, values ...interface{}) {
	if r.indent > 0 {
		w.Write(bytes.Repeat([]byte{'\t'}, r.indent))
	}
	fmt.Fprintf(w, fmtStr, values...)
	io.WriteString(w, "\n")
}

func (*renderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE document SYSTEM "CommonMark.dtd">
`)
}

func (*renderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}
