package main

// http://chris.eldredge.io/blog/2013/02/25/fun-with-nuget-rest-api/
// http://stackoverflow.com/questions/10231209/where-can-i-find-the-documentation-for-the-nuget-feed-api
//

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"regexp"
	"rekt"
	"strconv"
	"strings"
	"time"
)

const (
	NuGetRootDocument string = `<service xmlns="http://www.w3.org/2007/app" xmlns:atom="http://www.w3.org/2005/Atom" xml:base="https://www.nuget.org/api/v2/">
  <workspace>
    <atom:title>Default</atom:title>
    <collection href="Packages">
      <atom:title>Packages</atom:title>
    </collection>
  </workspace>
</service>`

	NuGetMetadata string = `<edmx:Edmx xmlns:edmx="http://schemas.microsoft.com/ado/2007/06/edmx" Version="1.0">
  <edmx:DataServices xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata" m:DataServiceVersion="2.0" m:MaxDataServiceVersion="2.0">
    <Schema xmlns="http://schemas.microsoft.com/ado/2006/04/edm" Namespace="NuGetGallery">
      <EntityType Name="V2FeedPackage" m:HasStream="true">
        <Key>
          <PropertyRef Name="Id"/>
          <PropertyRef Name="Version"/>
        </Key>
        <Property Name="Id" Type="Edm.String" Nullable="false" m:FC_TargetPath="SyndicationTitle" m:FC_ContentKind="text" m:FC_KeepInContent="false"/>
        <Property Name="Version" Type="Edm.String" Nullable="false"/>
        <Property Name="NormalizedVersion" Type="Edm.String"/>
        <Property Name="Authors" Type="Edm.String" m:FC_TargetPath="SyndicationAuthorName" m:FC_ContentKind="text" m:FC_KeepInContent="false"/>
        <Property Name="Copyright" Type="Edm.String"/>
        <Property Name="Created" Type="Edm.DateTime" Nullable="false"/>
        <Property Name="Dependencies" Type="Edm.String"/>
        <Property Name="Description" Type="Edm.String"/>
        <Property Name="DownloadCount" Type="Edm.Int32" Nullable="false"/>
        <Property Name="GalleryDetailsUrl" Type="Edm.String"/>
        <Property Name="IconUrl" Type="Edm.String"/>
        <Property Name="IsLatestVersion" Type="Edm.Boolean" Nullable="false"/>
        <Property Name="IsAbsoluteLatestVersion" Type="Edm.Boolean" Nullable="false"/>
        <Property Name="IsPrerelease" Type="Edm.Boolean" Nullable="false"/>
        <Property Name="Language" Type="Edm.String"/>
        <Property Name="LastUpdated" Type="Edm.DateTime" Nullable="false" m:FC_TargetPath="SyndicationUpdated" m:FC_ContentKind="text" m:FC_KeepInContent="false"/>
        <Property Name="Published" Type="Edm.DateTime" Nullable="false"/>
        <Property Name="PackageHash" Type="Edm.String"/>
        <Property Name="PackageHashAlgorithm" Type="Edm.String"/>
        <Property Name="PackageSize" Type="Edm.Int64" Nullable="false"/>
        <Property Name="ProjectUrl" Type="Edm.String"/>
        <Property Name="ReportAbuseUrl" Type="Edm.String"/>
        <Property Name="ReleaseNotes" Type="Edm.String"/>
        <Property Name="RequireLicenseAcceptance" Type="Edm.Boolean" Nullable="false"/>
        <Property Name="Summary" Type="Edm.String" m:FC_TargetPath="SyndicationSummary" m:FC_ContentKind="text" m:FC_KeepInContent="false"/>
        <Property Name="Tags" Type="Edm.String"/>
        <Property Name="Title" Type="Edm.String"/>
        <Property Name="VersionDownloadCount" Type="Edm.Int32" Nullable="false"/>
        <Property Name="MinClientVersion" Type="Edm.String"/>
        <Property Name="LastEdited" Type="Edm.DateTime"/>
        <Property Name="LicenseUrl" Type="Edm.String"/>
        <Property Name="LicenseNames" Type="Edm.String"/>
        <Property Name="LicenseReportUrl" Type="Edm.String"/>
      </EntityType>
      <EntityContainer Name="V2FeedContext" m:IsDefaultEntityContainer="true">
        <EntitySet Name="Packages" EntityType="NuGetGallery.V2FeedPackage"/>
        <FunctionImport Name="Search" ReturnType="Collection(NuGetGallery.V2FeedPackage)" EntitySet="Packages" m:HttpMethod="GET">
          <Parameter Name="searchTerm" Type="Edm.String"/>
          <Parameter Name="targetFramework" Type="Edm.String"/>
          <Parameter Name="includePrerelease" Type="Edm.Boolean"/>
        </FunctionImport>
        <FunctionImport Name="FindPackagesById" ReturnType="Collection(NuGetGallery.V2FeedPackage)" EntitySet="Packages" m:HttpMethod="GET">
          <Parameter Name="id" Type="Edm.String"/>
        </FunctionImport>
        <FunctionImport Name="GetUpdates" ReturnType="Collection(NuGetGallery.V2FeedPackage)" EntitySet="Packages" m:HttpMethod="GET">
          <Parameter Name="packageIds" Type="Edm.String"/>
          <Parameter Name="versions" Type="Edm.String"/>
          <Parameter Name="includePrerelease" Type="Edm.Boolean"/>
          <Parameter Name="includeAllVersions" Type="Edm.Boolean"/>
          <Parameter Name="targetFrameworks" Type="Edm.String"/>
          <Parameter Name="versionConstraints" Type="Edm.String"/>
        </FunctionImport>
      </EntityContainer>
    </Schema>
  </edmx:DataServices>
</edmx:Edmx>`

	NuGetSearchMalformedQuery string = `<?xml version="1.0" encoding="utf-8"?>
<m:error xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
  <m:code />
  <m:message xml:lang="en-US">Bad Request - Error in query syntax.</m:message>
</m:error>`

	NuGetPackageNotFound string = `<?xml version="1.0" encoding="utf-8"?>
<m:error xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <m:code />
    <m:message xml:lang="en-US">Resource not found for the segment 'Packages'.</m:message>
</m:error>`

	XmlMimeType string = "application/xml;charset=utf8"

	AtomMimeType string = "application/atom+xml;type=entry;charset=utf-8"
)

