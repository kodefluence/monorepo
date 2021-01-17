# Monorepo

Monorepo is a framework to reduce complexcity by implementing single source of truth structures which can be adventage because of several of things. In codelfuence monorepo we write all daily tools we need to brew our code.

## Ease of code reuse

Similar functionality or communication protocols can be abstracted into shared libraries and directly included by projects, without the need of a dependency package manager.


## Simplified dependency management

n a multiple repository environment where multiple projects depend on a third-party dependency, that dependency might be downloaded or built multiple times. In a monorepo the build can be easily optimized, as referenced dependencies all exist in the same codebase.

## Atomic commits

When projects that work together are contained in separate repositories, releases need to sync which versions of one project work with the other. And in large enough projects, managing compatible versions between dependencies can become dependency hell. In a monorepo this problem can be negated, since developers may change multiple projects atomically.

## Large-scale code refactoring

Since developers have access to the entire project, refactors can ensure that every piece of the project continues to function after a refactor.

## Collaboration across teams

In a monorepo that uses source dependencies (dependencies that are compiled from source), teams can improve projects being worked on by other teams. This leads to flexible code ownership.
