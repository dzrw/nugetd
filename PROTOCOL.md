PROTOCOL.md


/api/v2
=======

This project targets a subset of /api/v2/ rather than whatever design-by-committee disaster is going on in v3.  We're not saying that v2 isn't also an abomination, but at least it's locked-down and simplistic.


### Install a specific version of a package.

The client will send one or more search requests to find a package with a particular version.

```PowerShell
NuGet install "Newtonsoft.Json" -Version "7.0.0" -Source nugv2 -Verbosity detailed
```

a. GET https://www.nuget.org/api/v2/
b. GET https://www.nuget.org/api/v2/$metadata
1. GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='7.0')
2. GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='7.0.0')
3. GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='7.0.0.0')


### Install the latest version of a package

```PowerShell
NuGet install "Newtonsoft.Json" -Source mugv2 -Verbosity detailed
```

a. GET https://www.nuget.org/api/v2/
b. GET https://www.nuget.org/api/v2/$metadata
1. GET https://www.nuget.org/api/v2/FindPackagesById()?$filter=IsLatestVersion&$orderby=Version desc&$top=1&id='Newtonsoft.Json'
2. GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='8.0.2')


### List all packages

```PowerShell
NuGet list -Source nugv2 -Verbosity detailed
```

a. GET https://www.nuget.org/api/v2/
b. GET https://www.nuget.org/api/v2/$metadata
1. GET https://www.nuget.org/api/v2/Search()?$filter=IsLatestVersion&$orderby=Id&$skip=0&$top=30&searchTerm=''&targetFramework=''&includePrerelease=false

### Search for all packages containing a search term excluding pre-release packages

```PowerShell
NuGet list "Newtonsoft.Json" -Source nugv2 -Verbosity detailed
```

a. GET https://www.nuget.org/api/v2/
b. GET https://www.nuget.org/api/v2/$metadata
1. GET https://www.nuget.org/api/v2/Search()?$filter=IsLatestVersion&$orderby=Id&$skip=0&$top=30&searchTerm='Newtonsoft.Json'&targetFramework=''&includePrerelease=false


### Search for all packages containing a search term including pre-release packages

```PowerShell
NuGet list "Newtonsoft.Json" -Prerelease -Source nugv2 -Verbosity detailed
```

a. GET https://www.nuget.org/api/v2/
b. GET https://www.nuget.org/api/v2/$metadata
1. GET https://www.nuget.org/api/v2/Search()?$filter=IsAbsoluteLatestVersion&$orderby=Id&$skip=30&$top=30&searchTerm='Newtonsoft.Json'&targetFramework=''&includePrerelease=true


### Download all packages defined in a packages.config file.

The packages.config file:

```XML
<?xml version="1.0" encoding="utf-8"?>
<packages>
  <package id="Newtonsoft.Json" version="6.0.8" targetFramework="net45" />
  <package id="RestSharp" version="105.1.0" targetFramework="net45" />
</packages>
```

```PowerShell
mkdir packages
NuGet restore packages.config -Source nugv2 -NoCache -Verbosity detailed -PackagesDirectory packages
```

a. GET https://www.nuget.org/api/v2/
b. GET https://www.nuget.org/api/v2/$metadata
1. GET https://www.nuget.org/api/v2/Packages(Id='RestSharp',Version='105.1.0')
2. GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='6.0.8')


### Example: GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='7.0')

A package that doesn't exist.

Status: HTTP 404 Not Found

```XML
<?xml version="1.0" encoding="utf-8"?>
<m:error xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <m:code />
    <m:message xml:lang="en-US">Resource not found for the segment 'Packages'.</m:message>
</m:error>
```

### Example: GET https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='8.0.2')

Content-Type: application/atom+xml;type=feed;charset=utf-8

It's just an Atom document with an extra property bag object tacked on, and some goofy OData URLs.

