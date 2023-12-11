# This script should be run from the Makefile only.
set -eo pipefail

if [[ -n $(git status -s) ]]; then
  echo "There are uncommitted changes in the repository."
  exit 1
fi

# first release edge case
if [[ "$(gh release list --exclude-drafts)" == "" ]]; then
  start_commit="$(git rev-list --max-parents=0 HEAD)"
else
  old_version=$(gh release view --json tagName | jq -r ".tagName")
  start_commit=$(git show-ref --hash "refs/tags/$old_version")
fi

changelog="$(mktemp)"
printf "Full Changelog:\n" >> "$changelog"
git log --pretty='[%h] - %s' --abbrev-commit "$start_commit..HEAD" -- $directories | tee -a "$changelog"

release_version="$(jq -r .version .version.json)"
if gh release view "@audius-ctl/$release_version" &> /dev/null; then
  echo "Release $release_version already exists."
  echo "Please update the version in '.version.json' and push your changes."
  exit 1
fi
 
gh release create --target "$(git rev-parse HEAD)" -F "$changelog" "@audius-ctl/$release_version" bin/audius-ctl*
