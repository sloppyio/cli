# Changelog

## 1.13.1 (2018-12-17)

* fix(stats): show stats from projects with dashes in name

## 1.13.0 (2018-11-26)

* feat: support tcp ports, via port_mappings[].service_port

## 1.12.0 (2018-06-26)

* feat: Add Load Balancer properties to domain
* fix: update docker-compose support, to always set ssl:true when configuring a domain
* fix: sloppy-start now properly uses the `-p`/`-project` argument

## 1.11.1 (2018-04-18)

+ feat(docker-compose): support for deploying from docker-compose.yml files

## 1.10.0 (2018-02-28)

+ feat(docker-login-keystores): support reading native keystores for docker-login, see #34

## 1.9.0 (2018-01-19)
+ feat(log-filter): the logs can be filtered by date range

## 1.8.0 (2017-06-07)
+ feat(logging): configure external logging service

## 1.7.0 (2017-03-09)
+ chore(windows): add version-information resource
+ feat(projects): introduce force flag

## 1.6.0 (2017-01-27)
+ feat(update): add changelog uri
+ fix(version): set 5s timeout for http client
+ feat(decoder): allow users to escape variables
+ feat(change): if project does not exist, create it [GH-5]
+ feat(change): remove PROJECT argument [GH-5]

## 1.5.2 (2016-11-22)
+ fix(validation): Reject additional properties in sloppy.json
+ fix(logs): Disable timeout for reading log body

## 1.5.1 (2016-10-25)
+ chore(flags): correct pre-release label
+ docs(readme): Add blank lines between list items