```XML
<?xml version="1.0" encoding="utf-8"?>
<entry xml:base="https://www.nuget.org/api/v2/"
       xmlns="http://www.w3.org/2005/Atom"
       xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices"
       xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
  <id>https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='8.0.2')</id>
  <category term="NuGetGallery.V2FeedPackage" 
            scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" />
  <link rel="edit" title="V2FeedPackage" href="Packages(Id='Newtonsoft.Json',Version='8.0.2')" />
  <title type="text">Newtonsoft.Json</title>
  <summary type="text"></summary>
  <updated>2016-01-09T01:06:39Z</updated>
  <author>
    <name>James Newton-King</name>
  </author>
  <link rel="edit-media" title="V2FeedPackage" href="Packages(Id='Newtonsoft.Json',Version='8.0.2')/$value" />
  <content type="application/zip" src="https://www.nuget.org/api/v2/package/Newtonsoft.Json/8.0.2" />
  <m:properties>
    <d:Version>8.0.2</d:Version>
    <d:NormalizedVersion>8.0.2</d:NormalizedVersion>
    <d:Copyright m:null="true" />
    <d:Created m:type="Edm.DateTime">2016-01-09T01:06:39.39</d:Created>
    <d:Dependencies></d:Dependencies>
    <d:Description>Json.NET is a popular high-performance JSON framework for .NET</d:Description>
    <d:DownloadCount m:type="Edm.Int32">23441532</d:DownloadCount>
    <d:GalleryDetailsUrl>https://www.nuget.org/packages/Newtonsoft.Json/8.0.2</d:GalleryDetailsUrl>
    <d:IconUrl>http://www.newtonsoft.com/content/images/nugeticon.png</d:IconUrl>
    <d:IsLatestVersion m:type="Edm.Boolean">true</d:IsLatestVersion>
    <d:IsAbsoluteLatestVersion m:type="Edm.Boolean">true</d:IsAbsoluteLatestVersion>
    <d:IsPrerelease m:type="Edm.Boolean">false</d:IsPrerelease>
    <d:Language>en-US</d:Language>
    <d:Published m:type="Edm.DateTime">2016-01-09T01:06:39.39</d:Published>
    <d:PackageHash>e5yWmEfu68rmtG431zl9N/7PlNKQDIuiDW5MHlEFAZcecakcxrIGnKqrPAtWNILzK2oNanRB5cD150MYhECK3g==</d:PackageHash>
    <d:PackageHashAlgorithm>SHA512</d:PackageHashAlgorithm>
    <d:PackageSize m:type="Edm.Int64">1365056</d:PackageSize>
    <d:ProjectUrl>http://www.newtonsoft.com/json</d:ProjectUrl>
    <d:ReportAbuseUrl>https://www.nuget.org/package/ReportAbuse/Newtonsoft.Json/8.0.2</d:ReportAbuseUrl>
    <d:ReleaseNotes m:null="true" />
    <d:RequireLicenseAcceptance m:type="Edm.Boolean">false</d:RequireLicenseAcceptance>
    <d:Summary m:null="true" />
    <d:Tags>json</d:Tags>
    <d:Title>Json.NET</d:Title>
    <d:VersionDownloadCount m:type="Edm.Int32">303719</d:VersionDownloadCount>
    <d:MinClientVersion m:null="true" />
    <d:LastEdited m:type="Edm.DateTime" m:null="true" />
    <d:LicenseUrl>https://raw.github.com/JamesNK/Newtonsoft.Json/master/LICENSE.md</d:LicenseUrl>
    <d:LicenseNames m:null="true" />
    <d:LicenseReportUrl m:null="true" />
  </m:properties>
</entry>
```

### Example: GET https://www.nuget.org/api/v2/FindPackagesById()?$filter=IsLatestVersion&$orderby=Version%20desc&$top=1&id='jfjdk3'

No search results found.

```XML
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
```

### Example: GET https://www.nuget.org/api/v2/FindPackagesById()?$filter=IsLatestVersion&$orderby=Version%20desc&$top=1&id='Newtonsoft.Json'

Content-Type: application/atom+xml;type=feed;charset=utf-8


