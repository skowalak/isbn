# isbn

The `isbn` package inspects and verifies International Standard Book Numbers
(ISBNs) according to the [ISO standard][std].

> [!WARNING]
> This package will in its current version validate British Standard Book
> Numbers (SBN), ISBN-10 and ISBN-13 **only** according to their structure.
> 
> As of now it **will not** validate if an ISBN is actually contained in an
> ISBN range.
> 
> Furthermore it **will not** handle other but similar identifiers like EAN-13
> in general, Digital Object Identifier (DOI), ISBN-A, International Standard
> Music Number (ISMN), International Standard Serial Number (ISSN) or
> International Standard Link Identifier (ISLI).

## Contributions

Issues and PRs are welcome.

## License

The package is licensed under MIT license.

###### Install

```sh
go get github.com/skowalak/isbn
```

###### Documentation 

[![Go Reference](https://pkg.go.dev/badge/github.com/skowalak/isbn.svg)](https://pkg.go.dev/github.com/skowalak/isbn)


[std]: https://www.isbn-international.org/content/isbn-standard
