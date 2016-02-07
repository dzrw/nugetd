package rekt // my entire fucking day

import (
	"bytes"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"strconv"
	"time"
)

// NuGet's Atom structure is, basically, an abomination.
type nugetXmlWriter struct {
	b bytes.Buffer
}

func newWriter() *nugetXmlWriter {
	n := &nugetXmlWriter{b: bytes.Buffer{}}
	n.b.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	return n
}

func (n *nugetXmlWriter) Bytes() []byte {
	return n.b.Bytes()
}

func (n *nugetXmlWriter) el(key string) *nugetXmlWriter {
	s := fmt.Sprintf("\n<%s>", key)
	n.b.WriteString(s)
	return n
}

func (n *nugetXmlWriter) el_no_child_nodes(key string) *nugetXmlWriter {
	s := fmt.Sprintf("\n<%s />", key)
	n.b.WriteString(s)
	return n
}

func (n *nugetXmlWriter) el_with_text(key string, text string) *nugetXmlWriter {
	s := fmt.Sprintf("\n<%s>%s</%s>", key, text, key)
	n.b.WriteString(s)
	return n
}

func (n *nugetXmlWriter) el_with_text_or_empty(key string, text string) *nugetXmlWriter {
	var s string

	if len(text) == 0 {
		s = fmt.Sprintf("\n<%s />", key)
	} else {
		s = fmt.Sprintf("\n<%s>%s</%s>", key, text, key)
	}

	n.b.WriteString(s)
	return n
}

func (n *nugetXmlWriter) el_with_attr(key string) *nugetXmlWriter {
	s := fmt.Sprintf("\n<%s", key)
	n.b.WriteString(s)
	return n
}

func (n *nugetXmlWriter) fin() *nugetXmlWriter {
	n.b.WriteString(">")
	return n
}

func (n *nugetXmlWriter) fin_no_child_nodes() *nugetXmlWriter {
	n.b.WriteString(" />")
	return n
}

func (n *nugetXmlWriter) le(key string) *nugetXmlWriter {
	n.b.WriteString("</")
	n.b.WriteString(key)
	n.b.WriteString(">")
	return n
}

func (n *nugetXmlWriter) attr(key string, value string, omitempty bool) *nugetXmlWriter {
	if len(key) > 0 && (!omitempty || len(value) > 0) {
		s := fmt.Sprintf(" %s=\"%s\"", key, value)
		n.b.WriteString(s)
	}

	return n
}

func (n *nugetXmlWriter) text(s string) *nugetXmlWriter {
	n.b.WriteString(s)
	return n
}

func (w *nugetXmlWriter) atom_author(name, email string) {
	w.el("author")

	w.el_with_text_or_empty("name", name)

	if len(email) > 0 {
		w.el_with_text("email", email)
	}

	w.le("author")
}

func (w *nugetXmlWriter) nuget_link(rel, href, title string) {
	w.el_with_attr("link")
	w.attr("rel", rel, false)
	w.attr("href", href, false)
	w.attr("title", title, false)
	w.fin_no_child_nodes()
}

func (w *nugetXmlWriter) odatafy_string(key string, value string) {
	if len(value) == 0 {
		w.el_with_attr(key).attr("m:null", "true", false).fin_no_child_nodes()
		return
	}

	w.el(key).text(value).le(key)
}

func (w *nugetXmlWriter) odatafy_bool(key string, b bool) {
	w.el_with_attr(key).attr("m:type", "Edm.Boolean", false).fin()
	if b {
		w.text("true")
	} else {
		w.text("false")
	}
	w.le(key)
}

func (w *nugetXmlWriter) odatafy_int32(key string, i int) {
	w.el_with_attr(key).attr("m:type", "Edm.Int32", false).fin()
	w.text(strconv.Itoa(i))
	w.le(key)
}

func (w *nugetXmlWriter) odatafy_int64(key string, i int) {
	w.el_with_attr(key).attr("m:type", "Edm.Int64", false).fin()
	w.text(strconv.Itoa(i))
	w.le(key)
}

