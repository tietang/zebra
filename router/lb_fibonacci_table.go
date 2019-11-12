package router

import (
    "errors"
    "math"
)

func ceilingGet(tables [][]int64, key int64, smoothness bool, isTenMillisecond bool) int64 {
    index := 1
    if isTenMillisecond {
        index = 0
    }
    for i, v := range tables {
        k := v[index]
        v := v[2]
        if key <= k {
            if i == 0 || !smoothness {
                return v
            }
            kv0 := tables[i-1]
            k0 := kv0[index]
            v0 := kv0[2]
            diffk := k - k0
            diffv := v - v0
            dv := float64(diffv) * (float64(key-k0) / float64(diffk))
            //math.Ceil
            r := int64(math.Ceil(float64(v0) + dv))
            return r

        }
    }
    return 1
}

var FibonacciTable = [][]int64{
    {0, 1, 1000},
    {1, 10, 996},
    {2, 20, 992},
    {3, 30, 988},
    {5, 50, 981},
    {8, 80, 970},
    {13, 130, 951},
    {21, 210, 922},
    {34, 340, 874},
    {55, 550, 798},
    {89, 890, 680},
    {144, 1440, 508},
    {233, 2330, 285},
    {377, 3770, 48},
    {610, 6100, 10},
    {987, 9870, 6},
    {1597, 15970, 1},
}

//func GetCeiling(key int64) int64 {
//    value := ceilingGet(FibonacciTable, key, true, false)
//    return value
//}

func GetCeiling(key, x int64) int64 {
    if x%10 != 0 {
        panic(errors.New("value must be 10 times"))
    }
    value := ceilingGet(FibonacciTable, 10*key/x, true, false)
    return value
}

func GetCeilingByTenMillisecond(key int64) int64 {
    value := ceilingGet(FibonacciTable, key, true, true);
    return value;
}
