# Airnow GO Tech Test

---
### Setup
#### Requirements

* go v1.19
* make

#### Running
To compile `make build` and run with `bin/crawler -t [target url] -ti [timeout] -d [depth]`

#### Reviewing

All business logic is kept within the `internal` package.

#### CLI Flags
Usage of bin/crawler:

--d --depth int (int) --- specified nesting level for traversal (default 0)

--t --target (string) --- the url to crawl (default "https://crawler-test.com")

--ti --timeout (int) --- overall completion timeout in milliseconds (default 5000)

#### Example Usage
`./bin/crawler -t https://crawler-test.com -d 3 -ti 5000`

#### Testing

Run the test suite with `make test`

Generate coverage with `make cov`

#### Cleanup

Run `make clean` to remove coverage files and binaries

#### Makefile

`cov`       - Generate code coverage and run in browser

`test`      - Runs the full test suite

`clean`     - Removes coverage.out/coverage.out.tmp files, removes the `bin/crawler` binary

`build`     - Compiles the binary into bin/crawler

---
#### Decisions

I opted to store the links in a tree structure. This made sense to me as formatting the output of the links came for free:
the tree already has a parent-child relationship on insert, therefore it's trivial to print links with the correct indentation level.

I was conscious that I didn't have 100% code completion, I stopped at the point that it'd take too long for the
bounds of this exercise to chase the extra ~10% of coverage in the utils file. The pages that would induce these code branches
was hard to find, so I decided to move on.

Another decision was to not implement my own web server to run tests against. I've gone with using crawler-test.com, it seemed
to have a comprehensive amount of edge cases for free which allowed me to test my code more thoroughly in the allowed
time, rather than if I'd have tried to think and implement all of these edge cases in a stand-alone docker container.

I decided to remove duplicate nodes in the tree too - I thought this would be expected of the exercise as otherwise
it'd get stuck in an endless loop if one page linked back to it's parent, which is quite common.

The insertion method on the tree structure is front-loaded with checking that the child to be added is unique within the
tree. This was a conscious decision. I'm sure there'd be a way to optimise this, I don't think it's the most performant
implementation, however I thought optimising this for performance was overkill.

I decided on handling relative links and external links as I felt this was within the bounds of the project.
I also decided to discard anything that wasn't a 200 success code. I could have accepted just a http.Success, however I
felt that I'd also be swamped more edge cases to handle. I could be wrong on that, but that was my intuition when
implementing the get() function.

Layout-wise, I feel I've gone with a fairly common layout - `cmd` just holding main and the `internal` package takes
care of the actual logic and unit tests. I haven't extensively written integration tests, just a couple that I'm happy
that it works as expected. As a side-note, this could have been structured as a few .go files in the root of the project,
however as it's a CLI I felt it semantic to have an `internal` package for code that isn't meant for external parties to use.

---
#### Areas of Improvement

A less fragile test harness would be probably number one. I really don't like the idea that my tests are directly correlated
to my internet connection.

Concurrency, currently this is running in a single thread and the performance is slow. It's bottle-necked by internet 
connection and having to wait for a single query to return before starting the next. I had toyed around with making the 
Walk method concurrent and spawning recursive function calls as separate goroutines, however, again I thought this was 
also overkill for the project and would take me a bit more time to properly think about how to implement this as a 
concurrent and recursive algorithm properly without deadlocking.

It would have been nice to have some more thorough documentation, however as this is a CLI and not a library I felt it
was acceptable to lightly sprinkle some documentation instead of having lots of heavy explanations.