var registry *rekt.PackageRegistry

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	registry = openPackageRegistry("./packages.yml")
	log.Infof("Loaded %d NuGet package definitions.", len(registry.Packages))
}

func main() {
	r := gin.Default()

	// just a sanity check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// hard-coded root document
	r.GET("/api/v2/", func(c *gin.Context) {
		c.Data(200, XmlMimeType, []byte(NuGetRootDocument))
	})

	// where the (dubious) magic happens
	r.GET("/api/v2/:action", NuGetQueryDecoder(), NuGetFunctionDecoder(), rpcEndpoint)

	// static assests
	r.Static("packages", "./packages")

	r.Run() // listen and server on 0.0.0.0:8080
}

func openPackageRegistry(filename string) *rekt.PackageRegistry {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	registry, err := rekt.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	return registry
}

func rpcEndpoint(c *gin.Context) {
	op := c.MustGet("nuget.op").(string)
	query := c.MustGet("nuget.search").(*NuGetQuery)

	oplog := log.WithFields(log.Fields{
		"op":    op,
		"query": query,
	})

	switch op {
	case "$metadata":
		c.Data(200, XmlMimeType, []byte(NuGetMetadata))
		break
	case "Packages":
		data := c.MustGet("nuget.op.data").(*PackagesCommand)
		Packages(c, data, query, oplog.WithField("data", data))
	case "FindPackagesById":
		data := c.MustGet("nuget.op.data").(*FindPackagesByIdCommand)
		FindPackagesById(c, data, query, oplog.WithField("data", data))
	case "Search":
		data := c.MustGet("nuget.op.data").(*SearchCommand)
		Search(c, data, query, oplog.WithField("data", data))
	case "GetUpdates":
		data := c.MustGet("nuget.op.data").(*GetUpdatesCommand)
		GetUpdates(c, data, query, oplog.WithField("data", data))
	default:
		c.Status(400)
	}
}

/*
type Item struct {
	Title       string
	Link        *Link
	Author      *Author
	Description string // used as description in rss, summary in atom
	Id          string // used as guid in rss, id in atom
	Updated     time.Time
	Created     time.Time
}
*/

