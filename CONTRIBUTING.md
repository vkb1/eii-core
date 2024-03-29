# Contributing to open-edge-insights

## Pull Requests

Everybody can propose a pull request (PR) but only the
core-maintainers of the project can merge it.

### Forking the github project(s) for contributing

Please refer to the `The Github Flow` at https://git-scm.com/book/en/v2/GitHub-Contributing-to-a-Project
for contributing to the open-edge-insights projects. In this link, one can find a detailed explanation on 
how to fork the original projects, create pull request in original project to merge your forked project branch
and how to keep up with the original project (aka upstream remote)

### Commit and Pull Requests Hygiene

The following points will be looked at during the review.

1. **master** is the default branch, it is advised always to work on a new
   feature/bug fix developer branch by keeping your local **master** branch
   untouched. The below convention while naming our branches can be followed:
   * Feature branches - feature/idsid/feature_name
   * Bugfix branch - bugfix/idsid/jira_bug_ids

   More details on best git branching model
   [https://nvie.com/posts/a-successful-git-branching-model/](https://nvie.com/posts/a-successful-git-branching-model/)

2. Once your branch is created by following the above convention
   (`git checkout -b <branch_name>`) and it's ready with the changes,
   run the below commands:
   * `git commit -s` - commit the message just like the old fashioned
     way with one liner and the body. Commit message
     format below:

     ```sh
      <module_name>: one liner about the commit

      Additional details about the commit message

      Signed-off by: abc <abc@xyz.com>
     ```

     Here, `module_name` loosely refers to the folder name where the major
     change is been introduced.

   * `git push origin <branch_name>` - pushes your changes to the remote
     branch name (orgin/<branch_name>)
   * Refer the link https://git-scm.com/book/en/v2/GitHub-Contributing-to-a-Project on 
     creating the pull request into the original project
   * If one notices any conflicts after creating a pull request, just
     apply the latest master to your <branch_name> by running commands:
      * `git checkout master` - with your current branch being <branch_name>
      * `git pull` - pull latest remote master commits to your local master
        branch
      * `git checkout <branch_name>` - Get back to your branch
      * `git rebase master` - rebase your branch with latest master

3. After addressing code review comments, do a `git commit --amend` to amend
   the commit message and do a `git push -f origin <branch_name>` to 
   forcefully push your changes. This is needed as we want to maintain a single commit.

### Minimum requirements for a PR

The following are the minimal requirements that every PR should meet.

- **Pass Continuous Integration (pre-merge build)**: Every PR has to pass our CI

### Review Process

The reviewers may be busy. If they take long time to respond, feel free to
trigger them by additional comments in the PR thread.

It is the job of the developer that posts the PR to rebase the PR on the target
branch when the two diverge.

Below are some additional stuff that developers should adhere to:

* Please batch all your comments by adding to `review` by clicking on `Start a
  review` to start and add all further comments by clicking on `Add review comment`.
  Once done adding all the review comments, click on `Finish your review` and
  accordingly choose the appropriate option in the popup window to submit.

* In a pull request (i.e., on a feature/bugfix branch) one can have as many
  commits as possible. If all the commits are related to a single feature (eg:
  one has addressed review comments or fixed something etc.,), just ensure the
  `Title` and `Description` of the pull-request is up-to-date with respect to
  what is been added/modified/removed etc., This way, the maintainer
  during merging your pull request will squash all your commits into one and
  modify the squashed commit message to have both the `Title` and `Description`
  of the pull request.

  If pull-request has a single commit, then automatically it gets picked up in
  your pull-request description.

* Whenever one sees any failure in the pre-merge build, just check if “rebasing”
  your pull-request branch with latest master and re-push can fix the issue.
  If this is till failing, request the reviewers/maintainers to point out the 
  required compliance needed for passing the CI build for your PR.

### Merge code

Merging will be done as a 'Merge and Rebase' to avoid extraneous merge commits on the master branch
* Special case 'Squash and Merge' will be used sparingly and only as needed
  
