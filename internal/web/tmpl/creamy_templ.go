// Code generated by templ@v0.2.304 DO NOT EDIT.

package tmpl

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import (
	"fmt"
	"go.albinodrought.com/creamy-board/internal/repo"
)

// formatters:

func humanizeBytes(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func nameFile(thread *repo.Thread, post *repo.Post, file *repo.File) string {
	return fmt.Sprintf(
		"%v-%v-%v.%v",
		thread.ID,
		post.ID,
		file.Index,
		file.Extension,
	)
}

func linkBoard(board *repo.Board) templ.SafeURL {
	return templ.SafeURL(fmt.Sprintf("/%v/", board.Slug))
}
func linkFile(board *repo.Board, thread *repo.Thread, post *repo.Post, file *repo.File) templ.SafeURL {
	return templ.SafeURL(fmt.Sprintf(
		"/%v/src/%v",
		board.Slug,
		nameFile(thread, post, file),
	))
}
func linkThreadShow(board *repo.Board, thread *repo.Thread) templ.SafeURL {
	return templ.SafeURL(fmt.Sprintf("/%v/res/%v.html", board.Slug, thread.ID))
}
func linkPostShow(board *repo.Board, thread *repo.Thread, post *repo.Post) templ.SafeURL {
	return templ.SafeURL(fmt.Sprintf("/%v/res/%v.html#%v", board.Slug, thread.ID, post.ID))
}
func linkPostQuote(board *repo.Board, thread *repo.Thread, post *repo.Post) templ.SafeURL {
	return templ.SafeURL(fmt.Sprintf("/%v/res/%v.html#q%v", board.Slug, thread.ID, post.ID))
}

func textBoard(board *repo.Board) string {
	return fmt.Sprintf("/%v/ - %v", board.Slug, board.Title)
}
func textThread(board *repo.Board, thread *repo.Thread) string {
	return fmt.Sprintf("/%v/ - %v", board.Slug, thread.Subject)
}

// components:

func postHead(board *repo.Board, thread *repo.Thread, post *repo.Post) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<span")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__author\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_2 string = post.Author
		_, err = templBuffer.WriteString(templ.EscapeString(var_2))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<wbr>")
		if err != nil {
			return err
		}
		// Text
		var_3 := `&nbsp;`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<span")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__date\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_4 string = post.CreatedAt
		_, err = templBuffer.WriteString(templ.EscapeString(var_4))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<wbr>")
		if err != nil {
			return err
		}
		// Text
		var_5 := `&nbsp;`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__link-show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_6 templ.SafeURL = linkPostShow(board, thread, post)
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_6)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_7 := `No.`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		// Text
		var_8 := `&nbsp;`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__link-quote\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_9 templ.SafeURL = linkPostQuote(board, thread, post)
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_9)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_10 string = fmt.Sprintf("%v", post.ID)
		_, err = templBuffer.WriteString(templ.EscapeString(var_10))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func threadHead(board *repo.Board, thread *repo.Thread, post *repo.Post, inThread bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_11 := templ.GetChildren(ctx)
		if var_11 == nil {
			var_11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<span")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"thread__subject\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_12 string = thread.Subject
		_, err = templBuffer.WriteString(templ.EscapeString(var_12))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<wbr>")
		if err != nil {
			return err
		}
		// Text
		var_13 := `&nbsp;`
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		// TemplElement
		err = postHead(board, thread, post).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		// Text
		var_14 := `&nbsp;`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		// If
		if !inThread {
			// Element (standard)
			_, err = templBuffer.WriteString("<a")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"thread__link-reply\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" href=")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			var var_15 templ.SafeURL = linkThreadShow(board, thread)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_15)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Text
			var_16 := `[Reply]`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a>")
			if err != nil {
				return err
			}
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func fileFull(board *repo.Board, thread *repo.Thread, post *repo.Post, file *repo.File) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_17 := templ.GetChildren(ctx)
		if var_17 == nil {
			var_17 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"file\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"file__info\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_18 := `File:`
		_, err = templBuffer.WriteString(var_18)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_19 templ.SafeURL = linkFile(board, thread, post, file)
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_19)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" target=\"_blank\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_20 string = nameFile(thread, post, file)
		_, err = templBuffer.WriteString(templ.EscapeString(var_20))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		// Text
		var_21 := `&nbsp;(`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<span>")
		if err != nil {
			return err
		}
		// StringExpression
		var var_22 string = humanizeBytes(file.Bytes)
		_, err = templBuffer.WriteString(templ.EscapeString(var_22))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<wbr>")
		if err != nil {
			return err
		}
		// Text
		var_23 := `,`
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_24 templ.SafeURL = linkFile(board, thread, post, file)
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_24)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" download=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(file.OriginalName))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" title=\"Save as original filename\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_25 string = file.OriginalName
		_, err = templBuffer.WriteString(templ.EscapeString(var_25))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<wbr>")
		if err != nil {
			return err
		}
		// Text
		var_26 := `)`
		_, err = templBuffer.WriteString(var_26)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"file__img-container\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_27 templ.SafeURL = linkFile(board, thread, post, file)
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_27)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" target=\"_blank\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<img")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"file__img\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" src=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(string(linkFile(board, thread, post, file))))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

// GoExpression
func postFilesClass(files []repo.File) string {
	if len(files) <= 1 {
		return "post__files"
	}

	return "post__files post__files--multi"
}

