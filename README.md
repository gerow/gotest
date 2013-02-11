gotest
======

Some useful go testing tools.  (Actually only a deep equality checker/reporter for now).

How to use
======

If you have your gopath setup this is as easy as calling "go get github.com/gerow/gotest" and then importing "github.com/gerow/gotest" within your testing code.

Great, now assume you want to assert a deep equality between two things (let's call them foo and bar) and that this is within a test function so you have a *testing.T called t.  In order to assert deep equality between a and b simply call gotest.AssertDeepEqual(a, b, t).  The AssertDeepEqual call will recursively walk your "things" (interface{}) looking for differences.  On the first difference it spots it will begin making a series of calls to t.Errorf starting with the most specific difference and working its way up to the most general difference.
