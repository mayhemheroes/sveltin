package fuzz_sveltin_common

import (
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/sveltinio/sveltin/common"
)

func mayhemit(data []byte) int {

    if len(data) > 2 {
        num := int(data[0])
        data = data[1:]
        fuzzConsumer := fuzz.NewConsumer(data)
        
        switch num {
            
            case 0:
                var testStrings []string
                repeat, _ := fuzzConsumer.GetInt()

                for i := 0; i < repeat; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    testStrings = append(testStrings, temp)
                }

                findString, _ := fuzzConsumer.GetString()

                common.Contains(testStrings, findString)
                return 0

            case 1:
                var a []string
                var b []string
                repeatA, _ := fuzzConsumer.GetInt()
                repeatB, _ := fuzzConsumer.GetInt()

                for i := 0; i < repeatA; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    a = append(a, temp)
                }

                for i := 0; i < repeatB; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    b = append(b, temp)
                }

                common.Difference(a, b)
                return 0

            case 2:
                var testStrings []string
                repeat, _ := fuzzConsumer.GetInt()

                for i := 0; i < repeat; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    testStrings = append(testStrings, temp)
                }

                common.Unique(testStrings)
                return 0

            case 3:
                var a []string
                var b []string
                repeatA, _ := fuzzConsumer.GetInt()
                repeatB, _ := fuzzConsumer.GetInt()

                for i := 0; i < repeatA; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    a = append(a, temp)
                }

                for i := 0; i < repeatB; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    b = append(b, temp)
                }

                common.Union(a, b)
                return 0

            case 4:
                var testStrings []string
                repeat, _ := fuzzConsumer.GetInt()

                for i := 0; i < repeat; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    testStrings = append(testStrings, temp)
                }

                common.RemoveEmpty(testStrings)
                return 0

            case 5:
                m1 := make(map[string]string)
                m2 := make(map[string]string)

                fuzzConsumer.FuzzMap(m1)
                fuzzConsumer.FuzzMap(m2)

                common.UnionMap(m1, m2)
                return 0
        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}