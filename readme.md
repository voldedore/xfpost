# Parse a XenForo2 thread

## Download

```
go get -u github.com/voldedore/xfpost
```

## Usage

Get first page

```
xfpost get https://xenforo.com/t/thread-url.123/
```

Get from page 1 to page 5

```
xfpost get https://xenforo.com/t/thread-url.123/ -p 5
```

Write to specific file

```
xfpost get https://xenforo.com/t/thread-url.123/ -o 1-5.json
```

Help

```
xfpost -h
```

## Dependencies

This tool makes use of the following lib:

- github.com/PuerkitoBio/goquery
- github.com/spf13/cobra
