# Changelog

```log
0.0.7 - 2024/XX/XX
feat: container images
fix: call htmx.process() when Alpine.js changes the DOM
fix: duplicated ssl map blocks bug

0.0.6 - 2024/08/03
refactor: container profile frontend with HTMX and Alpine.js
feat: add storage performance units
fix!: renamed "cpuCores" => "millicores"
fix!: renamed "diskBytes" => "storageBytes"
fix!: renamed "inodes" => "storageInodes"
fix!: replaced uint to uint64 for profile, mapping and target's ids

0.0.5 - 2024/07/18
refactor: security events to activity records
feat: add zerolog as slog handler
fix: remove http custom logger
fix: remove orphan mappings files

0.0.4 - 2024/07/12
feat: add trail db service
feat: add security events
feat: limit failed logins attempts per ip
fix: rename all "get"(s) to "read"(s)

0.0.3 - 2024/06/28
feat: add async tasks
feat: make container deployment async

0.0.2 - 2024/06/21
refactor: use systemd to manage containers statuses
feat: sort registry result by pull count
feat: add transactional update timer overwrite
fix: DockerHubImageFactoryError when search for "rocket.chat"
fix: stop using profile id as prefix for container name
chore: update front
chore: update dependencies

0.0.1 - 2024/05/23
feat: initial release
```
