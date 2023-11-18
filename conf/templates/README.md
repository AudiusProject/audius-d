# config templates

collection of common templates that can be used for starting up audius-d

## default.\*.toml templates

These are a collection of defaults that are read in by audius-d and populate the conf structure should these values not already exist. This is to avoid a bunch of runtime logic for defaults. They are embedded into the go binary so do not need to exist at runtime.

There are three types of environments with specific behavior: dev, stage, and prod. CI as the goal of audius-d will extend from dev.
