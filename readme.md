## Syncing with Coroot Upstream

1. Add the upstream remote (if not already added):
   ```
   git remote add upstream https://github.com/coroot/coroot.git
   ```
2. Fetch the latest changes:
   ```
   git fetch upstream
   ```
3. Create a temporary branch for syncing:
   ```
   git checkout -b upstream-sync
   ```
4. Rebase or merge the changes:
   ```
   git rebase upstream/main
   ```
5. Merge into the `develop` branch:
   ```
   git checkout develop
   git merge upstream-sync
   ```
6. Push to the remote repository:
   ```
   git push origin develop
   ```
7. Delete the temporary branch:
   ```
   git branch -d upstream-sync
   ```
