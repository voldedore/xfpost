# Parse a XenForo2 thread

## Download

```
go get -u github.com/voldedore/xfpost
```

## Usage

Get first page

```
xfpost get https://voz.vn/t/caffe-tai-gia.2639/
```

Get from page 1 to page 5

```
xfpost get https://voz.vn/t/caffe-tai-gia.2639/ -p 5
```

## Dependencies

This tool makes use of the following lib:

- github.com/PuerkitoBio/goquery
- github.com/spf13/cobra

## TODO

[] Configuration for output file
[] Progress showing
