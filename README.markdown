# GoBom

A Golang implementation for Unicode's BOM (Byte Order Mark) detection.

BOM stands for Byte Order Mark. It is a standard by the Unicode organization
to understand the type of encoding is standing in front of us, by placing non
printable chars that explains if an encoding is UTF16 or UTF32, and what is the
endian that it uses.
The standard does not recommend to place a BOM to UTF8, but it supports that as
well.

This library was created for helping me detect if something contain a BOM, and
that's it. It does not do anything else, and there is no plan for anything other
then detecting it.

How does the library works?

The library can use io.Reader, and also "pure" byte slices and rune slices in
order to detect the type of BOM.

Please note:
If a BOM is not detected, then it will return "Unknown".
If a buffer is too small to detect BOM type it also returns "Unknown"