```XML
<?xml version="1.0" encoding="utf-8"?>
<feed xml:base="https://www.nuget.org/api/v2/" xmlns="http://www.w3.org/2005/Atom" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <id>https://www.nuget.org/api/v2/FindPackagesById</id>
    <title type="text">FindPackagesById</title>
    <updated>2016-02-07T01:18:35Z</updated>
    <link rel="self" title="FindPackagesById" href="FindPackagesById" />
    <entry>
        <id>https://www.nuget.org/api/v2/Packages(Id='Newtonsoft.Json',Version='8.0.2')</id>
        <category term="NuGetGallery.V2FeedPackage" scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" />
        <link rel="edit" title="V2FeedPackage" href="Packages(Id='Newtonsoft.Json',Version='8.0.2')" />
        <title type="text">Newtonsoft.Json</title>
        <summary type="text"></summary>
        <updated>2016-01-09T01:06:39Z</updated>
        <author>
            <name>James Newton-King</name>
        </author>
        <link rel="edit-media" title="V2FeedPackage" href="Packages(Id='Newtonsoft.Json',Version='8.0.2')/$value" />
        <content type="application/zip" src="https://www.nuget.org/api/v2/package/Newtonsoft.Json/8.0.2" />
        <m:properties>
            <d:Version>8.0.2</d:Version>
            <d:NormalizedVersion>8.0.2</d:NormalizedVersion>
            <d:Copyright m:null="true" />
            <d:Created m:type="Edm.DateTime">2016-01-09T01:06:39.39</d:Created>
            <d:Dependencies></d:Dependencies>
            <d:Description>Json.NET is a popular high-performance JSON framework for .NET</d:Description>
            <d:DownloadCount m:type="Edm.Int32">23456163</d:DownloadCount>
            <d:GalleryDetailsUrl>https://www.nuget.org/packages/Newtonsoft.Json/8.0.2</d:GalleryDetailsUrl>
            <d:IconUrl>http://www.newtonsoft.com/content/images/nugeticon.png</d:IconUrl>
            <d:IsLatestVersion m:type="Edm.Boolean">true</d:IsLatestVersion>
            <d:IsAbsoluteLatestVersion m:type="Edm.Boolean">true</d:IsAbsoluteLatestVersion>
            <d:IsPrerelease m:type="Edm.Boolean">false</d:IsPrerelease>
            <d:Language>en-US</d:Language>
            <d:Published m:type="Edm.DateTime">2016-01-09T01:06:39.39</d:Published>
            <d:PackageHash>e5yWmEfu68rmtG431zl9N/7PlNKQDIuiDW5MHlEFAZcecakcxrIGnKqrPAtWNILzK2oNanRB5cD150MYhECK3g==</d:PackageHash>
            <d:PackageHashAlgorithm>SHA512</d:PackageHashAlgorithm>
            <d:PackageSize m:type="Edm.Int64">1365056</d:PackageSize>
            <d:ProjectUrl>http://www.newtonsoft.com/json</d:ProjectUrl>
            <d:ReportAbuseUrl>https://www.nuget.org/package/ReportAbuse/Newtonsoft.Json/8.0.2</d:ReportAbuseUrl>
            <d:ReleaseNotes m:null="true" />
            <d:RequireLicenseAcceptance m:type="Edm.Boolean">false</d:RequireLicenseAcceptance>
            <d:Summary m:null="true" />
            <d:Tags>json</d:Tags>
            <d:Title>Json.NET</d:Title>
            <d:VersionDownloadCount m:type="Edm.Int32">304860</d:VersionDownloadCount>
            <d:MinClientVersion m:null="true" />
            <d:LastEdited m:type="Edm.DateTime" m:null="true" />
            <d:LicenseUrl>https://raw.github.com/JamesNK/Newtonsoft.Json/master/LICENSE.md</d:LicenseUrl>
            <d:LicenseNames m:null="true" />
            <d:LicenseReportUrl m:null="true" />
        </m:properties>
    </entry>
</feed>
```

### Example: GET https://www.nuget.org/api/v2/Search()?$filter=IsLatestVersion&$orderby=Id&$skip=0&$top=3&searchTerm='jfjdk3'&targetFramework=''&includePrerelease=false

No search results found.

```XML
<?xml version="1.0" encoding="utf-8"?>
<feed xml:base="https://www.nuget.org/api/v2/" xmlns="http://www.w3.org/2005/Atom" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <id>https://www.nuget.org/api/v2/Search</id>
    <title type="text">Search</title>
    <updated>2016-02-07T01:25:04Z</updated>
    <link rel="self" title="Search" href="Search" />
    <author>
        <name />
    </author>
</feed>
```

### Example: GET https://www.nuget.org/api/v2/Search()?$filter=IsLatestVersion&$orderby=Id&$skip=0&$top=3&searchTerm='Newtonsoft.Json'&targetFramework=''&includePrerelease=false

Content-Type: application/atom+xml;type=feed;charset=utf-8


