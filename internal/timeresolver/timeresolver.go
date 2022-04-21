package timeresolver

import "time"

type TimeResolver interface {
	Now() time.Time
}

type TimeResolverProd struct{}

func (t TimeResolverProd) Now() time.Time {
	return time.Now()
}

//Mocked timeresolver resolver
type TimeResolverMock struct {
	Time time.Time
}

func (t TimeResolverMock) Now() time.Time {
	return t.Time
}
