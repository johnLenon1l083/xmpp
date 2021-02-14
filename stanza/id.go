// Copyright 2020 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package stanza

import (
	"encoding/xml"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/internal/attr"
	"mellium.im/xmpp/internal/ns"
	"mellium.im/xmpp/jid"
)

// Namespaces used by this package, provided as a convenience.
const (
	// The namespace for unique and stable stanza and origin IDs.
	NSSid = "urn:xmpp:sid:0"
)

const idLen = 32

type ID struct {
	XMLName xml.Name `xml:"urn:xmpp:sid:0 stanza-id"`
	ID      string   `xml:"id,attr"`
	By      jid.JID  `xml:"by,attr"`
}

// TokenReader implements xmlstream.Marshaler.
func (id ID) TokenReader() xml.TokenReader {
	return xmlstream.Wrap(nil, xml.StartElement{
		Name: xml.Name{Space: NSSid, Local: "stanza-id"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: id.ID},
			{Name: xml.Name{Local: "by"}, Value: id.By.String()},
		},
	})
}

// WriteXML implements xmlstream.WriterTo.
func (id ID) WriteXML(w xmlstream.TokenWriter) (int, error) {
	return xmlstream.Copy(w, id.TokenReader())
}

type OriginID struct {
	XMLName xml.Name `xml:"urn:xmpp:sid:0 origin-id"`
	ID      string   `xml:"id,attr"`
}

// TokenReader implements xmlstream.Marshaler.
func (id OriginID) TokenReader() xml.TokenReader {
	return xmlstream.Wrap(nil, xml.StartElement{
		Name: xml.Name{Space: NSSid, Local: "origin-id"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: id.ID},
		},
	})
}

// WriteXML implements xmlstream.WriterTo.
func (id OriginID) WriteXML(w xmlstream.TokenWriter) (int, error) {
	return xmlstream.Copy(w, id.TokenReader())
}

func isStanza(name xml.Name) bool {
	return (name.Local == "iq" || name.Local == "message" || name.Local == "presence") &&
		(name.Space == ns.Client || name.Space == ns.Server)
}

// AddID returns an transformer that adds a random stanza ID to any stanzas that
// does not already have one.
func AddID(by jid.JID) xmlstream.Transformer {
	return xmlstream.InsertFunc(func(start xml.StartElement, level uint64, w xmlstream.TokenWriter) error {
		if isStanza(start.Name) && level == 1 {
			_, err := ID{
				ID: attr.RandomLen(idLen),
				By: by,
			}.WriteXML(w)
			return err
		}
		return nil
	})
}

var (
	addOriginID = xmlstream.InsertFunc(func(start xml.StartElement, level uint64, w xmlstream.TokenWriter) error {
		if isStanza(start.Name) && level == 1 {
			_, err := OriginID{
				ID: attr.RandomLen(idLen),
			}.WriteXML(w)
			return err
		}
		return nil
	})
)

// AddOriginID is an xmlstream.Transformer that adds an origin ID to any stanzas
// found in the input stream.
func AddOriginID(r xml.TokenReader) xml.TokenReader {
	return addOriginID(r)
}