var (
	pkgNewtonsoftJson = &rekt.Entry{
		Id:                   "https://www.example.com/api/v2/Packages(Id='Newtonsoft.Json',Version='8.0.2')",
		Title:                "Newtonsoft.Json",
		Updated:              time.Now(),
		AuthorName:           "James Newton-King",
		DownloadUrl:          "https://www.nuget.org/api/v2/package/Newtonsoft.Json/8.0.2",
		Version:              "8.0.2",
		NormalizedVersion:    "8.0.2",
		DownloadCount:        42,
		GalleryDetailsUrl:    "https://www.nuget.org/packages/Newtonsoft.Json/8.0.2",
		IconUrl:              "http://www.newtonsoft.com/content/images/nugeticon.png",
		IsLatestVersion:      true,
		Language:             "en-US",
		Published:            time.Now(),
		PackageHash:          "e5yWmEfu68rmtG431zl9N/7PlNKQDIuiDW5MHlEFAZcecakcxrIGnKqrPAtWNILzK2oNanRB5cD150MYhECK3g==",
		PackageHashAlgorithm: "SHA512",
		PackageSize:          1365056,
		Tags:                 "json",
		PackageTitle:         "Json.NET",
	}
)

func Packages(c *gin.Context, d *PackagesCommand, q *NuGetQuery, logger *log.Entry) {
	logger.Info("Packages")

	for _, pkg := range registry.Packages {
		if pkg.Title == d.Id && pkg.Version == d.Version {
			atom, err := pkg.ToAtom()
			if err != nil {
				c.Status(500)
				return
			}

			c.Data(200, AtomMimeType, atom)
			return
		}
	}

	c.Data(404, XmlMimeType, []byte(NuGetPackageNotFound))
}

/*
<?xml version="1.0" encoding="utf-8"?>
<feed xml:base="https://www.nuget.org/api/v2/" xmlns="http://www.w3.org/2005/Atom" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <id>https://www.nuget.org/api/v2/FindPackagesById</id>
    <title type="text">FindPackagesById</title>
    <updated>2016-02-07T01:28:53Z</updated>
    <link rel="self" title="FindPackagesById" href="FindPackagesById" />
    <author>
        <name />
    </author>
</feed>
*/

func FindPackagesById(c *gin.Context, d *FindPackagesByIdCommand, q *NuGetQuery, logger *log.Entry) {
	logger.Info("FindPackagesById")

	var match *rekt.Entry
	for _, pkg := range registry.Packages {
		if pkg.Title == d.Id && (q.Filter != "IsLatestVersion" || pkg.IsLatestVersion) {
			match = pkg
		}
	}

	feed := &rekt.Feed{
		Id:      "https://www.example.com/api/v2/FindPackagesById",
		Title:   "FindPackagesById",
		Updated: time.Now(),
		Entries: []*rekt.Entry{match},
	}

	atom, err := feed.ToAtom()
	if err != nil {
		c.Status(500)
		return
	}

	c.Data(200, AtomMimeType, atom)
	return
}

/*
<?xml version="1.0" encoding="utf-8"?>
<feed xml:base="https://www.nuget.org/api/v2/"
      xmlns="http://www.w3.org/2005/Atom"
      xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices"
      xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <id>https://www.nuget.org/api/v2/Search</id>
    <title type="text">Search</title>
    <updated>2016-02-07T01:25:04Z</updated>
    <link rel="self" title="Search" href="Search" />
    <author>
        <name />
    </author>
</feed>
*/

func Search(c *gin.Context, d *SearchCommand, q *NuGetQuery, logger *log.Entry) {
	logger.Info("Search")

	feed := &rekt.Feed{
		Id:      "https://www.example.com/api/v2/Search",
		Title:   "Search",
		Updated: time.Now(),
		Entries: registry.Packages, // just return all packages for every search
	}

	atom, err := feed.ToAtom()
	if err != nil {
		c.Status(500)
		return
	}

	c.Data(200, AtomMimeType, atom)
}

func GetUpdates(c *gin.Context, d *GetUpdatesCommand, q *NuGetQuery, logger *log.Entry) {
	logger.Warn("GetUpdates not implemented")
	c.Status(501)
}