func postFilesFull(board *repo.Board, thread *repo.Thread, post *repo.Post, files []repo.File) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_28 := templ.GetChildren(ctx)
		if var_28 == nil {
			var_28 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		// Element CSS
		var var_29 = []any{postFilesClass(files)}
		err = templ.RenderCSSItems(ctx, templBuffer, var_29...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_29).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// For
		for _, file := range post.Files {
			// TemplElement
			err = fileFull(board, thread, post, &file).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func postFull(board *repo.Board, thread *repo.Thread, post *repo.Post) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_30 := templ.GetChildren(ctx)
		if var_30 == nil {
			var_30 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__content\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__head\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// TemplElement
		err = postHead(board, thread, post).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// TemplElement
		err = postFilesFull(board, thread, post, post.Files).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__body\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_31 string = post.Body
		_, err = templBuffer.WriteString(templ.EscapeString(var_31))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func threadFull(board *repo.Board, thread *repo.Thread, main *repo.Post, other []repo.Post, inThread bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_32 := templ.GetChildren(ctx)
		if var_32 == nil {
			var_32 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"thread\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post post--op\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__content\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__head\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// TemplElement
		err = threadHead(board, thread, main, inThread).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// TemplElement
		err = postFilesFull(board, thread, main, main.Files).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"post__body\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_33 string = main.Body
		_, err = templBuffer.WriteString(templ.EscapeString(var_33))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"thread__posts\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// For
		for _, post := range other {
			// TemplElement
			err = postFull(board, thread, &post).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

// GoExpression
// wrappers:

func page(title string, description string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_34 := templ.GetChildren(ctx)
		if var_34 == nil {
			var_34 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// DocType
		_, err = templBuffer.WriteString(`<!doctype html>`)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<html")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" lang=\"en\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<head>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<meta")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" charset=\"utf-8\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<meta")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" http-equiv=\"X-UA-Compatible\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" content=\"IE=edge\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<meta")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" name=\"viewport\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" content=\"width=device-width,initial-scale=1.0\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<meta")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" name=\"description\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" content=\"{ description }\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" rel=\"icon\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"/favicon.ico\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<title>")
		if err != nil {
			return err
		}
		// StringExpression
		var var_35 string = title
		_, err = templBuffer.WriteString(templ.EscapeString(var_35))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<link")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=\"/css/main.css\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" rel=\"stylesheet\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</head>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<body>")
		if err != nil {
			return err
		}
		// Children
		err = var_34.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</body>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func wrapperBoard(board *repo.Board) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_36 := templ.GetChildren(ctx)
		if var_36 == nil {
			var_36 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<h1")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"board__title\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_37 string = textBoard(board)
		_, err = templBuffer.WriteString(templ.EscapeString(var_37))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h2")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"board__tag\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_38 string = board.Tagline
		_, err = templBuffer.WriteString(templ.EscapeString(var_38))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<hr>")
		if err != nil {
			return err
		}
		// Children
		err = var_36.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

// GoExpression
// entrypoints:

func ListBoards(boards []repo.Board) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_39 := templ.GetChildren(ctx)
		if var_39 == nil {
			var_39 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// TemplElement
		var_40 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<ul>")
			if err != nil {
				return err
			}
			// For
			for _, board := range boards {
				// Element (standard)
				_, err = templBuffer.WriteString("<a")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" href=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				var var_41 templ.SafeURL = linkBoard(&board)
				_, err = templBuffer.WriteString(templ.EscapeString(string(var_41)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_42 string = textBoard(&board)
				_, err = templBuffer.WriteString(templ.EscapeString(var_42))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</ul>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = page("Creamy Board", "the world is your oyster").Render(templ.WithChildren(ctx, var_40), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func ShowBoardAndRecents(brt *repo.BoardRecentThreads) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_43 := templ.GetChildren(ctx)
		if var_43 == nil {
			var_43 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// TemplElement
		var_44 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			// TemplElement
			var_45 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
				templBuffer, templIsBuffer := w.(*bytes.Buffer)
				if !templIsBuffer {
					templBuffer = templ.GetBuffer()
					defer templ.ReleaseBuffer(templBuffer)
				}
				// For
				for _, thread := range brt.RecentThreads {
					// TemplElement
					err = threadFull(&brt.Board, &thread.Thread, &thread.MainPost, thread.RecentPosts, false).Render(ctx, templBuffer)
					if err != nil {
						return err
					}
				}
				if !templIsBuffer {
					_, err = io.Copy(w, templBuffer)
				}
				return err
			})
			err = wrapperBoard(&brt.Board).Render(templ.WithChildren(ctx, var_45), templBuffer)
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = page(textBoard(&brt.Board), brt.Board.Tagline).Render(templ.WithChildren(ctx, var_44), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func ShowFullThread(bft *repo.BoardFullThread) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_46 := templ.GetChildren(ctx)
		if var_46 == nil {
			var_46 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// TemplElement
		var_47 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			// TemplElement
			var_48 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
				templBuffer, templIsBuffer := w.(*bytes.Buffer)
				if !templIsBuffer {
					templBuffer = templ.GetBuffer()
					defer templ.ReleaseBuffer(templBuffer)
				}
				// TemplElement
				err = threadFull(&bft.Board, &bft.FullThread.Thread, &bft.FullThread.MainPost, bft.FullThread.AllPosts, true).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				if !templIsBuffer {
					_, err = io.Copy(w, templBuffer)
				}
				return err
			})
			err = wrapperBoard(&bft.Board).Render(templ.WithChildren(ctx, var_48), templBuffer)
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = page(textThread(&bft.Board, &bft.FullThread.Thread), bft.FullThread.MainPost.Body).Render(templ.WithChildren(ctx, var_47), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}
