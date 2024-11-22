
# termo-beater

What is this?

This is a simple program to solve a question I've had: What is the best initial guess for [termo](https://term.ooo) (a portuguese clone of wordle)? 

The possible word list used for this is in `all_words.txt`, and is taken from the javascript files that 
comes with the website.


## Theory

Note: The program assumes the termo backend picks the word of the day at random uniformly.

TODO: Explain entropy and stuff

## Running the program

To run this program, just do `go run main.go`. The program will create a text file `output.txt` containing 
two columns, the words from the original word list, and the uncertaintly level of that word.
Better guesses have lower uncertainties, and worse guesses have higher uncertainties.

Running the following lines of [nushell](https://github.com/nushell/nushell) code, it is possible to find the best and worse words.

```
let results = 
    open output.txt 
    | from ssv -n 
    | rename word entropy 
    | update entropy { |it| $it.entropy | into float } 
    | sort-by entropy 

# Best word
$results | first

# Output:
# ┏━━━━━━━━━┳━━━━━━━┓
# ┃ word    ┃ tirao ┃
# ┃ entropy ┃ 7.36  ┃
# ┗━━━━━━━━━┻━━━━━━━┛

# Worst word
$results | last

# Output:
# ┏━━━━━━━━━┳━━━━━━━┓
# ┃ word    ┃ kayak ┃
# ┃ entropy ┃ 11.04 ┃
# ┗━━━━━━━━━┻━━━━━━━┛
```

So, if my implementation is correct (I think my mask computation routine is a bit sus), 
the best first guess is `tirão`, and the worst first guess is `kayak`.



