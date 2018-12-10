package utils

import (
    "time"
)

func DurationAdd(d1, d2 time.Duration) time.Duration {
    return time.Duration(d1.Nanoseconds() + d2.Nanoseconds())
}

func DurationMuti(d time.Duration, x int64) time.Duration {
    return time.Duration(x * d.Nanoseconds())
}

func DurationOneHalf(d time.Duration) time.Duration {
    return time.Duration(d.Nanoseconds() / 2)
}