func (w *nugetXmlWriter) odatafy_time(key string, t time.Time) {
	w.el_with_attr(key).attr("m:type", "Edm.DateTime", false).fin()
	w.text(t.Format(time.RFC3339))
	w.le(key)
}

func (w *nugetXmlWriter) odatafy_array(key string, a []string) {
	if len(a) == 0 {
		w.el_no_child_nodes(key)
	} else {
		// TODO: Find an example of the "d:Dependencies" element.
		w.el_with_text(key, "FIXME")
	}
}

type Feed struct {
	Id      string
	Title   string
	Updated time.Time
	Entries []*Entry
}

func (f *Feed) ToAtom() ([]byte, error) {
	w := newWriter()

	// <feed xml:base="https://www.nuget.org/api/v2/"
	//       xmlns="http://www.w3.org/2005/Atom"
	//       xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices"
	//       xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">

	w.el_with_attr("feed")
	w.attr("xml:base", "https://www.nuget.org/api/v2/", false)
	w.attr("xmlns", "http://www.w3.org/2005/Atom", false)
	w.attr("xmlns:d", "http://schemas.microsoft.com/ado/2007/08/dataservices", false)
	w.attr("xmlns:m", "http://schemas.microsoft.com/ado/2007/08/dataservices/metadata", false)
	w.fin()

	// <id>https://www.nuget.org/api/v2/FindPackagesById</id>
	w.el_with_text("id", f.Id)

	// <title type="text">FindPackagesById</title>
	w.el_with_attr("title").attr("type", "text", false).fin()
	w.text(f.Title)
	w.le("title")

	// <updated>2016-02-07T01:28:53Z</updated>
	w.el_with_text("updated", f.Updated.String())

	// <link rel="self" title="FindPackagesById" href="FindPackagesById" />
	w.nuget_link("self", f.Title, f.Title)

	// <author><name /></author>
	w.atom_author("", "")

	// all of the entries, if any
	for _, e := range f.Entries {
		w.el("entry")
		e.writeFields(w)
		w.le("entry")
	}

	w.le("feed")

	return w.Bytes(), nil
}

type PackageRegistry struct {
	Packages []*Entry `yaml:"packages"`
}

func Parse(data []byte) (*PackageRegistry, error) {
	var aux struct {
		Packages []*Entry `yaml:"packages"`
	}

	if err := yaml.Unmarshal(data, &aux); err != nil {
		return nil, err
	}

	return &PackageRegistry{aux.Packages}, nil
}

type Entry struct {
	Id                       string `yaml:"id"`
	Title                    string `yaml:"title"`
	Summary                  string `yaml:"summary"`
	Updated                  time.Time
	AuthorName               string `yaml:"authorName"`
	AuthorEmail              string `yaml:"authorEmail"`
	DownloadUrl              string `yaml:"downloadUrl"`
	Version                  string `yaml:"version"`
	NormalizedVersion        string
	Copyright                string `yaml:"copyright"`
	Created                  time.Time
	Dependencies             []string `yaml:"dependencies"`
	Description              string   `yaml:"description"`
	DownloadCount            int
	GalleryDetailsUrl        string `yaml:"galleryDetailsUrl"`
	IconUrl                  string `yaml:"iconUrl"`
	IsLatestVersion          bool   `yaml:"isLatestVersion"`
	IsAbsoluteLatestVersion  bool   `yaml:"isAbsoluteLatestVersion"`
	IsPrerelease             bool   `yaml:"isPrerelease"`
	Language                 string `yaml:"language"`
	Published                time.Time
	PackageHash              string `yaml:"packageHash"`
	PackageHashAlgorithm     string `yaml:"packageHashAlgorithm"`
	PackageSize              int    `yaml:"packageSize"` //int64
	ProjectUrl               string `yaml:"projectUrl"`
	ReportAbuseUrl           string
	ReleaseNotes             string
	RequireLicenseAcceptance bool
	PackageSummary           string `yaml:"packageSummary"`
	Tags                     string `yaml:"tags"`
	PackageTitle             string `yaml:"packageTitle"`
	VersionDownloadCount     int
	MinClientVersion         string `yaml:"minClientVersion"`
	LastEdited               time.Time
	LicenseUrl               string
	LicenseNames             string
	LicenseReportUrl         string
}