```XML
<?xml version="1.0" encoding="utf-8"?>
<feed xml:base="https://www.nuget.org/api/v2/" xmlns="http://www.w3.org/2005/Atom" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <id>https://www.nuget.org/api/v2/Search</id>
    <title type="text">Search</title>
    <updated>2016-02-07T01:22:33Z</updated>
    <link rel="self" title="Search" href="Search" />
    <entry>
        <id>https://www.nuget.org/api/v2/Packages(Id='AppNet.NET',Version='1.8.2.1')</id>
        <category term="NuGetGallery.V2FeedPackage" scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" />
        <link rel="edit" title="V2FeedPackage" href="Packages(Id='AppNet.NET',Version='1.8.2.1')" />
        <title type="text">AppNet.NET</title>
        <summary type="text"></summary>
        <updated>2015-05-03T18:40:27Z</updated>
        <author>
            <name>Sven Walther / lI' Ghun</name>
        </author>
        <link rel="edit-media" title="V2FeedPackage" href="Packages(Id='AppNet.NET',Version='1.8.2.1')/$value" />
        <content type="application/zip" src="https://www.nuget.org/api/v2/package/AppNet.NET/1.8.2.1" />
        <m:properties>
            <d:Version>1.8.2.1</d:Version>
            <d:NormalizedVersion>1.8.2.1</d:NormalizedVersion>
            <d:Copyright>Sven Walther</d:Copyright>
            <d:Created m:type="Edm.DateTime">2013-09-05T12:04:45.853</d:Created>
            <d:Dependencies>Newtonsoft.Json:4.5.11:</d:Dependencies>
            <d:Description>A complete .NET 4.0 App.net API implementation including posts, messages and the File API for example.
Find an example app on the homepage which is included in the source
It needs the current version of Newtonsoft JSON.NET</d:Description>
            <d:DownloadCount m:type="Edm.Int32">6923</d:DownloadCount>
            <d:GalleryDetailsUrl>https://www.nuget.org/packages/AppNet.NET/1.8.2.1</d:GalleryDetailsUrl>
            <d:IconUrl m:null="true" />
            <d:IsLatestVersion m:type="Edm.Boolean">true</d:IsLatestVersion>
            <d:IsAbsoluteLatestVersion m:type="Edm.Boolean">true</d:IsAbsoluteLatestVersion>
            <d:IsPrerelease m:type="Edm.Boolean">false</d:IsPrerelease>
            <d:Language>en-GB</d:Language>
            <d:Published m:type="Edm.DateTime">2013-09-05T12:04:45.853</d:Published>
            <d:PackageHash>xdvREcFp9pGHC6tVzz8l7XWg50l/V+ZNM5ediq1M3zxAVs6Fepx9q6i+eWZQzILM9ujy7upiVYg0yToJ5XjPDg==</d:PackageHash>
            <d:PackageHashAlgorithm>SHA512</d:PackageHashAlgorithm>
            <d:PackageSize m:type="Edm.Int64">118147</d:PackageSize>
            <d:ProjectUrl>https://github.com/liGhun/AppNet.NET</d:ProjectUrl>
            <d:ReportAbuseUrl>https://www.nuget.org/package/ReportAbuse/AppNet.NET/1.8.2.1</d:ReportAbuseUrl>
            <d:ReleaseNotes>Added fetch parameters to channel get method

prior in 1.8.x:
Added search API
Added Configurations API
Added general paramerts to Search API</d:ReleaseNotes>
            <d:RequireLicenseAcceptance m:type="Edm.Boolean">false</d:RequireLicenseAcceptance>
            <d:Summary m:null="true" />
            <d:Tags>App.net</d:Tags>
            <d:Title>A complete .NET library for the App.net API</d:Title>
            <d:VersionDownloadCount m:type="Edm.Int32">522</d:VersionDownloadCount>
            <d:MinClientVersion m:null="true" />
            <d:LastEdited m:type="Edm.DateTime" m:null="true" />
            <d:LicenseUrl>https://github.com/liGhun/AppNet.NET/blob/master/LICENSE.txt</d:LicenseUrl>
            <d:LicenseNames>BSD-3-Clause</d:LicenseNames>
            <d:LicenseReportUrl></d:LicenseReportUrl>
        </m:properties>
    </entry>
    <entry>
        <id>https://www.nuget.org/api/v2/Packages(Id='AcklenAvenue.Queueing.Serializers.JsonNet',Version='1.0.1.25')</id>
        <category term="NuGetGallery.V2FeedPackage" scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" />
        <link rel="edit" title="V2FeedPackage" href="Packages(Id='AcklenAvenue.Queueing.Serializers.JsonNet',Version='1.0.1.25')" />
        <title type="text">AcklenAvenue.Queueing.Serializers.JsonNet</title>
        <summary type="text"></summary>
        <updated>2015-09-21T18:48:59Z</updated>
        <author>
            <name>AcklenAvenue</name>
        </author>
        <link rel="edit-media" title="V2FeedPackage" href="Packages(Id='AcklenAvenue.Queueing.Serializers.JsonNet',Version='1.0.1.25')/$value" />
        <content type="application/zip" src="https://www.nuget.org/api/v2/package/AcklenAvenue.Queueing.Serializers.JsonNet/1.0.1.25" />
        <m:properties>
            <d:Version>1.0.1.25</d:Version>
            <d:NormalizedVersion>1.0.1.25</d:NormalizedVersion>
            <d:Copyright>Copyright Acklen Avenue, 2012-2015</d:Copyright>
            <d:Created m:type="Edm.DateTime">2015-09-21T18:48:59.297</d:Created>
            <d:Dependencies>Newtonsoft.Json:7.0.1:</d:Dependencies>
            <d:Description>Newtonsoft.Json serializer implementation for AcklenAvenue.Queueing.</d:Description>
            <d:DownloadCount m:type="Edm.Int32">892</d:DownloadCount>
            <d:GalleryDetailsUrl>https://www.nuget.org/packages/AcklenAvenue.Queueing.Serializers.JsonNet/1.0.1.25</d:GalleryDetailsUrl>
            <d:IconUrl>https://raw.githubusercontent.com/AcklenAvenue/acklenavenue.github.io/master/assets/img/acklenavenuelogo_block.png</d:IconUrl>
            <d:IsLatestVersion m:type="Edm.Boolean">true</d:IsLatestVersion>
            <d:IsAbsoluteLatestVersion m:type="Edm.Boolean">true</d:IsAbsoluteLatestVersion>
            <d:IsPrerelease m:type="Edm.Boolean">false</d:IsPrerelease>
            <d:Language m:null="true" />
            <d:Published m:type="Edm.DateTime">2015-09-21T18:48:59.297</d:Published>
            <d:PackageHash>jzPvn9VEpohVr+K4mqIboH3dLWVEvJg7jfPf82BwJaPqktWjSvS/BxQHh/TZFBZ42fJhd6RSvb23Ps85vkWJOg==</d:PackageHash>
            <d:PackageHashAlgorithm>SHA512</d:PackageHashAlgorithm>
            <d:PackageSize m:type="Edm.Int64">5216</d:PackageSize>
            <d:ProjectUrl>https://github.com/AcklenAvenue/AcklenAvenue.Queueing</d:ProjectUrl>
            <d:ReportAbuseUrl>https://www.nuget.org/package/ReportAbuse/AcklenAvenue.Queueing.Serializers.JsonNet/1.0.1.25</d:ReportAbuseUrl>
            <d:ReleaseNotes m:null="true" />
            <d:RequireLicenseAcceptance m:type="Edm.Boolean">false</d:RequireLicenseAcceptance>
            <d:Summary m:null="true" />
            <d:Tags m:null="true" />
            <d:Title>AcklenAvenue.Queuing.Serializers.JsonNet</d:Title>
            <d:VersionDownloadCount m:type="Edm.Int32">179</d:VersionDownloadCount>
            <d:MinClientVersion m:null="true" />
            <d:LastEdited m:type="Edm.DateTime" m:null="true" />
            <d:LicenseUrl>http://opensource.org/licenses/MIT</d:LicenseUrl>
            <d:LicenseNames m:null="true" />
            <d:LicenseReportUrl m:null="true" />
        </m:properties>
    </entry>
    <entry>
        <id>https://www.nuget.org/api/v2/Packages(Id='Bifrost.JSON',Version='1.0.0.32')</id>
        <category term="NuGetGallery.V2FeedPackage" scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" />
        <link rel="edit" title="V2FeedPackage" href="Packages(Id='Bifrost.JSON',Version='1.0.0.32')" />
        <title type="text">Bifrost.JSON</title>
        <summary type="text"></summary>
        <updated>2016-01-06T20:22:57Z</updated>
        <author>
            <name>Dolittle</name>
        </author>
        <link rel="edit-media" title="V2FeedPackage" href="Packages(Id='Bifrost.JSON',Version='1.0.0.32')/$value" />
        <content type="application/zip" src="https://www.nuget.org/api/v2/package/Bifrost.JSON/1.0.0.32" />
        <m:properties>
            <d:Version>1.0.0.32</d:Version>
            <d:NormalizedVersion>1.0.0.32</d:NormalizedVersion>
            <d:Copyright>Copyright 2008 - 2015</d:Copyright>
            <d:Created m:type="Edm.DateTime">2016-01-06T20:22:57.4</d:Created>
            <d:Dependencies>Newtonsoft.Json:[6.0.8, ):|Bifrost:[1.0.0.32, ):</d:Dependencies>
            <d:Description>Support for Newtonsoft JSON.net for Bifrost</d:Description>
            <d:DownloadCount m:type="Edm.Int32">5152</d:DownloadCount>
            <d:GalleryDetailsUrl>https://www.nuget.org/packages/Bifrost.JSON/1.0.0.32</d:GalleryDetailsUrl>
            <d:IconUrl>http://bifrost.dolittle.com/img/bifrost-icon.png</d:IconUrl>
            <d:IsLatestVersion m:type="Edm.Boolean">true</d:IsLatestVersion>
            <d:IsAbsoluteLatestVersion m:type="Edm.Boolean">true</d:IsAbsoluteLatestVersion>
            <d:IsPrerelease m:type="Edm.Boolean">false</d:IsPrerelease>
            <d:Language m:null="true" />
            <d:Published m:type="Edm.DateTime">2016-01-06T20:22:57.4</d:Published>
            <d:PackageHash>JnzhXzVvdZd0TaR2DZfGEPl8W2bn9n8uO3oDzRwhfl+RyntQot+jzZ78iPYin/kEhtBXmp4ckHLCrFMoys8nIw==</d:PackageHash>
            <d:PackageHashAlgorithm>SHA512</d:PackageHashAlgorithm>
            <d:PackageSize m:type="Edm.Int64">10743</d:PackageSize>
            <d:ProjectUrl>http://bifrost.dolittle.com/</d:ProjectUrl>
            <d:ReportAbuseUrl>https://www.nuget.org/package/ReportAbuse/Bifrost.JSON/1.0.0.32</d:ReportAbuseUrl>
            <d:ReleaseNotes>http://bifr.st/ReleaseNotes</d:ReleaseNotes>
            <d:RequireLicenseAcceptance m:type="Edm.Boolean">false</d:RequireLicenseAcceptance>
            <d:Summary m:null="true" />
            <d:Tags>MVVM LOB CQRS</d:Tags>
            <d:Title>Bifrost Newtonsoft JSON.net support</d:Title>
            <d:VersionDownloadCount m:type="Edm.Int32">54</d:VersionDownloadCount>
            <d:MinClientVersion m:null="true" />
            <d:LastEdited m:type="Edm.DateTime" m:null="true" />
            <d:LicenseUrl>http://github.com/dolittle/Bifrost/blob/master/MIT-LICENSE.txt</d:LicenseUrl>
            <d:LicenseNames m:null="true" />
            <d:LicenseReportUrl m:null="true" />
        </m:properties>
    </entry>
</feed>
```

