<div align="center">
 <h1> JSextractor</h1>  
<img src="https://github.com/ariary/JSextractor/blob/main/img/jse-small.png">
 <h4> JSextractor is a command line utility to gather JavaScript from webpages </h4>
 <p> Gather all JavaScript code of html webpages from command line by supplying them in stdin. JSextractor search for JavaScript in <code><script></code> tag (<code>src</code>  attribute or inline code) and in event handlers. </p>
  <p><strong><code>{ <a href="#install">Install it</a> ; <a href="#usage">Use it</a> } </code></strong></p>
</div>

## TL;DR
Fetch all js code from webpage:

```shell
export URL=[endpoint]
curl $URL -s | jse -gather-src -u $URL
```

## Demo
|⬇️ TUI demo|
|:---:| 
|![demo](https://github.com/ariary/JSextractor/blob/main/img/jse-tui.gif)|
|➡️ [command line demo](https://github.com/ariary/JSextractor/blob/main/img/jse.gif)|

## Usage
* [Gather JavaScript](#basic-usage)
* [Gather JavaScript for further analysis on it](#output-js-code-into-file)
* [Gather JavaScript from specific source](#restrict-harvest)
* [Launch Terminal based UI (TUI)](#launch-terminal-based-ui-tui)

### Basic usage

* Retrieve Javascript from an offline HTML file:
```shell
jse < [html_file]
```

* Alternatively, you could gather it following a curl command:
```shell
curl -s [url] | jse
```

### Output js code into file

This could be useful to perform further actions later on JavaScript, like scanning it or beautifying it. But in this case, the output **must** be a valid script:
```
curl [url] | jse -gather-src -u [url] > output.js
```
We use  `-gather-src` to retrieve code from `src` attribute (fetching the code). Otherwise it would return only the URL corresponding to the `src` value and thus making the output a non-valid JavaScript code.
When we use `-gather-src` we must also define the url from which we got the html page (this is used to fetch script hosting by te same site *e.g.* `src=/this/is/a/path.js`).

Also, all informative logs (line and source) are output to `stderr` to keep only js code in `stdout`

### Restrict harvest

`jse` search for js code from 3 sources by default. Sometimes, you only want code from a specific source. In this case you can disable other source gathering:
* `-ds`: don't look for js in src attribute of `<script>` tag
* `-de`: don't look for js from event handler attributes
* `-dt`: don't look for inline js of `<script>` tag

### Launch Terminal based UI (TUI)
`jse` comes w/ a TUI which makes JS browsing/gathering/analysis simpler.
 
From the TUI you browse all the scripts found, gather javascript from `src`attributes (`[CTRL] + G`). You can also change the url (`[CTRL] + U`);
To launch it add `-tui` falg (It is recommended to specify url also):
```
curl https://[domain]/path/to/file.html | jse -u  https://[domain]/path/to/file.html -tui
```

|view|shortcut|description|
|:---:|:---:|:---:|
|all|`CTRL+H`|Toogle help windows|
|all|`CTRL+C`|Exit TUI|
|scripts|`CTRL+G`|Gather script from source attribute for the current item|
|scripts|`CTRL+A`|Gather scripts from all source attributes|
|scripts|`CTRL+U`|Change tu url view|
|url|`TAB`|Change request method (`GET`or `cURL`)|
|url|`ENTER`|Perform request and gather js|
|url|`CTRL+Z`|Don't perform request, go back to scripts view|

[with 💥](https://github.com/ariary/bang/blob/main/EXAMPLES.md#jse)

## Install

### from release

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
* Beautify js? (jse is only for gathering js, maybe it is a better idea to not add other features)


## 🎁

The equivalent with WebAssembly: [`wasme`](https://gist.github.com/ariary/b5a6062b89a62ba33734ba4638b0b20b)
