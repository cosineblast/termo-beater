
package main

import (
    "fmt"
    "math"
    "os"
    "bufio"
    "log"
    "sync"
    "runtime"
)

type Color int

const (
    black = 0
    yellow = 1
    green = 2
)

type Mask [5]Color


func compare(answer string, guess string) Mask {

    var mask Mask

    for i := 0; i < 5; i++ {
        if answer[i] == guess[i] {
            mask[i] = green
        }
    }

    for i := 0; i < 5; i++ {
        if mask[i] == black {
            for j := 0; j < 5; j++ {
                if guess[i] == answer[j] && mask[j] == black {
                    mask[i] = yellow
                    break
                }
            }
        }
    }

    return mask
}

type WordDistribution map[string]float64

func computeEntropy(words WordDistribution) float64 {
    result := 0.0

    for _, value := range words {
        result += value * math.Log2(value)
    }

    return -result
}

// Given a random variable X, a guess g and a mask m,
// it returns Y such that for any answer w,
// P(Y = w) = P(X = w | compare(X, g) = m)
func restrictDistribution(words WordDistribution, guess string, mask Mask) WordDistribution {

    result := make(WordDistribution)
    total := 0.0

    for key, value := range words {
        if compare(key, guess) == mask {
            result[key] = value
            total += value
        }
    }

    // the sum may be less than one now, so we must adjust for that
    for key, value := range result {
        result[key] = value / total
    }

    return result
}


func measureResultingUncertainty(words WordDistribution, guess string, mask Mask) float64 {
    return computeEntropy(restrictDistribution(words, guess, mask))
}


func measureAverageResultingUncertainty(words WordDistribution, guess string) float64 {
    maskDistribution := make(map[Mask]float64)

    for word, probablity := range words {
        mask := compare(word, guess)

        maskDistribution[mask] += probablity
    }

    result := 0.0

    for mask, maskProbability := range maskDistribution {
        result += maskProbability * measureResultingUncertainty(words, guess, mask)
    }

    return result
}


func getWords(name string) ([]string, error) {

    f,err := os.Open(name)

    if err != nil {
        return nil, err
    }

    defer f.Close()

    var result []string

    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        result = append(result, scanner.Text())
    }

    return result, nil

}

func saveOutput(output chan struct {string; float64}, file *os.File) {
    for value := range output {
        fmt.Println("uncertainty(",value.string,") =", value.float64)
        file.WriteString(fmt.Sprintf("%s  %f\n", value.string, value.float64))
    }
}

func uniformDistribution(words []string) WordDistribution {
    distribution := make(WordDistribution)

    for _, word := range words {
        distribution[word] = 1/float64(len(words))
    }

    return distribution
}

const wordsFile = "all_words.txt"

func main() {
    fmt.Println("Loading words from", wordsFile, "...")

    words, err := getWords(wordsFile)

    if err != nil {
        log.Fatal(err)
        return
    }

    fmt.Println("Finished loading", len(words), "words =]")

    // we are going to assume all words are equally probable
    distribution := uniformDistribution(words)

    file, err := os.Create("output.txt")

    if err != nil {
        log.Fatal(err)
        return
    }

    defer file.Close()

    output := make(chan struct{string; float64}, 32)

    go saveOutput(output, file)

    todo := make(chan string, len(words))
    for _, word := range words {
        todo <- word
    }
    close(todo)

    var group sync.WaitGroup

    for i := 0; i < runtime.NumCPU(); i++ {
        group.Add(1)

        go func() {
            defer group.Done()

            for word := range todo {
                entropy := measureAverageResultingUncertainty(distribution, word)
                output <- struct {string; float64}{word, entropy}
            }
        }()
    }

    group.Wait()
}
