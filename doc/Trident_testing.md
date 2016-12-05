# Trident Testing

## Unit Testing

TBD

## CLI Testing

The CLI elements can be tested by the script found at: 

```
$base/doc/generatetestsetup.sh
```

This script uses ```tcli``` and ```tsetup``` to stage and populate the database with test data.

## Web Testing

Testing via the web interface is done with the [http://www.seleniumhq.org/projects/ide/](Selenium IDE suite).
It is a FireFox plugin that will handle most of the testing automatically. 

The following script is used to stage the database for testing: 

```
service trident stop
su postgres -c "/usr/sbin/tsetup --force-db-destroy setup_test_db"
service trident start
```

It is important to setup the uploaded files:
```
ln -s $base/test/files /tmp/trident-selenium
```

The DB test data can be found at ```$base/share/dbschemas/test_data.psql```
Test Selenium files are located at ```$base/test/```

 * ```$base/test/files``` - are files that get uploaded.
 * ```$base/test/case_*``` - are files that contain test cases.
 * ```$base/test/issue_*``` - are test cases based on open issues.
 * ```$base/test/selenium_suite``` - Is an ordered set of case_ and issue_ files

When starting Selenium IDE in Firefox, Use the ```Open Test Suite``` function to load ```$base/test/selenium_suite```.

When the suite is run ALL case_* cases should be green. issue_* cases that have not yet been resolved may be red.

