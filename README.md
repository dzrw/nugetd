# nugetd - A NuGet server without all of the bullsh!t

Your search for an open-source NuGet server is over.

## Why nugetd?

### Open Source

There are several NuGet servers out there that aren't free.  This one's free.

There's almost zero documentation on the NuGet v2 API out there because OData is a giant spec and the API is RPC-heavy.  We reverse engineered it.

### DevOps Friendly

Put the gun down, and step away from IIS.

It's 2016, yet time seems to have passed NuGet by. NuGet servers out there ask your operations teams to deploy Windows boxes or host on Azure or build some custom C# solution. That's silly.  A NuGet server isn't special - it serves zip files and has some (clumsy) search APIs, so we designed this server to run on your standard Linux box.  Or, heck, cross-compile it for Windows.  It's written in Go 1.5, so go nuts.

We're running it in a Docker container.  

### Just the Basics for Now

#### NuGet.exe

Today, the server supports the NuGet CLI commands below.  Later, we'll fill out the publishing capabilities, but it's just as easy to `scp` your packages up to the server from your build system.  So maybe we won't implement support for `NuGet push`.

- Install a specific package: `NuGet install "Newtonsoft.Json" -Version "8.0.2" -Source nugetd -Verbosity detailed`
- Restore packages: `NuGet restore packages.config -PackagesDirectory packages -Source nugetd -Verbosity detailed`
- List all packages: `NuGet list -Source nugetd -Verbosity detailed`

#### Visual Studio

Visual Studio seems to be able to show packages hosted from this server, but it doesn't seem like you can install packages from the GUI.  We haven't tried the PowerShell CLI.

#### /bin/bash

I hope to add a `curl`-friendly REST API in the future.
