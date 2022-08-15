# Airnow GO Tech Test

## Submission Guideline
Please fork this repository. When you have finished your work please provide a link to your own repository. You are also welcome to host this on Github if you don't have a Bitbucket account.

## Task
We would like you to prepare a simple command line tool for parsing nested links.

The application must accept the target source and nesting level as input parameters. The application must take the original page, find all the links and follow them until we get to the specified nesting level. One parent page is one level.  The application should print the list with children links tabbed by nested level.
Please use [goquery](https://github.com/PuerkitoBio/goquery) for getting links from the target URL. 

This task should take you somewhere from a few hours to a day. Please don't take any more time than this, we are more interested in your approach than completing the task.

## Outcomes
Some things we would like to see:

- Readme with a instruction how to install and run the application.
- Frequent commits.
- Few unit tests for the core logic.
