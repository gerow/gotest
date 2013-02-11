gotest
======

Some useful go testing tools.  (Actually only a deep equality checker/reporter for now).

How to install
======

If you have your gopath setup this is as easy as calling this from your terminal:

```bash
go get github.com/gerow/gotest
```

And then adding:

```go
import "github.com/gerow/gotest"
```
to your testing code.  If you want to update later just call this from your terminal:

```bash
go get -u github.com/gerow/gotest
```

How to use
==========
Excellent, you have gotest installed! Right now there's only one thing that it does: deep equality checking where differences are noted as calls to the Errorf method of testing.T. This is mostly aped from the actual DeepEqual implementation at reflect.DeepEqual.  The problem with that one is that it only tells you whether or not two things are equal without telling you exactly what isn't equal. This is great when your tests all pass, but when an equality check fails it's hard to know exactly what went wrong.

AssertDeepEqual
---------------
Great, now assume you want to assert a deep equality between two things (let's call them foo and bar) and that this is within a test function so you have a *testing.T called t.  In order to assert deep equality between a and b simply call:

```go
gotest.AssertDeepEqual(a, b, t)
```
The AssertDeepEqual call will recursively walk your "things" (interface{}) looking for differences.  On the first difference it spots it will begin making a series of calls to t.Errorf starting with the most specific difference and working its way up to the most general difference.

An actual use can be found [here](https://github.com/gerow/boop/blob/dev/src/boop/boop_config_test.go).  One function from there is reproduced below:
```go
func TestLoadConfigFromFile(t *testing.T) {
  t.Logf("Starting TestLoadConfigFromFileWithDefaults")
	const filename = "test.config.json"

	var exp Config

	exp.Port = 9180
	exp.OnlyAllowIps = []string{"192.168.1.1", "192.169.1.1"}
	exp.Commands = []Command{
		Command{
			"GET /here/is/get/path",
			"touch hello",
			[]string{"192.188.1.5", "192.162.1.55"},
			120},
		Command{
			"POST /here/is/post/path",
			"touch hello2",
			[]string{"192.188.2.5", "192.162.2.55"},
			320},
		Command{
			"DELETE /here/is/delete/path",
			"sleep 200",
			nil,
			0}}

	config, err := LoadConfigFromFile(filename)

	if err != nil {
		t.Errorf("Got error from LoadConfigFromFile:", err)
	}

	gotest.AssertDeepEqual(exp, *config, t)

}
```
