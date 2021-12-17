# JSextractor

Gather all JavaScript code from html webpages from command line.
`jse` take input from `stdin`, search for JavaScript in `<script>` tag (in `src` attribute and code) and in event handlers.

* [Gather JavaScript](#basic-usage)
* [Gather JavaScript for further analysis on it](#output-js-code-into-file)
* [Gather JavaScript from specific source](#restrict-harvest)

## Basic usage

* Retrieve Javascript from an offline HTML file:
```shell
cat [html_file] | jse
```

* Alternatively, you could gather it following a curl command:
```shell
curl [url] | jse
```

## Output js code into file

This could be useful for performing further actions later on JavaScript  like scanning it or beautifying it. But in this case, the output **must** be a valid script:
`cat [html_file] | jse -gather-src -d [domain_of_html]`

We use  `-gather-src` to retrieve code from `src` attribute (fetching the code). Otherwise it would return only the URL corresponding to the `src` value and thus making the output a non-valid JavaScript script.
When we use `-gather-src` we must also define the domain from which we got the html page, (this is used to fetch script hosting by te same site *e.g.* `src=/this/is/a/path.js`)

Also, all informative logs (line and source) are output to `stderr` to keep only js code in `stdout`

## Restrict harvest

`jse` search for js code from 3 sources by default. Sometimes, you only want code from a specific source. In this case you could disable other source gathering:
* `-ds`: don't look for js in src attribute of `<script>` tag
* `-de`: don't look for js from event handler attributes
* `-dt`: don't look for inline js of `<script>` tag")


## Install

### curl

```
curl -lO -L https://github.com/ariary/JSextractor/releases/latest/download/jse
```

### from code source

```shell
git clone https://github.com/ariary/JSextractor.git
make before.build
make build.jse
#install it in your $PATH
mv jse $HOME/.local/bin/
```


## Enhancement 🛣️

* Line counter is not working perfectly and must be improved