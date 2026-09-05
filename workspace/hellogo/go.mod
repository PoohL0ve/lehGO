module github.com/PoohL0ve/hellogo

go 1.27.1

// ../mystrings tells go to look in the parent directory of hellogo for the mystrings sibling directory.
replace github.com/PoohL0ve/mystrings v0.0.0 => ../mystrings
require github.com/PoohL0ve/mystrings v0.0.0
