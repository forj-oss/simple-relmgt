# Release Management

simple-relmgt is a GO program to help managing code releasing within Jenkins for Forjj projects.

The idea is to have a simple implementation of releasing a project code following a simple development process which can be re-used easily 

## Simple development release process

Currently `forj-oss` organization follow a basic github development process:

1. Fork the project
2. develop and test locally
3. submit a PR which start a Jenkins build 
4. Review and merge to Master

Usually, the master branch will deploy to github `latest` release

But for more stable product versioning, we need to introduce product version.
This was made manually by creating a github release, define the project version as decribed in the code, build and push the binary to github release.

Thanks to `simple-relmgt`, this extends this functionnality to Jenkins by updating few files, (VERSION & release notes) to deliver officially a final version.

## The release process

Releasing a version of the code is basically limited to :

- Update a Version file
- Create a versioned release note with proper information

Those files are updated/created and push using the basic development process.

`simple-relmgt` is called in the CI system when a PR is created and when we merge it.

Following are the steps used to release a project.

1. Defining a new release version

    Depending on the project code used, the version file can be `VERSION` or `version.go` or anything else.

    `simple-relmgt` will read that file and extract thanks to rules the latest version to deliver and release.

    The version must respect the [semver convention](https://semver.org/) `simple-relmgt` will test it.

    This file can be pushed and merged to the master with no more effect, except verifying version syntax.

2. Create a versioned release note file

    A release note is a Markdown file named as `<repo-root>/releases/release-<version>.md`. `simple-relmgt` can support different path or file name.
    But at least, the file name must contains the new Version to release.

    The release file has not real format to respect, except Markdown.

    It supports the [changelog convention](https://keepachangelog.com/en/1.0.0/)

    This file can be pushed and merged to master with no more effect. 
    But we can add a special effect to create the github release in *draft* state and which will be a copy of the release-note file.
    In that case, we call `simple-relmgt draft-it`

3. Prepare the release to be delivered.

    This step is the first one that influence the project Jenkins pipeline.

    When the versioned release note file defines a *publish date*, the PR will be considered as a *Releasing PR*

    `simple-relmgt` will detect the date, following a rule (regexp). The date must be formated as `YYYY/MM/DD`

    If the date is defined in the past, the PR can be merged, anytime. If the date is future, the *releasing PR* Jenkins pipeline will fail and a merge should not happen.

4. Delivering the release

    This is the last step.

    When the *Releasing PR* is valid (ie Jenkins pipeline succeed), the PR can be merged.
    Merging this PR will start the final release process.

    From that point, the Project jenkins pipeline will automatically deliver the release.

    The pipeline itself will do the following:

   1. create the git tag and push it. `simple-relmgt` can configure the upstream to use to push the code to. The pipeline call `simple-relmgt tag-it`
   2. create/update github release in draft state. The release description will be a copy of the *versioned release note* file.
   3. build and test the project
   4. if success, the pipeline will call `simple-relmgt release-it`. So, the github release will be updated to `released` or `pre-released`
      The release state is defined by another regexp detected in the release note. `simple-relmgt` will detect it. If not found, by default, it will set `released` except if we change the default use case.
    additionally, `simple-relmgt` will push some artifacts to the github release.

## Possible futur

For now, we thought this simple automated release process, will be good in most cases. But we may need to enhance it with [github deployment API](https://developer.github.com/v3/repos/deployments/).

We can also introduce, release branches, as today, we are releasing only from master branch

`simple-relmgt` currently uses [`github-release`](https://github.com/aktau/github-release) to manage github releases. We can replace it by running a GO plugin which will do the same.

If there is other ideas, comments, feel free to open a discussion issue.

Feel free to contribute to the project code as well, with a PR. I'm not necessarily a well GO developer. So, any hint, suggestions, proposal or code updates are welcome!

Forj Team