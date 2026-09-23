# bootstrap-library

## DEV Activities

```bash
TAG="0.1.1"
git checkout master
git pull --ff-only origin master
git tag -a v${TAG} -m "Release v${TAG}$"
git push origin v${TAG}
```