# Changelog

```log
0.1.1 - XXXX/XX/XX
refactor(front): mappings page with HTMX+Alpine.js

0.1.0 - 2025/03/24
feat: backups
fix: entrypoint parsing bug on create container

0.0.9 - 2024/11/27
refactor: overview frontend
refactor: login frontend
feat: GET /api/v1/marketplace/
feat: add account quota refresh to cli
fix: account resource usage bug when no containers
fix: improve sys-install verbosity
fix: broken unocss css reset

0.0.8 - 2024/09/30
refactor: use sri for activity records
feat: add system resource identifier
feat: add footer bar with resources usage
feat: add scheduled tasks popover
feat: improve host resource usage entity
fix: add and use blank metrics if container not running
fix: listen to delete custom events on container images
fix: refresh tasks popover on snapshot and archive

0.0.7 - 2024/09/23
feat: container images
feat: create security record for all write operations
feat: add account id to container profile
fix: call htmx.process() when Alpine.js changes the DOM
fix: duplicated ssl map blocks bug
fix: network parsing bug with newer versions of podman

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
