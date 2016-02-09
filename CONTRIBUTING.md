# How to contribute

We love contributions! Here's how you can submit changes to us:

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

In general you should:
* Make commits of logical units.
* Check for unnecessary whitespace with git diff --check before committing.
* choose sensible branch names
* write expressive PR descriptions


### Creating topic branches

All new work happens from the `develop` branch, so please create your topic branches there.

We encourage you to use the following format for your topic branch names.
The goals of this are:
* Be easy to read (I can tell what the branch is from a list of branches)
* Be tied to story management software
* Be easy to see if it’s a feature branch, bug fix etc
* Be easy to type


When naming your topic branches please use the following format:
  `workflow[/id]/descriptive_name_of_branch`

Here's a full example:
  `fix/jira-1234/publish_err_after_forceDisconnect`

The name should be short, but descriptive. Try to keep it less than five words.

Where workflow is one of:
 * **wip**: work in progress stuff, maybe a longer lived branch
 * **feat**: a new or updated feature
 * **fix**: a bug fix to an existing feature
 * **patch**: the same as fix, except that it targets `master` rather than `develop`

Where is a fully-qualified id of the story in an external system. This is optional if the story is fully managed within the PR.

  E.g `jira-1234`


### What should be in a pull request

At minimum you should always:
* explain what the PR does
* explain how it can be tested, if it requires testing
* link to any relevant stories

There are some good suggestions and a PR Template on [quickleft.com/blog/pull-request-templates-make-code-review-easier/](https://quickleft.com/blog/pull-request-templates-make-code-review-easier)
