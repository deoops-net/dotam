package main

import (
	"github.com/flosch/pongo2"
	log "github.com/sirupsen/logrus"
)

type Inode struct {
	Wrapper *pongo2.NodeWrapper
	Parser  *pongo2.Parser
}

func (i Inode) Execute(c *pongo2.ExecutionContext, t pongo2.TemplateWriter) *pongo2.Error {
	log.Debug("called Executor")
	i.Wrapper.Execute(c, t)
	return nil
}

// var SkipTag pongo2.TagParser = func () {}
func SkipTag2(doc *pongo2.Parser, start *pongo2.Token, args *pongo2.Parser) (tag pongo2.INodeTag, err *pongo2.Error) {
	wrapper, parser, err := doc.WrapUntilTag("endskip")
	log.Debug("called")
	log.Debug(doc)

	return Inode{Wrapper: wrapper, Parser: parser}, nil
}
