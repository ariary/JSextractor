# JSextractor

Gather all JavaScript code from html webpages from command line.
`js-extractor` take input from `stdin`, search for JavaScript in `<script>` tag, in event handler and in `src` attributes.. 
* Get all JavaScript in output
* Get all JavaScript into a file
* Get only JavaScript between `script` tag

## Output js code into file

This could be useful for performing further action later on JavaScript  like scanning it or beautifying it.
`cat [html_file] | js-exctractor -gather-src -d [domain_of_html]

We use  `-gather-src` to retrieve code from `src` attribute (fetching it). Otherwise it would return only the URL corresponding to the `src` value and thus making the output a non-valid JavaScript script.
When we use `-gather-src` we must also define the domain from which we got the html page, this to fetch `src=/this/is/a/path.js` script.

Also, all informative logs are output to `stderr` to keep only js code in `stdout`
