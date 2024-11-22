
# termo-beater

What is this?

This is a simple program to solve a question I've had: What is the best initial guess for [termo](https://term.ooo) (a portuguese clone of wordle)? 

I've taken the possible word list from the javascript files that comes with termo website, 
and I assume the termo backend chooses the word of the day uniformly.

To run this program, just do `go run main.go`. The program will create a text file `output.txt` containing 
two columns, the words from the original word list, and the uncertaintly level of that word.
Better guesses have lower uncertainties, and worse guesses have higher uncertainties.

## Theory

TODO: Explain entropy and stuff

## Running the program

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

So, if my implementation is correct (I think my mask routine is a bit sus), the best first guess is `tirão`, and the worst first guess is `kayak`.



