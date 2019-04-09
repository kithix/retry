# retry

Just another go retry package

### Usage
Call `retry.Do` with:
    - The thing you want retried, it must be of signature `func() error`
    - Some retry strategy `func(error) bool`. This is a way of determining what to do in the event of an error 

```
var data []byte
err := retry.Do(func() error{
    var err error
    data, err = Something()
    return err
}, retry.Limit(5))
```

#### Delay between attempts
```
var data []byte
err := retry.Do(func() error {
    var err error
    data, err = Something()
    return err
}, retry.WithWait(retry.Limit(5), 1*time.Second))
```

#### Custom retry strategies

A retry strategy is simply `func(error) bool`
```
var data []byte
err := retry.Do(func() error {
    var err error
    data, err = Something()
    return err
}, retry.WithLimit(func(err error) bool {
    if strings.Contains(err.Error(), "500") {
        fmt.Println("Error received:", err, "Retrying")
        return true
    }
    fmt.Println("Error received:", err, "Skipping")
    return false
}, 5))
```