## Syncing with Coroot Upstream

### 1. Fetch the latest changes from upstream

      git remote add upstream https://github.com/coroot/coroot.git
      git fetch upstream

### 2. Switch to your develop branch

      git checkout develop

### 3. Ensure your branch is up-to-date with origin develop

      git pull origin develop

### 4. Create a temporary branch to prepare upstream changes

      git checkout -b temp-upstream-sync upstream/main

### 5. Squash all upstream changes into a single commit

      git reset --soft upstream/main
      git commit -m "Squash upstream/main changes into a single commit"

### 6. Switch back to develop

      git checkout develop

### 7. Merge the single commit from temp-upstream-sync

      git merge --no-ff temp-upstream-sync -m "Merge squashed upstream changes into develop"

### 8. Push the updated develop branch

      git push origin develop

### 9. Cleanup the temporary branch

      git branch -d temp-upstream-sync

```

```
