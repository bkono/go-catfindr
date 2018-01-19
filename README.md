# go-catfindr

Go variation of [catfindr](github.com/bkono/catfindr) for lambda performance comparisons

## What's it do

See the matching README section of [catfindr](github.com/bkono/catfindr)

## Okay... but why a go version?

Two reasons:

1. I'm a fan of [code katas](http://codekata.com/). Writing brute force algos just isn't something that comes up in day to day work, so every once in awhile I like to take a problem I sort-of know, and attempt it from scratch all over again. Preferably in a different language.
2. A discussion came up that lead to a need for comparing Go & Java performance on lambda. This was built specifically for that purpose, as a completely arbitrary but kinda-sorta real world use case. Given it is an http endpoint, processes uploaded files, does some algo-y stuff, and encodes json, it √ all the boxes I needed.

## Since it was built to compare, what's the verdict? Give me the datas!

Let me start by saying neither is even slightly optimized. They were both written in one shot, kata style. Well, that isn't entirely fair, I did write all of the Go version in main.go at first. When I moved to finishing the part where I actually returned a json result, I decided to refactor for my own sanity.

Having said that, it actually moves this to a closer representation of a completely arbitrary business application. The results were *interesting*. Specifically because the java code performed like two different applications over the course of two test runs, done on consecutive days.

Data is as follows:

**go cold start**

```
INFO REPORT RequestId: ee78a98a-fd52-11e7-8036-21e4d9155ca4    Duration: 331.99 ms    Billed Duration: 400 ms     Memory Size: 512 MB    Max Memory Used: 40 MB
```

**go subsequent requests**

```
INFO REPORT RequestId: 0d4751e2-fd53-11e7-8036-21e4d9155ca4    Duration: 27.66 ms    Billed Duration: 100 ms     Memory Size: 512 MB    Max Memory Used: 40 M

.... ~1 min pause time ...

INFO REPORT RequestId: 3b0b048b-fd53-11e7-97d6-6127e7d46868    Duration: 118.83 ms    Billed Duration: 200 ms     Memory Size: 512 MB    Max Memory Used: 40 MB

INFO REPORT RequestId: 3e9fc3af-fd53-11e7-9989-e77a5c56ff18    Duration: 34.92 ms    Billed Duration: 100 ms     Memory Size: 512 MB    Max Memory Used: 40 MB
```

**java cold start, day 1**

```
INFO REPORT RequestId: e073f47e-fcc2-11e7-987f-8d0346a30a23    Duration: 14316.41 ms    Billed Duration: 14400 ms     Memory Size: 1536 MB    Max Memory Used: 378 MB
```

**java subsequent requests, day 1**

```
INFO REPORT RequestId: eca48867-fcc2-11e7-8803-b5dedd0e792b    Duration: 628.02 ms    Billed Duration: 700 ms     Memory Size: 1536 MB    Max Memory Used: 380 MB


INFO REPORT RequestId: f307d09c-fcc2-11e7-8803-b5dedd0e792b    Duration: 574.80 ms    Billed Duration: 600 ms     Memory Size: 1536 MB    Max Memory Used: 382 MB
```

**java cold start, day 2**

```
INFO REPORT RequestId: b15b68a2-fd53-11e7-93dc-9d3123d7ce00    Duration: 14344.75 ms    Billed Duration: 14400 ms     Memory Size: 1536 MB    Max Memory Used: 403 MB
```

**java subsequent requests, day 2** *note, this is where things get interesting*

```
INFO REPORT RequestId: cc9355e4-fd53-11e7-829a-8f40285afeee    Duration: 52.41 ms    Billed Duration: 100 ms     Memory Size: 1536 MB    Max Memory Used: 403 MB

INFO REPORT RequestId: ce90f9d3-fd53-11e7-8b95-8f9345cc15c4    Duration: 22.35 ms    Billed Duration: 100 ms     Memory Size: 1536 MB    Max Memory Used: 403 MB

... ~1 min pause time ...

INFO REPORT RequestId: 349e508b-fd54-11e7-81b7-fb49f5a7a7e9    Duration: 52.12 ms    Billed Duration: 100 ms     Memory Size: 1536 MB    Max Memory Used: 403 MB

```

As you can see in the above data, the go version booted very quickly, and held fairly consistent times under subsequent requests, with an interesting bump in response time whenever a minute or so of pause time existed between requests.

The java side had far more variance. On both days the cold start time was nothing short of atrocious, taking almost a solid 15s to launch (equating to ~36x the cold start cost of the Go variation). The fascinating part comes in when comparing day 1 and day 2 subsequent requests. On the first day, the process consumed a max memory allocation of 382 MB, while triggering a billed duration between 600 and 700 ms, effectively 6/7x the cost for subsequent Go requests. On day 2, with zero deployment or configuration changes, the first cold hit caused slightly more ram consumption, at 403 MB. Subsequent requests, however, managed to stay below the 100 ms cutoff, and were very competitive with the Go process. This will need further investigation to try to uncover the reasoning behind such a dramatic variation.

## Caveats

As with most benchmarks, this approach had more holes than swiss cheese. Despite that, it presents a decent starting point for getting an off-the-cuff feel for potential benefits and issues surrounding both technologies usage in lambda. **tldr:** Go consumes significantly less memory (no surprise), has a very respectable cold start time, and livable response times. Java showed horrendous launch times, significantly higher RAM usage, and wildly varied request times between test days. More investigation is needed before making any more WAGs for a conclusion.