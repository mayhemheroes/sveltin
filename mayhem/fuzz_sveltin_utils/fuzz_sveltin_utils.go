package fuzz_sveltin_utils

import (
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/sveltinio/sveltin/utils"
)

func mayhemit(data []byte) int {

    if len(data) > 2 {
        num := int(data[0])
        data = data[1:]
        fuzzConsumer := fuzz.NewConsumer(data)
        
        switch num {
            
            case 0:
                testURL, _ := fuzzConsumer.GetString()

                utils.IsValidURL(testURL)
                return 0

            case 1:
                testURL, _ := fuzzConsumer.GetString()

                utils.NewGitHubURLParser(testURL)
                return 0

            case 2:
                testNum, _ := fuzzConsumer.GetInt()

                utils.PlusOne(testNum)
                return 0

            case 3:
                testNum, _ := fuzzConsumer.GetInt()

                utils.MinusOne(testNum)
                return 0

            case 4:
                testNumX, _ := fuzzConsumer.GetInt()
                testNumY, _ := fuzzConsumer.GetInt()

                utils.Sum(testNumX, testNumY)
                return 0

            case 5:
                testName, _ := fuzzConsumer.GetString()

                utils.GetNPMClientInfo(testName)
                return 0

            case 6:
                testString, _ := fuzzConsumer.GetString()

                utils.IsEmpty(testString)
                return 0

            case 7:
                var testStrings []string
                repeat, _ := fuzzConsumer.GetInt()

                for i := 0; i < repeat; i++ {

                    temp, _ := fuzzConsumer.GetString()
                    testStrings = append(testStrings, temp)
                }

                utils.IsEmptySlice(testStrings)
                return 0

            case 8:
                testName, _ := fuzzConsumer.GetString()
                testBool, _ := fuzzConsumer.GetBool()

                utils.ToMDFile(testName, testBool)
                return 0

            case 9:
                testName, _ := fuzzConsumer.GetString()

                utils.ToLibFile(testName)
                return 0

            case 10:
                testTitle, _ := fuzzConsumer.GetString()

                utils.ToTitle(testTitle)
                return 0

            case 11:
                testURL, _ := fuzzConsumer.GetString()

                utils.ToURL(testURL)
                return 0

            case 12:
                testTxt, _ := fuzzConsumer.GetString()

                utils.ToSlug(testTxt)
                return 0

            case 13:
                testTxt, _ := fuzzConsumer.GetString()

                utils.ToSnakeCase(testTxt)
                return 0

            case 14:
                fullStr, _ := fuzzConsumer.GetString()
                replaceStr, _ := fuzzConsumer.GetString()

                utils.ToBasePath(fullStr, replaceStr)
                return 0

            case 15:
                testTxt, _ := fuzzConsumer.GetString()

                utils.ToVariableName(testTxt)
                return 0

            case 16:
                testTxt, _ := fuzzConsumer.GetString()

                utils.ReplaceIfNested(testTxt)
                return 0
                
            case 17:
                testTxt, _ := fuzzConsumer.GetString()

                utils.ConvertJSStringToStringArray(testTxt)
                return 0
        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}