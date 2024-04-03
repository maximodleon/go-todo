package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"todo"
)


var filename = "todo.json"

func main() {
  // get filename from env var if set
  if os.Getenv("TODO_FILENAME") != "" {
    filename = os.Getenv("TODO_FILENAME")
  }

  // Define the flags for the cli app
  add := flag.Bool("add", false, "Add task to the todo list")
  list := flag.Bool("list", false, "List all tasks")
  complete := flag.Int("complete", 0, "Item to be completed")
  delete := flag.Int("del", 0, "Item to be deleted")
  flag.Bool("v", false, "Verbose output")
  flag.Bool("skipDone", false, "Verbose output")


  flag.Usage = func() {
    fmt.Fprintf(flag.CommandLine.Output(),
      "%s tool. Developed for me\n", os.Args[0])

    fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024\n")
    fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
    fmt.Fprintln(flag.CommandLine.Output(), "\nTasks can be created using STDIN or using the -add flag")
    fmt.Fprintln(flag.CommandLine.Output(), "Example: echo 'item from STDIN' | todo -add")
    flag.PrintDefaults()
  }

  // parse the arguments
  flag.Parse()

  // define an items list
  l := &todo.List{}

  // read the file
  if err := l.Get(filename); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  switch {
    case *list:
      fmt.Print(l)
    case *complete > 0:
      // Mark the item as completed
      if err := l.Complete(*complete); err !=nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
      }

       // Save to file
      if err := l.Save(filename); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
      }
    case *add:
     // Add it to the list
     t, err := getTask(os.Stdin, flag.Args()...)
     if err != nil {
       fmt.Fprintln(os.Stderr, err)
       os.Exit(1)
     }

     tasks := strings.Split(t, "\n")

     for _, task := range tasks {
       if(len(task) != 0) {
        l.Add(strings.TrimSpace(task))
       }
     }

     if err := l.Save(filename); err != nil {
       fmt.Fprintln(os.Stderr, err)
       os.Exit(1)
     }
    case *delete > 0:
      // Mark the item as deleted
      if err := l.Delete(*delete); err !=nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
      }

       // Save to file
      if err := l.Save(filename); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
      }
    default:
     // invalid flag
     fmt.Fprintln(os.Stderr, "invalid flag")
     os.Exit(1)
  }
}

// task from: arguments or STDIN
func getTask (r io.Reader, args ...string) (string, error) {
  if len(args) > 0 {
    return strings.Join(args, " "), nil
  }

  s := bufio.NewScanner(r)
  output  := ""

  for {
    s.Scan()
    line := s.Text()

    if (len(line) == 0) {
      break;
    }

    output = output + line + "\n"
  }

  if err := s.Err(); err != nil {
    return "", err
  }

  //if len(s.Text()) == 0 {
   // return "", fmt.Errorf("Task cannot be blank")
 // }

  return output, nil
}
