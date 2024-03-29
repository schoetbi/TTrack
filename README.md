![travis-ci](https://travis-ci.com/schoetbi/TTrack.svg?token=49gyUuwE7PY3x5UsqSvU&branch=master)

# TTrack the command line time tracker

TTrack helps you to keep track of the time you spent on your projects and tasks

## Usage

### Begin a task

To start a task enter ``ttrack begin projA/task-123``. This command assumes that the starting time is "NOW". If you want to entere an explicit time you can do so by entering  ``ttrack begin projA/task-123 <date>``.

If you begin the next task the currently running task is automatically ended.

### End all tasks

If you stop working you can enter ``ttrack end`` to end all tasks without starting a new one. Now you can go home :-)

### Listing the tasks

To get a list of all tasks in a time period enter ``ttrack list``. To restrict the time there are two optional arguments. ``ttrack list [from] [end]``.

If you want to see the list only for todays tasks there is a shortcut: ``ttrack list t`` (t for today). ``ttrack list y`` show the tasks starting from yesterday.

### Creating reports

The report is created by entering ``ttrack report [<from>] [<to>]``.

You also can use two options:

- To group daily: ``ttrack report 11-1-2011 --daily``
- To group by project: ``ttrack report 11-1-2011 --project``

Both options can be combined to get a report of the time you spent on each project starting from 1st of November till now grouped by day.

## Date handling

For all dates and times TTrack tries to detect the format how you enter the date and time. 

For example TTrack detects ``11-03-2011`` as the 3rd November 2011. If you enter ``03.11.2001`` the date is also the third of November 2011 but in the German format.

### Enter time without a date

You can also ommit the date and only enter ``10:40``. In this case the date is assumed to be today.

### Time shortcodes

Instead of writing the full date and time you can enter a shortcode:

|Shortcode|Meaning|
| ------------- | ------------- |
| t | Today 00:00 |
| y | Yesterday 00:00|
| m | Start of current month |
| lm | Start of Last Month|


## Ideas behind TTrack

At the same time I wanted to keep track of both the time I spent on each project and also the
time I spent on each task within this project.

The tasks are entered in the form  ``<project>/<task>``

After all the work I like to have two kinds of report in a time period e.g. a month.

1. A report where I can see the time I spent on each project
2. A report for the time on each task

Both report types have the option to aggregate the work daily.

## Installation

You can find prebuild binaries for linux and windows under [Releases](https://github.com/schoetbi/TTrack/releases).

But you can also build ttrack yourself from source. To do so clone or download this repository and build it with ``go build``. Then install it with ``go install``.

## How data is stored

The data is stored in a [SQLITE](https://www.sqlite.org) database. The path of this file depends on the system you are using.

- Windows: ``C:\users\<username>\ttrack\ttrack.db``
- Linux: ``/home/<username>/ttrack/ttrack.db``
- Mac: ``No idea if someone knows the path please let me know``

## Some tips

- Put ttrack in your path. This way you can access it easily from every command prompt
- Under Windows you can execute ``ttrack end`` on system shutdown to finish all started tasks.

