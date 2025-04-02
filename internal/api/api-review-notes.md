# mwartell review notes for api.go

- There is a clever use of go:embed here to not have to ship a separate file. Also,
  thanks for teaching me that go:embed exists.

- Since you already have the wordlist.txt in memory, you don't need to pretend it is a
  file; a benchmark show that this takes 5.7ms to return one random word aand makes 96k
  memory allocations.

- I cannot figure why you are using a (degenerate) reservoir sampling algorithm for a
  sample of size 1. As you've written it, it is equivalent to just picking a random line
  from the file.

- there is no need to build a new random Source, the default Source is fine for casual
  randomness like you use here. Also, seeding need only be done at most once per
  application, not every time you want a random number.

- GetAcceptableWord is a poor name for a function that does nothing to Accept a word
  from the corpus. It also has a time-delay that should be much higher up in the
  application where giving time for the user to see feedback might be needed; this
  package has no concept of a user or a screen.

A simplified version of RandomWord takes only 2.6 nanoseconds to return a random word
and because of Go slice semantics, the elements of `words` point directoly to the
original string. The splitWords operation takes ~1.5ms but is only done once at startup.