func (e *Entry) ToAtom() ([]byte, error) {
	w := newWriter()

	// <entry xml:base="https://www.nuget.org/api/v2/"
	//       xmlns="http://www.w3.org/2005/Atom"
	//       xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices"
	//       xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">

	w.el_with_attr("entry")
	w.attr("xml:base", "https://www.nuget.org/api/v2/", false)
	w.attr("xmlns", "http://www.w3.org/2005/Atom", false)
	w.attr("xmlns:d", "http://schemas.microsoft.com/ado/2007/08/dataservices", false)
	w.attr("xmlns:m", "http://schemas.microsoft.com/ado/2007/08/dataservices/metadata", false)
	w.fin()

	e.writeFields(w)

	w.le("entry")

	return w.Bytes(), nil
}

func (e *Entry) writeFields(w *nugetXmlWriter) {
	// <id>https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='8.0.2')</id>
	w.el("id").text(e.Id).le("id")

	// <category term="NuGetGallery.V2FeedPackage" scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" />
	w.el_with_attr("category")
	w.attr("term", "NuGetGallery.V2FeedPackage", false)
	w.attr("scheme", "http://schemas.microsoft.com/ado/2007/08/dataservices/scheme", false)
	w.fin_no_child_nodes()

	// We aren't going to support publishing through this awful API.  We'll use a real publishing API.
	// <link rel="edit" title="V2FeedPackage" href="Packages(Id='Newtonsoft.Json',Version='8.0.2')" />
	// <link rel="edit-media" title="V2FeedPackage" href="Packages(Id='Newtonsoft.Json',Version='8.0.2')/$value" />

	// <title type="text">Newtonsoft.Json</title>
	w.el_with_attr("title").attr("type", "text", false).fin()
	w.text(e.Title)
	w.le("title")

	// <summary type="text"></summary>
	w.el_with_attr("summary").attr("type", "text", false).fin()
	w.text(e.Summary)
	w.le("summary")

	// <updated>2016-01-09T01:06:39Z</updated>
	w.el_with_text("updated", e.Updated.Format(time.RFC3339))

	// <author><name>James Newton-King</name></author>
	w.atom_author(e.AuthorName, e.AuthorEmail)

	// <content type="application/zip" src="https://www.nuget.org/api/v2/package/Newtonsoft.Json/8.0.2" />
	w.el_with_attr("content")
	w.attr("type", "application/zip", false)
	w.attr("src", e.DownloadUrl, false)
	w.fin_no_child_nodes()

	// Alright, now we're getting deeper into the XML vomit.
	// <m:properties>
	w.el("m:properties")

	// <d:Version>8.0.2</d:Version>
	w.odatafy_string("d:Version", e.Version)

	// <d:NormalizedVersion>8.0.2</d:NormalizedVersion>
	w.odatafy_string("d:NormalizedVersion", e.NormalizedVersion)

	// <d:Copyright m:null="true" />
	w.odatafy_string("d:Copyright", e.Copyright)

	// <d:Created m:type="Edm.DateTime">2016-01-09T01:06:39.39</d:Created>
	w.odatafy_time("d:Created", e.Created)

	// <d:Dependencies></d:Dependencies>
	w.odatafy_array("d:Dependencies", e.Dependencies)

	// <d:Description>Json.NET is a popular high-performance JSON framework for .NET</d:Description>
	w.odatafy_string("d:Description", e.Description)

	// <d:DownloadCount m:type="Edm.Int32">23441532</d:DownloadCount>
	w.odatafy_int32("d:DownloadCount", e.DownloadCount)

	// <d:GalleryDetailsUrl>https://www.nuget.org/packages/Newtonsoft.Json/8.0.2</d:GalleryDetailsUrl>
	w.odatafy_string("d:GalleryDetailsUrl", e.GalleryDetailsUrl)

	// <d:IconUrl>http://www.newtonsoft.com/content/images/nugeticon.png</d:IconUrl>
	w.odatafy_string("d:IconUrl", e.IconUrl)

	// <d:IsLatestVersion m:type="Edm.Boolean">true</d:IsLatestVersion>
	w.odatafy_bool("d:IsLatestVersion", e.IsLatestVersion)

	// <d:IsAbsoluteLatestVersion m:type="Edm.Boolean">true</d:IsAbsoluteLatestVersion>
	w.odatafy_bool("d:IsAbsoluteLatestVersion", e.IsAbsoluteLatestVersion)

	// <d:IsPrerelease m:type="Edm.Boolean">false</d:IsPrerelease>
	w.odatafy_bool("d:IsPrerelease", e.IsPrerelease)

	// <d:Language>en-US</d:Language>
	w.odatafy_string("d:Language", e.Language)

	// <d:Published m:type="Edm.DateTime">2016-01-09T01:06:39.39</d:Published>
	w.odatafy_time("d:Published", e.Published)

	// <d:PackageHash>e5yWmEfu68rmtG431zl9N/7PlNKQDIuiDW5MHlEFAZcecakcxrIGnKqrPAtWNILzK2oNanRB5cD150MYhECK3g==</d:PackageHash>
	w.odatafy_string("d:PackageHash", e.PackageHash)

	// <d:PackageHashAlgorithm>SHA512</d:PackageHashAlgorithm>
	w.odatafy_string("d:PackageHashAlgorithm", e.PackageHashAlgorithm)

	// <d:PackageSize m:type="Edm.Int64">1365056</d:PackageSize>
	w.odatafy_int64("d:PackageSize", e.PackageSize)

	// <d:ProjectUrl>http://www.newtonsoft.com/json</d:ProjectUrl>
	w.odatafy_string("d:ProjectUrl", e.ProjectUrl)

	// <d:ReportAbuseUrl>https://www.nuget.org/package/ReportAbuse/Newtonsoft.Json/8.0.2</d:ReportAbuseUrl>
	w.odatafy_string("d:ReportAbuseUrl", e.ReportAbuseUrl)

	// <d:ReleaseNotes m:null="true" />
	w.odatafy_string("d:ReleaseNotes", e.ReleaseNotes)

	// <d:RequireLicenseAcceptance m:type="Edm.Boolean">false</d:RequireLicenseAcceptance>
	w.odatafy_bool("d:RequireLicenseAcceptance", e.RequireLicenseAcceptance)

	// <d:Summary m:null="true" />
	w.odatafy_string("d:Summary", e.PackageSummary)

	// <d:Tags>json</d:Tags>
	w.odatafy_string("d:Tags", e.Tags)

	// <d:Title>Json.NET</d:Title>
	w.odatafy_string("d:Title", e.PackageTitle)

	// <d:VersionDownloadCount m:type="Edm.Int32">303719</d:VersionDownloadCount>
	w.odatafy_int32("d:VersionDownloadCount", e.VersionDownloadCount)

	// <d:MinClientVersion m:null="true" />
	w.odatafy_string("d:MinClientVersion", e.MinClientVersion)

	// <d:LastEdited m:type="Edm.DateTime" m:null="true" />
	w.odatafy_time("d:LastEdited", e.LastEdited)

	// <d:LicenseUrl>https://raw.github.com/JamesNK/Newtonsoft.Json/master/LICENSE.md</d:LicenseUrl>
	w.odatafy_string("d:LicenseUrl", e.LicenseUrl)

	// <d:LicenseNames m:null="true" />
	w.odatafy_string("d:LicenseNames", e.LicenseNames)

	// <d:LicenseReportUrl m:null="true" />
	w.odatafy_string("d:LicenseReportUrl", e.LicenseReportUrl)

	// Sit back, take a deep breath.  We're done with this shit.
	w.le("m:properties")
}
