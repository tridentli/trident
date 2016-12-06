# Trident

Trident is a Trusted Information Exchange Toolkit.

This is the upcoming new iteration of code that runs [Ops-Trust](https://www.ops-trust.net)
and a variety of other instances.

The Trident database and functionality are based on the
[Portal code](https://github.com/ops-trust/portal/),
but Trident is a full redesign and from scratch rewrite in [Go](https://golang.org/).

The underlying lirbary has been cared out into a dedicated project [Pitchfork](https://github.com/tridentli/pitchfork) for use in other systems.

Useful links:
 * [Trident Website](https://trident.li)
 * [Github](https://github.com/tridentli)
 * [Email](mailto:project@trident.li)
 * [Twitter](https://twitter.com/tridentli)

## License

Trident is placed under the [Apache Version 2.0 License](http://www.apache.org/licenses/).

## Components

This section briefly describes the components inside Trident.

### tridentd

```tridentd``` is the Trident Daemon. It runs in the background and
published a web-server on ```localhost``` and port ```8333```.
This HTTP-based web-server provides an ```/api/``` URL for
Trident API calls and the rest of the exported details consist
of the Trident web interface.

The tridentd web-service can be published to the general Internet
by using a HTTP proxy like Nginx or Apache and forwarding the
correct requests to it.

tridentd contains all the logic in Trident and is the tool that
talks to the Database.
All commands are routed through tridentd. Permission handling is
thus checked by tridentd.

Configuration for tridentd is stored in ```/etc/trident/trident.conf```.
This file should only be readable by tridentd which runs under the ```trident``` user account.

#### Trident API

The API is available under ```/api/```.

The API is always accessible from ```localhost```, it can be limited
for publishing to other hosts with the ```API_Enabled``` config option.

Authentication is done using [JSON Web Tokens](http://jwt.io/).
One can retrieve a token by using the ```system login``` command.

This is compatible with the [OAuth2](http://oauth.net/2/) protocol.

(XXX: Todo - further specification to follow)

#### Web CLI

The Web CLI under ```/cli/``` allows access to a simple HTML form that
can be used to type in custom commands, in a similar way of using
the ```tcli``` command.

#### OAuth2 / OpenID Connect

Trident allows external sites to verify and authenticate that the user
is an authorized user of the Trident installation.

#### Wiki

The Wiki built-in to Trident can operate either in no-javascript mode, where a
simple HTML textarea is used to edit the Markdown, or a
[very nice inline editor](http://epiceditor.com) that has direct Preview of what is being edited.

The Wiki format is the [standard Markdown](https://daringfireball.net/projects/markdown/) with [SmartyPants](v) [GitHub Flavored Markdown](https://help.github.com/articles/github-flavored-markdown/) extensions.

Before providing the rendered HTML to the user it is sanitized using [Blue Monday Sanitizer](https://github.com/microcosm-cc/bluemonday) thus avoiding inline javascript and other possible XSS options.

Due to use of two different Markdown engines (EpicEditor's "marked" and [Black Friday](https://github.com/russross/blackfriday)) and then filtering through Blue Monday some inconsistencies between the preview and the final output are to be expected. The "Source" tab can be used to identify such differences.

Differences between pages can be seen using the History tab.

A Table of Contents is automatically generated from the headers of a page.

The Child Pages tab can be seen used to list all the sub pages for that Page.

##### Formatting Examples

###### Headers

```
First Level Header (h1)
=======================

Double Underlined

Second Level Header (h2)
------------------------

Single Underlined

# First Level Header (h1)

First Paragraph text

## Second Level Header (h2)

Second Paragraph text

### Third Level Header (h3)

Third Paragraph text
```

###### Emphasis

```
*Italic characters*
_Italic characters_
**bold characters**
__bold characters__
~~strikethrough text~~
```

###### Tables
```
| Column One   | Column Left | Column Right |
|--------------|:------------|-------------:|
| Value 1      | Value One   | 10           |
| Value Two    | Value 2     | 200          |
| Value Three  | Value 3     | 300          |
```

###### Unordered Lists

```
 * Item One
 * Item Two
   * Item Two Sub 1
   * Item Two Sub 2
 * Item Three
```

###### Ordered Lists

```
1. Item 1
2. Item 2
    a. Item 2a
    b. Item 2b
3. Item 3
```

###### Code Examples

Either indent with four spaces or ```

```
    ```
    function test()
    {
        test
    }
    ```
```

With formatting:

```C
    #include <stdio.h>
    int main(void)
    {
         printf(stdout, "Hello World");
    }
```

###### Links

```
[WikiPage](WikiPage/)
[WikiPage](WikiPage/)
[External Site](https://example.net)
[Email Address](mailto:user@example.net)
```

##### Importing from other Wikis (FosWiki)

Trident has a tool named ```twikiexport``` that gathers files on a source host creating an archive
of the files relevant to transfer that Wiki to a Trident installation.

twikiexport currently only supports gathering files from FosWiki.

twikiexport works per Trust Group.

Usage of twikiexport:
```
twikiexport [-v] [-v] <wikiformat:foswiki> <path> <archivename>
```
for example:
```
twikiexport -v -v foswiki /path/to/foswiki/data/ wiki.tar.gz
```

Note that the FosWiki data directory is the one where 'data/Main' (wiki files) and 'pub/Main' directories exist.

The exporter ignores all the System files and only grabs files out of "Main".

The result will be the following message:
```
Done: 378 stored in archive (143 ignored)
```
With the archive created in /tmp/wiki.tar.gz.

One can import this archive in a running Trident system using:
```
tcli wiki import <trustgroup> <format:foswiki> <archivename> <destination-path>
```
for example:
```
tcli wiki import test foswiki /tmp/main.wiki /import/
```

Note that the destination path must be either empty or non-conflicting for
both wiki and file namespaces. We thus recommend setting the destination to
a name like ```/import/``` to avoid any conflicts.

When a conflict happens the import continues skipping the file.
Moving the conflict out of the way and rerunning the import will cause the
skipped file to be imported.

Importing wiki's require sysadmin privileges due to local file access.

###### Wiki archive format

The archive is a standard Tar file that is gzipped (.tar.gz).

It contains two directories, ```wiki``` for Wiki related files and ```files``` for attachments.

#### trident arguments

The following arguments can be provided in ```DAEMON_OPTS``` setting of ```/etc/default/trident```
or as arguments directlya to the tridentd binary.

| Argument          | Description                                                |
|-------------------|------------------------------------------------------------|
| --config          | Configuration File Directory                               |
| --syslog          | Log to syslog (also enabled when --daemonize is given)     |
| --daemonize       | Daemonize tridentd                                         |
| --pidfile         | PID File (useful in combo with daemonize)                  |
| --username        | Username to daemonize into                                 |

Following arguments should only used during testing or development and not on a production system:

| Argument          | Description                                                |
|-------------------|------------------------------------------------------------|
| --debug           | Enable verbose Debug output                                |
| --insecurecookies | Disables the HTTP requirement for cookies                  |
| --disabletwofactor| Disable Two Factor Authentication Check                    |
| --verbosedb       | Log all SQL Queries                                        |

### tcli

```tcli``` is the Trident CLI. It is in effect a simple HTTP client that speaks directly to the API of tridentd.

```tcli``` is pronounced as "Tickly".

As in the Web CLI typing ```tcli -help``` will provide details on available commands.

```
  # tcli -help
Usage of tcli:
  -r	Read an argument from the CLI, useful for passwords
  -server string
    	Server to talk to [env TRIDENT_SERVER] (default "http://localhost:8334")
  -tokenfile string
    	Token to use [env TRIDENT_TOKEN] (default "~/.trident_token")
  -v	Enable verbosity [env TRIDENT_VERBOSE]
```

```tcli``` stores its [JSON Web Token (JWT)](http://jwt.io/) authentication token in ```~/.trident_token```.
A custom token location can be configured using the ```TRIDENT_TOKEN``` environment variable allowing one to keep multiple tokens active (e.g. a normal user and one with sysadmin privileges).
Deleting the token effectively logs one out if one does not have another copy.

One can enable verbosity for ```tcli``` by setting the environment variable
```TRIDENT_VERBOSE=on```
this will show the HTTP URL and the HTTP Response Headers. Disable verbosity with:
```TRIDENT_VERBOSE=off```

### tsetup

```tsetup``` is the Trident Setup command. It helps in configuring the
PostgreSQL database and in upgrading schemas where needed.

This command must be run from the ```postgres``` user account or
another account having unix-loopback PostgreSQL administrative capabilities.

To add a user to the system use (as user ```postgres```):
```shell
tsetup adduser <username> <password>
```

```tsetup -help``` will provide details on available functions and arguments.

```
# tsetup -help
Note: No commands given
Usage: tsetup [<options>...] <cmd> [<arg>...]

 Options:
       --config <dir>
       --verbosedb
       --force-db-destroy
	--version
	--debug
	--help

 Command:
	help
	setup_db
	setup_test_db
	upgrade_db
	cleanup_db
	adduser <username> <password>
	setpassword <username> <password>
	sudo <username> [<cli commands>]
	version

Typically to be run from the 'postgres' account
that has access trusted access to PostgreSQL

The exit code will be zero when no problems are
encountered while non-zero (1 for simple errors,
others depending on the command)

```

Note that the ```username``` must be lowercase letters, and numbers, but the first character may not be a number.

This has to used to add an initial administrative (```sysadmin```) user after
which ```tcli``` or the Trident UI can be used to configure the rest of the system.

### Directories

| Directory                 | Permissions         | Description                                     |
|---------------------------|---------------------|-------------------------------------------------|
| ```/etc/trident/```       | 755 root:root       | Trident Configuration                           |
| ```/usr/share/trident/``` | 755 root:root       | Read-only files (templates, dbschemas, webroot) |
| ```/var/lib/trident/```   | 700 trident:trident | Trident intermediary files                      |


## Installation

Debian packaging is provided (use 'dpkg-buildpackage -b -uc -us' to generate a pacakge). FreeBSD packaging might follow at a later date.

Trident requires PostgreSQL 9.1+ as a database and Postfix for SMTP.

Trident prefers to be run behind Nginx or another HTTP proxy providing HTTPS access.
This avoids having Trident needing to know any SSL keys.

Trident should only be exposed to the outside world using HTTPS.
Thus do not send your details in cleartext HTTP.

One can also run Trident directly from source, see the Development section for more details.

After installing the Trident package one has to edit ```/etc/trident/trident.conf``` and
provide the correct database details (See Database Setup in the next sections).

The package automatically generates new RSA keys used for the [JSON Web Token (JWT)](http://jwt.io/)
that are used for authentication at package installation time ensuring that every
installation has unique JWTs.

### Quick and Dirty

The quick and dirty method of installing Trident:

```bash
apt-get install postgresql nginx postfix
dpkg -i pitchfork-data-VERSION.deb
dpkg -i trident-VERSION.deb
# edit /etc/trident/trident.conf
su - postgres -c "/usr/sbin/tsetup setup_db"
su - postgres -c "/usr/sbin/tsetup adduser USERNAME PASSWORD"
```

Then configure postfix + nginx as per details below.

Actual details about these commands can be found in the more specific sections of this document.

### Database Setup

Depending on having a local or remote database one has to follow the next two sections.

#### Local Database

Either have the PostgreSQL server package installed before installing the Trident package
or make sure that the ```postgres``` user is in the ```trident``` group with:
```shell
adduser postgres trident
```
This allows the postgres user to read ```/etc/trident/trident.conf``` and thus retrieve
the configuration settings needed for the database details.

As this is a local Database we can use peer authentication ```/etc/trident/trident.conf```
should have:
```
	"db_host": "/var/run/postgresql/",
	"db_port": "5432",
	"db_name": "trident",
	"db_user": "trident",
	"db_pass": "",
```

The package performs that if the postgres server is installed at the time the
Trident package is installed.


To create the database and the users that access it one can just run:
```shell
su - postgres -c "/usr/sbin/tsetup setup_db"
```

This will create the ```trident``` PostgreSQL user and the ```trident``` database.

#### Remote Database

Likely easiest is to temporarily install Trident on the
remote server and then run ```tsetup``` like normal and
then after ```setup_db``` to remove the package from the server.

In a nutshell what needs to be done:
 * Create the ```trident``` user on the remote server.
 * Create the ```trident``` database on the remote server.
 * Provide permissions for the user to access the database.
 * Check that ```pg_ha.conf``` contains the correct settings.
 * Ensure reachability of the port (firewalling, listener etc).
 * Configure ```/etc/trident/trident.conf``` correctly.
 * Run ```tsetup``` from the remote server as normal.

As this is a remote Database we must use md5 authentication ```/etc/trident/trident.conf```
should have:
```
	"db_host": "db.trident.example.net",
	"db_port": "5432",
	"db_name": "trident",
	"db_user": "trident",
	"db_pass": "trident",
```

#### Database Security

Trident's ```tsetup``` configures the postgres user to be the owner of the database.
The ```trident``` user only has access to ```SELECT```, ```UPDATE```,
```INSERT``` and ```DELETE``` on the various tables.

#### Database Versions

The ```share/dbschemas/``` directory contains files that describe the
Trident Database Schema.

Each DB change will be codified in a migration. Migrations with assigned
versions will be named ```DB_<version>.psql```. Where the version is the current
version of the database at the time that the given file should be applied.

For example if the current version is 3 and you would like to make a change to
the schema, set the name of the file to ```DB_3.psql```. Your update will set the
database version to 4 as it's last action.

The contents of the migration are in the form of a PSQL file. All activity
should be within a transaction (```BEGIN```, ```COMMIT```).
 his is so that if the update fails in any way, no update will occur.

Database permissions (read: ```GRANTS```) are managed by Trident itself,
these should not be given/revoked in these scripts.


### Webserver setup

Following sections detail how to configure the webserver (HTTP Proxy).

#### Nginx

The directories ```doc/conf/nginx/``` (source) or ```/etc/trident/nginx/```
(when installed from Debian package) contain Nginx configuration includes
for enabling Trident behind a Nginx server.

Drop the following in ```/etc/nginx/conf.d/trident.conf```
and of course, change as needed/wanted.

Items to change at minimum:
 * Hostname [```trident.example.net```]
 * SSL keys, referenced, but files are not included in this example
 * SSL Options
 * HTTP Key Pinning

Definitely verify the configuration with the Qualsys SSL Test:
  https://www.ssllabs.com/

Example ```/etc/nginx/conf.d/trident.conf``` or ```/etc/nginx/sites-available/``` depending on preference.
```nginx
# The Trident Daemon Upstream
include /etc/trident/nginx/trident-upstream.inc;

# Redirect all HTTP (80) traffic to HTTPS (443)
# Trident should only be exposed over HTTPS
server {
	listen 80 default_server;
	listen [::]:80 default_server;

        server_name _default_;

        rewrite ^ https://$host$request_uri permanent;
}

# The HTTPS (443) server that exposed Trident
server {
	listen 443 ssl;
	listen [::]:443 ssl;

	server_name trident.example.net;

	ssl_certificate		trident.crt;
	ssl_certificate_key	trident.key;
	ssl_prefer_server_ciphers on;

	# And other SSL options, recommended:
	# - ssl_dhparam
	# - ssl_protocols
	# - ssl_ciphers
	# See https://cipherli.st/ for details

	# STS header
	add_header Strict-Transport-Security "max-age=31536001";

	# HTTP Key Pinning
	add_header Public-Key-Pins "Public-Key-Pins: max-age=5184000; pin-sha256=\"...\""

	access_log /var/log/nginx/trident-access.log;

	# Include the config for making Trident work
	include /etc/trident/nginx/trident-server.inc;
}
```

### E-mail setup

Trident can be run behind both Postfix and Sendmail or likely other SMTP servers.

We provide details for configuring behind Postfix in the next section.

#### Postfix

Postfix is used as an inbound email setup.
This handles all SMTP specific things like STARTTLS, connection/rate limiting etc.

To make Postfix know about Trident one needs to add to /etc/aliases:
```aliases
trident-handler: "|/usr/sbin/trident-wrapper"
```

and to /etc/postfix/virtual something similar to:
```virtual
example.net                ----------------
mail-handler@example.net   trident-handler
@example.net               trident-handler
```

There are cases where, trident-handler might need to be trident-handler@localhost.

Critical elements of ```/etc/postfix/main.cf```:
```
alias_maps = hash:/etc/aliases
alias_database = hash:/etc/aliases
myhostname = portal.example.net
myorigin = portal.example.net
mydestination = localhost.example.net, localhost, portal.example.net
virtual_maps = hash:/etc/postfix/virtual
```

DOn't forget to update the aliases and virtual databases:

```
postmap /etc/postfix/virtual
newaliases
service postfix reload
```

## Reporting of Security Issues

Please contact project@trident.li directly, we'll deal with problems swiftly
and of course with proper attribution of any issues found.

Note that Trident is a community project, thus we do not offer any bounties
as there is no money involved in the project.

We do very much appreciate responsible disclosure, thus do please contact us
so that any issues can be properly addressed and existing installations updated.

## Problem solving

```tridentd``` per default logs to syslog, which is a good place to check for issues.

### Checklist

The following items should be checked:

 * Verify that the configuration is correct (```/etc/trident/trident.conf```)
 * Verify that the database is setup correctly and is accessible
 * Check syslog

### Debug mode

One can enable verbose debugging by adding the ```--debug``` option to ```DAEMON_OPTS```.

Passwords and Twofactor tokens are masked, other properties are all visible, thus be careful as it might reveal sensitive details.

### Running without HTTPS

For testing only, if a SSL certificate is not present or if one wants to inspect the HTTP traffic
between the client and the server one could forward proxy in nginx from port 80 instead of 443.
In this case add to ```DAEMON_OPTS``` the ```--insecurecookies``` option to disable the HTTPS requirement
for cookies and thus enabling the cookies to be stored in the browser.

## Trident Development

Following are some details about developing and improving Trident.

### Problem reports / Feature Requests

Please file issues on our [GitHub](https://github.com/tridentli/) account.

### Running Trident from source

This is a note for development purposes.

One requires ```golang``` (Go) 1.6+ to be installed for running.

After checking Trident out of git one needs to configure the ```GOPATH``` properly.

We store external libraries in ext/. To update that, use 'make deps' to fetch them.
Then to run  set ```GOPATH``` like this:
```shell
export GOPATH=$(pwd)/ext/_gopath/
```
You can then execute tridentd with:
```shell
src/cmd/tridentd/tridentd.go --disabletwofactor --insecurecookies --config=doc/conf/
```

Of course, don't forget to setup the database as detailed in the Database section.
And one might want to modify the config to adjust to taste.

### Repository Organization

We use git as a distributed revision control system.
Mainline is tracked in the master branch, development happens in sub-branches per feature.

The ```src``` directory contains all the source code, this directory is subdivided in:

| Directory | Description          |
|-----------|----------------------|
| ui        | Web User Interface (UI) components. Each module relates to a related component in ```lib``` and has templates for actually rendingering HTML in ```share/templates/<component>/``` |
| lib       | The per-component libraries of functions. These contain each a component head which links into the Menu system and thus also CLI code. Only code in this directory can execute SQL queries |

The ```share``` directory contains SQL Schemas in ```dbschemas```, the web root in ```webroot```
and Golang templates in ```templates```.

The ```ext``` directory contains external dependencies (update/fetch with the Makefile).

### Coding Style

We use the excellent [vim-go](https://github.com/fatih/vim-go) by Fatih Arslan to enforce formatting.

Emacs users might want to try [go-mode](https://github.com/dominikh/go-mode.el).
This provides a gofmt-before-save hook that can be installed by adding this line to your .emacs file:
```
(add-hook 'before-save-hook #'gofmt-before-save)
```

In general, please verify that whitespace is not affected while committing next to what is being committed is correct and properly tested.

We use a [Model View Controller (MVC)](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller)
approach to the code, though on two levels. The UI is strict-MVC.

### Releasing New Packages

Update ```debian/changelog``` with new date/timestamp and optional details, commit these details.

Then the below described sbuild procedure should be used to ensure vanilla builds that can be reproduced.

#### sbuild

On build.vm.trident.li we have an sbuild instance per instructions on https://wiki.debian.org/sbuild.
apt-cacher-ng and lintian are installed and configured for speed and verification.

~/.sbuildrc options:
```
$run_lintian = 1;
$lintian_opts = ['-i', '-I'];
$purge_session = 'successful';
$purge_build_deps = 'successful';
$purge_build_directory = 'successful';
```
The purge options are to allow debugging unsuccesful builds.

Check out a pristine git clone, then use:
```
sbuild -v -d sid --arch-all
```
to build the package.

If package building fails either check the 'Keeping session' line in the output or use the following to list sessions:
```
schroot -al
```

then start a shell inside the session with:
```
SESSION=sid-amd64-sbuild-<<UID>>
schroot -r -c ${SESSION}
```
providing the correct ```SESSION``` id to the variable from the details above.

The build root can be found inside the chroot in ```/build/```

Afterwards the session can be cleaned with:
```
schroot -e -c ${SESSION}
```

To generally clean sessions use:
```
schroot -e --all-sessions
```

#### Manual packages
Cheanup first:
```
make clean_ext
```
Make sure that all dependencies are up to date:
```
make ext
```

Build a new Debian package:
```
make pac
make vtests
```

Voila, a new package.

### Framework Reference

#### Member States

(XXX: this section is out of date, see state_mon and schema.psql for truth.)

The following member states exist:

| State      | Description                                                    |
|------------|----------------------------------------------------------------|
| nominated  | means somebody has nominated you but you don't know yet.       |
| vetted     | means you've been invouched and you still don't know about it  |
| approved   | will someday mean that admin@ has noted your vettedness and noted the absence of controversy about you. Right now you just go from vetted to approved immediately (criteria is identical.) |
| active     | means you've done everything you need to do and the system is not sending you any annoy-o-grams about your checklist |
| inactive   | means you used to be approved but lost your pgp key or lost a vouch or the vouch criteria was raised and now excludes you |
| blocked    | means somebody negvouched you and there's an investigation.    |
| idle       | means it's been X days (imagine "60") since you either logged into the UI or sent e-mail to one of the lists. |
| soonidle   | means you will soon be "idle" (we send mail warning of this so that you can log into the portal and prevent going idle.) |
| failed     | means your nomination timed out without reaching "vetted"      |

On this table the rows are membership states and the columns are member
capabilities/permissions.

|   state   | can_login | can_see | can_send | can_recv | blocked | hidden  |
|-----------|-----------|---------|----------|----------|---------|-------- |
| nominated | false     | false   | false    | false    | false   | false   |
| vetted    | false     | false   | false    | false    | false   | false   |
| approved  | true      | true    | false    | false    | false   | false   |
| active    | true      | true    | true     | true     | false   | false   |
| inactive  | true      | true    | true     | false    | false   | false   |
| blocked   | false     | false   | false    | false    | true    | true    |
| failed    | false     | false   | false    | false    | false   | true    |
| soonidle  | true      | true    | true     | true     | false   | false   |
| idle      | true      | true    | true     | false    | false   | false   |
| deceased  | false     | false   | false    | false    | true    | false   |

Transitions:

| From      | To         | Description                                                                                     |
|-----------|------------|-------------------------------------------------------------------------------------------------|
| NULL      | nominated  | When somebody nominates you, and mail is sent to vetting@ asking that folks check you out       |
| nominated | vetted     | When a cron job detects that you have enough invouches (target_invouches), and notifies admin@ about this) |
| vetted    | approved   | When an admin notes that there are no negvouches and manually slots you into "approved" status, and you finally hear for the first time that you are a member, or if that's not implemented yet, it's when a cron job notices that you've been vetted and automatically approves you |
| approved  | active     | when a cron job detects that you're approved but that you need to input a pgp (if that's required) and outvouch (if that's required) |
| active    | inactive   | When you lose your pgp key or it's suddenly required, or when you used to have enough invouches (min_invouches) but now you don't.) |
| inactive  | active     | When a cron job detects that you've outvouched and input a pgp key, and notifies by e-mail you about this |
| ANY       | blocked    | When an admin wants the system to camp onto your e-mail address and not allow further state changes or new nominations) |
| active    | soonidle   | When a cron job detects that you have not logged in or sent mail for some significant period of time, and sends you mail telling you that you will soon be idle.) |
| soonidle  | active     | When you log back into the UI or transmit to a mailinglist.                         |
| soonidle  | idle       | When you go a few more days without activity after being told you will soon be idle |
| idle      | active     | Same as soonidle -> active                                                          |

#### Member permissions

| Permission| Description                                                                           |
|-----------|---------------------------------------------------------------------------------------|
| can_login | means your password works at the main web portal UI                                   |
| can_see   | means you can see the membership list and other primary materials, including the wiki |
| can_send  | means you're allowed to send mail to the non-public-access mailing lists              |
| can_recv  | means you can receive mail to the subscription-checkbox mailing lists                 |
| blocked   | means you can't be nominated, nor log in, nor receive or send e-mail, nor be seen     |


## Credits

Following is a short list of credits related to the Trident Project:

 * Trident by [Jeroen Massar](http://jeroen.massar.ch)
 * Original [Ops-Trust Portal code](https://github.com/ops-trust/) by the [Ops-Trust Sysadmins](https://www.ops-trust.net): Paul Vixie, Ben April, Krassimir Tzvetanov, Chris Morrow, John Kristoff and Jeroen Massar.
 * [XKCD 936 'Password Strength' comic](https://xkcd.com/936/) by Randall Munroe of xkcd.com.
 * [Go Crypto library](https://github.com/kless/osutil/) by Jeramey Crawford and Jonas mg.
 * [Go UUID library](https://code.google.com/p/go-uuid/) by Google Inc.
 * [Go JSON Web Token (JWT)](https://github.com/dgrijalva/jwt-go) by Dave Grijalva.
 * [Go PostgreSQL driver 'pq'](https://github.com/lib/pq/) by Blake Mizerany and 'pq' Contributors.
 * [EpicEditor](http://epiceditor.com/) by Oscar Godson.
 * [Black Friday](https://github.com/russross/blackfriday) by Russ Ross
 * [Blue Monday](https://github.com/microcosm-cc/bluemonday) by David Kitchen.

Further details are available in the [Debian package copyright file](debian/copyright) in the source
or to be found in ```/usr/share/doc/trident/copyright``` for installed Trident package on Debian.