## Query String Analysis

```
GET https://www.nuget.org/api/v2/FindPackagesById()
?$filter=IsLatestVersion
&$orderby=Version desc
&$top=1&id='Newtonsoft.Json'

GET https://www.nuget.org/api/v2/Search()
?$filter=IsLatestVersion
&$orderby=Id
&$skip=0
&$top=30
&searchTerm=''
&targetFramework=''
&includePrerelease=false

GET https://www.nuget.org/api/v2/Search()
?$filter=IsLatestVersion
&$orderby=Id
&$skip=0
&$top=30
&searchTerm='Newtonsoft.Json'
&targetFramework=''
&includePrerelease=false

GET https://www.nuget.org/api/v2/Search()
?$filter=IsAbsoluteLatestVersion
&$orderby=Id
&$skip=30
&$top=30
&searchTerm='Newtonsoft.Json'
&targetFramework=''
&includePrerelease=true
```

Key | Observed Values | Remarks
--- | --- | ---
$filter | `IsLatestVersion`, `IsAbsoluteLatestVersion` | These look like boolean fields of the extended property bag.
$orderby | `Id`, `Version desc` | Does `Id` refer to the Atom entry id, or something magical? It's not in the bag.
$skip | `0`, `30`, `60` | pagination
$top | `1`, `30` | pagination
searchTerm | `''`, `'Newtonsoft.Json' |
targetFramework | `''` | .NET Framework (e.g. 3.5, 4.5)
includePrerelease | `true`, `false` | whether or not to include items with the `IsPrerelease` property set.
id | `Newtonsoft.Json` | Used when installing a package by name without specifying a version. 
