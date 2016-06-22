# Backgrounder [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/alileza/backgrounder) [![CircleCI](https://circleci.com/gh/alileza/backgrounder.svg?style=shield)](https://circleci.com/gh/alileza/backgrounder)

Simple goroutine manager.

### Example

`import "github.com/alileza/backgrounder"`

```go
var exampleResponse, godocResponse *http.Response

bg := backgrounder.New()

bg.Run(func() error {
    var err error
    exampleResponse, err = http.Get("http://example.com/")
    if err != nil{
        return err
    }
    
    return nil
})
bg.Run(func() error {
    var err error
    godocResponse, err = http.Get("http://godoc.org/")
    if err != nil{
        return err
    }
    
    return nil
})

// Adjust timeout by adding `time.Duration` as the first params of CatchErrs.
// e.g : bg.CatchErrs(time.Second * 5)
// default : 1m
errs := bg.CatchErrs()
if len(errs) != 0 {
    log.Fatal(errs)
}
```

#### Process Time Example
```go
var exampleResponse, godocResponse *http.Response

bg := backgrounder.New()

bg.RunProfile(func() error {
    var err error
    exampleResponse, err = http.Get("http://example.com/")
    if err != nil{
        return err
    }
    
    return nil
}, "get-example.com")
bg.RunProfile(func() error {
    var err error
    godocResponse, err = http.Get("http://godoc.org/")
    if err != nil{
        return err
    }
    
    return nil
}, "get-godoc.org")


errs := bg.CatchErrs()
if len(errs) != 0 {
    log.Fatal(errs)
}
fmt.Println(bg.GetProfiles())
// Output : map[get-godoc.org:2.279357612s get-example.com:579.31927ms]
```
