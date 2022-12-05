package annotation

import (
	"go/ast"
)

// SingleFileEntryVisitor 这部分和课堂演示差不多，但是我建议你们自己试着写一些
type SingleFileEntryVisitor struct {
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	if s != nil {
		return s.file.Get()
	}
	return File{}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {
	file, ok := node.(*ast.File)
	if ok {
		s.file = &fileVisitor{
			ans: newAnnotations(file, file.Doc),
		}
		return s.file
	}
	return s
}

type fileVisitor struct {
	ans     Annotations[*ast.File]
	types   []*typeVisitor
	visited bool
}

func (f *fileVisitor) Get() File {
	types := make([]Type, 0, len(f.types))
	for _, typ := range f.types {
		types = append(types, typ.Get())
	}
	return File{
		f.ans,
		types,
	}
}

func (f *fileVisitor) Visit(node ast.Node) ast.Visitor {
	typ, ok := node.(*ast.TypeSpec)
	if ok {
		tv := &typeVisitor {
			Annotations: newAnnotations(typ, typ.Doc),
			fields: make([]Field, 0, 0),
		}
		f.types = append(f.types, tv)
		return tv
	}

	return f
}

type File struct {
	Annotations[*ast.File]
	Types []Type
}

type typeVisitor struct {
	Annotations[*ast.TypeSpec]
	fields []Field
}

func (t *typeVisitor) Get() Type {
	typ := Type{
		Annotations: t.Annotations,
		Fields: t.fields,
	}
	return typ
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fd, ok := node.(*ast.Field)
	if ok {
		t.fields = append(t.fields, Field{newAnnotations(fd, fd.Doc)})
		return nil
	}
	return t
}

type Type struct {
	Annotations[*ast.TypeSpec]
	Fields []Field
}

type Field struct {
	Annotations[*ast.Field]
}
