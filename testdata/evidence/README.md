# :mag: evidence
Sample data for forensics processing

Forensics software need to be able to parse and process many different file formats. This repository contains samples of different file formats that can be used to test forensics software. Each file is accompanied by an entry in the [evidence.json](evidence.json) file with some metadata. 

Example entry for this README.md:
```
[
  …
  {
    "name": "README.md", 
    "mime": "text/plain", 
    "generator": "github.com"
  },
  …
]
```