var (
	reFunctionArgs = regexp.MustCompile(`(?P<Name>.+?)='(?P<Value>[^']*)'`)
)

// Parse NuGet's RPC-style routes into some kind of command object
func NuGetFunctionDecoder() gin.HandlerFunc {
	return func(c *gin.Context) {
		action := c.Param("action")

		if action == "$metadata" {
			c.Set("nuget.op", "$metadata")
		}

		// e.g. Packages(Id='Newtonsoft.Json',Version='7.0')
		if strings.HasPrefix(action, "Packages(") && strings.HasSuffix(action, ")") {
			action = strings.TrimPrefix(action, "Packages(")
			action = strings.TrimSuffix(action, ")")

			// e.g. Id='Newtonsoft.Json',Version='7.0'
			d := &PackagesCommand{}
			for _, m := range reFunctionArgs.FindAllStringSubmatch(action, -1) {

				// should probably fix the regexp to avoid this bug...
				k := strings.TrimPrefix(m[1], ",")
				v := m[2]

				switch k {
				case "Id":
					d.Id = v
					break
				case "Version":
					d.Version = v
					break
				}
			}

			c.Set("nuget.op", "Packages")
			c.Set("nuget.op.data", d)
		}

		// e.g. FindPackagesById()?$filter=IsLatestVersion&$orderby=Version%20desc&$top=1&id='Newtonsoft.Json'
		if strings.HasPrefix(action, "FindPackagesById()") {
			d := &FindPackagesByIdCommand{}

			d.Id = nuget_trim(c.Query("id"))

			c.Set("nuget.op", "FindPackagesById")
			c.Set("nuget.op.data", d)
		}

		// e.g. Search()?$filter=IsLatestVersion&$orderby=Id&$skip=0&$top=30&searchTerm=''&targetFramework=''&includePrerelease=false
		if strings.HasPrefix(action, "Search()") {
			d := &SearchCommand{}

			d.SearchTerm = nuget_trim(c.Query("searchTerm"))
			d.TargetFramework = nuget_trim(c.Query("targetFramework"))
			d.IncludePrerelease = c.DefaultQuery("includePrerelease", "false") == "true"

			c.Set("nuget.op", "Search")
			c.Set("nuget.op.data", d)
		}

		// I haven't actually seen an example for this one yet.  Just going off the $metadata.
		if strings.HasPrefix(action, "GetUpdates()") {
			d := &GetUpdatesCommand{}

			d.PackageIds = c.Query("packageIds")
			d.Versions = c.Query("versions")
			d.IncludePrerelease = c.DefaultQuery("includePrerelease", "false") == "true"
			d.IncludeAllVersions = c.DefaultQuery("includeAllVersions", "false") == "true"
			d.TargetFrameworks = c.Query("targetFrameworks")
			d.VersionConstraints = c.Query("versionConstraints")

			c.Set("nuget.op", "GetUpdates")
			c.Set("nuget.op.data", d)
		}
	}
}

func NuGetQueryDecoder() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := &NuGetQuery{}

		d.Filter = c.Query("$filter")
		d.OrderBy = c.Query("$orderby")
		d.Skip = IntQuery(c, "$skip", 0)
		d.Top = IntQuery(c, "$top", 1)

		c.Set("nuget.search", d)
	}
}

// NuGet's query strings are quoted. This isn't the dumbest thing I've had to workaround today.
func nuget_trim(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "'"), "'")
}

func IntQuery(c *gin.Context, param string, defaultValue int) int {
	v, ok := c.GetQuery(param)
	if !ok {
		return defaultValue
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}

	return i
}

type NuGetQuery struct {
	Filter  string `name:$filter`  // some kind of boolean expression
	OrderBy string `name:$orderby` // `foo` or `foo desc`
	Skip    int    `name:$skip`    // pagination
	Top     int    `name:$top`     // pagination
}

type PackagesCommand struct {
	Id      string
	Version string
}

type FindPackagesByIdCommand struct {
	Id string
}

type SearchCommand struct {
	SearchTerm        string
	TargetFramework   string
	IncludePrerelease bool
}

type GetUpdatesCommand struct {
	PackageIds         string
	Versions           string
	IncludePrerelease  bool
	IncludeAllVersions bool
	TargetFrameworks   string
	VersionConstraints string
}
