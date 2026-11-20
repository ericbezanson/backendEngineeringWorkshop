# Chapter 1: Go Basics + WikiData Stats

### What is go and why do we use it?
Go is an open-source, statically typed, compiled programming language. It was designed to combine the safety of languages like C++ with the readability of languages like Python.

#### Standout Features:

##### Built-In Concurrency: 
Very easy to run multiple tasks at once using "goroutines" (aka lightweight threads).

##### Massive community support/Standard Library:
A lot of out-of-the-box provided packages for networking, JSON parsing, cryptography, etc. All of Chapter 1, for example, was built using only standard library packages.. 

##### Performance:
Go compiles directly to machine code, which makes it extremely fast. It compiles to a single standalone binary file, which means it can be deployed easily to a server or Docker container without needing to install a runtime environment.

### How chapter 1 main.go functions:
Common Go convention is to have main.go be the entry point of an application. With our task of creating a stats readout for changes coming from the Wikimedia /stream/recentchange API being relatively small, I decided to keep all logic for Chapter 1, outside of testing, inside main.go.

With the exception of the go.mod configuration file, which effectively defines a Go module, I decided to keep this global and use it for all future chapters. I think this makes sense as it will be easier to manage one single dependency and adds the convenience of being able to run global commands across all chapters (e.g., tests with wildcard opperator "go test -race ./...").

The app itself immediately calls the fetchWikiData() function inside of main() when it runs, in a new goroutine using the go keyword. This pushes a continuous stream-fetching loop process into the background.

While in the background, fetchWikiData() will open a connection to Wikipedia's Server-Sent Events (SSE) stream, using bufio.Scanner (buffered I/O, reads data in small chunks) to wait for incoming lines of text. When data arrives, it strips the prefix, unmarshals (takes raw JSON, and byte-by-byte maps it directly into a go struct) the incoming JSON into a Go struct, and updates the statistics.

We use Mutex for concurrency control (added via "sync" and applied to AppStats). Because the background thread is constantly updating the "stats" in memory on one thread, the web server will try and read that same memory when the page is refreshed, leading to a potential race condition if that happens while new stats are being written to it. stats.mu.Lock() and Unlock() ensure that only one thread can touch the stats variable at a time.

Finally, when the /stats endpoint is hit, the statsHandler temporarily locks the data and counts the lengths of the in-memory maps (data structure, similar to Objects) to get distinct totals. It then translates it into JSON that the browser can read, understand, and display.

For testing we just make a fake request in a fake browser (httptest.NewRecorder()) and instead of starting the server we invoke the statsHandler function, hand it the fake request and validate the status code, header and paylod. 